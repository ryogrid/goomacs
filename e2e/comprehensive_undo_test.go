package e2e

import (
	"testing"
	"time"
)

func TestComprehensiveUndo(t *testing.T) {
	// --- Undo of character insertion ---
	t.Run("Undo_CharInsertion", func(t *testing.T) {
		content := "hello"
		path := createTestFile(t, "undo_char.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to end and insert a character
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("X")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "helloX")

		// Undo should remove the inserted character
		h.SendKeys("C-_")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineEquals(t, 0, "hello")
		h.AssertMessageLine(t, "Undo")
	})

	// --- Undo of deletion (C-d) ---
	t.Run("Undo_DeleteChar", func(t *testing.T) {
		content := "abcde"
		path := createTestFile(t, "undo_del.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("abcde", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Delete first character with C-d
		h.SendKeys("C-d")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineEquals(t, 0, "bcde")

		// Undo should restore 'a'
		h.SendKeys("C-_")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineEquals(t, 0, "abcde")
		h.AssertCursorAt(t, 0, 0)
	})

	// --- Undo of kill line (C-k) ---
	t.Run("Undo_KillLine", func(t *testing.T) {
		content := "kill this line"
		path := createTestFile(t, "undo_kill.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("kill this line", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Kill the line
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineEquals(t, 0, "")

		// Undo should restore the line
		h.SendKeys("C-_")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "kill this line")
	})

	// --- Undo of region kill (C-w) ---
	t.Run("Undo_RegionKill", func(t *testing.T) {
		content := "ABCDEF"
		path := createTestFile(t, "undo_region.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("ABCDEF", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Set mark, move forward 3 chars, kill region "ABC"
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
	})

	// --- Undo of yank (C-y) ---
	t.Run("Undo_Yank", func(t *testing.T) {
		content := "hello world"
		path := createTestFile(t, "undo_yank.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("hello world", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Kill "hello world" first to put it in kill ring
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineEquals(t, 0, "")

		// Yank it back
		h.SendKeys("C-y")
		time.Sleep(200 * time.Millisecond)
		h.AssertLineContains(t, 0, "hello world")

		// Undo the yank — should remove the yanked text
		h.SendKeys("C-_")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineEquals(t, 0, "")
	})

	// --- Undo of newline insertion (Enter) ---
	t.Run("Undo_NewlineInsertion", func(t *testing.T) {
		content := "oneline"
		path := createTestFile(t, "undo_newline.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("oneline", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to middle and insert newline
		for i := 0; i < 3; i++ {
			h.SendKeys("C-f")
			time.Sleep(30 * time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("Enter")
		time.Sleep(100 * time.Millisecond)

		// Verify the line was split
		h.AssertLineContains(t, 0, "one")
		h.AssertLineContains(t, 1, "line")

		// Undo should join the lines back
		h.SendKeys("C-_")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "oneline")
	})

	// --- Undo with empty stack (no changes) ---
	t.Run("Undo_EmptyStack_NoCrash", func(t *testing.T) {
		h := newHarness(t)
		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Undo with no changes should show message, not crash
		h.SendKeys("C-_")
		time.Sleep(100 * time.Millisecond)
		h.AssertMessageLine(t, "No further undo information")

		// Editor should still be functional
		h.SendKeys("a")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "a")
	})

	// --- Undo after save (undo is still available) ---
	t.Run("Undo_AfterSave", func(t *testing.T) {
		content := "original"
		path := createTestFile(t, "undo_save.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("original", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Type some text
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("!")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "original!")

		// Save with C-x C-s
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-s")
		time.Sleep(200 * time.Millisecond)
		h.AssertMessageLine(t, "Saved")

		// Undo should still work after save
		h.SendKeys("C-_")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineEquals(t, 0, "original")
		h.AssertMessageLine(t, "Undo")

		// Status bar should show [Modified] since we undid past the save point
		h.AssertStatusBar(t, "[Modified]")
	})

	// --- Multiple undos in sequence ---
	t.Run("Undo_MultipleSequential", func(t *testing.T) {
		content := "base"
		path := createTestFile(t, "undo_multi.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("base", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Insert 3 characters one by one
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("1")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("2")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("3")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "base123")

		// Undo 3 times should restore to "base"
		h.SendKeys("C-_")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "base12")

		h.SendKeys("C-_")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "base1")

		h.SendKeys("C-_")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineEquals(t, 0, "base")
	})

	// --- Undo of backspace ---
	t.Run("Undo_Backspace", func(t *testing.T) {
		content := "test"
		path := createTestFile(t, "undo_bs.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("test", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to end and backspace
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("BSpace")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineEquals(t, 0, "tes")

		// Undo should restore the deleted char
		h.SendKeys("C-_")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineEquals(t, 0, "test")
	})

	// --- Undo restores cursor position ---
	t.Run("Undo_RestoresCursorPosition", func(t *testing.T) {
		content := "abcd"
		path := createTestFile(t, "undo_cursor.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("abcd", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to col 2 and insert 'X'
		h.SendKeys("C-f")
		time.Sleep(30 * time.Millisecond)
		h.SendKeys("C-f")
		time.Sleep(100 * time.Millisecond)
		h.AssertCursorAt(t, 0, 2)

		h.SendKeys("X")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "abXcd")
		h.AssertCursorAt(t, 0, 3)

		// Undo should restore cursor to col 2 (the position before insertion)
		h.SendKeys("C-_")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineEquals(t, 0, "abcd")
		h.AssertCursorAt(t, 0, 2)
	})

	// --- Undo of kill line at end of line (join) ---
	t.Run("Undo_KillLineJoin", func(t *testing.T) {
		content := "first\nsecond"
		path := createTestFile(t, "undo_killjoin.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("first", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Move to end of first line and C-k (join with next)
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("C-k")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "firstsecond")

		// Undo should restore the two separate lines
		h.SendKeys("C-_")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "first")
		h.AssertLineContains(t, 1, "second")
	})

	// --- Multiple undo then new edit (undo stack truncated) ---
	t.Run("Undo_ThenNewEdit", func(t *testing.T) {
		content := "abc"
		path := createTestFile(t, "undo_newedit.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("abc", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Insert X at end
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("X")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "abcX")

		// Undo
		h.SendKeys("C-_")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineEquals(t, 0, "abc")

		// Now insert Y instead
		h.SendKeys("C-e")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("Y")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "abcY")

		// Undo should undo Y, getting back to "abc"
		h.SendKeys("C-_")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineEquals(t, 0, "abc")
	})

	// --- 10 edits, 10 undos ---
	t.Run("Undo_TenEdits_TenUndos", func(t *testing.T) {
		content := ""
		path := createTestFile(t, "undo_ten.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("undo_ten.txt", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Insert digits 0-9
		for i := 0; i < 10; i++ {
			h.SendKeys(string(rune('0' + i)))
			time.Sleep(50 * time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "0123456789")

		// Undo 10 times
		for i := 0; i < 10; i++ {
			h.SendKeys("C-_")
			time.Sleep(50 * time.Millisecond)
		}
		time.Sleep(100 * time.Millisecond)

		// Buffer should be empty again
		h.AssertLineEquals(t, 0, "")
	})

	// --- Undo of multi-line region kill restores all lines ---
	t.Run("Undo_MultiLineRegionKill", func(t *testing.T) {
		content := "line1\nline2\nline3\nline4"
		path := createTestFile(t, "undo_mlregion.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("line1", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Set mark at beginning, move down 2 lines
		h.SendKeys("C-Space")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(100 * time.Millisecond)

		// Kill region (line1\nline2\n)
		h.SendKeys("C-w")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "line3")

		// Undo should restore all killed lines
		h.SendKeys("C-_")
		time.Sleep(100 * time.Millisecond)
		h.AssertLineContains(t, 0, "line1")
		h.AssertLineContains(t, 1, "line2")
		h.AssertLineContains(t, 2, "line3")
		h.AssertLineContains(t, 3, "line4")
	})

	// --- Undo message appears ---
	t.Run("Undo_MessageShown", func(t *testing.T) {
		content := "msg"
		path := createTestFile(t, "undo_msg.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("msg", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// Make an edit
		h.SendKeys("a")
		time.Sleep(100 * time.Millisecond)

		// Undo should show "Undo" message
		h.SendKeys("C-_")
		time.Sleep(100 * time.Millisecond)
		h.AssertMessageLine(t, "Undo")
	})

	// --- Undo empty stack message ---
	t.Run("Undo_EmptyStackMessage", func(t *testing.T) {
		content := "text"
		path := createTestFile(t, "undo_emptymsg.txt", content)
		h := newHarness(t, path)
		if err := h.WaitForContent("text", 5*time.Second); err != nil {
			t.Fatalf("file did not open: %v", err)
		}

		// No edits, just undo — should show "No further undo information"
		h.SendKeys("C-_")
		time.Sleep(100 * time.Millisecond)
		h.AssertMessageLine(t, "No further undo information")
	})
}
