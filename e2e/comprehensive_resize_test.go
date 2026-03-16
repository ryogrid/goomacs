package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestComprehensiveResize(t *testing.T) {

	// Helper to create a multi-line file for resize tests.
	makeResizeFile := func(t *testing.T, lines int) string {
		t.Helper()
		dir, err := os.MkdirTemp("/tmp", "g")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		t.Cleanup(func() { os.RemoveAll(dir) })
		var content strings.Builder
		for i := 1; i <= lines; i++ {
			fmt.Fprintf(&content, "line %d content here\n", i)
		}
		path := filepath.Join(dir, "r.go")
		if err := os.WriteFile(path, []byte(content.String()), 0o644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		return path
	}

	// resizeAndRedraw resizes the terminal and sends C-g (a no-op in normal mode)
	// to force a full redraw, since tmux resize-window may not immediately
	// propagate content to the capture buffer.
	resizeAndRedraw := func(h *Harness, width, height int) {
		h.Resize(width, height)
		time.Sleep(300 * time.Millisecond)
		h.SendKeys("C-g")
		time.Sleep(300 * time.Millisecond)
	}

	// --- Resize from 80x24 to 40x12 ---
	t.Run("Shrink_80x24_to_40x12", func(t *testing.T) {
		path := makeResizeFile(t, 20)
		h := newHarness(t, path)

		if err := h.WaitForContent("line 1", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Verify initial state at 80x24
		h.AssertScreenContains(t, "line 1 content here")
		h.AssertStatusBar(t, "r.go")

		// Resize to 40x12 and force redraw
		resizeAndRedraw(h, 40, 12)

		// After resize, content should still be visible
		screen := h.CapturePane()
		if !strings.Contains(screen, "line 1") {
			t.Errorf("after shrink: expected 'line 1' visible, got:\n%s", screen)
		}

		// Status bar should still be present (at new row h.height-2)
		lines := h.Capture()
		statusRow := h.height - 2
		if statusRow >= 0 && statusRow < len(lines) {
			if !strings.Contains(lines[statusRow], "r.go") {
				t.Errorf("after shrink: status bar should contain filename, got %q", lines[statusRow])
			}
		}

		// Editor still functional — can type
		h.SendKeys("X")
		time.Sleep(200 * time.Millisecond)
		h.AssertScreenContains(t, "Xline 1")
	})

	// --- Resize from 80x24 to 200x50 ---
	t.Run("Expand_80x24_to_200x50", func(t *testing.T) {
		path := makeResizeFile(t, 20)
		h := newHarness(t, path)

		if err := h.WaitForContent("line 1", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Resize to 200x50 and force redraw
		resizeAndRedraw(h, 200, 50)

		// Content should be visible — more lines visible now
		screen := h.CapturePane()
		if !strings.Contains(screen, "line 1") {
			t.Errorf("after expand: expected 'line 1' visible, got:\n%s", screen)
		}
		// With 50 rows, all 20 lines should be visible
		if !strings.Contains(screen, "line 20") {
			t.Errorf("after expand: expected 'line 20' visible with larger terminal, got:\n%s", screen)
		}

		// Status bar should still be present
		lines := h.Capture()
		statusRow := h.height - 2
		if statusRow >= 0 && statusRow < len(lines) {
			if !strings.Contains(lines[statusRow], "r.go") {
				t.Errorf("after expand: status bar should contain filename, got %q", lines[statusRow])
			}
		}

		// Editor still functional
		h.SendKeys("Y")
		time.Sleep(200 * time.Millisecond)
		h.AssertScreenContains(t, "Yline 1")
	})

	// --- Resize with split windows (C-x 2) ---
	t.Run("Resize_WithVerticalSplit", func(t *testing.T) {
		path := makeResizeFile(t, 20)
		h := newHarness(t, path)

		if err := h.WaitForContent("line 1", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Split vertically
		h.SendKeys("C-x")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("2")
		time.Sleep(300 * time.Millisecond)

		// Verify two status bars before resize
		lines := h.Capture()
		statusCount := 0
		for _, line := range lines {
			if strings.Contains(line, "r.go") && strings.Contains(line, "Line ") {
				statusCount++
			}
		}
		if statusCount < 2 {
			t.Errorf("before resize: expected 2 status bars with 'r.go', found %d", statusCount)
		}

		// Resize to 80x40 (larger) and force redraw
		resizeAndRedraw(h, 80, 40)

		// Both windows should still show content and status bars
		lines = h.Capture()
		statusCount = 0
		for _, line := range lines {
			if strings.Contains(line, "r.go") && strings.Contains(line, "Line ") {
				statusCount++
			}
		}
		if statusCount < 2 {
			t.Errorf("after resize to 80x40: expected 2 status bars, found %d", statusCount)
		}

		// Content should still be visible
		screen := h.CapturePane()
		if !strings.Contains(screen, "line 1") {
			t.Errorf("after resize with split: expected 'line 1' visible, got:\n%s", screen)
		}

		// Resize smaller and force redraw
		resizeAndRedraw(h, 60, 16)

		// Both windows should still exist with smaller heights
		lines = h.Capture()
		statusCount = 0
		for _, line := range lines {
			if strings.Contains(line, "r.go") && strings.Contains(line, "Line ") {
				statusCount++
			}
		}
		if statusCount < 2 {
			t.Errorf("after resize to 60x16: expected 2 status bars, found %d", statusCount)
		}

		// Editor still functional
		h.SendKeys("Z")
		time.Sleep(200 * time.Millisecond)
		h.AssertScreenContains(t, "Zline 1")
	})

	// --- Resize during search mode ---
	t.Run("Resize_DuringSearch", func(t *testing.T) {
		path := makeResizeFile(t, 20)
		h := newHarness(t, path)

		if err := h.WaitForContent("line 1", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Enter search mode
		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("line 5")
		time.Sleep(300 * time.Millisecond)

		// Verify search is active
		if err := h.WaitForContent("I-search:", 2*time.Second); err != nil {
			t.Fatalf("search mode not active: %v", err)
		}

		// Resize during search — do NOT send C-g (that would cancel search)
		// Instead, resize and then send C-s (search forward again) to trigger a redraw
		// while staying in search mode
		h.Resize(60, 18)
		time.Sleep(300 * time.Millisecond)
		// Send C-s to advance search (also triggers redraw) — search remains active
		h.SendKeys("C-s")
		time.Sleep(300 * time.Millisecond)

		// Search state should be preserved — message line should still show search prompt
		lines := h.Capture()
		msgLine := lines[len(lines)-1]
		if !strings.Contains(msgLine, "I-search:") {
			t.Errorf("after resize during search: expected search prompt on message line, got %q", msgLine)
		}
		if !strings.Contains(msgLine, "line 5") {
			t.Errorf("after resize during search: expected query 'line 5' in prompt, got %q", msgLine)
		}

		// Accept search and verify cursor moved to match
		h.SendKeys("Enter")
		time.Sleep(200 * time.Millisecond)

		// Content should show line 5
		h.AssertScreenContains(t, "line 5 content here")
	})

	// --- Resize during minibuffer input ---
	t.Run("Resize_DuringMinibuffer", func(t *testing.T) {
		path := makeResizeFile(t, 10)
		h := newHarness(t, path)

		if err := h.WaitForContent("line 1", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Enter minibuffer via C-x C-f
		h.SendKeys("C-x")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("C-f")
		time.Sleep(200 * time.Millisecond)

		// Type some input
		h.SendKeys("somefile.txt")
		time.Sleep(200 * time.Millisecond)

		// Verify minibuffer is active
		if err := h.WaitForContent("Find file:", 2*time.Second); err != nil {
			t.Fatalf("minibuffer not active: %v", err)
		}

		// Resize during minibuffer — send C-e (move to end, no-op) to trigger redraw
		// while staying in minibuffer mode
		h.Resize(100, 30)
		time.Sleep(300 * time.Millisecond)
		h.SendKeys("C-e")
		time.Sleep(300 * time.Millisecond)

		// Minibuffer should still show the prompt and input
		lines := h.Capture()
		msgLine := lines[len(lines)-1]
		if !strings.Contains(msgLine, "Find file:") {
			t.Errorf("after resize: expected 'Find file:' prompt, got %q", msgLine)
		}
		if !strings.Contains(msgLine, "somefile.txt") {
			t.Errorf("after resize: expected input 'somefile.txt' preserved, got %q", msgLine)
		}

		// Cancel and verify editor still functional
		h.SendKeys("C-g")
		time.Sleep(200 * time.Millisecond)

		h.AssertScreenContains(t, "line 1 content here")
	})
}
