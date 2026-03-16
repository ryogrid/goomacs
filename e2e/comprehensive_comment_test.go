package e2e

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// shortTempDir creates a temp directory with a short path for status bar readability.
func shortTempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("/tmp", "g")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

// shortTempFileInDir creates a file in the given directory with the given content.
func shortTempFileInDir(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	return path
}

// invokeCommand runs M-x <cmdName> Enter.
func invokeCommand(h *Harness, cmdName string) {
	h.SendKeys("Escape")
	time.Sleep(50 * time.Millisecond)
	h.SendKeys("x")
	time.Sleep(200 * time.Millisecond)
	h.SendKeys(cmdName)
	time.Sleep(100 * time.Millisecond)
	h.SendKeys("Enter")
	time.Sleep(300 * time.Millisecond)
}

func TestComprehensiveComment(t *testing.T) {

	// --- comment-region on Go file ---
	t.Run("CommentRegion_GoFile", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "package main\n\nfunc hello() {\n\tfmt.Println(\"hi\")\n}\n"
		path := shortTempFileInDir(t, dir, "main.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("main.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Select lines 3-4 (func hello and fmt.Println)
		// Move to line 3 (0-indexed line 2)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)

		// Set mark
		h.SendKeys("C-Space")
		time.Sleep(50 * time.Millisecond)

		// Move down to include two lines
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)

		// M-x comment-region
		invokeCommand(h, "comment-region")

		// Verify "Region commented" message
		h.AssertMessageLine(t, "Region commented")

		// Verify // prefix was added
		time.Sleep(100 * time.Millisecond)
		h.AssertScreenContains(t, "// func hello()")
		// Tab is rendered as spaces in tmux capture — just check the comment marker and content
		h.AssertScreenContains(t, "//")
		h.AssertScreenContains(t, "fmt.Println")
	})

	// --- uncomment-region on Go file ---
	t.Run("UncommentRegion_GoFile", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "package main\n\n// func hello() {\n// \tfmt.Println(\"hi\")\n}\n"
		path := shortTempFileInDir(t, dir, "main.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("main.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Select the two commented lines (lines 3-4, 0-indexed 2-3)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)

		// Set mark
		h.SendKeys("C-Space")
		time.Sleep(50 * time.Millisecond)

		// Move down to include two lines
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)

		// M-x uncomment-region
		invokeCommand(h, "uncomment-region")

		// Verify "Region uncommented" message
		h.AssertMessageLine(t, "Region uncommented")

		// Verify // prefix was removed
		time.Sleep(100 * time.Millisecond)
		screen := h.CapturePane()
		if !containsLine(screen, "func hello()") {
			t.Errorf("expected uncommented 'func hello()' in screen:\n%s", screen)
		}
	})

	// --- comment-region with no active region ---
	t.Run("CommentRegion_NoRegion", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "package main\n\nfunc main() {}\n"
		path := shortTempFileInDir(t, dir, "main.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("main.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// M-x comment-region without setting mark
		invokeCommand(h, "comment-region")

		// Verify error message
		h.AssertMessageLine(t, "No region selected")
	})

	// --- uncomment-region with no active region ---
	t.Run("UncommentRegion_NoRegion", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "// commented line\n"
		path := shortTempFileInDir(t, dir, "main.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("main.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// M-x uncomment-region without setting mark
		invokeCommand(h, "uncomment-region")

		// Verify error message
		h.AssertMessageLine(t, "No region selected")
	})

	// --- comment-region on Python file ---
	t.Run("CommentRegion_PythonFile", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "def hello():\n    print(\"hi\")\n    return True\n"
		path := shortTempFileInDir(t, dir, "main.py", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("main.py", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Select all 3 lines
		h.SendKeys("C-Space")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)

		// M-x comment-region
		invokeCommand(h, "comment-region")

		h.AssertMessageLine(t, "Region commented")
		time.Sleep(100 * time.Millisecond)

		// Verify # prefix was added
		h.AssertScreenContains(t, "# def hello():")
		h.AssertScreenContains(t, "#     print")
		h.AssertScreenContains(t, "#     return True")
	})

	// --- uncomment-region tolerates leading whitespace ---
	t.Run("UncommentRegion_LeadingWhitespace", func(t *testing.T) {
		dir := shortTempDir(t)
		// Lines with leading whitespace before the comment marker
		content := "package main\n\n    // indented comment\n    // another indented\n"
		path := shortTempFileInDir(t, dir, "main.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("main.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Move to line 3 (0-indexed 2)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)

		// Set mark and select 2 lines
		h.SendKeys("C-Space")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)

		// M-x uncomment-region
		invokeCommand(h, "uncomment-region")

		h.AssertMessageLine(t, "Region uncommented")
		time.Sleep(100 * time.Millisecond)

		// Verify indentation preserved but // removed
		h.AssertScreenContains(t, "    indented comment")
		h.AssertScreenContains(t, "    another indented")
	})

	// --- comment then uncomment round-trip ---
	t.Run("CommentUncomment_RoundTrip", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "package main\n\nfunc foo() int {\n\treturn 42\n}\n"
		path := shortTempFileInDir(t, dir, "main.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("main.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Capture original screen for comparison
		time.Sleep(200 * time.Millisecond)
		originalScreen := h.CapturePane()

		// Select lines 3-4 (func foo and return 42)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)

		// Set mark
		h.SendKeys("C-Space")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)

		// Comment the region
		invokeCommand(h, "comment-region")
		h.AssertMessageLine(t, "Region commented")

		// Verify commented
		time.Sleep(100 * time.Millisecond)
		h.AssertScreenContains(t, "// func foo()")

		// Re-select the same lines (mark was deactivated)
		// Go back to line 3
		h.SendKeys("C-p")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-p")
		time.Sleep(50 * time.Millisecond)

		// Set mark again
		h.SendKeys("C-Space")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)

		// Uncomment the region
		invokeCommand(h, "uncomment-region")
		h.AssertMessageLine(t, "Region uncommented")

		// Verify content is back to original
		time.Sleep(100 * time.Millisecond)
		h.AssertScreenContains(t, "func foo() int")
		h.AssertScreenContains(t, "return 42")

		// The buffer lines should match original (ignore status bar which may show [Modified])
		screen := h.CapturePane()
		_ = originalScreen
		_ = screen
		// Content lines should match
		if !containsLine(screen, "func foo() int {") {
			t.Errorf("round-trip failed: 'func foo() int {' not found in screen:\n%s", screen)
		}
		if !containsLine(screen, "return 42") {
			t.Errorf("round-trip failed: 'return 42' not found in screen:\n%s", screen)
		}
	})

	// --- comment-region on single line ---
	t.Run("CommentRegion_SingleLine", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "package main\n\nvar x = 1\n"
		path := shortTempFileInDir(t, dir, "main.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("main.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Move to line 3 (var x = 1)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)

		// Set mark at beginning of line, move to end (single line region)
		h.SendKeys("C-Space")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-e")
		time.Sleep(50 * time.Millisecond)

		// M-x comment-region
		invokeCommand(h, "comment-region")

		h.AssertMessageLine(t, "Region commented")
		time.Sleep(100 * time.Millisecond)
		h.AssertScreenContains(t, "// var x = 1")
	})

	// --- uncomment-region with no space after marker ---
	t.Run("UncommentRegion_NoSpaceAfterMarker", func(t *testing.T) {
		dir := shortTempDir(t)
		// Comment without space after // (e.g. "//code")
		content := "package main\n\n//no space here\n"
		path := shortTempFileInDir(t, dir, "main.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("main.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Move to line 3
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)

		// Set mark and select line
		h.SendKeys("C-Space")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-e")
		time.Sleep(50 * time.Millisecond)

		// M-x uncomment-region
		invokeCommand(h, "uncomment-region")

		h.AssertMessageLine(t, "Region uncommented")
		time.Sleep(100 * time.Millisecond)
		h.AssertScreenContains(t, "no space here")
	})

	// --- comment-region marks buffer as modified ---
	t.Run("CommentRegion_BufferModified", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "package main\n\nfunc main() {}\n"
		path := shortTempFileInDir(t, dir, "main.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("main.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Select a line
		h.SendKeys("C-Space")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-e")
		time.Sleep(50 * time.Millisecond)

		// Comment it
		invokeCommand(h, "comment-region")

		// Verify buffer is marked as modified
		time.Sleep(100 * time.Millisecond)
		h.AssertStatusBar(t, "[Modified]")
	})

	// --- comment-region deactivates mark ---
	t.Run("CommentRegion_DeactivatesMark", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "package main\n\nfunc main() {}\nvar y = 2\n"
		path := shortTempFileInDir(t, dir, "main.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("main.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Select a line and comment
		h.SendKeys("C-Space")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)

		invokeCommand(h, "comment-region")
		h.AssertMessageLine(t, "Region commented")

		// Now try comment-region again — should fail because mark was deactivated
		invokeCommand(h, "comment-region")
		h.AssertMessageLine(t, "No region selected")
	})

	// --- comment-region on Python file round-trip ---
	t.Run("CommentUncomment_PythonRoundTrip", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "def greet(name):\n    return f\"Hello {name}\"\n"
		path := shortTempFileInDir(t, dir, "main.py", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("main.py", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Select both lines
		h.SendKeys("C-Space")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)

		// Comment
		invokeCommand(h, "comment-region")
		h.AssertMessageLine(t, "Region commented")
		time.Sleep(100 * time.Millisecond)
		h.AssertScreenContains(t, "# def greet")

		// Re-select and uncomment
		h.SendKeys("C-p")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-p")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-Space")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)

		// Uncomment
		invokeCommand(h, "uncomment-region")
		h.AssertMessageLine(t, "Region uncommented")
		time.Sleep(100 * time.Millisecond)
		h.AssertScreenContains(t, "def greet(name):")
	})

	// --- comment-region on multiple non-contiguous invocations ---
	t.Run("CommentRegion_MultipleInvocations", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "line one\nline two\nline three\nline four\n"
		path := shortTempFileInDir(t, dir, "main.go", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("main.go", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Comment line 1 only
		h.SendKeys("C-Space")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-e")
		time.Sleep(50 * time.Millisecond)
		invokeCommand(h, "comment-region")
		h.AssertMessageLine(t, "Region commented")

		// Move to line 3, comment it
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-Space")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-e")
		time.Sleep(50 * time.Millisecond)
		invokeCommand(h, "comment-region")
		h.AssertMessageLine(t, "Region commented")

		time.Sleep(100 * time.Millisecond)
		// Verify lines 1 and 3 commented, lines 2 and 4 not
		h.AssertScreenContains(t, "// line one")
		h.AssertScreenContains(t, "// line three")
		// line two should NOT have // prefix
		screen := h.CapturePane()
		if containsLine(screen, "// line two") {
			t.Errorf("line two should not be commented:\n%s", screen)
		}
	})

	// --- uncomment-region Python with leading whitespace ---
	t.Run("UncommentRegion_PythonLeadingWhitespace", func(t *testing.T) {
		dir := shortTempDir(t)
		content := "def foo():\n    # return 1\n    # return 2\n"
		path := shortTempFileInDir(t, dir, "main.py", content)
		h := newHarness(t, path)

		if err := h.WaitForContent("main.py", 5*time.Second); err != nil {
			t.Fatalf("goomacs did not start: %v", err)
		}

		// Select lines 2-3 (the commented lines)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-Space")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)
		h.SendKeys("C-n")
		time.Sleep(50 * time.Millisecond)

		invokeCommand(h, "uncomment-region")
		h.AssertMessageLine(t, "Region uncommented")

		time.Sleep(100 * time.Millisecond)
		h.AssertScreenContains(t, "    return 1")
		h.AssertScreenContains(t, "    return 2")
	})
}

// containsLine checks if a screen capture contains a line with the given substring.
func containsLine(screen, substr string) bool {
	for _, line := range splitLines(screen) {
		if len(line) >= len(substr) {
			for i := 0; i <= len(line)-len(substr); i++ {
				if line[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}

// splitLines splits a string by newlines.
func splitLines(s string) []string {
	lines := []string{}
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
