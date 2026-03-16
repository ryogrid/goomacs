# E2E No-Op and Boundary Condition Test Cases

Structured test case definitions for no-op keys and boundary conditions across all editor modes.

---

## No-Op Test Cases

No-op test cases verify that pressing keys with no binding in a given mode does **not** crash, corrupt state, or change the buffer. After every no-op input, the buffer content and cursor position must remain unchanged.

### Normal Mode No-Ops

Keys that are not bound in normal mode. Control keys not listed in the switch statement fall through silently (no action, no error).

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-NOOP-001 | Normal mode, buffer with text "hello", cursor at (0,2) | C-q | Buffer unchanged, cursor at (0,2), no crash |
| TC-NOOP-002 | Normal mode, buffer with text "hello", cursor at (0,2) | C-o | Buffer unchanged, cursor at (0,2), no crash |
| TC-NOOP-003 | Normal mode, buffer with text "hello", cursor at (0,2) | C-t | Buffer unchanged, cursor at (0,2), no crash |
| TC-NOOP-004 | Normal mode, buffer with text "hello", cursor at (0,0) | M-z (Alt+z, unbound Alt key) | Buffer unchanged, cursor at (0,0), no crash |
| TC-NOOP-005 | Normal mode, buffer with text "hello", cursor at (0,0) | M-a (Alt+a, unbound Alt key) | Buffer unchanged, cursor at (0,0), no crash |
| TC-NOOP-006 | Normal mode, buffer with text "hello", cursor at (0,0) | M-b (Alt+b, unbound Alt key) | Buffer unchanged, cursor at (0,0), no crash |
| TC-NOOP-007 | Normal mode, empty buffer, cursor at (0,0) | C-q | Buffer unchanged (1 empty line), cursor at (0,0), no crash |

### Search Mode No-Ops

In search mode, only these keys are handled: C-s, C-r, C-g, Enter, Backspace, and printable runes. All other keys **exit search** and are re-posted for normal processing.

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-NOOP-008 | Search forward mode (C-s pressed), query "hel" | C-x | Search exits, C-x prefix activated (message shows "C-x-") |
| TC-NOOP-009 | Search forward mode (C-s pressed), query "hel" | C-v | Search exits, page scrolls down |
| TC-NOOP-010 | Search forward mode (C-s pressed), query empty | C-n | Search exits, cursor moves down one line |
| TC-NOOP-011 | Search backward mode (C-r pressed), query "hel" | C-a | Search exits, cursor moves to beginning of line |
| TC-NOOP-012 | Search forward mode (C-s pressed), query empty | C-g | Search cancelled, cursor returns to original position, message "Quit" |

### Minibuffer Mode No-Ops

In minibuffer mode, handled keys: Enter, C-g, Esc, Backspace, Left/C-b, Right/C-f, C-a, C-e, C-d, C-k, Tab, printable runes. All other keys are ignored (minibuffer input preserved).

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-NOOP-013 | Minibuffer active (C-x C-f), input "test" | C-n | Minibuffer stays active, input "test" unchanged |
| TC-NOOP-014 | Minibuffer active (C-x C-f), input "test" | C-p | Minibuffer stays active, input "test" unchanged |
| TC-NOOP-015 | Minibuffer active (C-x C-f), input "test" | C-v | Minibuffer stays active, input "test" unchanged |
| TC-NOOP-016 | Minibuffer active (M-x), input "find" | C-w | Minibuffer stays active, input "find" unchanged |
| TC-NOOP-017 | Minibuffer active (C-x C-f), input "test" | C-y | Minibuffer stays active, input "test" unchanged |
| TC-NOOP-018 | Minibuffer active (C-x C-f), input "" (empty) | C-n | Minibuffer stays active, empty input unchanged |

### Confirm Mode No-Ops

