package e2e

import (
	"strings"
	"testing"
	"time"
)

func TestComprehensiveKillYank(t *testing.T) {
	// --- C-k: kill to end of line ---
	t.Run("Ck_KillToEndOfLine", func(t *testing.T) {
		content := "hello world"
		path := createTestFile(t, "ck_eol.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello world", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to col 5 (after "hello")
		for i := 0; i < 5; i++ {
			h.SendKeys("C-f")
			time.Sleep(30 * time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 5)

		// C-k kills " world"
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 5)
		h.AssertLineEquals(t, 0, "hello")
	})

	t.Run("Ck_AtEndOfLine_JoinsNextLine", func(t *testing.T) {
		content := "first\nsecond"
		path := createTestFile(t, "ck_join.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("first", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to end of first line
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 5)

		// C-k at end of line joins with next
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 5)
		h.AssertLineContains(t, 0, "firstsecond")
	})

	t.Run("Ck_OnEmptyLine", func(t *testing.T) {
		content := "hello\n\nworld"
		path := createTestFile(t, "ck_empty.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to empty line (line 1)
		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)

		// C-k on empty line kills the newline (joins with next line)
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)
		// Now line 1 should be "world"
		h.AssertLineContains(t, 1, "world")
	})

	// --- Consecutive kills accumulate in kill ring ---
	t.Run("Ck_ConsecutiveKills_Accumulate", func(t *testing.T) {
		content := "line one\nline two\nline three"
		path := createTestFile(t, "ck_consec.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("line one", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// C-k twice: first kills "line one", second kills the newline
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)

		// Now cursor should be at start of what was "line two"
		h.AssertCursorAt(t, 0, 0)
		h.AssertLineContains(t, 0, "line two")

		// Move to end of buffer and yank — should paste "line one\n" (both kills accumulated)
		h.SendKeys("C-e")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("Enter")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-y")
		time.Sleep(200 * time.Millisecond)

		// The yanked text should contain "line one" followed by the joined newline
		h.AssertScreenContains(t, "line one")
	})

	// --- Non-kill key breaks consecutive kill chain ---
	t.Run("Ck_NonKillKey_BreaksChain", func(t *testing.T) {
		content := "aaa\nbbb\nccc"
		path := createTestFile(t, "ck_break.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("aaa", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// C-k to kill "aaa"
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)

		// C-f (non-kill) breaks the chain — moves past the newline to "bbb"
		h.SendKeys("C-f")
		time.Sleep(100 * time.Millisecond)

		// C-k again kills "bb" (from col 1 to end since C-f moved past newline to start of bbb, then forward 1)
		// Actually after first C-k: line 0 is empty "", line 1 is "bbb", line 2 is "ccc"
		// C-f at end of empty line wraps to start of "bbb" (line 1, col 0)
		// Actually wait — after C-k on "aaa", line 0 becomes "" (empty), cursor at (0,0)
		// Then C-k on empty line joins with next: line 0 becomes "bbb", cursor at (0,0)
		// Hmm, no. Let me re-think.
		// After first C-k: "aaa" is killed, line 0 is now empty "", lines: ["", "bbb", "ccc"], cursor at (0,0)
		// C-f on empty line at col 0: wraps to line 1, col 0 → cursor at (1,0)
		// C-k at (1,0) kills "bbb" → line 1 becomes "", cursor at (1,0)
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)

		// Move to end and yank — should only get "bbb" (second kill), not "aaa" + "bbb"
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-e")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-y")
		time.Sleep(200 * time.Millisecond)

		// The yank should paste only "bbb" (the second, separate kill)
		h.AssertScreenContains(t, "cccbbb")
	})

	// --- C-w: kill region (single line) ---
	t.Run("Cw_SingleLineRegion", func(t *testing.T) {
		content := "hello world"
		path := createTestFile(t, "cw_single.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello world", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Set mark at col 0
		h.SendKeys("C-Space")
		time.Sleep(100 * time.Millisecond)

		// Move to col 5
		for i := 0; i < 5; i++ {
			h.SendKeys("C-f")
			time.Sleep(30 * time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 5)

		// C-w kills "hello"
		h.SendKeys("C-w")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)
		h.AssertLineContains(t, 0, " world")
	})

	// --- C-w: multi-line region ---
	t.Run("Cw_MultiLineRegion", func(t *testing.T) {
		content := "line one\nline two\nline three"
		path := createTestFile(t, "cw_multi.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("line one", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Set mark at beginning
		h.SendKeys("C-Space")
		time.Sleep(100 * time.Millisecond)

		// Move down two lines to start of line three
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)

		// C-w kills from start of "line one" to col 0 of "line three"
		// That means "line one\nline two\n" is killed
		h.SendKeys("C-w")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)
		h.AssertLineContains(t, 0, "line three")
	})

	// --- C-w: zero-width region (mark == point) ---
	t.Run("Cw_ZeroWidthRegion", func(t *testing.T) {
		content := "hello"
		path := createTestFile(t, "cw_zero.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Set mark at current position
		h.SendKeys("C-Space")
		time.Sleep(100 * time.Millisecond)

		// C-w with zero-width region: no text killed, buffer unchanged
		h.SendKeys("C-w")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)
		h.AssertLineContains(t, 0, "hello")
	})

	// --- M-w: copy region without deleting ---
	t.Run("Mw_CopiesWithoutDeleting", func(t *testing.T) {
		content := "hello world"
		path := createTestFile(t, "mw_copy.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello world", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Set mark at col 0
		h.SendKeys("C-Space")
		time.Sleep(100 * time.Millisecond)

		// Move to col 5
		for i := 0; i < 5; i++ {
			h.SendKeys("C-f")
			time.Sleep(30 * time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond)

		// M-w copies "hello" without deleting
		h.SendKeys("Escape")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("w")
		time.Sleep(100 * time.Millisecond)

		// Buffer should be unchanged
		h.AssertLineContains(t, 0, "hello world")
		h.AssertMessageLine(t, "Region copied")

		// Move to end of line and yank
		h.SendKeys("C-e")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-y")
		time.Sleep(200 * time.Millisecond)

		// "hello" should be pasted at end
		h.AssertLineContains(t, 0, "hello worldhello")
	})

	// --- C-y: yank single line ---
	t.Run("Cy_YankSingleLine", func(t *testing.T) {
		content := "hello world"
		path := createTestFile(t, "cy_single.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello world", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Kill to end of line
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineEquals(t, 0, "")

		// Yank it back
		h.SendKeys("C-y")
		time.Sleep(200 * time.Millisecond)
		h.AssertLineContains(t, 0, "hello world")
	})

	// --- C-y: yank multi-line ---
	t.Run("Cy_YankMultiLine", func(t *testing.T) {
		content := "aaa\nbbb\nccc"
		path := createTestFile(t, "cy_multi.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("aaa", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Set mark, move down 2 lines, C-w to kill multi-line region
		h.SendKeys("C-Space")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("C-w")
		time.Sleep(100 * time.Millisecond)

		// Should just have "ccc" left
		h.AssertLineContains(t, 0, "ccc")

		// Move to end and yank
		h.SendKeys("C-e")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("Enter")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-y")
		time.Sleep(200 * time.Millisecond)

		// The multi-line kill should be yanked back
		h.AssertScreenContains(t, "aaa")
		h.AssertScreenContains(t, "bbb")
	})

	// --- C-y: yank into empty buffer ---
	t.Run("Cy_YankIntoEmptyBuffer", func(t *testing.T) {
		content := "hello"
		path := createTestFile(t, "cy_empty.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Kill the line content
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineEquals(t, 0, "")

		// Now we have an empty buffer (single empty line) with "hello" in kill ring
		// Yank into it
		h.SendKeys("C-y")
		time.Sleep(200 * time.Millisecond)
		h.AssertLineContains(t, 0, "hello")
	})

	// --- C-y: empty kill ring ---
	t.Run("Cy_EmptyKillRing_NoChange", func(t *testing.T) {
		h := newHarness(t)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Type some text
		h.SendKeys("a")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("b")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("c")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "abc")
		h.AssertCursorAt(t, 0, 3)

		// C-y with empty kill ring: no change
		h.SendKeys("C-y")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "abc")
		h.AssertCursorAt(t, 0, 3)
	})

	// --- Region + yank interaction ---
	t.Run("RegionKillThenYankElsewhere", func(t *testing.T) {
		content := "aaaa\nbbbb\ncccc"
		path := createTestFile(t, "region_yank.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("aaaa", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Set mark at beginning of line 0
		h.SendKeys("C-Space")
		time.Sleep(100 * time.Millisecond)

		// Move to end of "aaaa"
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)

		// C-w kills "aaaa"
		h.SendKeys("C-w")
		time.Sleep(100 * time.Millisecond)
		// Line 0 should now be empty (just the remainder after "aaaa")
		lines := h.Capture()
		line0 := strings.TrimRight(lines[0], " ")
		if line0 != "" {
			t.Logf("line 0 after kill: %q (expected empty)", line0)
		}

		// Move to end of line 2 (was line 2, now shifted)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)

		// Yank "aaaa" at end of cccc
		h.SendKeys("C-y")
		time.Sleep(200 * time.Millisecond)
		h.AssertScreenContains(t, "ccccaaaa")
	})

	// --- C-k then C-y round-trip ---
	t.Run("Ck_Cy_RoundTrip", func(t *testing.T) {
		content := "the quick brown fox"
		path := createTestFile(t, "ck_cy_rt.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("the quick", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Kill the entire line content
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineEquals(t, 0, "")

		// Yank it back
		h.SendKeys("C-y")
		time.Sleep(200 * time.Millisecond)
		h.AssertLineContains(t, 0, "the quick brown fox")
	})

	// --- Multiple C-k on multi-line: kill all then yank all ---
	t.Run("Ck_MultipleLines_YankAll", func(t *testing.T) {
		content := "AAA\nBBB"
		path := createTestFile(t, "ck_multi_yank.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("AAA", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// C-k 4 times: kill "AAA", kill newline, kill "BBB", (no more to kill)
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)

		// Buffer should be a single empty line now
		lines := h.Capture()
		line0 := strings.TrimRight(lines[0], " ")
		if line0 != "" {
			t.Errorf("expected empty buffer after killing all, got %q", line0)
		}

		// Yank should restore everything (consecutive kills accumulated)
		h.SendKeys("C-y")
		time.Sleep(200 * time.Millisecond)
		h.AssertScreenContains(t, "AAA")
		h.AssertScreenContains(t, "BBB")
	})

	// --- C-w with mark set backward (point before mark) ---
	t.Run("Cw_BackwardRegion", func(t *testing.T) {
		content := "hello world"
		path := createTestFile(t, "cw_backward.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello world", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to col 5
		for i := 0; i < 5; i++ {
			h.SendKeys("C-f")
			time.Sleep(30 * time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond)

		// Set mark at col 5
		h.SendKeys("C-Space")
		time.Sleep(100 * time.Millisecond)

		// Move backward to col 0
		h.SendKeys("C-a")
		time.Sleep(100 * time.Millisecond)

		// C-w kills from col 0 to col 5 (backward region)
		h.SendKeys("C-w")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)
		h.AssertLineContains(t, 0, " world")
	})

	// --- C-w then C-y: kill region and yank it back ---
	t.Run("Cw_Cy_RoundTrip", func(t *testing.T) {
		content := "ABCDEF"
		path := createTestFile(t, "cw_cy.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("ABCDEF", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Set mark at col 1
		h.SendKeys("C-f")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-Space")
		time.Sleep(100 * time.Millisecond)

		// Move to col 4
		for i := 0; i < 3; i++ {
			h.SendKeys("C-f")
			time.Sleep(30 * time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond)

		// C-w kills "BCD"
		h.SendKeys("C-w")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "AEF")

		// Yank "BCD" back at current position
		h.SendKeys("C-y")
		time.Sleep(200 * time.Millisecond)
		h.AssertLineContains(t, 0, "ABCDEF")
	})

	// --- M-w then C-y: copy then paste ---
	t.Run("Mw_Cy_CopyThenPaste", func(t *testing.T) {
		content := "xyz"
		path := createTestFile(t, "mw_cy.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("xyz", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Set mark at col 0, move to col 3 (select all "xyz")
		h.SendKeys("C-Space")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)

		// M-w copies
		h.SendKeys("Escape")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("w")
		time.Sleep(100 * time.Millisecond)

		// Buffer unchanged
		h.AssertLineContains(t, 0, "xyz")

		// Yank at end
		h.SendKeys("C-e")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-y")
		time.Sleep(200 * time.Millisecond)
		h.AssertLineContains(t, 0, "xyzxyz")
	})

	// --- C-w with no mark active ---
	t.Run("Cw_NoMarkActive", func(t *testing.T) {
		content := "untouched"
		path := createTestFile(t, "cw_nomark.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("untouched", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// C-w without setting mark first: no change (mark not active by default)
		h.SendKeys("C-w")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "untouched")
		h.AssertCursorAt(t, 0, 0)
	})

	// --- Consecutive C-k then yank verifies accumulation ---
	t.Run("Ck_ThreeConsecutive_YankAll", func(t *testing.T) {
		content := "one\ntwo\nthree\nfour"
		path := createTestFile(t, "ck_three.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("one", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Kill "one" (content), newline, "two" (content), newline — 4 consecutive C-k
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)

		// Now at start of "three"
		h.AssertLineContains(t, 0, "three")

		// Move to end and yank: should get "one\ntwo\n" (accumulated)
		h.SendKeys("C-e")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-y")
		time.Sleep(200 * time.Millisecond)

		h.AssertScreenContains(t, "one")
		h.AssertScreenContains(t, "two")
	})

	// --- Yank replaces killed content at different position ---
	t.Run("KillAtOnePlaceYankAtAnother", func(t *testing.T) {
		content := "start\nmiddle\nend"
		path := createTestFile(t, "kill_yank_pos.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("start", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Kill "start" from line 0
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)

		// Move to end of "end" (now line 1 because line 0 is empty, line 1 is middle, line 2 is end)
		// Actually after C-k: line 0 is "" (empty), cursor at (0,0)
		// Navigate: C-n to line 1 (middle), C-n to line 2 (end), C-e to end of "end"
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)

		// Yank "start" at end of "end"
		h.SendKeys("C-y")
		time.Sleep(200 * time.Millisecond)
		h.AssertLineContains(t, 2, "endstart")
	})
}
