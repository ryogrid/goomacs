package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestComprehensiveNoop(t *testing.T) {

	// --- C-x followed by invalid second key ---
	// In goomacs, unmatched rune keys in C-x prefix are consumed (not re-dispatched)
	t.Run("CxPrefix_InvalidSecondKey_a", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "hello world\nsecond line\n"
		path := shortTempFileInDir(t, dir, "t.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("t.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// C-x then 'a' — prefix cancelled, 'a' consumed (not inserted)
		h.SendKeys("C-x")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("a")
		time.Sleep(200 * time.Millisecond)

		// Buffer should be unchanged — 'a' was not inserted
		h.AssertScreenContains(t, "hello world")
		h.AssertCursorAt(t, 0, 0)

		// Editor still functional — can type after prefix cancellation
		h.SendKeys("Z")
		time.Sleep(100 * time.Millisecond)
		h.AssertScreenContains(t, "Zhello world")
	})

	t.Run("CxPrefix_InvalidSecondKey_z", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "hello world\n"
		path := shortTempFileInDir(t, dir, "t.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("t.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// C-x then 'z' — prefix cancelled, 'z' consumed (not inserted)
		h.SendKeys("C-x")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("z")
		time.Sleep(200 * time.Millisecond)

		// Buffer unchanged
		h.AssertScreenContains(t, "hello world")
		h.AssertCursorAt(t, 0, 0)
	})

	t.Run("CxPrefix_InvalidControlKey", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "hello world\nsecond line\n"
		path := shortTempFileInDir(t, dir, "t.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("t.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// C-x then unbound control key (C-n is not a valid C-x subcommand)
		// In the C-x prefix handler, unmatched keys are consumed
		h.SendKeys("C-x")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(200 * time.Millisecond)

		// Cursor should stay at (0,0) — C-n was consumed by prefix handler
		h.AssertCursorAt(t, 0, 0)

		// Editor is still functional — type a char
		h.SendKeys("X")
		time.Sleep(100 * time.Millisecond)
		h.AssertScreenContains(t, "X")
	})

	t.Run("CxPrefix_EditorStillFunctional", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "line one\nline two\n"
		path := shortTempFileInDir(t, dir, "t.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("t.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Send several invalid C-x sequences
		for i := 0; i < 5; i++ {
			h.SendKeys("C-x")
			time.Sleep(50 * time.Millisecond)
			h.SendKeys("C-g") // C-g after C-x — cancel prefix
			time.Sleep(50 * time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond)

		// Editor still functional — cursor should be at start
		h.AssertCursorAt(t, 0, 0)
		// Can still type
		h.SendKeys("Z")
		time.Sleep(100 * time.Millisecond)
		h.AssertScreenContains(t, "Zline one")
	})

	// --- Unbound control keys in normal mode ---
	t.Run("UnboundControlKeys_NormalMode", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "test content\nsecond line\n"
		path := shortTempFileInDir(t, dir, "t.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("t.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		time.Sleep(100 * time.Millisecond)
		beforeScreen := h.CapturePane()

		// C-o and C-t are unbound — should be no-ops (or ignored)
		// Note: C-q may be intercepted by the terminal (flow control)
		h.SendKeys("C-o")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("C-t")
		time.Sleep(100 * time.Millisecond)

		afterScreen := h.CapturePane()

		// Buffer content should be unchanged (compare content lines, not status bar)
		beforeLines := strings.Split(beforeScreen, "\n")
		afterLines := strings.Split(afterScreen, "\n")
		// Compare first content line
		if len(beforeLines) > 0 && len(afterLines) > 0 {
			if strings.TrimSpace(beforeLines[0]) != strings.TrimSpace(afterLines[0]) {
				t.Errorf("content changed after unbound keys: before=%q after=%q", beforeLines[0], afterLines[0])
			}
		}

		// Editor still functional
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 12) // end of "test content"
	})

	// --- Rapid key repeat: 100 C-f presses ---
	t.Run("RapidKeyRepeat_100_Cf", func(t *testing.T) {
		dir := shortTempDir(t)
		// "short\n" loads as single line "short" (trailing newline trimmed by file reader)
		content := "short\nsecond\n"
		path := shortTempFileInDir(t, dir, "t.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("t.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Send 100 C-f presses rapidly using SendKeysRepeat
		h.SendKeysRepeat("C-f", 100)
		time.Sleep(500 * time.Millisecond)

		// Cursor should be at end of buffer — "short" (5 chars) + "second" (6 chars)
		// Buffer has 2 lines: ["short", "second"]; end is (1, 6)
		row, col, err := h.CursorPosition()
		if err != nil {
			t.Fatalf("failed to get cursor: %v", err)
		}
		if row != 1 || col != 6 {
			t.Errorf("expected cursor at (1,6) after 100 C-f, got (%d,%d)", row, col)
		}

		// Editor still functional
		h.SendKeys("a")
		time.Sleep(100 * time.Millisecond)
		h.AssertScreenContains(t, "seconda")
	})

	// --- Rapid mode switching: C-s then C-g repeated 20 times ---
	t.Run("RapidModeSwitching_SearchCancel", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "hello world\nsearch target\nthird line\n"
		path := shortTempFileInDir(t, dir, "t.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("t.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Rapid C-s then C-g 20 times
		for i := 0; i < 20; i++ {
			h.SendKeys("C-s")
			time.Sleep(30 * time.Millisecond)
			h.SendKeys("C-g")
			time.Sleep(30 * time.Millisecond)
		}
		time.Sleep(200 * time.Millisecond)

		// Cursor should still be at (0,0) — no state corruption
		h.AssertCursorAt(t, 0, 0)

		// Buffer unchanged
		h.AssertScreenContains(t, "hello world")
		h.AssertScreenContains(t, "search target")

		// Editor still functional
		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)
	})

	// --- C-g in Normal mode (deactivate mark) ---
	t.Run("Cg_NormalMode_DeactivateMark", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "line one\nline two\n"
		path := shortTempFileInDir(t, dir, "t.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("t.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Set mark then cancel
		h.SendKeys("C-Space")
		time.Sleep(100 * time.Millisecond)
		h.AssertMessageLine(t, "Mark set")

		h.SendKeys("C-g")
		time.Sleep(100 * time.Millisecond)

		// Mark deactivated, cursor unchanged at (0,0)
		h.AssertCursorAt(t, 0, 0)
	})

	// --- C-g in Search mode (cancel search) ---
	t.Run("Cg_SearchMode_Cancel", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "hello world\ntarget line\n"
		path := shortTempFileInDir(t, dir, "t.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("t.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Start search, type query, then cancel
		h.SendKeys("C-s")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("target")
		time.Sleep(200 * time.Millisecond)

		// Should have found "target" on line 1
		h.AssertMessageLine(t, "I-search: target")

		// Cancel — cursor returns to original position
		h.SendKeys("C-g")
		time.Sleep(100 * time.Millisecond)

		h.AssertCursorAt(t, 0, 0)
		h.AssertMessageLine(t, "Quit")
	})

	// --- C-g in Minibuffer mode (cancel input) ---
	t.Run("Cg_MinibufferMode_Cancel", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "hello\n"
		path := shortTempFileInDir(t, dir, "t.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("t.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Open find-file minibuffer
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-f")
		time.Sleep(200 * time.Millisecond)

		h.AssertMessageLine(t, "Find file:")

		// Cancel with C-g
		h.SendKeys("C-g")
		time.Sleep(100 * time.Millisecond)

		h.AssertMessageLine(t, "Quit")
	})

	// --- C-g in C-x prefix mode (cancel prefix) ---
	t.Run("Cg_CxPrefixMode_Cancel", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "hello\n"
		path := shortTempFileInDir(t, dir, "t.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("t.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// C-x prefix
		h.SendKeys("C-x")
		time.Sleep(100 * time.Millisecond)
		h.AssertMessageLine(t, "C-x-")

		// C-g cancels the prefix (falls through to normal C-g handler)
		h.SendKeys("C-g")
		time.Sleep(100 * time.Millisecond)

		// Editor back to normal, cursor unchanged
		h.AssertCursorAt(t, 0, 0)
	})

	// --- Confirm mode: invalid keys do nothing ---
	t.Run("ConfirmMode_InvalidKeys", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "hello\n"
		path := shortTempFileInDir(t, dir, "t.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("t.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Modify buffer so kill-buffer asks for confirmation
		h.SendKeys("X")
		time.Sleep(100 * time.Millisecond)

		// C-x k to trigger confirm prompt
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("k")
		time.Sleep(200 * time.Millisecond)
		// Accept default (just press Enter in kill-buffer minibuffer)
		h.SendKeys("Enter")
		time.Sleep(200 * time.Millisecond)

		// Should show "Buffer modified; kill anyway? (y/n)"
		h.AssertMessageLine(t, "kill anyway")

		// Send invalid keys — confirm mode stays active but message is cleared
		// (message = "" runs at top of key handler before confirm check)
		// The important thing is confirm mode is NOT exited:
		h.SendKeys("a")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("C-f")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("Enter")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("x")
		time.Sleep(100 * time.Millisecond)

		// Verify confirm mode is still active by pressing 'n' — should cancel
		h.SendKeys("n")
		time.Sleep(100 * time.Millisecond)
		h.AssertMessageLine(t, "Cancelled")

		// Buffer should still be open (confirm was cancelled, not confirmed)
		h.AssertStatusBar(t, "t.go")
	})

	// --- Confirm mode: y confirms ---
	t.Run("ConfirmMode_YConfirms", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "hello\n"
		path := shortTempFileInDir(t, dir, "t.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("t.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Modify buffer
		h.SendKeys("X")
		time.Sleep(100 * time.Millisecond)

		// C-x k + Enter to trigger confirm
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("k")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("Enter")
		time.Sleep(200 * time.Millisecond)

		h.AssertMessageLine(t, "kill anyway")

		// Press y to confirm kill
		h.SendKeys("y")
		time.Sleep(200 * time.Millisecond)

		// Buffer should be killed — now on *scratch*
		h.AssertScreenContains(t, "*scratch*")
	})

	// --- Confirm mode: C-g quits ---
	t.Run("ConfirmMode_CgQuits", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "hello\n"
		path := shortTempFileInDir(t, dir, "t.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("t.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Modify buffer
		h.SendKeys("X")
		time.Sleep(100 * time.Millisecond)

		// C-x k + Enter to trigger confirm
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("k")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("Enter")
		time.Sleep(200 * time.Millisecond)

		h.AssertMessageLine(t, "kill anyway")

		// C-g to quit
		h.SendKeys("C-g")
		time.Sleep(100 * time.Millisecond)
		h.AssertMessageLine(t, "Quit")

		// Buffer should still be open
		h.AssertStatusBar(t, "t.go")
	})

	// --- Buffer content and cursor unchanged after no-op sequence ---
	t.Run("NoopSequence_BufferUnchanged", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "line one\nline two\nline three\n"
		path := shortTempFileInDir(t, dir, "t.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("t.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Move to a known position
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-f")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-f")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-f")
		time.Sleep(100 * time.Millisecond)

		// Cursor at (1, 3) — "lin|e two"
		h.AssertCursorAt(t, 1, 3)

		// Capture before
		beforeScreen := h.CapturePane()

		// Send a bunch of no-ops: C-o, C-t, C-x + C-g (cancel prefix)
		h.SendKeys("C-o")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-t")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-g")
		time.Sleep(100 * time.Millisecond)

		// Cursor should remain at same position
		h.AssertCursorAt(t, 1, 3)

		// Content should be unchanged
		afterScreen := h.CapturePane()
		beforeContentLines := extractContentLines(beforeScreen, 3)
		afterContentLines := extractContentLines(afterScreen, 3)
		for i := 0; i < len(beforeContentLines) && i < len(afterContentLines); i++ {
			if beforeContentLines[i] != afterContentLines[i] {
				t.Errorf("content line %d changed: before=%q after=%q", i, beforeContentLines[i], afterContentLines[i])
			}
		}
	})

	// --- Rapid C-f at buffer end is no-op ---
	t.Run("RapidCf_AtBufferEnd", func(t *testing.T) {
		dir := shortTempDir(t)
		// "ab\n" loads as single line "ab" (trailing newline trimmed)
		content := "ab\n"
		path := shortTempFileInDir(t, dir, "t.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("t.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Move to end of buffer
		h.SendKeys("Escape")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys(">")
		time.Sleep(200 * time.Millisecond)

		// Now send 50 more C-f — should all be no-ops
		h.SendKeysRepeat("C-f", 50)
		time.Sleep(300 * time.Millisecond)

		// Cursor at end of single-line buffer: (0, 2)
		h.AssertCursorAt(t, 0, 2)
	})

	// --- Grep mode: typing blocked as read-only ---
	t.Run("GrepMode_TypingBlocked", func(t *testing.T) {
		dir := setupGrepFixtureForNoop(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		invokeFindGrep(t, h, "grep -rnH hello ./main.go | sort")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		// Try typing various keys — all should be blocked as read-only
		h.SendKeys("x")
		time.Sleep(100 * time.Millisecond)
		h.AssertMessageLine(t, "Buffer is read-only")

		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)
		h.AssertMessageLine(t, "Buffer is read-only")

		h.SendKeys("BSpace")
		time.Sleep(100 * time.Millisecond)
		h.AssertMessageLine(t, "Buffer is read-only")
	})

	// --- Multiple C-x prefix sequences don't corrupt state ---
	t.Run("MultipleCxPrefix_NoCorruption", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "test line\n"
		path := shortTempFileInDir(t, dir, "t.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("t.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// C-x C-x C-x... rapid prefix presses — each should show "C-x-" and be cancelled by next
		for i := 0; i < 10; i++ {
			h.SendKeys("C-x")
			time.Sleep(30 * time.Millisecond)
		}
		// Last one leaves us in C-x prefix mode — cancel it
		h.SendKeys("C-g")
		time.Sleep(200 * time.Millisecond)

		// Editor should be functional
		h.SendKeys("A")
		time.Sleep(100 * time.Millisecond)
		h.AssertScreenContains(t, "Atest line")
	})

	// --- Mixed mode transitions don't corrupt state ---
	t.Run("MixedModeTransitions", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "hello\nworld\n"
		path := shortTempFileInDir(t, dir, "t.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("t.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Search mode → cancel
		h.SendKeys("C-s")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-g")
		time.Sleep(50 * time.Millisecond)

		// Minibuffer → cancel
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-f")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("C-g")
		time.Sleep(50 * time.Millisecond)

		// C-x prefix → cancel
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-g")
		time.Sleep(50 * time.Millisecond)

		// Set mark → cancel
		h.SendKeys("C-Space")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-g")
		time.Sleep(100 * time.Millisecond)

		// Cursor should still be at (0,0) and buffer unchanged
		h.AssertCursorAt(t, 0, 0)
		h.AssertScreenContains(t, "hello")
		h.AssertScreenContains(t, "world")
	})
}

// extractContentLines returns the first n lines from a screen capture (trimmed).
func extractContentLines(screen string, n int) []string {
	lines := strings.Split(screen, "\n")
	result := make([]string, 0, n)
	for i := 0; i < n && i < len(lines); i++ {
		result = append(result, strings.TrimRight(lines[i], " "))
	}
	return result
}

// setupGrepFixtureForNoop creates a temp dir with a single file for grep testing.
func setupGrepFixtureForNoop(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("/tmp", "g")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	content := "package main\n\nfunc main() {\n\tfmt.Println(\"hello\")\n}\n"
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	return dir
}

