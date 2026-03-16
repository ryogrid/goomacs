package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// shortTempFile creates a file in a short temp dir to avoid status bar truncation.
func shortTempFile(t *testing.T, name, content string) string {
	t.Helper()
	dir, err := os.MkdirTemp("/tmp", "g")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	return path
}

func TestComprehensiveBuffer(t *testing.T) {
	// --- C-x C-f with existing file ---
	t.Run("FindFile_ExistingFile", func(t *testing.T) {
		path := shortTempFile(t, "exist.txt", "hello from file")
		h := newHarness(t)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		openFile(h, path)

		h.AssertScreenContains(t, "hello from file")
		h.AssertStatusBar(t, "exist.txt")
	})

	// --- C-x C-f with nonexistent file ---
	t.Run("FindFile_NonexistentFile", func(t *testing.T) {
		dir, err := os.MkdirTemp("/tmp", "g")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		t.Cleanup(func() { os.RemoveAll(dir) })
		newPath := filepath.Join(dir, "newfile.txt")

		h := newHarness(t)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		openFile(h, newPath)

		h.AssertMessageLine(t, "(New file)")
		h.AssertStatusBar(t, "newfile.txt")
	})

	// --- C-x C-f with same file already open: switches to existing buffer ---
	t.Run("FindFile_AlreadyOpen_NoDuplicate", func(t *testing.T) {
		path := shortTempFile(t, "dup.txt", "unique content")
		h := newHarness(t)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		openFile(h, path)
		h.AssertScreenContains(t, "unique content")

		// Switch back to scratch
		switchBuffer(h, "*scratch*")
		time.Sleep(300 * time.Millisecond)
		h.AssertStatusBar(t, "*scratch*")

		// Open same file again
		openFile(h, path)

		h.AssertScreenContains(t, "unique content")
		h.AssertMessageLine(t, "Switch to buffer")

		// Verify no duplicate in buffer list
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-b")
		time.Sleep(300 * time.Millisecond)

		screen := h.CapturePane()
		count := strings.Count(screen, "dup.txt")
		if count != 1 {
			t.Errorf("expected exactly 1 occurrence of 'dup.txt' in buffer list, got %d\n%s", count, screen)
		}
	})

	// --- C-x b with default (previous buffer) ---
	t.Run("SwitchBuffer_Default_PreviousBuffer", func(t *testing.T) {
		path := shortTempFile(t, "sw.txt", "file content here")
		h := newHarness(t)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		openFile(h, path)
		h.AssertStatusBar(t, "sw.txt")

		// C-x b Enter (empty input) switches to previous buffer (*scratch*)
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("b")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("Enter")
		time.Sleep(300 * time.Millisecond)

		h.AssertStatusBar(t, "*scratch*")

		// Switch back
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("b")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("Enter")
		time.Sleep(300 * time.Millisecond)

		h.AssertStatusBar(t, "sw.txt")
		h.AssertScreenContains(t, "file content here")
	})

	// --- C-x b with typed buffer name ---
	t.Run("SwitchBuffer_TypedName", func(t *testing.T) {
		path := shortTempFile(t, "named.txt", "named buffer content")
		h := newHarness(t)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		openFile(h, path)
		h.AssertStatusBar(t, "named.txt")

		switchBuffer(h, "*scratch*")
		time.Sleep(300 * time.Millisecond)
		h.AssertStatusBar(t, "*scratch*")

		switchBuffer(h, path)
		time.Sleep(300 * time.Millisecond)
		h.AssertStatusBar(t, "named.txt")
		h.AssertScreenContains(t, "named buffer content")
	})

	// --- C-x k on unmodified buffer: closes without confirmation ---
	t.Run("KillBuffer_Unmodified_NoConfirmation", func(t *testing.T) {
		path := shortTempFile(t, "kill.txt", "to be killed")
		h := newHarness(t)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		openFile(h, path)
		h.AssertStatusBar(t, "kill.txt")

		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("k")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("Enter")
		time.Sleep(300 * time.Millisecond)

		h.AssertStatusBar(t, "*scratch*")
		h.AssertMessageLine(t, "Killed buffer")
	})

	// --- C-x k on modified buffer: y confirms close ---
	t.Run("KillBuffer_Modified_YConfirms", func(t *testing.T) {
		path := shortTempFile(t, "km.txt", "original")
		h := newHarness(t)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		openFile(h, path)
		h.AssertStatusBar(t, "km.txt")

		// Modify
		h.SendKeys("X")
		time.Sleep(100 * time.Millisecond)
		h.AssertStatusBar(t, "[Modified]")

		// C-x k Enter
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("k")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("Enter")
		time.Sleep(300 * time.Millisecond)

		h.AssertMessageLine(t, "Buffer modified; kill anyway? (y/n)")

		h.SendKeys("y")
		time.Sleep(300 * time.Millisecond)

		h.AssertStatusBar(t, "*scratch*")
		h.AssertMessageLine(t, "Killed buffer")
	})

	// --- C-x k on modified buffer: n cancels close ---
	t.Run("KillBuffer_Modified_NCancels", func(t *testing.T) {
		path := shortTempFile(t, "kn.txt", "keep me")
		h := newHarness(t)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		openFile(h, path)
		h.AssertStatusBar(t, "kn.txt")

		// Modify
		h.SendKeys("Z")
		time.Sleep(100 * time.Millisecond)
		h.AssertStatusBar(t, "[Modified]")

		// C-x k Enter
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("k")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("Enter")
		time.Sleep(300 * time.Millisecond)

		h.AssertMessageLine(t, "Buffer modified; kill anyway? (y/n)")

		h.SendKeys("n")
		time.Sleep(300 * time.Millisecond)

		h.AssertStatusBar(t, "kn.txt")
		h.AssertScreenContains(t, "Zkeep me")
	})

	// --- C-x C-b: buffer list shows all open buffers with correct markers ---
	t.Run("BufferList_ShowsAllBuffers", func(t *testing.T) {
		path1 := shortTempFile(t, "b1.txt", "file one")
		path2 := shortTempFile(t, "b2.txt", "file two")

		h := newHarness(t)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		openFile(h, path1)
		h.AssertStatusBar(t, "b1.txt")

		openFile(h, path2)
		h.AssertStatusBar(t, "b2.txt")

		// Modify buf2
		h.SendKeys("M")
		time.Sleep(100 * time.Millisecond)

		// Open buffer list
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-b")
		time.Sleep(300 * time.Millisecond)

		h.AssertStatusBar(t, "*Buffer List*")

		screen := h.CapturePane()

		if !strings.Contains(screen, "*scratch*") {
			t.Errorf("buffer list should contain *scratch*\n%s", screen)
		}
		if !strings.Contains(screen, "b1.txt") {
			t.Errorf("buffer list should contain b1.txt\n%s", screen)
		}
		if !strings.Contains(screen, "b2.txt") {
			t.Errorf("buffer list should contain b2.txt\n%s", screen)
		}

		// buf2 should have '*' for modified
		lines := h.Capture()
		foundBuf2Modified := false
		for _, line := range lines {
			if strings.Contains(line, "b2.txt") && strings.Contains(line, "*") {
				foundBuf2Modified = true
				break
			}
		}
		if !foundBuf2Modified {
			t.Errorf("b2.txt should have '*' modified marker\nscreen:\n%s", screen)
		}
	})

	// --- C-x C-b then Enter: switches to selected buffer ---
	t.Run("BufferList_EnterSwitchesToBuffer", func(t *testing.T) {
		path := shortTempFile(t, "tg.txt", "target buffer content")
		h := newHarness(t)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		openFile(h, path)
		h.AssertStatusBar(t, "tg.txt")

		// Open buffer list (buffers: *scratch* at 0, tg.txt at 1)
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-b")
		time.Sleep(300 * time.Millisecond)
		h.AssertStatusBar(t, "*Buffer List*")

		// Cursor at (0,0) -> line 0 = *scratch*. Press Enter.
		h.SendKeys("Enter")
		time.Sleep(300 * time.Millisecond)

		h.AssertStatusBar(t, "*scratch*")
	})

	// --- C-x C-b then navigate and switch ---
	t.Run("BufferList_NavigateAndSwitch", func(t *testing.T) {
		path1 := shortTempFile(t, "nav.txt", "navigate target")
		path2 := shortTempFile(t, "s2.txt", "second file")

		h := newHarness(t)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		openFile(h, path1)
		h.AssertStatusBar(t, "nav.txt")

		openFile(h, path2)
		h.AssertStatusBar(t, "s2.txt")

		// Buffer list (buffers: *scratch*(0), nav.txt(1), s2.txt(2))
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-b")
		time.Sleep(300 * time.Millisecond)
		h.AssertStatusBar(t, "*Buffer List*")

		// Move to line 1 (nav.txt) and press Enter
		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("Enter")
		time.Sleep(300 * time.Millisecond)

		h.AssertStatusBar(t, "nav.txt")
		h.AssertScreenContains(t, "navigate target")
	})
}

// typeString sends each character of a string individually with small delays.
func typeString(h *Harness, s string) {
	for _, ch := range s {
		h.SendKeys(string(ch))
		time.Sleep(20 * time.Millisecond)
	}
}

// openFile opens a file via C-x C-f minibuffer.
func openFile(h *Harness, path string) {
	h.SendKeys("C-x")
	time.Sleep(50 * time.Millisecond)
	h.SendKeys("C-f")
	time.Sleep(200 * time.Millisecond)
	typeString(h, path)
	time.Sleep(100 * time.Millisecond)
	h.SendKeys("Enter")
	time.Sleep(300 * time.Millisecond)
}

// switchBuffer switches to a buffer via C-x b minibuffer.
func switchBuffer(h *Harness, name string) {
	h.SendKeys("C-x")
	time.Sleep(50 * time.Millisecond)
	h.SendKeys("b")
	time.Sleep(200 * time.Millisecond)
	typeString(h, name)
	time.Sleep(100 * time.Millisecond)
	h.SendKeys("Enter")
}
