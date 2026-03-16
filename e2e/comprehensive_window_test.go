package e2e

import (
	"strings"
	"testing"
	"time"
)

// countWindowStatusBars counts lines that look like status bars and contain the given name.
// A status bar contains "Line " (position indicator) and the given name.
// Active bars use spaces (reverse video), inactive use dashes.
func countWindowStatusBars(lines []string, height int, name string) int {
	count := 0
	for i := 0; i < height-1; i++ {
		line := lines[i]
		if strings.Contains(line, name) && strings.Contains(line, "Line ") {
			count++
		}
	}
	return count
}

func TestComprehensiveWindow(t *testing.T) {
	// --- C-x 2: vertical split, both windows show same buffer, two status lines ---
	t.Run("VerticalSplit_BothShowSameBuffer", func(t *testing.T) {
		path := shortTempFile(t, "vsplit.txt", "vertical split content\nline two\nline three")
		h := newHarness(t, path)
		if err := h.WaitForContent("vsplit.txt", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// C-x 2 to split vertically
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("2")
		time.Sleep(300 * time.Millisecond)

		lines := h.Capture()

		// Both windows should show the same buffer content
		screen := h.CapturePane()
		count := strings.Count(screen, "vertical split content")
		if count < 2 {
			t.Errorf("expected buffer content in both windows, found %d occurrences\nScreen:\n%s", count, screen)
		}

		// Should have 2 status bars — count lines containing both filename and "Line "
		statusCount := countWindowStatusBars(lines, h.height, "vsplit.txt")
		if statusCount < 2 {
			t.Errorf("expected 2 status bars showing vsplit.txt, got %d\nScreen:\n%s", statusCount, screen)
		}
	})

	// --- C-x 3: horizontal split, both windows side-by-side, vertical separator ---
	t.Run("HorizontalSplit_SideBySide", func(t *testing.T) {
		path := shortTempFile(t, "hsplit.txt", "horizontal content")
		h := newHarness(t, path)
		if err := h.WaitForContent("hsplit.txt", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// C-x 3 to split horizontally
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("3")
		time.Sleep(300 * time.Millisecond)

		screen := h.CapturePane()

		// Verify vertical separator drawn
		if !strings.Contains(screen, "│") {
			t.Errorf("expected vertical separator '│' after horizontal split\nScreen:\n%s", screen)
		}

		// Both sides should show the buffer content
		count := strings.Count(screen, "horizontal content")
		if count < 2 {
			t.Errorf("expected buffer content in both panes, found %d occurrences\nScreen:\n%s", count, screen)
		}
	})

	// --- C-x o: cycles through windows ---
	t.Run("CycleWindows_CxO", func(t *testing.T) {
		path := shortTempFile(t, "cyc.txt", "cycle content")
		h := newHarness(t, path)
		if err := h.WaitForContent("cyc.txt", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Split vertically
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("2")
		time.Sleep(300 * time.Millisecond)

		// Capture initial state
		before := h.CapturePane()

		// C-x o to switch
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("o")
		time.Sleep(300 * time.Millisecond)

		// Screen should change (active/inactive status bars swap)
		after := h.CapturePane()
		if before == after {
			t.Error("screen did not change after C-x o; expected active window indicator to swap")
		}

		// C-x o again should return to original
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("o")
		time.Sleep(300 * time.Millisecond)

		restored := h.CapturePane()
		if restored != before {
			t.Error("screen did not return to original after two C-x o cycles")
		}
	})

	// --- C-x 0: closes current window ---
	t.Run("CloseWindow_Cx0", func(t *testing.T) {
		path := shortTempFile(t, "cls.txt", "close me")
		h := newHarness(t, path)
		if err := h.WaitForContent("cls.txt", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Split
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("2")
		time.Sleep(300 * time.Millisecond)

		lines := h.Capture()
		if countStatusBars(lines, h.height) < 2 {
			t.Fatal("split did not create 2 windows")
		}

		// C-x 0 to close current window
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("0")
		time.Sleep(300 * time.Millisecond)

		lines = h.Capture()
		if countStatusBars(lines, h.height) != 1 {
			t.Errorf("expected 1 status bar after C-x 0, got %d\nScreen:\n%s",
				countStatusBars(lines, h.height), h.CapturePane())
		}

		// Content should still be visible
		h.AssertScreenContains(t, "close me")
	})

	// --- C-x 0 with only one window ---
	t.Run("CloseWindow_OnlyOne_NoOp", func(t *testing.T) {
		h := newHarness(t)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		before := h.CapturePane()

		// C-x 0 with single window — should be a no-op
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("0")
		time.Sleep(300 * time.Millisecond)

		after := h.CapturePane()

		// Verify still one window, no crash
		lines := h.Capture()
		if countStatusBars(lines, h.height) != 1 {
			t.Errorf("expected 1 status bar, got %d", countStatusBars(lines, h.height))
		}

		// Editor should still be functional — type something
		h.SendKeys("x")
		time.Sleep(100 * time.Millisecond)
		h.AssertScreenContains(t, "x")

		_ = before
		_ = after
	})

	// --- C-x 1: closes all other windows ---
	t.Run("CloseOtherWindows_Cx1", func(t *testing.T) {
		path := shortTempFile(t, "one.txt", "only me remains")
		h := newHarness(t, path)
		if err := h.WaitForContent("one.txt", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Split twice to create 3 windows
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("2")
		time.Sleep(300 * time.Millisecond)

		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("2")
		time.Sleep(300 * time.Millisecond)

		lines := h.Capture()
		if countStatusBars(lines, h.height) < 2 {
			t.Fatal("expected multiple windows after splits")
		}

		// C-x 1 to close all other windows
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("1")
		time.Sleep(300 * time.Millisecond)

		lines = h.Capture()
		if countStatusBars(lines, h.height) != 1 {
			t.Errorf("expected 1 status bar after C-x 1, got %d\nScreen:\n%s",
				countStatusBars(lines, h.height), h.CapturePane())
		}

		h.AssertScreenContains(t, "only me remains")
	})

	// --- Editing in split windows: edit in one, visible in other showing same buffer ---
	t.Run("EditInSplitWindow_PropagatesContent", func(t *testing.T) {
		path := shortTempFile(t, "prop.txt", "original text")
		h := newHarness(t, path)
		if err := h.WaitForContent("original text", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Split vertically — both windows show same buffer
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("2")
		time.Sleep(300 * time.Millisecond)

		// Type text in the active window (top window after split)
		h.SendKeys("Z")
		time.Sleep(200 * time.Millisecond)

		// Both windows should show the edited content
		screen := h.CapturePane()
		count := strings.Count(screen, "Zoriginal text")
		if count < 2 {
			t.Errorf("expected edited content visible in both windows, found %d occurrences\nScreen:\n%s", count, screen)
		}
	})

	// --- Independent scroll: different buffers in split windows scroll independently ---
	t.Run("IndependentScroll_SplitWindows", func(t *testing.T) {
		// Create two files — one with enough lines to scroll, one short
		var longContent strings.Builder
		for i := 1; i <= 50; i++ {
			if i > 1 {
				longContent.WriteByte('\n')
			}
			longContent.WriteString(strings.Repeat("x", 10))
		}
		longLines := strings.Split(longContent.String(), "\n")
		longLines[0] = "TOPMARKER"
		longPath := shortTempFile(t, "long.txt", strings.Join(longLines, "\n"))
		shortPath := shortTempFile(t, "short.txt", "SHORT CONTENT\nline 2")

		h := newHarness(t, longPath)
		if err := h.WaitForContent("TOPMARKER", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Open the short file too
		openFile(h, shortPath)
		time.Sleep(200 * time.Millisecond)
		h.AssertScreenContains(t, "SHORT CONTENT")

		// Switch back to the long file
		switchBuffer(h, longPath)
		time.Sleep(300 * time.Millisecond)
		h.AssertScreenContains(t, "TOPMARKER")

		// Split vertically — both windows show long.txt
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("2")
		time.Sleep(300 * time.Millisecond)

		// Switch the bottom window to short.txt
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("o")
		time.Sleep(200 * time.Millisecond)

		switchBuffer(h, shortPath)
		time.Sleep(300 * time.Millisecond)

		// Verify: one window shows SHORT CONTENT, one shows TOPMARKER
		screen := h.CapturePane()
		if !strings.Contains(screen, "SHORT CONTENT") {
			t.Errorf("expected SHORT CONTENT in one window\nScreen:\n%s", screen)
		}
		if !strings.Contains(screen, "TOPMARKER") {
			t.Errorf("expected TOPMARKER in other window\nScreen:\n%s", screen)
		}

		// Switch back to the long.txt window and scroll down
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("o")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("C-v")
		time.Sleep(300 * time.Millisecond)

		// After scrolling long.txt, SHORT CONTENT should still be visible in the other window
		screen = h.CapturePane()
		if !strings.Contains(screen, "SHORT CONTENT") {
			t.Errorf("expected SHORT CONTENT still visible after scrolling other window\nScreen:\n%s", screen)
		}
		// TOPMARKER should be scrolled away from the long.txt window
		if strings.Contains(screen, "TOPMARKER") {
			t.Errorf("expected TOPMARKER to be scrolled away, but still visible\nScreen:\n%s", screen)
		}
	})
}
