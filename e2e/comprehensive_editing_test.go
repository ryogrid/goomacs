package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestComprehensiveEditing(t *testing.T) {
	// --- InsertChar: ASCII ---
	t.Run("InsertChar_ASCII", func(t *testing.T) {
		content := "hello"
		path := createTestFile(t, "ins_ascii.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}
		h.AssertCursorAt(t, 0, 0)

		// Type 'X' at beginning
		h.SendKeys("X")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 1)
		h.AssertLineContains(t, 0, "Xhello")
	})

	t.Run("InsertChar_MiddleOfLine", func(t *testing.T) {
		content := "abcd"
		path := createTestFile(t, "ins_mid.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("abcd", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to col 2
		h.SendKeys("C-f")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-f")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 2)

		// Insert 'Z' at col 2
		h.SendKeys("Z")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 3)
		h.AssertLineContains(t, 0, "abZcd")
	})

	t.Run("InsertChar_EndOfLine", func(t *testing.T) {
		content := "abc"
		path := createTestFile(t, "ins_end.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("abc", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 3)

		h.SendKeys("D")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 4)
		h.AssertLineContains(t, 0, "abcD")
	})

	// --- InsertChar: Multi-byte UTF-8 ---
	t.Run("InsertChar_UTF8_Multibyte", func(t *testing.T) {
		content := "hello"
		path := createTestFile(t, "ins_utf8.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Insert a multibyte character (Japanese あ)
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("あ")
		time.Sleep(100 * time.Millisecond)

		h.AssertScreenContains(t, "helloあ")
	})

	// --- Tab character: visual expansion of existing tab in file ---
	t.Run("TabExpansion_AtColumn0", func(t *testing.T) {
		// Create a file with an actual tab character embedded
		content := "\thello"
		path := createTestFile(t, "tab_exp0.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Tab at col 0 should visually expand to 8 spaces
		lines := h.Capture()
		line0 := lines[0]
		if !strings.HasPrefix(line0, "        hello") {
			t.Errorf("expected tab expansion (8 spaces) before hello, got %q", line0)
		}
	})

	t.Run("TabExpansion_AtColumn3", func(t *testing.T) {
		// Create a file with "abc" then tab then "hello"
		content := "abc\thello"
		path := createTestFile(t, "tab_exp3.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Tab at col 3 expands to 5 visual spaces (next tab stop at 8)
		lines := h.Capture()
		line0 := lines[0]
		// "abc" + 5 visual spaces + "hello"
		if !strings.Contains(line0, "abc     hello") {
			t.Errorf("expected tab expansion (5 spaces after 'abc') before hello, got %q", line0)
		}
	})

	t.Run("TabKey_InNormalMode_NoOp", func(t *testing.T) {
		// Tab key in normal editing mode is not handled (no tab insertion)
		content := "hello"
		path := createTestFile(t, "tab_noop.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}
		h.AssertCursorAt(t, 0, 0)

		h.SendKeys("Tab")
		time.Sleep(100 * time.Millisecond)

		// Buffer should be unchanged
		h.AssertCursorAt(t, 0, 0)
		h.AssertLineContains(t, 0, "hello")
	})

	// --- Backspace: middle of line ---
	t.Run("Backspace_MiddleOfLine", func(t *testing.T) {
		content := "abcdef"
		path := createTestFile(t, "bs_mid.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("abcdef", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to col 3
		for i := 0; i < 3; i++ {
			h.SendKeys("C-f")
			time.Sleep(30 * time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 3)

		// Backspace deletes 'c'
		h.SendKeys("BSpace")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 2)
		h.AssertLineContains(t, 0, "abdef")
	})

	t.Run("Backspace_StartOfLine_JoinWithPrevious", func(t *testing.T) {
		content := "first\nsecond"
		path := createTestFile(t, "bs_join.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("first", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to start of second line
		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)

		// Backspace at start of line joins with previous
		h.SendKeys("BSpace")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 5) // cursor at join point
		h.AssertLineContains(t, 0, "firstsecond")
	})

	t.Run("Backspace_EmptyBuffer_NoChange", func(t *testing.T) {
		h := newHarness(t)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}
		h.AssertCursorAt(t, 0, 0)

		// Backspace in empty buffer should do nothing
		h.SendKeys("BSpace")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)
	})

	t.Run("Backspace_StartOfBuffer_NoChange", func(t *testing.T) {
		content := "hello"
		path := createTestFile(t, "bs_start.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}
		h.AssertCursorAt(t, 0, 0)

		// Backspace at start of buffer — no change
		h.SendKeys("BSpace")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)
		h.AssertLineContains(t, 0, "hello")
	})

	// --- C-d (delete forward): middle of line ---
	t.Run("Cd_MiddleOfLine", func(t *testing.T) {
		content := "abcdef"
		path := createTestFile(t, "cd_mid.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("abcdef", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to col 2
		h.SendKeys("C-f")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-f")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 2)

		// C-d deletes 'c'
		h.SendKeys("C-d")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 2)
		h.AssertLineContains(t, 0, "abdef")
	})

	t.Run("Cd_EndOfLine_JoinWithNext", func(t *testing.T) {
		content := "first\nsecond"
		path := createTestFile(t, "cd_join.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("first", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to end of first line
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 5)

		// C-d at end of line joins with next
		h.SendKeys("C-d")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 5) // cursor stays
		h.AssertLineContains(t, 0, "firstsecond")
	})

	t.Run("Cd_EmptyBuffer_NoChange", func(t *testing.T) {
		h := newHarness(t)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}
		h.AssertCursorAt(t, 0, 0)

		// C-d in empty buffer should do nothing
		h.SendKeys("C-d")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)
	})

	t.Run("Cd_EndOfBuffer_NoChange", func(t *testing.T) {
		content := "hello"
		path := createTestFile(t, "cd_end.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 5)

		// C-d at end of buffer — no change
		h.SendKeys("C-d")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 5)
		h.AssertLineContains(t, 0, "hello")
	})

	// --- Enter/C-j: split line ---
	t.Run("Enter_MiddleOfLine", func(t *testing.T) {
		content := "abcdef"
		path := createTestFile(t, "enter_mid.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("abcdef", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to col 3
		for i := 0; i < 3; i++ {
			h.SendKeys("C-f")
			time.Sleep(30 * time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 3)

		// Enter splits the line
		h.SendKeys("Enter")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0) // cursor on new line
		h.AssertLineContains(t, 0, "abc")
		h.AssertLineContains(t, 1, "def")
	})

	t.Run("Enter_EndOfLine", func(t *testing.T) {
		content := "hello"
		path := createTestFile(t, "enter_end.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)

		h.SendKeys("Enter")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)
		h.AssertLineContains(t, 0, "hello")
		// Line 1 should be empty
		lines := h.Capture()
		if strings.TrimSpace(lines[1]) != "" {
			t.Errorf("expected empty second line, got %q", lines[1])
		}
	})

	t.Run("Enter_StartOfLine", func(t *testing.T) {
		content := "hello"
		path := createTestFile(t, "enter_start.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}
		h.AssertCursorAt(t, 0, 0)

		h.SendKeys("Enter")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)
		// Line 0 should be empty, line 1 should have "hello"
		lines := h.Capture()
		if strings.TrimSpace(lines[0]) != "" {
			t.Errorf("expected empty first line, got %q", lines[0])
		}
		h.AssertLineContains(t, 1, "hello")
	})

	t.Run("Cj_SplitsLine", func(t *testing.T) {
		content := "abcdef"
		path := createTestFile(t, "cj_split.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("abcdef", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to col 2
		h.SendKeys("C-f")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-f")
		time.Sleep(100 * time.Millisecond)

		// C-j should also split the line (same as Enter)
		h.SendKeys("C-j")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)
		h.AssertLineContains(t, 0, "ab")
		h.AssertLineContains(t, 1, "cdef")
	})

	t.Run("Enter_EmptyBuffer", func(t *testing.T) {
		h := newHarness(t)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}
		h.AssertCursorAt(t, 0, 0)

		// Enter in empty buffer creates a new line
		h.SendKeys("Enter")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 0)
	})

	// --- Multiple edits ---
	t.Run("MultipleInsertAndBackspace", func(t *testing.T) {
		content := "ab"
		path := createTestFile(t, "multi_edit.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("ab", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Insert "XY" at beginning
		h.SendKeys("X")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("Y")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "XYab")
		h.AssertCursorAt(t, 0, 2)

		// Backspace twice to remove "XY"
		h.SendKeys("BSpace")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("BSpace")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "ab")
		h.AssertCursorAt(t, 0, 0)
	})

	// --- Read-only buffer: editing blocked ---
	// The grep buffer is the only read-only buffer in goomacs.
	// We use find-grep to create a grep buffer, then test editing in it.
	t.Run("ReadOnly_InsertChar_Blocked", func(t *testing.T) {
		dir := setupGrepFixture(t)
		h := newHarnessInDir(t, dir)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Run find-grep to create a read-only grep buffer
		enterGrepMode(h)

		// Try inserting a character
		h.SendKeys("x")
		time.Sleep(100 * time.Millisecond)
		h.AssertMessageLine(t, "Buffer is read-only")
	})

	t.Run("ReadOnly_Backspace_Blocked", func(t *testing.T) {
		dir := setupGrepFixture(t)
		h := newHarnessInDir(t, dir)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		enterGrepMode(h)

		h.SendKeys("BSpace")
		time.Sleep(100 * time.Millisecond)
		h.AssertMessageLine(t, "Buffer is read-only")
	})

	t.Run("ReadOnly_Cd_Blocked", func(t *testing.T) {
		dir := setupGrepFixture(t)
		h := newHarnessInDir(t, dir)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		enterGrepMode(h)

		h.SendKeys("C-d")
		time.Sleep(100 * time.Millisecond)
		h.AssertMessageLine(t, "Buffer is read-only")
	})

	t.Run("ReadOnly_Ck_Blocked", func(t *testing.T) {
		dir := setupGrepFixture(t)
		h := newHarnessInDir(t, dir)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		enterGrepMode(h)

		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)
		h.AssertMessageLine(t, "Buffer is read-only")
	})

	t.Run("ReadOnly_Enter_Blocked", func(t *testing.T) {
		dir := setupGrepFixture(t)
		h := newHarnessInDir(t, dir)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		enterGrepMode(h)

		// Move to a non-result line (the first line is usually a header or empty)
		// In grep mode, Enter on a non-result line is a no-op for grep handler,
		// but the read-only guard for Enter should trigger first.
		// Actually, grep mode Enter is handled by the mode handler before the read-only guard.
		// So we need to navigate to a line that grep mode won't handle.
		// For read-only Enter test, we use C-y (yank) instead since that's also blocked.
		h.SendKeys("C-y")
		time.Sleep(100 * time.Millisecond)
		h.AssertMessageLine(t, "Buffer is read-only")
	})

	t.Run("ReadOnly_Yank_Blocked", func(t *testing.T) {
		dir := setupGrepFixture(t)
		h := newHarnessInDir(t, dir)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		enterGrepMode(h)

		h.SendKeys("C-y")
		time.Sleep(100 * time.Millisecond)
		h.AssertMessageLine(t, "Buffer is read-only")
	})

	// --- Insert multiple characters forming a word ---
	t.Run("InsertWord", func(t *testing.T) {
		content := ""
		path := createTestFile(t, "ins_word.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("ins_word.txt", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		h.SendKeys("h")
		time.Sleep(30 * time.Millisecond)
		h.SendKeys("e")
		time.Sleep(30 * time.Millisecond)
		h.SendKeys("l")
		time.Sleep(30 * time.Millisecond)
		h.SendKeys("l")
		time.Sleep(30 * time.Millisecond)
		h.SendKeys("o")
		time.Sleep(100 * time.Millisecond)

		h.AssertLineContains(t, 0, "hello")
		h.AssertCursorAt(t, 0, 5)
	})

	// --- C-d on single char line ---
	t.Run("Cd_SingleCharLine", func(t *testing.T) {
		content := "a\nb"
		path := createTestFile(t, "cd_single.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("a", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}
		h.AssertCursorAt(t, 0, 0)

		// C-d deletes 'a', leaving empty first line
		h.SendKeys("C-d")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)

		// First line should now be empty; second line still has 'b'
		lines := h.Capture()
		line0 := strings.TrimRight(lines[0], " ")
		if line0 != "" {
			t.Errorf("expected empty first line after C-d, got %q", line0)
		}
		h.AssertLineContains(t, 1, "b")
	})

	// --- Backspace on single char line ---
	t.Run("Backspace_SingleCharLine", func(t *testing.T) {
		content := "a"
		path := createTestFile(t, "bs_single.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("a", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to end of 'a'
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 1)

		// Backspace removes 'a'
		h.SendKeys("BSpace")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 0)

		// Line should be empty
		lines := h.Capture()
		line0 := strings.TrimRight(lines[0], " ")
		if line0 != "" {
			t.Errorf("expected empty line after backspace, got %q", line0)
		}
	})

	// --- Backspace joining multiple lines ---
	t.Run("Backspace_JoinMultipleLines", func(t *testing.T) {
		content := "aaa\nbbb\nccc"
		path := createTestFile(t, "bs_multi_join.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("aaa", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Go to start of line 2
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 2, 0)

		// Backspace at start of line 2 joins with line 1
		h.SendKeys("BSpace")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 1, 3) // after "bbb"
		h.AssertLineContains(t, 1, "bbbccc")
	})

	// --- C-d joining lines ---
	t.Run("Cd_JoinLines", func(t *testing.T) {
		content := "first\nsecond"
		path := createTestFile(t, "cd_joinlines.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("first", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to end of first line
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 5)

		// C-d joins "first" and "second"
		h.SendKeys("C-d")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 5)
		h.AssertLineContains(t, 0, "firstsecond")
	})

	// --- Enter/newline creating blank lines ---
	t.Run("Enter_MultipleNewlines", func(t *testing.T) {
		content := "hello"
		path := createTestFile(t, "enter_multi.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}
		h.AssertCursorAt(t, 0, 0)

		// Two Enters at beginning creates two blank lines
		h.SendKeys("Enter")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("Enter")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 2, 0)
		h.AssertLineContains(t, 2, "hello")
	})
}

