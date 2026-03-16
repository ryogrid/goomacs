package e2e

import (
	"strings"
	"testing"
	"time"
)

func TestComprehensiveInteraction(t *testing.T) {

	// --- Search while region is active: verify mark behavior after search exit ---
	t.Run("SearchWhileRegionActive", func(t *testing.T) {
		content := "alpha beta gamma delta"
		path := createTestFile(t, "interact_search_region.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("alpha", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Set mark at col 0
		h.SendKeys("C-Space")
		time.Sleep(100 * time.Millisecond)

		// Move forward to col 5
		for i := 0; i < 5; i++ {
			h.SendKeys("C-f")
			time.Sleep(30 * time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond)

		// Enter search mode (C-s)
		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("gamma")
		time.Sleep(300 * time.Millisecond)

		// Accept search with Enter — cursor stays at match
		h.SendKeys("Enter")
		time.Sleep(200 * time.Millisecond)

		// Cursor should be at "gamma" position (col 11)
		h.AssertLineContains(t, 0, "alpha beta gamma delta")

		// Editor should still be functional — type a character
		h.SendKeys("X")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "Xgamma")
	})

	// --- Undo after yank: C-y to paste, C-_ to undo, verify yanked text removed ---
	t.Run("UndoAfterYank", func(t *testing.T) {
		content := "hello world"
		path := createTestFile(t, "interact_undo_yank.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello world", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Kill "hello world" to fill kill ring
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineEquals(t, 0, "")

		// Move to end (already there) and yank
		h.SendKeys("C-y")
		time.Sleep(200 * time.Millisecond)
		h.AssertLineContains(t, 0, "hello world")

		// Undo the yank — should remove the pasted text
		h.SendKeys("C-_")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineEquals(t, 0, "")
		h.AssertMessageLine(t, "Undo")
	})

	// --- Undo after kill region: C-w to kill, C-_ to undo, verify region text restored ---
	t.Run("UndoAfterKillRegion", func(t *testing.T) {
		content := "ABCDEF"
		path := createTestFile(t, "interact_undo_region.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("ABCDEF", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Set mark, move forward 3, C-w to kill "ABC"
		h.SendKeys("C-Space")
		time.Sleep(100 * time.Millisecond)
		for i := 0; i < 3; i++ {
			h.SendKeys("C-f")
			time.Sleep(30 * time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("C-w")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineEquals(t, 0, "DEF")

		// Undo should restore "ABCDEF"
		h.SendKeys("C-_")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "ABCDEF")
		h.AssertMessageLine(t, "Undo")
	})

	// --- C-x C-f during search mode: exits search and opens find-file prompt ---
	t.Run("FindFileDuringSearch", func(t *testing.T) {
		content := "search target here"
		path := createTestFile(t, "interact_search_ff.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("search target", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Start forward search
		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("target")
		time.Sleep(300 * time.Millisecond)
		h.AssertMessageLine(t, "I-search: target")

		// C-x C-f — unrecognized key in search exits search, then C-x is processed
		// Search mode: C-x is unrecognized → exits search, re-posts C-x
		// Then C-x prefix mode is active, C-f opens find-file
		h.SendKeys("C-x")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("C-f")
		time.Sleep(300 * time.Millisecond)

		// Should now be in minibuffer with "Find file:" prompt
		h.AssertMessageLine(t, "Find file:")
		// Cancel out
		h.SendKeys("C-g")
		time.Sleep(100 * time.Millisecond)
	})

	// --- Split window then kill buffer: C-x 2, C-x k, window shows next buffer ---
	t.Run("SplitWindowThenKillBuffer", func(t *testing.T) {
		path1 := shortTempFile(t, "splitkill1.txt", "buffer one content")
		path2 := shortTempFile(t, "splitkill2.txt", "buffer two content")
		h := newHarness(t, path1)
		if err := h.WaitForContent("splitkill1.txt", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Open second file
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-f")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys(path2)
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("Enter")
		time.Sleep(300 * time.Millisecond)
		if err := h.WaitForContent("buffer two", 5*time.Second); err != nil {
			t.Fatalf("second file did not open: %v", err)
		}

		// Split window (C-x 2)
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("2")
		time.Sleep(300 * time.Millisecond)

		// Both windows should show buffer two
		screen := h.CapturePane()
		count := strings.Count(screen, "buffer two content")
		if count < 2 {
			t.Errorf("expected buffer two in both windows, found %d\nScreen:\n%s", count, screen)
		}

		// Kill the buffer (unmodified, so no confirmation)
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("k")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("Enter")
		time.Sleep(300 * time.Millisecond)

		// After killing buffer two, should switch to buffer one
		if err := h.WaitForContent("buffer one", 5*time.Second); err != nil {
			t.Errorf("expected to switch to buffer one after kill: %v", err)
		}
	})

	// --- Comment-region then undo: M-x comment-region, C-_ undo, comments removed ---
	t.Run("CommentRegionThenUndo", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "func main() {\n\tfmt.Println(\"hi\")\n}\n"
		path := shortTempFileInDir(t, dir, "undo_comment.go", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("undo_comment.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Select first two lines
		h.SendKeys("C-Space")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)

		// M-x comment-region
		invokeCommand(h, "comment-region")
		h.AssertMessageLine(t, "Region commented")

		// Verify comments were added
		h.AssertLineContains(t, 0, "// func main")

		// Undo should remove the comments
		h.SendKeys("C-_")
		time.Sleep(200 * time.Millisecond)

		// Verify original content is restored (no // prefix)
		lines := h.Capture()
		line0 := strings.TrimRight(lines[0], " ")
		if strings.HasPrefix(line0, "//") {
			t.Errorf("expected comments removed after undo, got %q", line0)
		}
		h.AssertLineContains(t, 0, "func main")
	})

	// --- Save then undo: C-x C-s, C-_ undo, buffer shows [Modified] again ---
	t.Run("SaveThenUndo", func(t *testing.T) {
		path := shortTempFile(t, "su.txt", "original text")
		h := newHarness(t, path)
		if err := h.WaitForContent("original text", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Type some text to modify the buffer
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("!")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "original text!")
		h.AssertStatusBar(t, "[Modified]")

		// Save with C-x C-s
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-s")
		time.Sleep(300 * time.Millisecond)
		h.AssertMessageLine(t, "Saved")

		// Status bar should NOT show [Modified] after save
		lines := h.Capture()
		statusBar := lines[h.height-2]
		if strings.Contains(statusBar, "[Modified]") {
			t.Errorf("expected no [Modified] after save, got: %s", statusBar)
		}

		// Undo the edit
		h.SendKeys("C-_")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineEquals(t, 0, "original text")
		h.AssertMessageLine(t, "Undo")

		// Status bar should show [Modified] again since we undid past the save point
		h.AssertStatusBar(t, "[Modified]")
	})
}