In confirm mode (y/n prompt), only 'y', 'n', and C-g are handled. All other keys are silently ignored — the prompt remains.

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-NOOP-019 | Confirm mode ("Buffer modified; kill anyway?") | 'a' | Prompt remains, still in confirm mode |
| TC-NOOP-020 | Confirm mode ("Buffer modified; kill anyway?") | 'x' | Prompt remains, still in confirm mode |
| TC-NOOP-021 | Confirm mode ("Buffer modified; kill anyway?") | Enter | Prompt remains, still in confirm mode |
| TC-NOOP-022 | Confirm mode ("Buffer modified; kill anyway?") | C-f | Prompt remains, still in confirm mode |
| TC-NOOP-023 | Confirm mode ("Buffer modified; kill anyway?") | C-a | Prompt remains, still in confirm mode |
| TC-NOOP-024 | Confirm mode ("Buffer modified; kill anyway?") | Backspace | Prompt remains, still in confirm mode |
| TC-NOOP-025 | Confirm mode ("Buffer modified; kill anyway?") | Space | Prompt remains, still in confirm mode |

### C-x Prefix Mode No-Ops

After pressing C-x, only specific second keys are recognized: C-s, C-c, C-b, C-f, and runes b/2/3/o/0/1/k. Any other key cancels the prefix with no action.

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-NOOP-026 | Normal mode, C-x pressed (prefix active) | 'a' | Prefix cancelled, buffer unchanged, no crash |
| TC-NOOP-027 | Normal mode, C-x pressed (prefix active) | 'z' | Prefix cancelled, buffer unchanged, no crash |
| TC-NOOP-028 | Normal mode, C-x pressed (prefix active) | C-q | Prefix cancelled, buffer unchanged, no crash |
| TC-NOOP-029 | Normal mode, C-x pressed (prefix active) | C-n | Prefix cancelled, buffer unchanged (note: C-n not handled in prefix, falls through to rune handler — actually it's a control key that isn't matched, so the `switch ev.Key()` in the prefix handler doesn't match it, prefix is cleared, and the switch in the outer handler processes it as C-n) |
| TC-NOOP-030 | Normal mode, C-x pressed (prefix active) | 'x' | Prefix cancelled, buffer unchanged, no crash |
| TC-NOOP-031 | Normal mode, C-x pressed (prefix active) | C-d | Prefix cancelled; C-d falls through to normal dispatch (deletes char at cursor) |

### Grep Mode No-Ops

In grep mode (buffer ReadOnly=true, Mode="grep"), the grepModeHandler handles: Enter, n, p, g, q, M-n, M-p. Other rune keys fall through to ReadOnly check which blocks editing. Movement control keys (C-f, C-b, C-n, C-p) work normally.

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-NOOP-032 | Grep buffer active, cursor on result line | 'x' | "Buffer is read-only" message, buffer unchanged |
| TC-NOOP-033 | Grep buffer active, cursor on result line | 'a' | "Buffer is read-only" message, buffer unchanged |
| TC-NOOP-034 | Grep buffer active, cursor on result line | 'z' | "Buffer is read-only" message, buffer unchanged |
| TC-NOOP-035 | Grep buffer active, cursor on result line | Backspace | "Buffer is read-only" message, buffer unchanged |
| TC-NOOP-036 | Grep buffer active, cursor on result line | C-k | "Buffer is read-only" message, buffer unchanged |
| TC-NOOP-037 | Grep buffer active, cursor on result line | C-d | "Buffer is read-only" message, buffer unchanged |
| TC-NOOP-038 | Grep buffer active, cursor on non-result line (header/empty) | Enter | "No grep result on this line" message, cursor stays |

### Read-Only Buffer No-Ops (Non-Grep)

The *Buffer List* buffer is not explicitly ReadOnly in the code (it's refreshed each time), but grep buffers are. This section covers editing attempts on any ReadOnly buffer.

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-NOOP-039 | Read-only buffer active | InsertChar 'a' | "Buffer is read-only" message, buffer unchanged |
| TC-NOOP-040 | Read-only buffer active | C-y (yank) | "Buffer is read-only" message, buffer unchanged |
| TC-NOOP-041 | Read-only buffer active | Enter | "Buffer is read-only" message, buffer unchanged |

---

## Boundary Condition Test Cases

Boundary condition test cases verify correct behavior at the extremes of data structures and editor state.

### Empty Buffer Boundaries

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-BOUND-001 | Empty buffer (1 empty line), cursor at (0,0) | Backspace | No change — already at start of only line |
| TC-BOUND-002 | Empty buffer (1 empty line), cursor at (0,0) | C-d (delete forward) | No change — nothing to delete |
| TC-BOUND-003 | Empty buffer (1 empty line), cursor at (0,0) | C-k (kill line) | No change — empty line, nothing to kill (lastKill flag set) |
| TC-BOUND-004 | Empty buffer (1 empty line), cursor at (0,0) | C-f (forward) | No change — already at end of buffer |
| TC-BOUND-005 | Empty buffer (1 empty line), cursor at (0,0) | C-b (backward) | No change — already at start of buffer |
| TC-BOUND-006 | Empty buffer (1 empty line), cursor at (0,0) | C-n (down) | No change — only one line |
| TC-BOUND-007 | Empty buffer (1 empty line), cursor at (0,0) | C-p (up) | No change — already at first line |
| TC-BOUND-008 | Empty buffer (1 empty line), cursor at (0,0) | C-a (beginning of line) | No change — already at column 0 |
| TC-BOUND-009 | Empty buffer (1 empty line), cursor at (0,0) | C-e (end of line) | No change — empty line, end == column 0 |
| TC-BOUND-010 | Empty buffer (1 empty line), cursor at (0,0) | M-< (beginning of buffer) | No change — already at (0,0) |
| TC-BOUND-011 | Empty buffer (1 empty line), cursor at (0,0) | M-> (end of buffer) | No change — already at end of only line |
| TC-BOUND-012 | Empty buffer (1 empty line), cursor at (0,0) | C-v (page down) | Scroll offset stays 0, cursor stays (0,0) |
| TC-BOUND-013 | Empty buffer (1 empty line), cursor at (0,0) | M-v (page up) | Scroll offset stays 0, cursor stays (0,0) |
| TC-BOUND-014 | Empty buffer (1 empty line), cursor at (0,0) | Enter (newline) | Buffer now has 2 lines (both empty), cursor at (1,0) |

### Single Character Buffer Boundaries

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-BOUND-015 | Buffer with "a" (1 line, 1 char), cursor at (0,0) | C-d | Buffer becomes empty line, cursor at (0,0) |
| TC-BOUND-016 | Buffer with "a" (1 line, 1 char), cursor at (0,1) | Backspace | Buffer becomes empty line, cursor at (0,0) |
| TC-BOUND-017 | Buffer with "a" (1 line, 1 char), cursor at (0,0) | C-d, then C-_ (undo) | Buffer restored to "a", cursor at (0,0) |
| TC-BOUND-018 | Buffer with "a" (1 line, 1 char), cursor at (0,1) | C-f | No change — at end of only line |
| TC-BOUND-019 | Buffer with "a" (1 line, 1 char), cursor at (0,0) | C-b | No change — at start of buffer |

### Movement at Buffer Boundaries

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-BOUND-020 | 5-line buffer, cursor at (0,0) (first line, first col) | C-p | Cursor stays at (0,0) |
| TC-BOUND-021 | 5-line buffer, cursor at (4, end) (last line, last col) | C-n | Cursor stays at (4, end) |
| TC-BOUND-022 | 5-line buffer, cursor at (0,0) | C-b | Cursor stays at (0,0) |
| TC-BOUND-023 | 5-line buffer, cursor at (4, end) | C-f | Cursor stays at (4, end) |
| TC-BOUND-024 | Line 0: "short", Line 1: "this is longer", cursor at (1,12) | C-p | Cursor clamped to (0, 5) — length of "short" |
| TC-BOUND-025 | Line 0: "this is longer", Line 1: "short", cursor at (0,12) | C-n | Cursor clamped to (1, 5) — length of "short" |
| TC-BOUND-026 | 5-line buffer, cursor at end of line 2 | C-f | Cursor wraps to (3, 0) — next line, column 0 |
| TC-BOUND-027 | 5-line buffer, cursor at (3, 0) | C-b | Cursor wraps to (2, end) — previous line, last column |

### Page Up/Down Boundary Conditions

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-BOUND-028 | 30-line buffer, 24-row terminal (viewHeight=22), cursor at line 20 | C-v (page down) | Scroll offset advances, cursor moves to stay in view; less than full page remaining |
| TC-BOUND-029 | 30-line buffer, 24-row terminal, scroll at top | C-v, C-v | Second C-v reaches near end of buffer; scroll offset capped at max |
| TC-BOUND-030 | 30-line buffer, scroll at line 5 | M-v (page up) | Scroll offset decreases by viewHeight |
| TC-BOUND-031 | 30-line buffer, scroll at line 2 | M-v (page up) | Scroll offset clamped to 0 (less than full page to scroll) |

### Undo Stack Boundaries

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-BOUND-032 | Empty buffer, no edits made | C-_ (undo) | "No further undo information" message, no crash |
| TC-BOUND-033 | Buffer with 1 edit on undo stack | C-_ | Edit undone, undo stack now empty |
| TC-BOUND-034 | Buffer with 1 edit on undo stack | C-_, C-_ | First undo succeeds, second shows "No further undo information" |
| TC-BOUND-035 | Perform 101 character insertions (each with SaveUndo) | C-_ x100 | Undoes 100 operations (maxUndoEntries); 101st C-_ shows "No further undo information" |
| TC-BOUND-036 | Perform 10 edits, undo all 10 | Check that undone buffer matches original state exactly |

### Consecutive Kill Accumulation

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-BOUND-037 | 3-line buffer: "aaa\nbbb\nccc", cursor at (0,0) | C-k, C-k (kill line content, kill newline) | Kill ring has one entry: "aaa\n" (consecutive kills accumulate) |
| TC-BOUND-038 | 3-line buffer: "aaa\nbbb\nccc", cursor at (0,0) | C-k, C-k, C-k, C-k | Kill ring entry accumulates: "aaa\nbbb\n" (killed both lines + newlines) |
| TC-BOUND-039 | 3-line buffer: "aaa\nbbb\nccc", cursor at (0,0) | C-k, C-f, C-k | Two separate kill ring entries — C-f breaks consecutive chain |
| TC-BOUND-040 | 50-line buffer, cursor at (0,0) | C-k x100 (alternating content + newline kills) | All kills accumulated in one entry, buffer ends as single empty line, no crash |

### Kill Ring Empty

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-BOUND-041 | Fresh buffer, kill ring empty | C-y (yank) | No change — nothing to yank |
| TC-BOUND-042 | Fresh buffer, kill ring empty | M-w (copy region with mark active at same pos) | Zero-width region, nothing copied, kill ring still empty |

### Region Boundaries

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-BOUND-043 | Buffer "hello", cursor at (0,2) | C-Space (set mark), C-w (kill region) | Zero-width region (mark == point), mark deactivated, nothing killed |
| TC-BOUND-044 | Buffer "hello", no mark set | C-w (kill region) | No action — MarkActive is false |
| TC-BOUND-045 | Buffer "hello", no mark set | M-w (copy region) | No action — MarkActive is false |

### Window Management Boundaries

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-BOUND-046 | Single window | C-x 0 (close window) | No action — cannot close only window |
| TC-BOUND-047 | Single window | C-x 1 (close other windows) | No action — no other windows to close |
| TC-BOUND-048 | Single window | C-x o (other window) | No action — only one window |
| TC-BOUND-049 | Two vertical windows | C-x 3 (horizontal split) | Message: "Cannot split horizontally while in vertical split mode" |
| TC-BOUND-050 | Two horizontal windows | C-x 2 (vertical split) | Message: "Cannot split vertically while in horizontal split mode" |

### Search Boundaries

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-BOUND-051 | Buffer "hello world", cursor at (0,0) | C-s (enter search), Enter immediately (empty query) | Search exits, cursor stays at (0,0) |
| TC-BOUND-052 | Buffer "hello world", cursor at (0,0) | C-s, type "xyz" (no match) | Message: "Failing I-search: xyz", cursor unchanged from search start |
| TC-BOUND-053 | Buffer with "target" on line 3 of 10, cursor at line 8 | C-s, type "target" | Wrap-around: finds "target" on line 3 |
| TC-BOUND-054 | Buffer with "target" on line 8 of 10, cursor at line 3 | C-r, type "target" | Backward wrap-around: finds "target" on line 8 |
| TC-BOUND-055 | Buffer "aaa", cursor at (0,1) | C-s, type "a", C-s, C-s | Forward search repeats: finds at col 1, col 2, wraps to col 0 |

### Backspace in Search

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-BOUND-056 | Buffer "hello", search forward, query "hel" (match at 0,0) | Backspace | Query becomes "he", re-searches from original position, still matches |
| TC-BOUND-057 | Buffer "hello", search forward, query "h" | Backspace | Query becomes empty, cursor returns to original position |
| TC-BOUND-058 | Buffer "hello", search forward, query empty | Backspace | No change — nothing to delete from empty query |

### Large Content Boundaries

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-BOUND-059 | Buffer with 1000 lines | M-> (end of buffer) | Cursor at (999, end of last line), no performance issue |
| TC-BOUND-060 | Buffer with 1000 lines, cursor at end | M-< (beginning of buffer) | Cursor at (0,0) |
| TC-BOUND-061 | Buffer with line of 500 characters | C-e | Cursor at column 500 |
| TC-BOUND-062 | Buffer with line of 500 characters | C-a then C-e then C-a | Cursor at (row, 0) after round-trip |

### Minimum Window Size

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-BOUND-063 | Small terminal (40x6), two vertical windows | Verify both windows render | Each window gets ~2 rows of text + status line |
| TC-BOUND-064 | Terminal resized to 40x4 with two vertical windows | Verify no crash | Windows should still render with minimal space |

### File Operation Boundaries

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-BOUND-065 | Buffer with Filename="" (no file) | C-x C-s | Message: "No file name" |
| TC-BOUND-066 | *scratch* buffer, unmodified | C-x k | Buffer killed without confirmation, new *scratch* created |
| TC-BOUND-067 | Only one buffer open, kill it | C-x k, Enter (default) | Buffer killed, new *scratch* buffer auto-created |

### Tab Completion Boundaries

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-BOUND-068 | Minibuffer (C-x C-f), input "nonexistent_dir/foo" | Tab | No completion — directory doesn't exist, input unchanged |
| TC-BOUND-069 | Minibuffer (C-x C-f), input "" (empty) | Tab | Completes based on current directory entries |
| TC-BOUND-070 | Minibuffer (M-x), input "zzz" (no matching commands) | Tab | No completion — no commands start with "zzz", input unchanged |
| TC-BOUND-071 | Minibuffer (M-x), input "find-grep" (exact single match) | Tab | Completes to "find-grep" |
| TC-BOUND-072 | Minibuffer (M-x), input "com" | Tab | Shows both "comment-region" and... no, "uncomment-region" doesn't start with "com". Only "comment-region" matches → completes fully |

### Quit Boundary

| ID | Precondition | Input Sequence | Expected Post-Condition |
|----|-------------|----------------|------------------------|
| TC-BOUND-073 | Modified buffer | C-x C-c | Warning message: "Modified buffers exist; exit anyway? (C-x C-c to confirm)" |
| TC-BOUND-074 | Modified buffer, warning shown | C-x C-c again | Editor exits |
| TC-BOUND-075 | Modified buffer, warning shown | Any other key (e.g., C-f) | Warning cleared, editor continues normally |

---

## Summary Statistics

| Category | Test Case Count |
|----------|----------------|
| No-op: Normal mode | 7 (TC-NOOP-001 to TC-NOOP-007) |
| No-op: Search mode | 5 (TC-NOOP-008 to TC-NOOP-012) |
| No-op: Minibuffer mode | 6 (TC-NOOP-013 to TC-NOOP-018) |
| No-op: Confirm mode | 7 (TC-NOOP-019 to TC-NOOP-025) |
| No-op: C-x Prefix mode | 6 (TC-NOOP-026 to TC-NOOP-031) |
| No-op: Grep mode | 7 (TC-NOOP-032 to TC-NOOP-038) |
| No-op: Read-only buffer | 3 (TC-NOOP-039 to TC-NOOP-041) |
| **Total No-Op** | **41** |
| Boundary: Empty buffer | 14 (TC-BOUND-001 to TC-BOUND-014) |
| Boundary: Single char buffer | 5 (TC-BOUND-015 to TC-BOUND-019) |
| Boundary: Movement | 8 (TC-BOUND-020 to TC-BOUND-027) |
| Boundary: Page up/down | 4 (TC-BOUND-028 to TC-BOUND-031) |
| Boundary: Undo stack | 5 (TC-BOUND-032 to TC-BOUND-036) |
| Boundary: Kill accumulation | 4 (TC-BOUND-037 to TC-BOUND-040) |
| Boundary: Kill ring empty | 2 (TC-BOUND-041 to TC-BOUND-042) |
| Boundary: Region | 3 (TC-BOUND-043 to TC-BOUND-045) |
| Boundary: Window management | 5 (TC-BOUND-046 to TC-BOUND-050) |
| Boundary: Search | 8 (TC-BOUND-051 to TC-BOUND-058) |
| Boundary: Large content | 4 (TC-BOUND-059 to TC-BOUND-062) |
| Boundary: Min window size | 2 (TC-BOUND-063 to TC-BOUND-064) |
| Boundary: File operations | 3 (TC-BOUND-065 to TC-BOUND-067) |
| Boundary: Tab completion | 5 (TC-BOUND-068 to TC-BOUND-072) |
| Boundary: Quit | 3 (TC-BOUND-073 to TC-BOUND-075) |
| **Total Boundary** | **75** |
| **Grand Total** | **116** |

## Notes

- **No-op behavior by mode:**
  - **Normal mode:** Unbound control keys (C-q, C-o, C-t) simply do nothing — no message, no error.
  - **Search mode:** Unrecognized keys **exit search** and are re-posted for normal processing. This is intentional Emacs behavior.
  - **Minibuffer mode:** Unrecognized keys are **silently ignored** — the minibuffer stays active with input preserved.
  - **Confirm mode:** Only 'y', 'n', and C-g are accepted. All other keys are **silently ignored**.
  - **C-x prefix:** After C-x, unrecognized second keys cancel the prefix. The key is then processed by the normal switch — control keys may trigger their normal action, while unrecognized runes (like 'a', 'z') fall through to InsertChar.
  - **Grep mode:** Unbound rune keys fall through to ReadOnly check → "Buffer is read-only". Navigation keys (C-f, C-b, C-n, C-p) work normally.

- **Key distinction: C-x prefix fallthrough.** When C-x prefix is cancelled by an unmatched key, the key handling continues to the normal switch statement. This means C-x C-n would cancel prefix AND move cursor down. C-x a would cancel prefix AND insert 'a'. This is the actual codebase behavior — the `prefixCx = false` is set, but the key event still falls through.

- **maxUndoEntries = 100** in buffer.go. Testing undo limit requires inserting 101 operations and verifying only the last 100 can be undone.