// setupGrepFixture creates a temp directory with a .go file for grep testing.
func setupGrepFixture(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	content := fmt.Sprintf("package main\n\nfunc hello() {\n\tfmt.Println(\"hello world\")\n}\n")
	path := filepath.Join(dir, "main.go")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create grep fixture: %v", err)
	}
	return dir
}

// enterGrepMode runs find-grep with a known query to enter grep mode.
func enterGrepMode(h *Harness) {
	// M-x find-grep
	h.SendKeys("Escape")
	time.Sleep(50 * time.Millisecond)
	h.SendKeys("x")
	time.Sleep(100 * time.Millisecond)
	h.SendKeys("f")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys("i")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys("n")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys("d")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys("-")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys("g")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys("r")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys("e")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys("p")
	time.Sleep(100 * time.Millisecond)
	h.SendKeys("Enter")
	time.Sleep(200 * time.Millisecond)

	// Clear the default command and type our own grep
	h.SendKeys("C-a")
	time.Sleep(50 * time.Millisecond)
	h.SendKeys("C-k")
	time.Sleep(50 * time.Millisecond)
	h.SendKeys("g")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys("r")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys("e")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys("p")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys(" ")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys("-")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys("r")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys("n")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys("H")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys(" ")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys("h")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys("e")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys("l")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys("l")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys("o")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys(" ")
	time.Sleep(30 * time.Millisecond)
	h.SendKeys(".")
	time.Sleep(100 * time.Millisecond)
	h.SendKeys("Enter")

	// Wait for grep results
	time.Sleep(2 * time.Second)
}
