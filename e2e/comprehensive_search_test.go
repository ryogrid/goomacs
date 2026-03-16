package e2e

import (
	"strings"
	"testing"
	"time"
)

func TestComprehensiveSearch(t *testing.T) {
	// --- Search with no matches ---
	t.Run("ForwardSearch_NoMatches", func(t *testing.T) {
		content := "hello world\nfoo bar\nbaz qux"
		path := createTestFile(t, "search_nomatch.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello world", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("zzzzz")
		time.Sleep(300 * time.Millisecond)

		h.AssertMessageLine(t, "Failing I-search: zzzzz")

		// Cancel to restore
		h.SendKeys("C-g")
		time.Sleep(100 * time.Millisecond)
	})

	t.Run("BackwardSearch_NoMatches", func(t *testing.T) {
		content := "hello world\nfoo bar\nbaz qux"
		path := createTestFile(t, "search_nomatch_back.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello world", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to end
		h.SendKeys("Escape")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys(">")
		time.Sleep(200 * time.Millisecond)

		h.SendKeys("C-r")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("zzzzz")
		time.Sleep(300 * time.Millisecond)

		h.AssertMessageLine(t, "Failing I-search backward: zzzzz")

		h.SendKeys("C-g")
		time.Sleep(100 * time.Millisecond)
	})

	// --- Forward search wrap-around ---
	t.Run("ForwardSearch_WrapAround", func(t *testing.T) {
		content := "line one\nline two\ntarget here\nline four\nline five\nline six\nline seven\nline eight"
		path := createTestFile(t, "search_wrap.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("line one", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move cursor to line 5 (past "target here" on line 2)
		for i := 0; i < 5; i++ {
			h.SendKeys("C-n")
			time.Sleep(50 * time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 5, 0)

		// Forward search for "target" — should wrap around to line 2
		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("target")
		time.Sleep(300 * time.Millisecond)

		h.AssertMessageLine(t, "I-search: target")
		h.AssertCursorAt(t, 2, 0)

		// Accept
		h.SendKeys("Enter")
		time.Sleep(100 * time.Millisecond)
	})

	// --- Backward search wrap-around ---
	t.Run("BackwardSearch_WrapAround", func(t *testing.T) {
		content := "line one\nline two\nline three\nline four\nline five\ntarget here"
		path := createTestFile(t, "search_wrap_back.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("line one", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Cursor starts at (0,0), which is before "target here" on last line
		// Backward search should wrap to end and find "target" on line 5
		h.SendKeys("C-r")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("target")
		time.Sleep(300 * time.Millisecond)

		h.AssertMessageLine(t, "I-search backward: target")
		h.AssertCursorAt(t, 5, 0)

		h.SendKeys("Enter")
		time.Sleep(100 * time.Millisecond)
	})

	// --- Search cancel (C-g) restores position ---
	t.Run("SearchCancel_RestoresPosition", func(t *testing.T) {
		content := "aaa\nbbb\nccc\nddd"
		path := createTestFile(t, "search_cancel.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("aaa", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to line 1, col 1
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-f")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 1)

		// Search for "ddd" (moves to line 3)
		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("ddd")
		time.Sleep(300 * time.Millisecond)
		h.AssertCursorAt(t, 3, 0)

		// Cancel with C-g — cursor should return to (1, 1)
		h.SendKeys("C-g")
		time.Sleep(200 * time.Millisecond)
		h.AssertCursorAt(t, 1, 1)
		h.AssertMessageLine(t, "Quit")
	})

	// --- Search accept (Enter) keeps position ---
	t.Run("SearchAccept_KeepsPosition", func(t *testing.T) {
		content := "first\nsecond\nthird"
		path := createTestFile(t, "search_accept.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("first", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("third")
		time.Sleep(300 * time.Millisecond)
		h.AssertCursorAt(t, 2, 0)

		// Accept with Enter
		h.SendKeys("Enter")
		time.Sleep(200 * time.Millisecond)

		// Cursor stays at (2, 0)
		h.AssertCursorAt(t, 2, 0)

		// Message line should be clear (no I-search)
		lines := h.Capture()
		msgLine := strings.TrimSpace(lines[h.height-1])
		if strings.Contains(msgLine, "I-search") {
			t.Errorf("message line should not contain I-search after accept, got %q", msgLine)
		}
	})

	// --- Search exit via unrecognized key: key is re-processed ---
	t.Run("SearchExit_UnrecognizedKey_Reprocessed", func(t *testing.T) {
		content := "hello world\nsecond line"
		path := createTestFile(t, "search_exit_key.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello world", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Search for "second" — moves to (1, 0)
		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("second")
		time.Sleep(300 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)

		// Press C-e: exits search and moves to end of line
		h.SendKeys("C-e")
		time.Sleep(200 * time.Millisecond)

		// Should now be at end of "second line" (col 11)
		h.AssertCursorAt(t, 1, 11)

		// Verify search mode exited (no I-search prompt)
		lines := h.Capture()
		msgLine := strings.TrimSpace(lines[h.height-1])
		if strings.Contains(msgLine, "I-search") {
			t.Errorf("search should have exited, but message line shows: %q", msgLine)
		}
	})

	t.Run("SearchExit_Ca_MovesToBeginning", func(t *testing.T) {
		content := "hello world\nfoo bar baz"
		path := createTestFile(t, "search_exit_ca.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello world", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Search for "bar" — moves to (1, 4)
		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("bar")
		time.Sleep(300 * time.Millisecond)
		h.AssertCursorAt(t, 1, 4)

		// Press C-a: exits search and moves to beginning of line
		h.SendKeys("C-a")
		time.Sleep(200 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)
	})

	// --- Search with empty query ---
	t.Run("SearchEmptyQuery_NoChange", func(t *testing.T) {
		content := "hello world"
		path := createTestFile(t, "search_empty_q.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello world", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.AssertCursorAt(t, 0, 0)

		// Start search, immediately accept with Enter (empty query)
		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)
		h.AssertMessageLine(t, "I-search:")

		h.SendKeys("Enter")
		time.Sleep(200 * time.Millisecond)

		// Cursor should stay at (0, 0) — no crash
		h.AssertCursorAt(t, 0, 0)
		h.AssertLineContains(t, 0, "hello world")
	})

	t.Run("SearchEmptyQuery_CsAgain_NoCrash", func(t *testing.T) {
		content := "hello world"
		path := createTestFile(t, "search_empty_cs.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello world", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Start search with empty query, then press C-s again (no query to search)
		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)

		// Should not crash, cancel to exit
		h.SendKeys("C-g")
		time.Sleep(200 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)
		h.AssertLineContains(t, 0, "hello world")
	})

	// --- Backspace in search re-searches from original position ---
	t.Run("SearchBackspace_ResearchesFromOriginal", func(t *testing.T) {
		content := "apple banana\ncherry date\napricot fig"
		path := createTestFile(t, "search_bspace.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("apple banana", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Start search
		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)

		// Type "ap" — matches "apple" at (0, 0)
		h.SendKeys("a")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("p")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)
		h.AssertMessageLine(t, "I-search: ap")

		// Type "r" — "apr" matches "apricot" at (2, 0)
		h.SendKeys("r")
		time.Sleep(200 * time.Millisecond)
		h.AssertCursorAt(t, 2, 0)
		h.AssertMessageLine(t, "I-search: apr")

		// Backspace — removes "r", query becomes "ap", re-searches from ORIGINAL position
		// Should find "apple" at (0, 0) again
		h.SendKeys("BSpace")
		time.Sleep(200 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)
		h.AssertMessageLine(t, "I-search: ap")

		h.SendKeys("C-g")
		time.Sleep(100 * time.Millisecond)
	})

	t.Run("SearchBackspace_AllCharsDeleted_RestoresOriginal", func(t *testing.T) {
		content := "aaa\nbbb"
		path := createTestFile(t, "search_bspace_all.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("aaa", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to line 1
		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)

		// Start search, type "a" (moves to line 0 via wrap)
		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("a")
		time.Sleep(200 * time.Millisecond)

		// Backspace removes "a", query empty — cursor returns to original (1, 0)
		h.SendKeys("BSpace")
		time.Sleep(200 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)
		h.AssertMessageLine(t, "I-search:")

		h.SendKeys("C-g")
		time.Sleep(100 * time.Millisecond)
	})

	// --- Switching direction: C-s then C-r within same search ---
	t.Run("SwitchDirection_ForwardThenBackward", func(t *testing.T) {
		content := "apple one\nbanana two\napple three\ncherry four\napple five"
		path := createTestFile(t, "search_switch_dir.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("apple one", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Start forward search for "apple"
		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("apple")
		time.Sleep(300 * time.Millisecond)
		// First match: "apple" at (0, 0)
		h.AssertCursorAt(t, 0, 0)

		// C-s to advance to next match: (2, 0)
		h.SendKeys("C-s")
		time.Sleep(300 * time.Millisecond)
		h.AssertCursorAt(t, 2, 0)

		// C-s again: (4, 0)
		h.SendKeys("C-s")
		time.Sleep(300 * time.Millisecond)
		h.AssertCursorAt(t, 4, 0)

		// Now switch to backward with C-r: should find previous "apple" at (2, 0)
		h.SendKeys("C-r")
		time.Sleep(300 * time.Millisecond)
		h.AssertCursorAt(t, 2, 0)
		h.AssertMessageLine(t, "I-search backward: apple")

		h.SendKeys("Enter")
		time.Sleep(100 * time.Millisecond)
	})

	t.Run("SwitchDirection_BackwardThenForward", func(t *testing.T) {
		content := "xxx\nfoo\nxxx\nfoo\nxxx"
		path := createTestFile(t, "search_switch_dir2.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("xxx", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to end
		h.SendKeys("Escape")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys(">")
		time.Sleep(200 * time.Millisecond)

		// Start backward search for "foo"
		h.SendKeys("C-r")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("foo")
		time.Sleep(300 * time.Millisecond)
		// Should find "foo" at line 3
		h.AssertCursorAt(t, 3, 0)

		// C-r again to find previous "foo" at line 1
		h.SendKeys("C-r")
		time.Sleep(300 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)

		// Switch to forward with C-s: should find next "foo" at line 3
		h.SendKeys("C-s")
		time.Sleep(300 * time.Millisecond)
		h.AssertCursorAt(t, 3, 0)

		h.SendKeys("Enter")
		time.Sleep(100 * time.Millisecond)
	})

	// --- Forward search next (C-s C-s) ---
	t.Run("ForwardSearch_NextMatch", func(t *testing.T) {
		content := "word one\nword two\nword three"
		path := createTestFile(t, "search_next.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("word one", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("word")
		time.Sleep(300 * time.Millisecond)
		// First match at (0, 0)
		h.AssertCursorAt(t, 0, 0)

		// C-s to next
		h.SendKeys("C-s")
		time.Sleep(300 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)

		// C-s to next
		h.SendKeys("C-s")
		time.Sleep(300 * time.Millisecond)
		h.AssertCursorAt(t, 2, 0)

		// C-s wraps around to first match
		h.SendKeys("C-s")
		time.Sleep(300 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)

		h.SendKeys("Enter")
		time.Sleep(100 * time.Millisecond)
	})

	// --- Backward search previous (C-r C-r) ---
	t.Run("BackwardSearch_PreviousMatch", func(t *testing.T) {
		content := "word one\nword two\nword three"
		path := createTestFile(t, "search_prev.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("word one", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to end
		h.SendKeys("Escape")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys(">")
		time.Sleep(200 * time.Millisecond)

		h.SendKeys("C-r")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("word")
		time.Sleep(300 * time.Millisecond)
		// Last match: "word" at (2, 0)
		h.AssertCursorAt(t, 2, 0)

		// C-r to previous
		h.SendKeys("C-r")
		time.Sleep(300 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)

		// C-r to previous
		h.SendKeys("C-r")
		time.Sleep(300 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)

		// C-r wraps around to last match
		h.SendKeys("C-r")
		time.Sleep(300 * time.Millisecond)
		h.AssertCursorAt(t, 2, 0)

		h.SendKeys("Enter")
		time.Sleep(100 * time.Millisecond)
	})

	// --- Search for text in the middle of a line ---
	t.Run("ForwardSearch_MiddleOfLine", func(t *testing.T) {
		content := "the quick brown fox"
		path := createTestFile(t, "search_middle.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("the quick", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("brown")
		time.Sleep(300 * time.Millisecond)
		// "brown" starts at col 10
		h.AssertCursorAt(t, 0, 10)

		h.SendKeys("Enter")
		time.Sleep(100 * time.Millisecond)
	})

	// --- Incremental search: cursor moves as you type ---
	t.Run("IncrementalSearch_CursorUpdatesPerChar", func(t *testing.T) {
		content := "abcdef\nabc123\nabcxyz"
		path := createTestFile(t, "search_incr.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("abcdef", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)

		// Type "a" — matches at (0, 0)
		h.SendKeys("a")
		time.Sleep(150 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)

		// Type "b" — "ab" still matches at (0, 0)
		h.SendKeys("b")
		time.Sleep(150 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)

		// Type "c" — "abc" still at (0, 0)
		h.SendKeys("c")
		time.Sleep(150 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)

		// Type "x" — "abcx" matches "abcxyz" at (2, 0)
		h.SendKeys("x")
		time.Sleep(200 * time.Millisecond)
		h.AssertCursorAt(t, 2, 0)
		h.AssertMessageLine(t, "I-search: abcx")

		h.SendKeys("C-g")
		time.Sleep(100 * time.Millisecond)
	})

	// --- Search then edit: search exits, edit works ---
	t.Run("SearchThenEdit_Works", func(t *testing.T) {
		content := "hello world"
		path := createTestFile(t, "search_edit.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello world", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Search for "world" — cursor at (0, 6)
		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("world")
		time.Sleep(300 * time.Millisecond)
		h.AssertCursorAt(t, 0, 6)

		// Accept search
		h.SendKeys("Enter")
		time.Sleep(200 * time.Millisecond)

		// Now type text — should insert at cursor position
		h.SendKeys("X")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "hello Xworld")
	})

	// --- Multiple searches in sequence ---
	t.Run("MultipleSearches_Sequential", func(t *testing.T) {
		content := "alpha\nbeta\ngamma"
		path := createTestFile(t, "search_multi.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("alpha", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// First search for "beta"
		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("beta")
		time.Sleep(300 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)
		h.SendKeys("Enter")
		time.Sleep(200 * time.Millisecond)

		// Second search for "gamma"
		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("gamma")
		time.Sleep(300 * time.Millisecond)
		h.AssertCursorAt(t, 2, 0)
		h.SendKeys("Enter")
		time.Sleep(100 * time.Millisecond)
	})

	// --- Search cancel does not modify buffer ---
	t.Run("SearchCancel_BufferUnchanged", func(t *testing.T) {
		content := "do not change"
		path := createTestFile(t, "search_nomod.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("do not change", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("change")
		time.Sleep(300 * time.Millisecond)
		h.SendKeys("C-g")
		time.Sleep(200 * time.Millisecond)

		// Buffer content should be unchanged
		h.AssertLineContains(t, 0, "do not change")
		h.AssertCursorAt(t, 0, 0)

		// Status bar should not show [Modified]
		lines := h.Capture()
		statusBar := lines[h.height-2]
		if strings.Contains(statusBar, "Modified") {
			t.Errorf("buffer should not be modified after search cancel, status: %q", statusBar)
		}
	})

	// --- Backward search from middle of buffer ---
	t.Run("BackwardSearch_FromMiddle", func(t *testing.T) {
		content := "first target\nsecond line\nthird line"
		path := createTestFile(t, "search_back_mid.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("first target", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to line 2
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 2, 0)

		// Backward search for "target"
		h.SendKeys("C-r")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("target")
		time.Sleep(300 * time.Millisecond)

		// Should find at (0, 6)
		h.AssertCursorAt(t, 0, 6)
		h.AssertMessageLine(t, "I-search backward: target")

		h.SendKeys("Enter")
		time.Sleep(100 * time.Millisecond)
	})

	// --- Failing search then backspace restores ---
	t.Run("FailingSearch_BackspaceRestores", func(t *testing.T) {
		content := "apple pie"
		path := createTestFile(t, "search_fail_bs.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("apple pie", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)

		// Type "applez" — no match
		h.SendKeys("applez")
		time.Sleep(300 * time.Millisecond)
		h.AssertMessageLine(t, "Failing I-search: applez")

		// Backspace removes 'z' — "apple" matches at (0, 0)
		h.SendKeys("BSpace")
		time.Sleep(200 * time.Millisecond)
		h.AssertMessageLine(t, "I-search: apple")
		h.AssertCursorAt(t, 0, 0)

		h.SendKeys("Enter")
		time.Sleep(100 * time.Millisecond)
	})
}
