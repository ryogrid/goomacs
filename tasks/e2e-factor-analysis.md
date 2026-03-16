# E2E Factor Analysis Table

Systematic enumeration of all testable factors for the goomacs editor, organized by category with priority assignments.

## Priority Legend

| Priority | Meaning | Rationale |
|----------|---------|-----------|
| **Critical** | Must test | Core functionality; failures cause crashes, data loss, or wrong behavior |
| **Important** | Should test | Significant user-facing behavior; failures cause confusion or minor wrong behavior |
| **Edge-case** | Nice to have | Unusual combinations; failures are cosmetic or unlikely to occur in practice |

---

## Factor 1: Editor Mode

**Category:** Mode State
**Priority:** Critical

The editor mode determines which key bindings are active and how input is dispatched.

| Level | Description | Source |
|-------|-------------|--------|
| Normal | Default editing mode; all standard bindings active | main.go event loop default path |
| Search Forward | Incremental forward search (C-s); search-specific bindings | `searchMode=true, searchForward=true` |
| Search Backward | Incremental backward search (C-r); search-specific bindings | `searchMode=true, searchForward=false` |
| Minibuffer | Text input mode for prompts (find-file, switch-buffer, M-x, goto-line) | `minibufferMode=true` |
| Confirm | Yes/no confirmation dialog (e.g., kill modified buffer, quit with unsaved) | `confirmMode=true` |
| C-x Prefix | Waiting for second key after C-x press | `prefixCx=true` |
| Grep | Special buffer mode for grep results; grep-specific bindings | `buf.Mode="grep"` via modeHandlers |

**Rationale:** Mode determines the entire input dispatch path. Every key binding test must specify which mode is active. Critical because wrong mode dispatch causes all keys to behave incorrectly.

---

## Factor 2: Cursor Row Position

**Category:** Cursor State
**Priority:** Critical

| Level | Description | Boundary Condition |
|-------|-------------|-------------------|
| First row (0) | Cursor at top of buffer | C-p, C-b at start-of-line, M-v at top |
| Middle row | Cursor in middle of buffer | Normal operation baseline |
| Last row | Cursor at bottom of buffer | C-n, C-f at end-of-line, C-v at bottom |
| Beyond viewport | Cursor scrolled off-screen | Scroll offset interaction with movement |

**Rationale:** Row boundaries affect movement wrapping, scroll behavior, and page up/down edge cases.

---

## Factor 3: Cursor Column Position

**Category:** Cursor State
**Priority:** Critical

| Level | Description | Boundary Condition |
|-------|-------------|-------------------|
| Column 0 (start of line) | Cursor at beginning of line | C-b wraps to previous line, Backspace joins lines |
| Middle of line | Cursor in middle of text | Normal operation baseline |
| End of line | Cursor at last character | C-f wraps to next line, C-d joins lines |
| Past short line (clamped) | After C-n/C-p from longer line to shorter line | Column clamping behavior |

**Rationale:** Column boundaries affect character-level movement, wrapping between lines, and editing operations like backspace/delete at edges.

---

## Factor 4: Buffer Content Type

**Category:** Buffer State
**Priority:** Critical

| Level | Description | Key Characteristics |
|-------|-------------|-------------------|
| Empty buffer | Zero lines or single empty line | All movement/editing at boundary |
| Single line | One line of text | No vertical movement possible |
| Multi-line (small) | 2-30 lines | Standard editing scenario |
| Multi-line (large, 1000+ lines) | Very long file | Scroll behavior, page up/down, performance |
| Wide lines (500+ chars) | Lines exceeding terminal width | Horizontal scroll, rendering |
| Mixed content (tabs, UTF-8) | Contains tabs, multi-byte characters | Tab expansion, character width handling |

**Rationale:** Buffer content determines which operations are valid and which hit boundary conditions. Empty and single-line buffers are critical edge cases.

---

## Factor 5: Modification State

**Category:** Buffer State
**Priority:** Important

| Level | Description | Effect |
|-------|-------------|--------|
| Unmodified | Buffer matches on-disk file (or freshly created) | Clean status bar, no save prompt on close |
| Modified | Buffer has unsaved changes | `[Modified]` in status bar, confirmation on C-x k / C-x C-c |

**Rationale:** Modification state affects save/close workflows and status bar display. Important for kill-buffer and quit confirmation flows.

---

## Factor 6: Mark State

**Category:** Region/Selection
**Priority:** Critical

| Level | Description | Effect |
|-------|-------------|--------|
| No mark set | `MarkActive=false`, no prior mark | C-w/M-w have no region to operate on |
| Mark active, same position | `MarkActive=true`, mark == cursor | Zero-width region for C-w/M-w |
| Mark active, different position | `MarkActive=true`, mark != cursor | Region defined for kill/copy/comment |

