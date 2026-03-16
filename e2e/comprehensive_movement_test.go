package e2e

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestComprehensiveMovement(t *testing.T) {
	// --- C-f (forward char) ---
	t.Run("Cf_MiddleOfLine", func(t *testing.T) {
		content := "abcdef"
		path := createTestFile(t, "cf.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("abcdef", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}
		h.AssertCursorAt(t, 0, 0)

		h.SendKeys("C-f")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 1)

		h.SendKeys("C-f")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 2)
	})

	t.Run("Cf_EndOfLine_WrapsToNextLine", func(t *testing.T) {
		content := "ab\ncd"
		path := createTestFile(t, "cf_wrap.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("ab", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to end of first line (col 2)
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 2)

		// C-f at end of line should wrap to beginning of next line
		h.SendKeys("C-f")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)
	})

	t.Run("Cf_EndOfBuffer_NoOp", func(t *testing.T) {
		content := "ab"
		path := createTestFile(t, "cf_end.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("ab", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to end of buffer
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 2)

		// C-f at end of buffer should do nothing
		h.SendKeys("C-f")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 2)
	})

	// --- C-b (backward char) ---
	t.Run("Cb_MiddleOfLine", func(t *testing.T) {
		content := "abcdef"
		path := createTestFile(t, "cb.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("abcdef", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 6)

		h.SendKeys("C-b")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 5)
	})

	t.Run("Cb_StartOfLine_WrapsToPrevLine", func(t *testing.T) {
		content := "ab\ncd"
		path := createTestFile(t, "cb_wrap.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("ab", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to beginning of second line
		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)

		// C-b should wrap to end of previous line
		h.SendKeys("C-b")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 2)
	})

	t.Run("Cb_StartOfBuffer_NoOp", func(t *testing.T) {
		content := "ab"
		path := createTestFile(t, "cb_start.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("ab", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.AssertCursorAt(t, 0, 0)

		// C-b at start of buffer should do nothing
		h.SendKeys("C-b")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)
	})

	// --- C-n (next line) ---
	t.Run("Cn_BasicDown", func(t *testing.T) {
		content := "line one\nline two\nline three"
		path := createTestFile(t, "cn.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("line one", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.AssertCursorAt(t, 0, 0)

		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)

		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 2, 0)
	})

	t.Run("Cn_ColumnClamping", func(t *testing.T) {
		// Long line then short line
		content := "this is a longer line\nshort"
		path := createTestFile(t, "cn_clamp.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("this is", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to end of first line (col 21)
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 21)

		// C-n should clamp to end of short line (col 5)
		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 5)
	})

	t.Run("Cn_LastLine_NoOp", func(t *testing.T) {
		content := "only line"
		path := createTestFile(t, "cn_end.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("only line", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.AssertCursorAt(t, 0, 0)

		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)
	})

	// --- C-p (previous line) ---
	t.Run("Cp_BasicUp", func(t *testing.T) {
		content := "line one\nline two\nline three"
		path := createTestFile(t, "cp.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("line one", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move down two lines
		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 2, 0)

		// Move up
		h.SendKeys("C-p")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)
	})

	t.Run("Cp_ColumnClamping", func(t *testing.T) {
		// Short line then long line — move to end of long, go up
		content := "short\nthis is a longer line"
		path := createTestFile(t, "cp_clamp.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("short", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Go to second line end
		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 21)

		// C-p should clamp to end of short line (col 5)
		h.SendKeys("C-p")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 5)
	})

	t.Run("Cp_FirstLine_NoOp", func(t *testing.T) {
		content := "only line"
		path := createTestFile(t, "cp_start.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("only line", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.AssertCursorAt(t, 0, 0)

		h.SendKeys("C-p")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)
	})

	// --- C-a (beginning of line) ---
	t.Run("Ca_BeginningOfLine", func(t *testing.T) {
		content := "hello world"
		path := createTestFile(t, "ca.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 11)

		h.SendKeys("C-a")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)
	})

	t.Run("Ca_AlreadyAtBeginning", func(t *testing.T) {
		content := "hello"
		path := createTestFile(t, "ca_start.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.AssertCursorAt(t, 0, 0)

		h.SendKeys("C-a")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)
	})

	// --- C-e (end of line) ---
	t.Run("Ce_EndOfLine", func(t *testing.T) {
		content := "hello world"
		path := createTestFile(t, "ce.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.AssertCursorAt(t, 0, 0)

		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 11)
	})

	t.Run("Ce_EmptyLine", func(t *testing.T) {
		content := "first\n\nthird"
		path := createTestFile(t, "ce_empty.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("first", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to empty second line
		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)

		// C-e on empty line stays at 0
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)
	})

	// --- C-v (page down) ---
	t.Run("Cv_PageDown", func(t *testing.T) {
		var lines []string
		for i := 1; i <= 60; i++ {
			lines = append(lines, fmt.Sprintf("Line %d", i))
		}
		path := createTestFile(t, "cv.txt", strings.Join(lines, "\n"))
		h := newHarness(t, path)
		if err := h.WaitForContent("Line 1", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.AssertScreenContains(t, "Line 1")

		h.SendKeys("C-v")
		time.Sleep(300 * time.Millisecond)

		// After page down, Line 1 should not be at the top
		captured := h.Capture()
		line0 := strings.TrimSpace(captured[0])
		if line0 == "Line 1" {
			t.Error("C-v did not scroll down: Line 1 still at top")
		}
	})

	t.Run("Cv_NearEnd_PartialPage", func(t *testing.T) {
		// Create 30 lines, move to line 25, then page down — less than full page remaining
		var lines []string
		for i := 1; i <= 30; i++ {
			lines = append(lines, fmt.Sprintf("Line %d", i))
		}
		path := createTestFile(t, "cv_partial.txt", strings.Join(lines, "\n"))
		h := newHarness(t, path)
		if err := h.WaitForContent("Line 1", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Jump to near end: M-> then go up a few lines
		h.SendKeys("Escape")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys(">")
		time.Sleep(200 * time.Millisecond)

		// Now at the end. Go back up ~5 lines
		for i := 0; i < 5; i++ {
			h.SendKeys("C-p")
			time.Sleep(50 * time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond)

		// C-v should not crash even with less than a page remaining
		h.SendKeys("C-v")
		time.Sleep(300 * time.Millisecond)

		// Should see Line 30 on screen
		h.AssertScreenContains(t, "Line 30")
	})

	// --- M-v (page up) ---
	t.Run("Mv_PageUp", func(t *testing.T) {
		var lines []string
		for i := 1; i <= 60; i++ {
			lines = append(lines, fmt.Sprintf("Line %d", i))
		}
		path := createTestFile(t, "mv.txt", strings.Join(lines, "\n"))
		h := newHarness(t, path)
		if err := h.WaitForContent("Line 1", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Page down first
		h.SendKeys("C-v")
		time.Sleep(300 * time.Millisecond)

		// Page up
		h.SendKeys("Escape")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("v")
		time.Sleep(300 * time.Millisecond)

		// Should see Line 1 again
		h.AssertScreenContains(t, "Line 1")
	})

	t.Run("Mv_AtTop_NoOp", func(t *testing.T) {
		var lines []string
		for i := 1; i <= 30; i++ {
			lines = append(lines, fmt.Sprintf("Line %d", i))
		}
		path := createTestFile(t, "mv_top.txt", strings.Join(lines, "\n"))
		h := newHarness(t, path)
		if err := h.WaitForContent("Line 1", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// M-v at top should not crash
		h.SendKeys("Escape")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("v")
		time.Sleep(300 * time.Millisecond)

		// Line 1 should still be visible at top
		h.AssertScreenContains(t, "Line 1")
		h.AssertCursorAt(t, 0, 0)
	})

	// --- M-< (beginning of buffer) ---
	t.Run("MltBeginningOfBuffer", func(t *testing.T) {
		var lines []string
		for i := 1; i <= 40; i++ {
			lines = append(lines, fmt.Sprintf("Line %d", i))
		}
		path := createTestFile(t, "mlt.txt", strings.Join(lines, "\n"))
		h := newHarness(t, path)
		if err := h.WaitForContent("Line 1", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Go to end first
		h.SendKeys("Escape")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys(">")
		time.Sleep(300 * time.Millisecond)

		// M-< to jump to beginning
		h.SendKeys("Escape")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("<")
		time.Sleep(300 * time.Millisecond)

		h.AssertCursorAt(t, 0, 0)
		h.AssertScreenContains(t, "Line 1")
	})

	// --- M-> (end of buffer) ---
	t.Run("MgtEndOfBuffer", func(t *testing.T) {
		var lines []string
		for i := 1; i <= 40; i++ {
			lines = append(lines, fmt.Sprintf("Line %d", i))
		}
		path := createTestFile(t, "mgt.txt", strings.Join(lines, "\n"))
		h := newHarness(t, path)
		if err := h.WaitForContent("Line 1", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.AssertCursorAt(t, 0, 0)

		// M-> to jump to end
		h.SendKeys("Escape")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys(">")
		time.Sleep(300 * time.Millisecond)

		// Should see Line 40 on screen
		h.AssertScreenContains(t, "Line 40")
	})

	// --- Arrow keys ---
	t.Run("ArrowRight", func(t *testing.T) {
		content := "abcdef"
		path := createTestFile(t, "arrow_r.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("abcdef", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.SendKeys("Right")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 1)

		h.SendKeys("Right")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 2)
	})

	t.Run("ArrowLeft", func(t *testing.T) {
		content := "abcdef"
		path := createTestFile(t, "arrow_l.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("abcdef", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)

		h.SendKeys("Left")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 5)

		h.SendKeys("Left")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 4)
	})

	t.Run("ArrowDown", func(t *testing.T) {
		content := "line one\nline two\nline three"
		path := createTestFile(t, "arrow_d.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("line one", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.SendKeys("Down")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)

		h.SendKeys("Down")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 2, 0)
	})

	t.Run("ArrowUp", func(t *testing.T) {
		content := "line one\nline two\nline three"
		path := createTestFile(t, "arrow_u.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("line one", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 2, 0)

		h.SendKeys("Up")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)

		h.SendKeys("Up")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)
	})

	// --- Boundary positions ---
	t.Run("Cf_WrapMultipleLines", func(t *testing.T) {
		content := "a\nb\nc"
		path := createTestFile(t, "cf_multi.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("a", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// From (0,0), C-f to (0,1=end), C-f to (1,0), C-f to (1,1=end), C-f to (2,0)
		h.SendKeys("C-f") // (0,1)
		time.Sleep(50 * time.Millisecond)
		h.AssertCursorAt(t, 0, 1)

		h.SendKeys("C-f") // wrap to (1,0)
		time.Sleep(50 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)

		h.SendKeys("C-f") // (1,1)
		time.Sleep(50 * time.Millisecond)
		h.AssertCursorAt(t, 1, 1)

		h.SendKeys("C-f") // wrap to (2,0)
		time.Sleep(50 * time.Millisecond)
		h.AssertCursorAt(t, 2, 0)
	})

	t.Run("Cb_WrapMultipleLines", func(t *testing.T) {
		content := "a\nb\nc"
		path := createTestFile(t, "cb_multi.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("a", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Go to end of buffer
		h.SendKeys("Escape")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys(">")
		time.Sleep(200 * time.Millisecond)
		h.AssertCursorAt(t, 2, 1)

		h.SendKeys("C-b") // (2,0)
		time.Sleep(50 * time.Millisecond)
		h.AssertCursorAt(t, 2, 0)

		h.SendKeys("C-b") // wrap to (1,1)
		time.Sleep(50 * time.Millisecond)
		h.AssertCursorAt(t, 1, 1)

		h.SendKeys("C-b") // (1,0)
		time.Sleep(50 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)

		h.SendKeys("C-b") // wrap to (0,1)
		time.Sleep(50 * time.Millisecond)
		h.AssertCursorAt(t, 0, 1)
	})

	t.Run("Cn_Cp_ColumnClampRoundTrip", func(t *testing.T) {
		// Moving down to a short line clamps, moving back up preserves clamped col
		content := "this is a longer line\nhi\nthis is a longer line"
		path := createTestFile(t, "clamp_rt.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("this is", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Go to col 10 on first line
		for i := 0; i < 10; i++ {
			h.SendKeys("C-f")
			time.Sleep(30 * time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 10)

		// Move down — clamps to col 2 (end of "hi")
		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 2)

		// Move down again — stays at col 2 (not restored to 10)
		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 2, 2)
	})

	// --- Empty buffer movement ---
	t.Run("EmptyBuffer_AllMovements", func(t *testing.T) {
		h := newHarness(t)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// All movements should be no-ops in empty buffer
		h.SendKeys("C-f")
		time.Sleep(50 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)

		h.SendKeys("C-b")
		time.Sleep(50 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)

		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)

		h.SendKeys("C-p")
		time.Sleep(50 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)

		h.SendKeys("C-a")
		time.Sleep(50 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)

		h.SendKeys("C-e")
		time.Sleep(50 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)
	})

	// --- Multi-line C-a and C-e ---
	t.Run("Ca_Ce_OnDifferentLines", func(t *testing.T) {
		content := "first line\nsecond line here\nthird"
		path := createTestFile(t, "ca_ce_ml.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("first", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Go to second line, middle
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		for i := 0; i < 5; i++ {
			h.SendKeys("C-f")
			time.Sleep(30 * time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 5)

		// C-a to beginning of second line
		h.SendKeys("C-a")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)

		// C-e to end of second line
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 16)
	})

	// --- Page up/down near boundaries ---
	t.Run("Cv_Mv_MultiplePagesDown", func(t *testing.T) {
		var lines []string
		for i := 1; i <= 100; i++ {
			lines = append(lines, fmt.Sprintf("Line %d", i))
		}
		path := createTestFile(t, "cv_multi.txt", strings.Join(lines, "\n"))
		h := newHarness(t, path)
		if err := h.WaitForContent("Line 1", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Page down 3 times
		for i := 0; i < 3; i++ {
			h.SendKeys("C-v")
			time.Sleep(200 * time.Millisecond)
		}

		// Line 1 should not be visible
		screen := h.CapturePane()
		if strings.Contains(screen, "Line 1\n") {
			firstLine := strings.TrimSpace(h.Capture()[0])
			if firstLine == "Line 1" {
				t.Error("after 3 page-downs, Line 1 should not be at top")
			}
		}

		// Page up 3 times to go back
		for i := 0; i < 3; i++ {
			h.SendKeys("Escape")
			time.Sleep(50 * time.Millisecond)
			h.SendKeys("v")
			time.Sleep(200 * time.Millisecond)
		}

		// Line 1 should be visible again
		h.AssertScreenContains(t, "Line 1")
	})

	// --- Arrow key wrapping (same as C-f/C-b behavior) ---
	t.Run("ArrowRight_WrapsToNextLine", func(t *testing.T) {
		content := "ab\ncd"
		path := createTestFile(t, "arr_r_wrap.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("ab", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 2)

		// Right arrow at end of line wraps
		h.SendKeys("Right")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)
	})

	t.Run("ArrowLeft_WrapsToPrevLine", func(t *testing.T) {
		content := "ab\ncd"
		path := createTestFile(t, "arr_l_wrap.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("ab", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Go to beginning of second line
		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)

		// Left at start of line wraps to end of previous line
		h.SendKeys("Left")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 2)
	})

	t.Run("ArrowDown_ColumnClamping", func(t *testing.T) {
		content := "this is a longer line\nhi"
		path := createTestFile(t, "arr_d_clamp.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("this is", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 21)

		h.SendKeys("Down")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 2)
	})

	t.Run("ArrowUp_ColumnClamping", func(t *testing.T) {
		content := "hi\nthis is a longer line"
		path := createTestFile(t, "arr_u_clamp.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hi", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 21)

		h.SendKeys("Up")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 2)
	})
}
