package e2e

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestComprehensiveGrep(t *testing.T) {
	// --- Invocation ---
	t.Run("Invoke_DefaultPrompt", func(t *testing.T) {
		dir := setupGrepFixtures(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// M-x find-grep
		h.SendKeys("Escape")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("x")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("find-grep")
		time.Sleep(100 * time.Millisecond)
		h.SendKeys("Enter")
		time.Sleep(200 * time.Millisecond)

		// Verify minibuffer shows the default command with the prompt
		h.AssertMessageLine(t, "Run find-grep:")
	})

	// --- Enter on valid grep result ---
	t.Run("Enter_ValidResult", func(t *testing.T) {
		dir := setupGrepFixtures(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		invokeFindGrep(t, h, "grep -rnH hello ./hello.go | sort")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		// Press Enter on the first result to jump to source
		h.SendKeys("Enter")
		time.Sleep(500 * time.Millisecond)

		// Verify we jumped to the source file
		h.AssertStatusBar(t, "hello.go")
		// Verify content is visible
		h.AssertScreenContains(t, "hello")
	})

	// --- Enter on non-result line ---
	t.Run("Enter_NonResultLine", func(t *testing.T) {
		dir := setupGrepFixtures3Files(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Use a command that produces a non-parseable header line before grep results
		invokeFindGrep(t, h, "echo '=== Results ===' && grep -rnH marker . | sort")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		// Cursor starts at line 0 which is "=== Results ===" (not a grep result)
		h.SendKeys("Enter")
		time.Sleep(300 * time.Millisecond)

		// Should still be in grep buffer with message
		h.AssertStatusBar(t, "*grep*")
		h.AssertMessageLine(t, "No grep result on this line")
	})

	// --- n/p navigation ---
	t.Run("Navigation_N", func(t *testing.T) {
		dir := setupGrepFixtures3Files(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		invokeFindGrep(t, h, "grep -rnH marker . | sort")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		row1, _, err := h.CursorPosition()
		if err != nil {
			t.Fatalf("failed to get cursor: %v", err)
		}

		// Press n to go to next result
		h.SendKeys("n")
		time.Sleep(200 * time.Millisecond)

		row2, _, err := h.CursorPosition()
		if err != nil {
			t.Fatalf("failed to get cursor: %v", err)
		}

		if row2 <= row1 {
			t.Errorf("n did not advance cursor: was row %d, now row %d", row1, row2)
		}
	})

	t.Run("Navigation_P", func(t *testing.T) {
		dir := setupGrepFixtures3Files(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		invokeFindGrep(t, h, "grep -rnH marker . | sort")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		// Move to second result
		h.SendKeys("n")
		time.Sleep(200 * time.Millisecond)

		row1, _, err := h.CursorPosition()
		if err != nil {
			t.Fatalf("failed to get cursor: %v", err)
		}

		// Press p to go back
		h.SendKeys("p")
		time.Sleep(200 * time.Millisecond)

		row2, _, err := h.CursorPosition()
		if err != nil {
			t.Fatalf("failed to get cursor: %v", err)
		}

		if row2 >= row1 {
			t.Errorf("p did not move cursor back: was row %d, now row %d", row1, row2)
		}
	})

	// --- n at last result ---
	t.Run("N_AtLastResult", func(t *testing.T) {
		dir := setupGrepFixtures3Files(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		invokeFindGrep(t, h, "grep -rnH marker . | sort")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		// Navigate to last result (4 results total, start at first)
		for i := 0; i < 3; i++ {
			h.SendKeys("n")
			time.Sleep(200 * time.Millisecond)
		}

		row1, _, err := h.CursorPosition()
		if err != nil {
			t.Fatalf("failed to get cursor: %v", err)
		}

		// n at last result
		h.SendKeys("n")
		time.Sleep(200 * time.Millisecond)

		h.AssertMessageLine(t, "No more results")

		row2, _, err := h.CursorPosition()
		if err != nil {
			t.Fatalf("failed to get cursor: %v", err)
		}
		if row1 != row2 {
			t.Errorf("cursor moved after n at last result: was row %d, now row %d", row1, row2)
		}
	})

	// --- p at first result ---
	t.Run("P_AtFirstResult", func(t *testing.T) {
		dir := setupGrepFixtures3Files(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		invokeFindGrep(t, h, "grep -rnH marker . | sort")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		row1, _, err := h.CursorPosition()
		if err != nil {
			t.Fatalf("failed to get cursor: %v", err)
		}

		// p at first result
		h.SendKeys("p")
		time.Sleep(200 * time.Millisecond)

		h.AssertMessageLine(t, "No more results")

		row2, _, err := h.CursorPosition()
		if err != nil {
			t.Fatalf("failed to get cursor: %v", err)
		}
		if row1 != row2 {
			t.Errorf("cursor moved after p at first result: was row %d, now row %d", row1, row2)
		}
	})

	// --- M-n / M-p file navigation ---
	t.Run("FileNav_MN", func(t *testing.T) {
		dir := setupGrepFixtures3Files(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		invokeFindGrep(t, h, "grep -rnH marker . | sort")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		// Cursor on first result (a.txt)
		row1, _, err := h.CursorPosition()
		if err != nil {
			t.Fatalf("failed to get cursor: %v", err)
		}
		h.AssertLineContains(t, row1, "a.txt")

		// M-n to jump to next file (b.txt)
		h.SendKeys("M-n")
		time.Sleep(300 * time.Millisecond)

		row2, _, err := h.CursorPosition()
		if err != nil {
			t.Fatalf("failed to get cursor: %v", err)
		}
		h.AssertLineContains(t, row2, "b.txt")

		// M-n again to jump to c.txt
		h.SendKeys("M-n")
		time.Sleep(300 * time.Millisecond)

		row3, _, err := h.CursorPosition()
		if err != nil {
			t.Fatalf("failed to get cursor: %v", err)
		}
		h.AssertLineContains(t, row3, "c.txt")
	})

	t.Run("FileNav_MP", func(t *testing.T) {
		dir := setupGrepFixtures3Files(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		invokeFindGrep(t, h, "grep -rnH marker . | sort")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		// Navigate to c.txt first
		h.SendKeys("M-n")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("M-n")
		time.Sleep(300 * time.Millisecond)

		row1, _, err := h.CursorPosition()
		if err != nil {
			t.Fatalf("failed to get cursor: %v", err)
		}
		h.AssertLineContains(t, row1, "c.txt")

		// M-p to jump back to b.txt
		h.SendKeys("M-p")
		time.Sleep(300 * time.Millisecond)

		row2, _, err := h.CursorPosition()
		if err != nil {
			t.Fatalf("failed to get cursor: %v", err)
		}
		h.AssertLineContains(t, row2, "b.txt")

		// M-p to jump back to a.txt
		h.SendKeys("M-p")
		time.Sleep(300 * time.Millisecond)

		row3, _, err := h.CursorPosition()
		if err != nil {
			t.Fatalf("failed to get cursor: %v", err)
		}
		h.AssertLineContains(t, row3, "a.txt")
	})

	t.Run("FileNav_MN_AtLastFile", func(t *testing.T) {
		dir := setupGrepFixtures3Files(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		invokeFindGrep(t, h, "grep -rnH marker . | sort")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		// Navigate to last file (c.txt)
		h.SendKeys("M-n")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("M-n")
		time.Sleep(300 * time.Millisecond)

		row1, _, err := h.CursorPosition()
		if err != nil {
			t.Fatalf("failed to get cursor: %v", err)
		}

		// M-n at last file
		h.SendKeys("M-n")

		if err := h.WaitForContent("No more files", 5*time.Second); err != nil {
			t.Fatalf("No more files message did not appear: %v", err)
		}

		row2, _, err := h.CursorPosition()
		if err != nil {
			t.Fatalf("failed to get cursor: %v", err)
		}
		if row1 != row2 {
			t.Errorf("cursor moved after M-n at last file: was row %d, now row %d", row1, row2)
		}
	})

	t.Run("FileNav_MP_AtFirstFile", func(t *testing.T) {
		dir := setupGrepFixtures3Files(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		invokeFindGrep(t, h, "grep -rnH marker . | sort")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		row1, _, err := h.CursorPosition()
		if err != nil {
			t.Fatalf("failed to get cursor: %v", err)
		}

		// M-p at first file
		h.SendKeys("M-p")

		if err := h.WaitForContent("No more files", 5*time.Second); err != nil {
			t.Fatalf("No more files message did not appear: %v", err)
		}

		row2, _, err := h.CursorPosition()
		if err != nil {
			t.Fatalf("failed to get cursor: %v", err)
		}
		if row1 != row2 {
			t.Errorf("cursor moved after M-p at first file: was row %d, now row %d", row1, row2)
		}
	})

	// --- g (refresh) ---
	t.Run("Refresh_G", func(t *testing.T) {
		dir := setupGrepFixtures3Files(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		invokeFindGrep(t, h, "grep -rnH marker . | sort")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		// Verify initial results (4 results)
		h.AssertStatusBar(t, "Line 1/4")

		// Add a new matching line to a fixture file
		newPath := filepath.Join(dir, "a.txt")
		if err := os.WriteFile(newPath, []byte("line1\nmarker alpha\nline3\nmarker alpha2\n"), 0o644); err != nil {
			t.Fatalf("failed to update fixture: %v", err)
		}

		// Press g to refresh
		h.SendKeys("g")

		// Wait for refreshed results (now 5 lines)
		if err := h.WaitForContent("Line 1/5", 10*time.Second); err != nil {
			t.Fatalf("refresh did not show new results: %v", err)
		}

		h.AssertScreenContains(t, "marker alpha2")
	})

	// --- q (quit) ---
	t.Run("Quit_Q", func(t *testing.T) {
		dir := setupGrepFixtures(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		invokeFindGrep(t, h, "grep -rnH hello .")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		// Press q to close grep buffer
		h.SendKeys("q")
		time.Sleep(300 * time.Millisecond)

		// Should be back to *scratch*
		h.AssertStatusBar(t, "*scratch*")
	})

	// --- Read-only: typing characters blocked ---
	t.Run("ReadOnly_InsertBlocked", func(t *testing.T) {
		dir := setupGrepFixtures(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		invokeFindGrep(t, h, "grep -rnH hello .")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		// Try typing a regular character
		h.SendKeys("x")
		time.Sleep(200 * time.Millisecond)

		h.AssertMessageLine(t, "Buffer is read-only")
		// Still in grep buffer
		h.AssertStatusBar(t, "*grep*")
	})

	t.Run("ReadOnly_BackspaceBlocked", func(t *testing.T) {
		dir := setupGrepFixtures(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		invokeFindGrep(t, h, "grep -rnH hello .")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		h.SendKeys("BSpace")
		time.Sleep(200 * time.Millisecond)

		h.AssertMessageLine(t, "Buffer is read-only")
	})

	t.Run("ReadOnly_CkBlocked", func(t *testing.T) {
		dir := setupGrepFixtures(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		invokeFindGrep(t, h, "grep -rnH hello .")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		h.SendKeys("C-k")
		time.Sleep(200 * time.Millisecond)

		h.AssertMessageLine(t, "Buffer is read-only")
	})

	// --- Enter on correct line number ---
	t.Run("Enter_CorrectLineNumber", func(t *testing.T) {
		dir := setupGrepFixtures(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Grep for a specific line in hello.go — "hello world" is on line 4
		invokeFindGrep(t, h, "grep -rnH 'hello world' ./hello.go")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		// Press Enter to jump to source
		h.SendKeys("Enter")
		time.Sleep(500 * time.Millisecond)

		// Should be in hello.go at line 4 (0-indexed: row 3)
		h.AssertStatusBar(t, "hello.go")
		h.AssertStatusBar(t, "Line 4/")
	})

	// --- Navigate between n results across files ---
	t.Run("Navigation_NP_MultiFile", func(t *testing.T) {
		dir := setupGrepFixtures3Files(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		invokeFindGrep(t, h, "grep -rnH marker . | sort")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		// Should have 4 results: a.txt:2, b.txt:1, b.txt:3, c.txt:3
		h.AssertStatusBar(t, "Line 1/4")

		// Navigate forward through all results
		h.SendKeys("n")
		time.Sleep(200 * time.Millisecond)
		h.AssertStatusBar(t, "Line 2/4")

		h.SendKeys("n")
		time.Sleep(200 * time.Millisecond)
		h.AssertStatusBar(t, "Line 3/4")

		h.SendKeys("n")
		time.Sleep(200 * time.Millisecond)
		h.AssertStatusBar(t, "Line 4/4")

		// Navigate back
		h.SendKeys("p")
		time.Sleep(200 * time.Millisecond)
		h.AssertStatusBar(t, "Line 3/4")

		h.SendKeys("p")
		time.Sleep(200 * time.Millisecond)
		h.AssertStatusBar(t, "Line 2/4")

		h.SendKeys("p")
		time.Sleep(200 * time.Millisecond)
		h.AssertStatusBar(t, "Line 1/4")
	})

	// --- Enter then return to grep buffer ---
	t.Run("Enter_ThenReturnToGrep", func(t *testing.T) {
		dir := setupGrepFixtures(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		invokeFindGrep(t, h, "grep -rnH hello ./hello.go")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		// Jump to source
		h.SendKeys("Enter")
		time.Sleep(500 * time.Millisecond)
		h.AssertStatusBar(t, "hello.go")

		// Switch back to previous buffer (grep) with C-x b Enter
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("b")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("Enter")
		time.Sleep(300 * time.Millisecond)

		h.AssertStatusBar(t, "*grep*")
	})

	// --- Grep with no matches ---
	t.Run("NoMatches", func(t *testing.T) {
		dir := setupGrepFixtures(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		invokeFindGrep(t, h, "grep -rnH zzzznonexistent .")

		if err := h.WaitForContent("No matches found", 10*time.Second); err != nil {
			t.Fatalf("no-match message did not appear: %v", err)
		}

		// Should still be on *scratch* (no grep buffer created)
		h.AssertStatusBar(t, "*scratch*")
	})

	// --- Enter on multiple different results opens correct files ---
	t.Run("Enter_MultipleFiles", func(t *testing.T) {
		dir := setupGrepFixtures3Files(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		invokeFindGrep(t, h, "grep -rnH marker . | sort")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		// First result is a.txt — jump to it
		h.SendKeys("Enter")
		time.Sleep(500 * time.Millisecond)
		h.AssertStatusBar(t, "a.txt")

		// Go back to grep buffer
		h.SendKeys("C-x")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("b")
		time.Sleep(200 * time.Millisecond)
		h.SendKeys("Enter")
		time.Sleep(300 * time.Millisecond)
		h.AssertStatusBar(t, "*grep*")

		// Navigate to b.txt result and jump
		h.SendKeys("M-n")
		time.Sleep(300 * time.Millisecond)
		h.SendKeys("Enter")
		time.Sleep(500 * time.Millisecond)
		h.AssertStatusBar(t, "b.txt")
	})

	// --- Refresh with removed match ---
	t.Run("Refresh_RemovedMatch", func(t *testing.T) {
		dir := setupGrepFixtures3Files(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		invokeFindGrep(t, h, "grep -rnH marker . | sort")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		// Initially 4 results
		h.AssertStatusBar(t, "Line 1/4")

		// Remove all matches from c.txt
		cPath := filepath.Join(dir, "c.txt")
		if err := os.WriteFile(cPath, []byte("line1\nline2\nline3\n"), 0o644); err != nil {
			t.Fatalf("failed to update fixture: %v", err)
		}

		// Press g to refresh
		h.SendKeys("g")

		// Wait for refreshed results (now 3 lines)
		if err := h.WaitForContent("Line 1/3", 10*time.Second); err != nil {
			t.Fatalf("refresh did not show updated results: %v", err)
		}
	})

	// --- Functional after quit and re-invoke ---
	t.Run("Quit_ThenReinvoke", func(t *testing.T) {
		dir := setupGrepFixtures(t)
		h := newHarnessInDir(t, dir)

		if err := h.WaitForContent("*scratch*", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// First grep
		invokeFindGrep(t, h, "grep -rnH hello .")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not appear: %v", err)
		}

		// Quit
		h.SendKeys("q")
		time.Sleep(300 * time.Millisecond)
		h.AssertStatusBar(t, "*scratch*")

		// Re-invoke grep with different search
		invokeFindGrep(t, h, "grep -rnH greet .")

		if err := h.WaitForContent("*grep*", 10*time.Second); err != nil {
			t.Fatalf("grep buffer did not reappear: %v", err)
		}

		h.AssertScreenContains(t, "greet")
	})
}