**Rationale:** Mark state is critical for region operations (C-w, M-w, comment-region, uncomment-region). Zero-width region is a key edge case.

---

## Factor 7: Kill Ring State

**Category:** Buffer State
**Priority:** Important

| Level | Description | Effect |
|-------|-------------|--------|
| Empty kill ring | No prior kills/copies | C-y does nothing |
| Non-empty kill ring | At least one entry | C-y pastes last entry |

**Rationale:** Kill ring state affects yank behavior and consecutive kill accumulation.

---

## Factor 8: Undo Stack State

**Category:** Buffer State
**Priority:** Important

| Level | Description | Effect |
|-------|-------------|--------|
| Empty undo stack | No prior operations | C-_ does nothing |
| Partially filled (1-99) | Some operations recorded | C-_ undoes last operation |
| Full (100 entries) | At maxUndoEntries limit | Oldest entries dropped on new operations |

**Rationale:** Undo stack has a hard limit of 100 entries (`maxUndoEntries` in buffer.go). Testing the limit and empty stack are important edge cases.

---

## Factor 9: Read-Only Flag

**Category:** Buffer State
**Priority:** Important

| Level | Description | Effect |
|-------|-------------|--------|
| Read-write | Normal editing buffer | All editing operations allowed |
| Read-only | `ReadOnly=true` (e.g., *Buffer List*, *grep*) | InsertChar, Backspace, C-d, C-k, C-y, Enter blocked with "Buffer is read-only" message |

**Rationale:** Read-only state must correctly block all editing operations while allowing navigation and M-x commands.

---

## Factor 10: Key Input Class

**Category:** Input
**Priority:** Critical

| Level | Description | Examples |
|-------|-------------|---------|
| Printable character (ASCII) | Standard letters, digits, symbols | a-z, 0-9, punctuation |
| Printable character (UTF-8 multi-byte) | Non-ASCII characters | Japanese characters, emoji |
| Control key (bound) | Ctrl combinations with defined bindings | C-f, C-b, C-n, C-p, C-a, C-e, C-k, C-d, C-w, C-y, C-_ |
| Control key (unbound) | Ctrl combinations without bindings | C-q, C-o, C-t |
| Meta/Alt key | Alt combinations | M-v, M-<, M->, M-w, M-x, M-n, M-p |
| Special key | Named keys | Enter, Backspace, Tab, Arrow keys, Esc |
| No-op / irrelevant key | Keys that should have no effect in current mode | Unbound keys in grep mode, random keys in confirm mode |

**Rationale:** Different key classes exercise different input dispatch paths. Testing unbound and no-op keys is critical for robustness.

---

## Factor 11: File Backing State

**Category:** File System
**Priority:** Important

| Level | Description | Example |
|-------|-------------|---------|
| File-backed (existing) | Buffer loaded from existing file | `NewBufferFromFile("test.go")` |
| File-backed (new/nonexistent) | Filename set but no file on disk | C-x C-f with new path |
| Special buffer (no file) | Transient buffer with special name | `*scratch*`, `*Buffer List*`, `*grep*` |

**Rationale:** File backing affects save operations, status bar display, and highlighter initialization.

---

## Factor 12: Highlighter State

**Category:** Rendering
**Priority:** Edge-case

| Level | Description | Condition |
|-------|-------------|-----------|
| Active (lexer matched) | Syntax highlighting enabled | File extension matches known language (e.g., .go, .py) |
| Inactive (no lexer) | No syntax highlighting | Unknown extension or special buffer |

**Rationale:** Highlighter affects rendering but not editing behavior. Edge-case priority because highlighting failures are cosmetic.

---

## Factor 13: Window Count

**Category:** Window Configuration
**Priority:** Important

| Level | Description | State |
|-------|-------------|-------|
| Single window | One window occupies full screen | `len(windows)==1` |
| Two windows (vertical split) | Top/bottom layout via C-x 2 | `splitMode="vertical"`, `len(windows)==2` |
| Two windows (horizontal split) | Side-by-side layout via C-x 3 | `splitMode="horizontal"`, `len(windows)==2` |

**Rationale:** Window configuration affects layout calculation, status bar rendering, and window management commands (C-x 0, C-x 1, C-x o).

---

## Factor 14: Active Window

**Category:** Window Configuration
**Priority:** Important

| Level | Description | State |
|-------|-------------|-------|
| First window (index 0) | Top or left window active | `activeWindowIdx=0` |
| Second window (index 1) | Bottom or right window active | `activeWindowIdx=1` (only when split) |

**Rationale:** Active window determines which buffer receives input and which status bar shows active styling (reverse video vs dashes).

---

## Factor 15: Terminal Size

**Category:** Terminal Configuration
**Priority:** Important

| Level | Description | Dimensions |
|-------|-------------|-----------|
| Standard (80x24) | Default terminal size | Normal operation baseline |
| Small (40x12) | Reduced terminal | Tests minimum viable rendering, window constraints |
| Large (200x50) | Expanded terminal | Tests rendering at scale, no artifacts |

**Rationale:** Terminal size affects viewport calculations, page up/down distance, scroll behavior, and window split constraints. Resize events (SIGWINCH) must be handled correctly.

---

## Factor 16: Consecutive Kill State

**Category:** Buffer State
**Priority:** Important

| Level | Description | Effect |
|-------|-------------|--------|
| First kill (lastKill=false) | No preceding kill operation | C-k creates new kill ring entry |
| Consecutive kill (lastKill=true) | Previous operation was also a kill | C-k appends to existing kill ring entry |

**Rationale:** Consecutive kill accumulation is a core Emacs behavior. The `lastKill` flag determines whether kills append or create new entries.

---

## Factor 17: Search State

**Category:** Search
**Priority:** Important

| Level | Description | State |
|-------|-------------|-------|
| No active search | Normal mode | `searchMode=false` |
| Search with match | Query has matching text | `searchHasMatch=true` |
| Search failing | Query has no match | `searchHasMatch=false`, "Failing I-search" message |
| Search with empty query | C-s/C-r pressed but no chars typed | `searchQuery` is empty |

**Rationale:** Search state affects message line display, cursor positioning, and key behavior (C-s/C-r cycle through matches).

---

## Factor 18: Buffer Mode (Local)

**Category:** Buffer Configuration
**Priority:** Important

| Level | Description | Bindings |
|-------|-------------|----------|
| Normal (empty string) | Standard editing mode | Default key dispatch |
| Grep mode | Grep results buffer | n/p navigation, Enter jumps to source, g refreshes, q quits |

**Rationale:** Buffer-local mode overrides key dispatch via `modeHandlers` map. Grep mode has its own complete set of bindings.

---

## Factor 19: Minibuffer Context

**Category:** Minibuffer
**Priority:** Important

| Level | Description | Tab Completion |
|-------|-------------|---------------|
| Find file (C-x C-f) | File path input | File system path completion |
| Switch buffer (C-x b) | Buffer name input | No special completion |
| Kill buffer (C-x k) | Buffer name input | No special completion |
| M-x command | Command name input | Command name completion via `FindCommandsByPrefix` |
| Goto line (M-g) | Line number input | No completion |

**Rationale:** Different minibuffer contexts have different tab completion behaviors and callbacks. Testing each context ensures the minibuffer subsystem works correctly across all entry points.

---

## Summary Statistics

| Category | Factor Count | Total Levels |
|----------|-------------|-------------|
| Mode State | 1 | 7 |
| Cursor State | 2 | 8 |
| Buffer State | 6 | 17 |
| Region/Selection | 1 | 3 |
| Input | 1 | 7 |
| File System | 1 | 3 |
| Rendering | 1 | 2 |
| Window Configuration | 2 | 5 |
| Terminal Configuration | 1 | 3 |
| Search | 1 | 4 |
| Buffer Configuration | 1 | 2 |
| Minibuffer | 1 | 5 |
| **Total** | **19 factors** | **66 levels** |

## Priority Distribution

| Priority | Factor Count | Factors |
|----------|-------------|---------|
| **Critical** | 5 | Editor Mode, Cursor Row Position, Cursor Column Position, Buffer Content Type, Mark State, Key Input Class |
| **Important** | 12 | Modification State, Kill Ring State, Undo Stack State, Read-Only Flag, File Backing State, Window Count, Active Window, Terminal Size, Consecutive Kill State, Search State, Buffer Mode, Minibuffer Context |
| **Edge-case** | 1 | Highlighter State |

## Notes

- **Full combinatorial space:** 19 factors with 3-7 levels each yields an astronomically large test space (~10^15 combinations). 3-way combinatorial coverage (t=3) reduces this to a manageable 200-400 test cases.
- **Infeasible combinations** must be pruned: e.g., Grep mode + Confirm mode cannot coexist; Read-only buffer + editing operations are tested separately as no-op cases; Minibuffer mode + window commands are invalid.
- **Factor interdependencies:** Some factors are only relevant in specific modes (e.g., Search State only matters when Editor Mode = Search; Minibuffer Context only matters when Editor Mode = Minibuffer).
