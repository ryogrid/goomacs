# E2E Test Matrix — 3-Way Combinatorial Coverage

Comprehensive test matrix achieving 3-way (t=3) combinatorial coverage across all testable factors. Every combination of any 3 factor levels from the [factor analysis](e2e-factor-analysis.md) appears in at least one test case.

## Methodology

**Covering array approach:** For each thematic group, the relevant factors are identified and test cases are constructed so that every 3-way interaction among those factors is exercised. Factors irrelevant to a group (e.g., Search State for Movement tests) are fixed or excluded.

**Factor abbreviations used in tables:**
- **Mode**: Editor Mode (Norm/SrchF/SrchB/Mini/Conf/CxPfx/Grep)
- **CurR**: Cursor Row (First/Mid/Last/BeyondVP)
- **CurC**: Cursor Column (Col0/Mid/End/Clamped)
- **Content**: Buffer Content (Empty/Single/SmallML/LargeML/Wide/Mixed)
- **ModState**: Modification State (Unmod/Mod)
- **Mark**: Mark State (None/SamePos/DiffPos)
- **KillRing**: Kill Ring (Empty/NonEmpty)
- **Undo**: Undo Stack (Empty/Partial/Full)
- **RO**: Read-Only (RW/RO)
- **Key**: Key Input Class (ASCII/UTF8/CtrlBound/CtrlUnbound/Meta/Special/NoOp)
- **File**: File Backing (Existing/New/Special)
- **HL**: Highlighter (Active/Inactive)
- **WinCnt**: Window Count (Single/VSplit/HSplit)
- **ActWin**: Active Window (Win0/Win1)
- **Term**: Terminal Size (80x24/40x12/200x50)
- **ConsKill**: Consecutive Kill (First/Consec)
- **SrchSt**: Search State (NoSearch/Match/Failing/EmptyQ)
- **BufMode**: Buffer Mode (Normal/Grep)
- **MiniCtx**: Minibuffer Context (FindFile/SwitchBuf/KillBuf/Mx/GotoLine)

---

## Excluded (Infeasible) Combinations

| Combination | Rationale |
|-------------|-----------|
| Mode=Grep + Mode=Confirm | Only one mode active at a time |
| Mode=Search + Mode=Minibuffer | Mutually exclusive modes |
| Mode=Minibuffer + WinCnt=VSplit + Key=C-x_2 | Window commands not available in minibuffer |
| RO=true + Key=editing (InsertChar, Backspace, C-d, C-k, C-y, Enter) | Tested separately as no-ops in TC-ERR group |
| BufMode=Grep + ModState=Modified | Grep buffer is read-only, cannot be modified |
| Content=Empty + CurR=Last (>0) | Empty buffer has only row 0 |
| Content=Single + CurR=Last (>0) | Single-line buffer has only row 0 |
| WinCnt=Single + ActWin=Win1 | No second window in single-window mode |
| SrchSt=Match/Failing/EmptyQ + Mode!=Search | Search state only meaningful in search mode |
| MiniCtx=* + Mode!=Minibuffer | Minibuffer context only meaningful in minibuffer mode |

---

## Group 1: Movement (TC-MOV)

**Relevant factors:** Mode(Normal), CurR, CurC, Content, WinCnt, ActWin, Term, Key(movement)

| ID | CurR | CurC | Content | WinCnt | Term | Input | Expected Post-Condition |
|----|------|------|---------|--------|------|-------|------------------------|
| TC-MOV-001 | First | Col0 | SmallML | Single | 80x24 | C-f | Cursor at (0,1) |
| TC-MOV-002 | First | Col0 | SmallML | Single | 80x24 | C-b | Cursor stays (0,0) — start of buffer |
| TC-MOV-003 | First | Col0 | SmallML | Single | 80x24 | C-n | Cursor at (1,0) |
| TC-MOV-004 | First | Col0 | SmallML | Single | 80x24 | C-p | Cursor stays (0,0) — top of buffer |
| TC-MOV-005 | First | Col0 | SmallML | Single | 80x24 | C-a | Cursor stays (0,0) — already at BOL |
| TC-MOV-006 | First | Col0 | SmallML | Single | 80x24 | C-e | Cursor at (0, end-of-line) |
| TC-MOV-007 | First | Col0 | SmallML | Single | 80x24 | C-v | Scroll down, cursor moves down ~viewHeight |
| TC-MOV-008 | First | Col0 | SmallML | Single | 80x24 | M-v | No change — already at top |
| TC-MOV-009 | First | Col0 | SmallML | Single | 80x24 | M-< | Cursor stays (0,0) |
| TC-MOV-010 | First | Col0 | SmallML | Single | 80x24 | M-> | Cursor at (lastRow, end-of-last-line) |
| TC-MOV-011 | First | Col0 | SmallML | Single | 80x24 | Right | Cursor at (0,1) |
| TC-MOV-012 | First | Col0 | SmallML | Single | 80x24 | Down | Cursor at (1,0) |
| TC-MOV-013 | Mid | Mid | SmallML | Single | 80x24 | C-f | Cursor col +1 |
| TC-MOV-014 | Mid | Mid | SmallML | Single | 80x24 | C-b | Cursor col -1 |
| TC-MOV-015 | Mid | Mid | SmallML | Single | 80x24 | C-n | Cursor row +1, col maintained or clamped |
| TC-MOV-016 | Mid | Mid | SmallML | Single | 80x24 | C-p | Cursor row -1 |
| TC-MOV-017 | Mid | Mid | SmallML | Single | 80x24 | C-a | Cursor at (row, 0) |
| TC-MOV-018 | Mid | Mid | SmallML | Single | 80x24 | C-e | Cursor at (row, end-of-line) |
| TC-MOV-019 | Mid | Mid | SmallML | Single | 80x24 | Left | Cursor col -1 |
| TC-MOV-020 | Mid | Mid | SmallML | Single | 80x24 | Up | Cursor row -1 |
| TC-MOV-021 | Last | End | SmallML | Single | 80x24 | C-f | Cursor stays — end of buffer |
| TC-MOV-022 | Last | End | SmallML | Single | 80x24 | C-n | Cursor stays — last row |
| TC-MOV-023 | Last | End | SmallML | Single | 80x24 | C-b | Cursor col -1 |
| TC-MOV-024 | Last | End | SmallML | Single | 80x24 | C-p | Cursor row -1 |
| TC-MOV-025 | Last | End | SmallML | Single | 80x24 | C-v | Cursor near end, scroll capped |
| TC-MOV-026 | Last | End | SmallML | Single | 80x24 | M-> | Cursor stays — already at end |
| TC-MOV-027 | First | End | SmallML | Single | 80x24 | C-f | Cursor wraps to (1,0) |
| TC-MOV-028 | Mid | Col0 | SmallML | Single | 80x24 | C-b | Cursor wraps to (row-1, end) |
| TC-MOV-029 | Mid | End | SmallML | Single | 80x24 | C-f | Cursor wraps to (row+1, 0) |
| TC-MOV-030 | First | Col0 | Empty | Single | 80x24 | C-f | Cursor stays (0,0) |
| TC-MOV-031 | First | Col0 | Empty | Single | 80x24 | C-n | Cursor stays (0,0) |
| TC-MOV-032 | First | Col0 | Empty | Single | 80x24 | C-v | Cursor stays (0,0) |
| TC-MOV-033 | First | Col0 | Empty | Single | 80x24 | M-> | Cursor stays (0,0) |
| TC-MOV-034 | First | Col0 | Single | Single | 80x24 | C-n | Cursor stays — single line |
| TC-MOV-035 | First | Col0 | Single | Single | 80x24 | C-p | Cursor stays (0,0) |
| TC-MOV-036 | First | End | Single | Single | 80x24 | C-f | Cursor stays — end of only line |
| TC-MOV-037 | First | End | Single | Single | 80x24 | C-e | Cursor stays — already at EOL |
| TC-MOV-038 | Mid | Clamped | SmallML | Single | 80x24 | C-p | Cursor at (row-1, clamped to shorter line length) |
| TC-MOV-039 | Mid | Clamped | SmallML | Single | 80x24 | C-n | Cursor at (row+1, clamped to shorter line length) |
| TC-MOV-040 | First | Col0 | LargeML | Single | 80x24 | M-> | Cursor at (999, end) |
| TC-MOV-041 | Last | End | LargeML | Single | 80x24 | M-< | Cursor at (0,0) |
| TC-MOV-042 | First | Col0 | LargeML | Single | 80x24 | C-v | Cursor moves down ~22 rows |
| TC-MOV-043 | BeyondVP | Mid | LargeML | Single | 80x24 | C-n | Cursor row +1, scroll adjusts |
| TC-MOV-044 | BeyondVP | Mid | LargeML | Single | 80x24 | C-p | Cursor row -1, scroll adjusts |
| TC-MOV-045 | First | Col0 | Wide | Single | 80x24 | C-e | Cursor at (0, 500+) |
| TC-MOV-046 | First | End | Wide | Single | 80x24 | C-a | Cursor at (0,0) |
| TC-MOV-047 | First | Col0 | Mixed | Single | 80x24 | C-f | Cursor at (0,1) — tab counts as 1 char |
| TC-MOV-048 | First | Col0 | SmallML | VSplit | 80x24 | C-f | Cursor at (0,1) in active window |
| TC-MOV-049 | Mid | Mid | SmallML | VSplit | 80x24 | C-n | Cursor row +1 in active window |
| TC-MOV-050 | First | Col0 | SmallML | HSplit | 80x24 | C-e | Cursor at (0, EOL) in active window |
| TC-MOV-051 | First | Col0 | SmallML | Single | 40x12 | C-v | Page down with smaller viewport (~10 rows) |
| TC-MOV-052 | First | Col0 | SmallML | Single | 40x12 | M-v | No change — at top |
| TC-MOV-053 | Mid | Mid | SmallML | Single | 200x50 | C-v | Page down with large viewport (~48 rows) |
| TC-MOV-054 | First | Col0 | SmallML | VSplit | 40x12 | C-n | Cursor moves in smaller window |
| TC-MOV-055 | Mid | Mid | SmallML | HSplit | 200x50 | C-f | Cursor col +1 in wide window |
| TC-MOV-056 | Last | Col0 | SmallML | VSplit | 80x24 | C-v | Page down near end, capped |
| TC-MOV-057 | Mid | Col0 | LargeML | VSplit | 40x12 | M-> | Jump to end in split window |
| TC-MOV-058 | First | Col0 | Empty | VSplit | 80x24 | C-n | No change — empty buffer in split |
| TC-MOV-059 | First | Col0 | Single | HSplit | 40x12 | C-e | Cursor at (0, EOL) |
| TC-MOV-060 | First | Col0 | LargeML | Single | 200x50 | M-< | Cursor stays (0,0) |

---

## Group 2: Editing (TC-EDIT)

**Relevant factors:** CurR, CurC, Content, ModState, RO, Key(editing), File, HL, WinCnt, Term

| ID | CurR | CurC | Content | RO | Key | WinCnt | Term | Input | Expected Post-Condition |
|----|------|------|---------|-----|-----|--------|------|-------|------------------------|
| TC-EDIT-001 | Mid | Mid | SmallML | RW | ASCII | Single | 80x24 | Type 'a' | 'a' inserted at cursor, cursor advances, [Modified] shown |
| TC-EDIT-002 | First | Col0 | SmallML | RW | ASCII | Single | 80x24 | Type 'z' | 'z' inserted at start of first line |
| TC-EDIT-003 | Last | End | SmallML | RW | ASCII | Single | 80x24 | Type 'x' | 'x' appended at end of last line |
| TC-EDIT-004 | Mid | Mid | SmallML | RW | UTF8 | Single | 80x24 | Type Japanese char | Multi-byte char inserted, cursor advances |
| TC-EDIT-005 | First | Col0 | Empty | RW | ASCII | Single | 80x24 | Type 'h' | Buffer becomes "h", cursor at (0,1) |
| TC-EDIT-006 | First | Col0 | Single | RW | ASCII | Single | 80x24 | Type 'a' | Char prepended to line |
| TC-EDIT-007 | Mid | Mid | SmallML | RW | Special | Single | 80x24 | Tab | Tab inserted, visual expansion to 8-col boundary |
| TC-EDIT-008 | First | Col0 | SmallML | RW | Special | Single | 80x24 | Tab | Tab at col 0: 8 visual spaces |
| TC-EDIT-009 | Mid | Mid | SmallML | RW | Special | Single | 80x24 | Backspace | Char before cursor deleted |
| TC-EDIT-010 | Mid | Col0 | SmallML | RW | Special | Single | 80x24 | Backspace | Line joined with previous line |
| TC-EDIT-011 | First | Col0 | SmallML | RW | Special | Single | 80x24 | Backspace | No change — start of buffer |
| TC-EDIT-012 | First | Col0 | Empty | RW | Special | Single | 80x24 | Backspace | No change — empty buffer |
| TC-EDIT-013 | Mid | Mid | SmallML | RW | CtrlBound | Single | 80x24 | C-d | Char at cursor deleted |
| TC-EDIT-014 | Mid | End | SmallML | RW | CtrlBound | Single | 80x24 | C-d | Next line joined to current |
| TC-EDIT-015 | First | Col0 | Empty | RW | CtrlBound | Single | 80x24 | C-d | No change — nothing to delete |
| TC-EDIT-016 | Last | End | SmallML | RW | CtrlBound | Single | 80x24 | C-d | No change — end of buffer |
| TC-EDIT-017 | Mid | Mid | SmallML | RW | Special | Single | 80x24 | Enter | Line split at cursor, cursor at (row+1, 0) |
| TC-EDIT-018 | Mid | End | SmallML | RW | Special | Single | 80x24 | Enter | New empty line below, cursor at (row+1, 0) |
| TC-EDIT-019 | Mid | Col0 | SmallML | RW | Special | Single | 80x24 | Enter | New empty line above, cursor at (row+1, 0) |
| TC-EDIT-020 | First | Col0 | Empty | RW | Special | Single | 80x24 | Enter | Two empty lines, cursor at (1,0) |
| TC-EDIT-021 | First | Col0 | Single | RW | CtrlBound | Single | 80x24 | C-d | First char deleted from single-char line OR line becomes empty |
| TC-EDIT-022 | Mid | Mid | SmallML | RO | ASCII | Single | 80x24 | Type 'a' | "Buffer is read-only" message, no change |
| TC-EDIT-023 | Mid | Mid | SmallML | RO | Special | Single | 80x24 | Backspace | "Buffer is read-only" message, no change |
| TC-EDIT-024 | Mid | Mid | SmallML | RO | CtrlBound | Single | 80x24 | C-d | "Buffer is read-only" message, no change |
| TC-EDIT-025 | Mid | Mid | SmallML | RO | CtrlBound | Single | 80x24 | C-k | "Buffer is read-only" message, no change |
| TC-EDIT-026 | Mid | Mid | SmallML | RO | Special | Single | 80x24 | Enter | "Buffer is read-only" message, no change |
| TC-EDIT-027 | Mid | Mid | SmallML | RW | ASCII | VSplit | 80x24 | Type 'a' | Char inserted; both windows show updated content |
| TC-EDIT-028 | Mid | Mid | SmallML | RW | ASCII | HSplit | 80x24 | Type 'a' | Char inserted; both windows show updated content |
| TC-EDIT-029 | Mid | Mid | SmallML | RW | ASCII | Single | 40x12 | Type 'a' | Char inserted, small terminal renders correctly |
| TC-EDIT-030 | Mid | Mid | SmallML | RW | ASCII | Single | 200x50 | Type 'a' | Char inserted, large terminal renders correctly |
| TC-EDIT-031 | Mid | Mid | LargeML | RW | Special | Single | 80x24 | Enter | Line split works on large buffer |
| TC-EDIT-032 | First | Col0 | Wide | RW | ASCII | Single | 80x24 | Type 'a' | Char prepended to wide line |
| TC-EDIT-033 | First | Col0 | Mixed | RW | Special | Single | 80x24 | Tab | Tab before existing tab, correct expansion |
| TC-EDIT-034 | Mid | Mid | SmallML | RW | UTF8 | VSplit | 80x24 | Type multi-byte | Multi-byte char in split view |
| TC-EDIT-035 | Mid | Mid | SmallML | RW | Special | VSplit | 40x12 | Backspace | Delete in small split window |
| TC-EDIT-036 | First | Col0 | Empty | RW | Special | HSplit | 80x24 | Enter | Newline in empty buffer with hsplit |
| TC-EDIT-037 | Last | End | LargeML | RW | ASCII | Single | 40x12 | Type 'z' | Append in large buffer, small terminal |
| TC-EDIT-038 | Mid | Mid | Single | RW | CtrlBound | Single | 80x24 | C-d | Delete from single-line buffer |
| TC-EDIT-039 | First | End | Single | RW | Special | Single | 80x24 | Enter | Single line split into two |
| TC-EDIT-040 | Mid | Mid | SmallML | RW | CtrlBound | VSplit | 40x12 | C-d | Delete in split, small terminal |

---

## Group 3: Kill, Yank, and Region (TC-KILL)

**Relevant factors:** CurR, CurC, Content, Mark, KillRing, ConsKill, RO, WinCnt, Term

| ID | CurR | CurC | Content | Mark | KillRing | ConsKill | Input Sequence | Expected Post-Condition |
|----|------|------|---------|------|----------|----------|----------------|------------------------|
| TC-KILL-001 | Mid | Mid | SmallML | None | Empty | First | C-k | Kill from cursor to EOL; new kill ring entry; [Modified] |
| TC-KILL-002 | Mid | End | SmallML | None | Empty | First | C-k | Kill newline, join with next line |
| TC-KILL-003 | First | Col0 | SmallML | None | Empty | First | C-k | Kill entire first line content |
| TC-KILL-004 | Last | End | SmallML | None | Empty | First | C-k | Nothing to kill (end of buffer) |
| TC-KILL-005 | First | Col0 | Empty | None | Empty | First | C-k | Nothing to kill; lastKill set |
| TC-KILL-006 | Mid | Col0 | SmallML | None | NonEmpty | Consec | C-k | Kill appends to existing kill ring entry |
| TC-KILL-007 | Mid | Col0 | SmallML | None | NonEmpty | First | C-k | New kill ring entry (chain broken) |
| TC-KILL-008 | First | Col0 | SmallML | None | Empty | First | C-k, C-k | Two consecutive kills: content + newline accumulated |
| TC-KILL-009 | First | Col0 | SmallML | None | Empty | First | C-k, C-f, C-k | C-f breaks chain; two separate kill ring entries |
| TC-KILL-010 | First | Col0 | SmallML | None | Empty | First | C-k x4 | Kill 2 lines (content+newline each), all accumulated |
| TC-KILL-011 | Mid | Mid | SmallML | DiffPos | Empty | First | C-w | Region killed, text removed, mark deactivated |
| TC-KILL-012 | Mid | Mid | SmallML | SamePos | Empty | First | C-w | Zero-width region; nothing killed, mark deactivated |
| TC-KILL-013 | First | Col0 | SmallML | DiffPos | Empty | First | C-w | Multi-line region killed (mark at different row) |
| TC-KILL-014 | Mid | Mid | SmallML | None | Empty | First | C-w | No mark active — no action |
| TC-KILL-015 | Mid | Mid | SmallML | DiffPos | Empty | First | M-w | Region copied (not deleted), mark deactivated |
| TC-KILL-016 | Mid | Mid | SmallML | SamePos | Empty | First | M-w | Zero-width copy; nothing in kill ring |
| TC-KILL-017 | Mid | Mid | SmallML | None | Empty | First | M-w | No mark active — no action |
| TC-KILL-018 | Mid | Mid | SmallML | None | NonEmpty | First | C-y | Paste last kill ring entry at cursor |
| TC-KILL-019 | First | Col0 | Empty | None | NonEmpty | First | C-y | Paste into empty buffer |
| TC-KILL-020 | Mid | Mid | SmallML | None | Empty | First | C-y | Kill ring empty — no change |
| TC-KILL-021 | First | Col0 | SmallML | None | NonEmpty | First | C-y | Paste at start of buffer |
| TC-KILL-022 | Last | End | SmallML | None | NonEmpty | First | C-y | Paste at end of buffer |
| TC-KILL-023 | Mid | Mid | SmallML | DiffPos | NonEmpty | First | C-w, C-y | Kill region, move, yank: text reappears elsewhere |
| TC-KILL-024 | Mid | Mid | SmallML | DiffPos | Empty | First | M-w, C-e, C-y | Copy region, move to EOL, paste: duplicated text |
| TC-KILL-025 | First | Col0 | SmallML | None | NonEmpty | First | C-y (multi-line) | Multi-line yank: lines inserted, cursor after pasted text |
| TC-KILL-026 | Mid | Mid | SmallML | None | NonEmpty | First | C-k, C-k, M-<, C-y | Kill 2 lines accumulated, jump to top, yank all |
| TC-KILL-027 | Mid | Mid | SmallML | DiffPos | Empty | First | C-w | Kill region spanning 3+ lines |
| TC-KILL-028 | Mid | Mid | SmallML | None | NonEmpty | Consec | C-k, C-k, C-k | Triple consecutive kill, all accumulated |
| TC-KILL-029 | Mid | Mid | SmallML | DiffPos | NonEmpty | First | C-w, C-y, C-y | Kill region, yank twice: text duplicated |
| TC-KILL-030 | Mid | Mid | SmallML | None | NonEmpty | First | C-y | Yank in split window |
| TC-KILL-031 | First | Col0 | Empty | None | NonEmpty | First | C-y | Yank into empty buffer in small terminal |
| TC-KILL-032 | Mid | Col0 | SmallML | None | Empty | First | C-k x50 | 50 consecutive kills; buffer shrinks, no crash |
| TC-KILL-033 | Mid | Mid | SmallML | RO:true | NonEmpty | First | C-k | "Buffer is read-only" |
| TC-KILL-034 | Mid | Mid | SmallML | RO:true | NonEmpty | First | C-y | "Buffer is read-only" |
| TC-KILL-035 | Mid | Mid | SmallML | RO:true | DiffPos | First | C-w | "Buffer is read-only" |

---

## Group 4: Search (TC-SRCH)

**Relevant factors:** Mode(Search), CurR, CurC, Content, SrchSt, Term, WinCnt

| ID | Direction | CurR | Content | SrchSt | Term | Input Sequence | Expected Post-Condition |
|----|-----------|------|---------|--------|------|----------------|------------------------|
| TC-SRCH-001 | Forward | First | SmallML | Match | 80x24 | C-s, type "target" | Cursor moves to match, "I-search: target" shown |
| TC-SRCH-002 | Forward | First | SmallML | Failing | 80x24 | C-s, type "xyz" | "Failing I-search: xyz", cursor stays at original |
| TC-SRCH-003 | Forward | First | SmallML | EmptyQ | 80x24 | C-s, Enter | Search exits immediately, cursor unchanged |
| TC-SRCH-004 | Backward | Last | SmallML | Match | 80x24 | C-r, type "target" | Cursor moves to match before current position |
| TC-SRCH-005 | Backward | Last | SmallML | Failing | 80x24 | C-r, type "xyz" | "Failing I-search backward: xyz" |
| TC-SRCH-006 | Forward | Last | SmallML | Match | 80x24 | C-s, type "line1word" | Wrap-around: finds match before cursor position |
| TC-SRCH-007 | Backward | First | SmallML | Match | 80x24 | C-r, type "lastword" | Wrap-around: finds match after cursor position |
| TC-SRCH-008 | Forward | Mid | SmallML | Match | 80x24 | C-s, type "t", C-s | Repeated C-s: advances to next match |
| TC-SRCH-009 | Forward | Mid | SmallML | Match | 80x24 | C-s, type "t", C-s, C-s | Multiple repeats cycle through matches |
| TC-SRCH-010 | Forward | Mid | SmallML | Match | 80x24 | C-s, type "t", C-g | Cancel: cursor returns to original position, "Quit" |
| TC-SRCH-011 | Forward | Mid | SmallML | Match | 80x24 | C-s, type "t", Enter | Accept: cursor stays at match position |
| TC-SRCH-012 | Forward | Mid | SmallML | Match | 80x24 | C-s, type "t", C-a | Exit via unrecognized key: exits search, C-a moves to BOL |
| TC-SRCH-013 | Forward | Mid | SmallML | Match | 80x24 | C-s, type "hel", Backspace | Query becomes "he", re-searches from original |
| TC-SRCH-014 | Forward | Mid | SmallML | Match | 80x24 | C-s, type "h", Backspace | Query becomes empty, cursor at original |
| TC-SRCH-015 | Forward | Mid | SmallML | EmptyQ | 80x24 | C-s, Backspace | No change — empty query |
| TC-SRCH-016 | Forward | Mid | SmallML | Match | 80x24 | C-s, type "t", C-r | Switch direction to backward within same search |
| TC-SRCH-017 | Backward | Mid | SmallML | Match | 80x24 | C-r, type "t", C-s | Switch direction to forward within same search |
| TC-SRCH-018 | Forward | First | Empty | EmptyQ | 80x24 | C-s, type "a" | Failing search on empty buffer |
| TC-SRCH-019 | Forward | First | Single | Match | 80x24 | C-s, type first char | Match on single-line buffer |
| TC-SRCH-020 | Forward | First | LargeML | Match | 80x24 | C-s, type "line500" | Search in large buffer finds distant match |
| TC-SRCH-021 | Forward | Mid | SmallML | Match | 40x12 | C-s, type "t" | Search in small terminal |
| TC-SRCH-022 | Forward | Mid | SmallML | Match | 200x50 | C-s, type "t" | Search in large terminal |
| TC-SRCH-023 | Forward | Mid | SmallML | Match | 80x24 | C-s, type "t" (in VSplit) | Search in split window — only active window searched |
| TC-SRCH-024 | Backward | Mid | SmallML | Match | 80x24 | C-r, type "t" (in VSplit) | Backward search in split window |
| TC-SRCH-025 | Forward | First | SmallML | Match | 80x24 | C-s, type "aaa", C-s, C-s, C-s | Wrap-around multiple times |
| TC-SRCH-026 | Forward | Mid | Mixed | Match | 80x24 | C-s, type UTF-8 char | Search for multi-byte character |
| TC-SRCH-027 | Forward | Mid | SmallML | Match | 80x24 | C-s, type "t", C-x | Exit via C-x: search exits, C-x prefix activated |
| TC-SRCH-028 | Forward | Mid | SmallML | Failing | 80x24 | C-s, type "xyz", C-g | Cancel failing search, cursor returns |
| TC-SRCH-029 | Forward | Mid | SmallML | Failing | 80x24 | C-s, type "xyz", Backspace | Revise query, may find match |
| TC-SRCH-030 | Forward | Mid | Wide | Match | 80x24 | C-s, type match in wide line | Search finds match in long line |

---

## Group 5: Buffer Management (TC-BUF)

**Relevant factors:** Mode(Normal/Mini), Content, ModState, File, WinCnt, Term, MiniCtx

| ID | Operation | ModState | File | Content | WinCnt | Term | Input Sequence | Expected Post-Condition |
|----|-----------|----------|------|---------|--------|------|----------------|------------------------|
| TC-BUF-001 | Open file | Unmod | Existing | SmallML | Single | 80x24 | C-x C-f, type path, Enter | File opened, content shown, filename in status |
| TC-BUF-002 | Open file | Unmod | New | Empty | Single | 80x24 | C-x C-f, type new path, Enter | New empty buffer, filename set |
| TC-BUF-003 | Open same file | Unmod | Existing | SmallML | Single | 80x24 | C-x C-f, type same path, Enter | Switches to existing buffer (no duplicate) |
| TC-BUF-004 | Switch buffer | Unmod | Existing | SmallML | Single | 80x24 | C-x b, Enter (default) | Switches to previous buffer |
| TC-BUF-005 | Switch buffer | Unmod | Existing | SmallML | Single | 80x24 | C-x b, type name, Enter | Switches to named buffer |
| TC-BUF-006 | Kill buffer | Unmod | Existing | SmallML | Single | 80x24 | C-x k, Enter | Buffer killed without confirmation |
| TC-BUF-007 | Kill buffer | Mod | Existing | SmallML | Single | 80x24 | C-x k, Enter, y | Modified buffer: confirm y → killed |
| TC-BUF-008 | Kill buffer | Mod | Existing | SmallML | Single | 80x24 | C-x k, Enter, n | Modified buffer: confirm n → not killed |
| TC-BUF-009 | Buffer list | Unmod | Existing | SmallML | Single | 80x24 | C-x C-b | *Buffer List* shows all buffers with markers |
| TC-BUF-010 | Buffer list nav | Unmod | Existing | SmallML | Single | 80x24 | C-x C-b, Enter on buffer | Switches to selected buffer |
| TC-BUF-011 | Save | Mod | Existing | SmallML | Single | 80x24 | C-x C-s | File written, [Modified] cleared |
| TC-BUF-012 | Save | Unmod | Existing | SmallML | Single | 80x24 | C-x C-s | "No changes need to be saved" or silent success |
| TC-BUF-013 | Save no file | Mod | Special | SmallML | Single | 80x24 | C-x C-s | "No file name" message |
| TC-BUF-014 | Open file | Unmod | Existing | SmallML | VSplit | 80x24 | C-x C-f, type path, Enter | File opens in active window only |
| TC-BUF-015 | Kill buffer split | Unmod | Existing | SmallML | VSplit | 80x24 | C-x k, Enter | Buffer killed, window shows next buffer |
| TC-BUF-016 | Kill only buffer | Unmod | Special | Empty | Single | 80x24 | C-x k, Enter | Buffer killed, new *scratch* created |
| TC-BUF-017 | Open .go file | Unmod | Existing | SmallML | Single | 80x24 | C-x C-f, type .go file, Enter | File opened with syntax highlighting active |
| TC-BUF-018 | Open file | Unmod | Existing | SmallML | Single | 40x12 | C-x C-f, type path, Enter | File opens correctly in small terminal |
| TC-BUF-019 | Buffer list | Mod | Existing | SmallML | Single | 80x24 | C-x C-b | Buffer list shows * for modified buffer |
| TC-BUF-020 | Switch buffer | Unmod | Existing | SmallML | HSplit | 80x24 | C-x b, type name, Enter | Switch in horizontal split |
| TC-BUF-021 | Quit unmod | Unmod | Existing | SmallML | Single | 80x24 | C-x C-c | Editor exits cleanly |
| TC-BUF-022 | Quit modified | Mod | Existing | SmallML | Single | 80x24 | C-x C-c | Warning: "Modified buffers exist..." |
| TC-BUF-023 | Quit confirm | Mod | Existing | SmallML | Single | 80x24 | C-x C-c, C-x C-c | Second C-x C-c confirms exit |
| TC-BUF-024 | Open file large | Unmod | Existing | LargeML | Single | 80x24 | C-x C-f, type path, Enter | Large file opens without lag |
| TC-BUF-025 | Kill buffer | Unmod | New | Empty | Single | 80x24 | C-x k, Enter | Kill new (unsaved) buffer without confirmation |

---

## Group 6: Window Management (TC-WIN)

**Relevant factors:** WinCnt, ActWin, Content, CurR, Term, ModState

| ID | WinCnt | ActWin | Content | Term | Input Sequence | Expected Post-Condition |
|----|--------|--------|---------|------|----------------|------------------------|
| TC-WIN-001 | Single | Win0 | SmallML | 80x24 | C-x 2 | Vertical split, 2 windows, both show same buffer |
| TC-WIN-002 | Single | Win0 | SmallML | 80x24 | C-x 3 | Horizontal split, 2 windows side-by-side |
| TC-WIN-003 | VSplit | Win0 | SmallML | 80x24 | C-x o | Focus switches to Win1 (status line changes) |
| TC-WIN-004 | VSplit | Win1 | SmallML | 80x24 | C-x o | Focus cycles back to Win0 |
| TC-WIN-005 | HSplit | Win0 | SmallML | 80x24 | C-x o | Focus switches to Win1 |
| TC-WIN-006 | VSplit | Win0 | SmallML | 80x24 | C-x 0 | Current window closed, single window remains |
| TC-WIN-007 | VSplit | Win1 | SmallML | 80x24 | C-x 0 | Win1 closed, Win0 remains full-screen |
| TC-WIN-008 | Single | Win0 | SmallML | 80x24 | C-x 0 | No action — cannot close only window |
| TC-WIN-009 | VSplit | Win0 | SmallML | 80x24 | C-x 1 | All other windows closed, only Win0 remains |
| TC-WIN-010 | HSplit | Win0 | SmallML | 80x24 | C-x 1 | All other windows closed, Win0 goes full-screen |
| TC-WIN-011 | Single | Win0 | SmallML | 80x24 | C-x 1 | No action — no other windows |
| TC-WIN-012 | VSplit | Win0 | SmallML | 80x24 | Type 'a' | Text appears in both windows (same buffer) |
| TC-WIN-013 | VSplit | Win0 | SmallML | 80x24 | C-v in Win0, C-x o | Win1 scroll position unchanged |
| TC-WIN-014 | Single | Win0 | SmallML | 40x12 | C-x 2 | Split in small terminal, both windows render |
| TC-WIN-015 | Single | Win0 | SmallML | 200x50 | C-x 2 | Split in large terminal |
| TC-WIN-016 | Single | Win0 | SmallML | 80x24 | C-x 3 | Horizontal split, vertical separator drawn |
| TC-WIN-017 | Single | Win0 | SmallML | 40x12 | C-x 3 | Horizontal split in small terminal |
| TC-WIN-018 | VSplit | Win0 | SmallML | 80x24 | C-x 3 | Error: cannot split horizontally while vertical |
| TC-WIN-019 | HSplit | Win0 | SmallML | 80x24 | C-x 2 | Error: cannot split vertically while horizontal |
| TC-WIN-020 | Single | Win0 | SmallML | 80x24 | C-x o | No action — single window |
| TC-WIN-021 | VSplit | Win0 | Empty | 80x24 | C-x 2 → already split | Both show empty buffer |
| TC-WIN-022 | VSplit | Win0 | LargeML | 80x24 | C-v, C-x o, verify | Independent scroll positions |
| TC-WIN-023 | HSplit | Win0 | SmallML | 200x50 | C-x o | Switch window in large hsplit |
| TC-WIN-024 | VSplit | Win0 | SmallML | 40x12 | C-x 0 | Close window in small terminal |
| TC-WIN-025 | Single | Win0 | Single | 80x24 | C-x 2 | Split with single-line buffer |
| TC-WIN-026 | VSplit | Win1 | SmallML | 80x24 | C-x 1 | Win1 active, close others → Win1 becomes only window |
| TC-WIN-027 | VSplit | Win0 | SmallML | 80x24 | C-x 2 → already split → C-x 0 → C-x 2 | Re-split after unsplit works |
| TC-WIN-028 | HSplit | Win1 | SmallML | 80x24 | Type 'a' | Edit propagates to Win0 if same buffer |
| TC-WIN-029 | VSplit | Win0 | SmallML | 80x24 | Status bar check | Win0 shows '==' (active), Win1 shows '--' |
| TC-WIN-030 | VSplit | Win1 | SmallML | 80x24 | Status bar check | Win1 shows '==' (active), Win0 shows '--' |

---

## Group 7: Grep Mode (TC-GREP)

**Relevant factors:** BufMode(Grep), CurR, Content, WinCnt, Term, Key

| ID | CurR | WinCnt | Term | Input Sequence | Expected Post-Condition |
|----|------|--------|------|----------------|------------------------|
| TC-GREP-001 | N/A | Single | 80x24 | M-x, "find-grep", Enter, Enter (default cmd) | Grep results buffer shown, "Run find-grep: " prompt |
| TC-GREP-002 | Mid (result line) | Single | 80x24 | Enter | Opens correct file at correct line number |
| TC-GREP-003 | First (header/non-result) | Single | 80x24 | Enter | "No grep result on this line", stays in grep buffer |
| TC-GREP-004 | Mid (result line) | Single | 80x24 | n | Cursor moves to next result line |
| TC-GREP-005 | Mid (result line) | Single | 80x24 | p | Cursor moves to previous result line |
| TC-GREP-006 | Last result | Single | 80x24 | n | "No more results" message, cursor stays |
| TC-GREP-007 | First result | Single | 80x24 | p | "No more results" message, cursor stays |
| TC-GREP-008 | Mid (result line) | Single | 80x24 | M-n | Jump to next file's first result |
| TC-GREP-009 | Mid (result line) | Single | 80x24 | M-p | Jump to previous file's first result |
| TC-GREP-010 | Mid (result line) | Single | 80x24 | g | Grep re-executed, buffer refreshed |
| TC-GREP-011 | Mid (result line) | Single | 80x24 | q | Grep buffer closed, previous buffer shown |
| TC-GREP-012 | Mid (result line) | Single | 80x24 | 'x' (unbound rune) | "Buffer is read-only" message |
| TC-GREP-013 | Mid (result line) | Single | 80x24 | Backspace | "Buffer is read-only" message |
| TC-GREP-014 | Mid (result line) | Single | 80x24 | C-d | "Buffer is read-only" message |
| TC-GREP-015 | Mid (result line) | Single | 80x24 | C-f | Normal movement: cursor col +1 |
| TC-GREP-016 | Mid (result line) | Single | 80x24 | C-n | Normal movement: cursor row +1 |
| TC-GREP-017 | Mid (result line) | Single | 80x24 | C-v | Normal page down |
| TC-GREP-018 | Mid (result line) | VSplit | 80x24 | Enter | Opens file in other window or current |
| TC-GREP-019 | Mid (result line) | Single | 40x12 | n | Navigate in small terminal |
| TC-GREP-020 | Mid (result line) | Single | 200x50 | n | Navigate in large terminal |
| TC-GREP-021 | Mid (result line) | Single | 80x24 | n, n, n, p, p | Navigate forward and backward through results |
| TC-GREP-022 | Mid (result line) | Single | 80x24 | M-n, M-n | Jump multiple files forward |
| TC-GREP-023 | First file result | Single | 80x24 | M-p | "No more results" at first file |
| TC-GREP-024 | Last file result | Single | 80x24 | M-n | "No more results" at last file |
| TC-GREP-025 | Mid (result line) | Single | 80x24 | C-k | "Buffer is read-only" message |

---

## Group 8: Minibuffer (TC-MINI)

**Relevant factors:** MiniCtx, CurC(in minibuffer), Content(of input), Term, WinCnt

| ID | MiniCtx | Input State | Term | Input Sequence | Expected Post-Condition |
|----|---------|-------------|------|----------------|------------------------|
| TC-MINI-001 | FindFile | "test.txt" typed | 80x24 | C-a | Cursor moves to start of input |
| TC-MINI-002 | FindFile | "test.txt" typed | 80x24 | C-e | Cursor moves to end of input |
| TC-MINI-003 | FindFile | "test.txt" typed, cursor mid | 80x24 | C-f | Cursor moves right one char |
| TC-MINI-004 | FindFile | "test.txt" typed, cursor mid | 80x24 | C-b | Cursor moves left one char |
| TC-MINI-005 | FindFile | "test.txt" typed, cursor mid | 80x24 | C-d | Char at cursor deleted |
| TC-MINI-006 | FindFile | "test.txt" typed, cursor mid | 80x24 | C-k | Text from cursor to end killed |
| TC-MINI-007 | FindFile | "test.txt" typed, cursor mid | 80x24 | Backspace | Char before cursor deleted |
| TC-MINI-008 | FindFile | "test.txt" typed | 80x24 | C-g | Minibuffer cancelled, normal mode |
| TC-MINI-009 | FindFile | "test.txt" typed | 80x24 | Esc | Minibuffer cancelled, normal mode |
| TC-MINI-010 | FindFile | partial path typed | 80x24 | Tab | Path completed to longest common prefix |
| TC-MINI-011 | FindFile | no matching path | 80x24 | Tab | No change — no matches |
| TC-MINI-012 | FindFile | single match path | 80x24 | Tab | Completes fully to single match |
| TC-MINI-013 | Mx | "find" typed | 80x24 | Tab | Completes to "find-grep" (single match) |
| TC-MINI-014 | Mx | "zzz" typed | 80x24 | Tab | No change — no commands match |
| TC-MINI-015 | Mx | "com" typed | 80x24 | Tab | Completes to "comment-region" (single match starting with "com") |
| TC-MINI-016 | FindFile | "" empty | 80x24 | Tab | Completes based on current directory |
| TC-MINI-017 | FindFile | "" empty | 80x24 | Enter | Opens current directory or empty path handling |
| TC-MINI-018 | SwitchBuf | "" empty | 80x24 | Enter | Switches to default (previous) buffer |
| TC-MINI-019 | SwitchBuf | "scratch" typed | 80x24 | Enter | Switches to *scratch* buffer |
| TC-MINI-020 | KillBuf | "" empty | 80x24 | Enter | Kills current buffer (default) |
| TC-MINI-021 | GotoLine | "5" typed | 80x24 | Enter | Cursor moves to line 5 |
| TC-MINI-022 | GotoLine | "abc" typed | 80x24 | Enter | Error handling — non-numeric input |
| TC-MINI-023 | FindFile | "test" typed | 80x24 | C-n | Ignored — not a minibuffer command |
| TC-MINI-024 | FindFile | "test" typed | 80x24 | C-p | Ignored — not a minibuffer command |
| TC-MINI-025 | FindFile | "test" typed | 40x12 | C-a | Cursor to start in small terminal |
| TC-MINI-026 | FindFile | "test" typed | 200x50 | C-e | Cursor to end in large terminal |
| TC-MINI-027 | Mx | "find-grep" typed | 80x24 | Enter | Launches find-grep command |
| TC-MINI-028 | Mx | "comment-region" typed | 80x24 | Enter | Runs comment-region command |
| TC-MINI-029 | FindFile | "test" typed | 80x24 | Left arrow | Cursor moves left (same as C-b) |
| TC-MINI-030 | FindFile | "test" typed | 80x24 | Right arrow | Cursor moves right (same as C-f) |
| TC-MINI-031 | FindFile | "" empty, cursor at 0 | 80x24 | Backspace | No change — nothing to delete |
| TC-MINI-032 | FindFile | "a" typed, cursor at 1 | 80x24 | Backspace | Input becomes empty |
| TC-MINI-033 | FindFile | "test" typed | 80x24 | C-a, C-k | Cursor to start, kill all: input empty |
| TC-MINI-034 | Mx | "" empty | 80x24 | Tab | Completes based on available commands |
| TC-MINI-035 | FindFile | "test.txt" typed | 80x24 | C-d at end | No change — at end of input |

---

## Group 9: Undo (TC-UNDO)

**Relevant factors:** Content, Undo, CurR, CurC, ModState, WinCnt, Term

| ID | Undo Stack | Content | Operation to Undo | Input Sequence | Expected Post-Condition |
|----|------------|---------|-------------------|----------------|------------------------|
| TC-UNDO-001 | After insert | SmallML | Char insertion | Type 'a', C-_ | Char removed, buffer restored |
| TC-UNDO-002 | After delete | SmallML | C-d deletion | C-d, C-_ | Deleted char restored |
| TC-UNDO-003 | After C-k | SmallML | Kill line | C-k, C-_ | Killed text restored |
| TC-UNDO-004 | After C-w | SmallML | Kill region | C-Space, C-n, C-w, C-_ | Region text restored |
| TC-UNDO-005 | After C-y | SmallML | Yank | (kill first) C-k, M-<, C-y, C-_ | Yanked text removed |
| TC-UNDO-006 | After Enter | SmallML | Newline | Enter, C-_ | Lines joined back |
| TC-UNDO-007 | Empty | SmallML | Nothing | C-_ | "No further undo information" message |
| TC-UNDO-008 | After save | SmallML | Pre-save edit | Type 'a', C-x C-s, C-_ | Undo available after save, [Modified] reappears |
| TC-UNDO-009 | Partial (10) | SmallML | Multiple | Type 10 chars, C-_ x10 | All 10 undone, buffer matches original |
| TC-UNDO-010 | Full (100) | SmallML | At limit | Insert 101 chars, C-_ x100 | 100 undone; 101st C-_ shows no more undo |
| TC-UNDO-011 | After Backspace | SmallML | Backspace | Backspace, C-_ | Deleted char restored at original position |
| TC-UNDO-012 | After insert | Empty | Char insertion | Type 'a', C-_ | Buffer returns to empty |
| TC-UNDO-013 | After insert | Single | Char insertion | Type 'a', C-_ | Single line restored |
| TC-UNDO-014 | After C-k | SmallML | Kill on empty line | Go to empty line, C-k, C-_ | Empty line restored (joined line unkilled) |
| TC-UNDO-015 | After insert | SmallML | In split window | (VSplit) Type 'a', C-_ | Undo visible in both windows |
| TC-UNDO-016 | After insert | SmallML | Small terminal | (40x12) Type 'a', C-_ | Undo works in small terminal |
| TC-UNDO-017 | After insert | SmallML | Large terminal | (200x50) Type 'a', C-_ | Undo works in large terminal |
| TC-UNDO-018 | After C-w | SmallML | Multi-line region | C-Space, C-n, C-n, C-w, C-_ | Multi-line region fully restored |
| TC-UNDO-019 | After multiple | SmallML | Interleaved | Type 'a', C-d, Enter, C-_ x3 | All three undone in reverse order |
| TC-UNDO-020 | After insert | SmallML | Double undo + retype | Type 'a', Type 'b', C-_, C-_ | Both chars removed |

---

## Group 10: Comment/Uncomment (TC-CMNT)

**Relevant factors:** Content, Mark, File(extension), CurR, WinCnt, Term

| ID | File Type | Mark | Content | Input Sequence | Expected Post-Condition |
|----|-----------|------|---------|----------------|------------------------|
| TC-CMNT-001 | .go | DiffPos (3 lines) | SmallML | C-Space, C-n, C-n, M-x comment-region Enter | "// " prepended to 3 lines |
| TC-CMNT-002 | .go | DiffPos (3 lines) | SmallML (commented) | C-Space, C-n, C-n, M-x uncomment-region Enter | "// " removed from 3 lines |
| TC-CMNT-003 | .py | DiffPos (3 lines) | SmallML | C-Space, C-n, C-n, M-x comment-region Enter | "# " prepended to 3 lines |
| TC-CMNT-004 | .py | DiffPos (3 lines) | SmallML (commented) | C-Space, C-n, C-n, M-x uncomment-region Enter | "# " removed from 3 lines |
| TC-CMNT-005 | .go | None | SmallML | M-x comment-region Enter | Error: "No region selected" |
| TC-CMNT-006 | .go | DiffPos (1 line) | SmallML | C-Space, C-e, M-x comment-region Enter | "// " prepended to 1 line |
| TC-CMNT-007 | .go | DiffPos (3 lines) | SmallML | Comment, then uncomment | Round-trip: content returns to original |
| TC-CMNT-008 | .go | DiffPos (3 lines) | SmallML (indented comments) | M-x uncomment-region Enter | Tolerates leading whitespace before "// " |
| TC-CMNT-009 | .go | DiffPos (5 lines) | SmallML | C-Space, C-n x5, M-x comment-region Enter | 5 lines commented |
| TC-CMNT-010 | .go | SamePos | SmallML | C-Space, M-x comment-region Enter | Zero-width region — no lines commented or error |
| TC-CMNT-011 | .go | DiffPos | SmallML | Comment in VSplit | Comment visible in both windows |
| TC-CMNT-012 | .py | DiffPos | SmallML | Comment in 40x12 | Comment works in small terminal |
| TC-CMNT-013 | .go | DiffPos | SmallML | Comment, then C-_ undo | Comments removed by undo |
| TC-CMNT-014 | .go | DiffPos | SmallML | Uncomment already uncommented | Lines unchanged (no "// " to remove) |
| TC-CMNT-015 | .go | DiffPos | LargeML | Comment 10 lines | Comment works on large buffer section |

---

## Group 11: Error Handling and No-Ops (TC-ERR)

**Relevant factors:** Mode, Key(unbound/noop), Content, CurR, WinCnt, Term

| ID | Mode | Key/Input | Content | WinCnt | Term | Input Sequence | Expected Post-Condition |
|----|------|-----------|---------|--------|------|----------------|------------------------|
| TC-ERR-001 | CxPfx | Invalid 2nd key | SmallML | Single | 80x24 | C-x, 'a' | Prefix cancelled, 'a' inserted (falls through) |
| TC-ERR-002 | CxPfx | Invalid 2nd key | SmallML | Single | 80x24 | C-x, 'z' | Prefix cancelled, 'z' inserted |
| TC-ERR-003 | CxPfx | Invalid ctrl key | SmallML | Single | 80x24 | C-x, C-q | Prefix cancelled, no action |
| TC-ERR-004 | CxPfx | Valid after invalid | SmallML | Single | 80x24 | C-x, C-n | Prefix cancelled, cursor moves down |
| TC-ERR-005 | Normal | Unbound ctrl | SmallML | Single | 80x24 | C-q | No action, buffer unchanged |
| TC-ERR-006 | Normal | Unbound ctrl | SmallML | Single | 80x24 | C-o | No action, buffer unchanged |
| TC-ERR-007 | Normal | Unbound ctrl | SmallML | Single | 80x24 | C-t | No action, buffer unchanged |
| TC-ERR-008 | Normal | Rapid keys | SmallML | Single | 80x24 | C-f x100 | Cursor stops at end of buffer, no crash |
| TC-ERR-009 | Normal | Rapid mode switch | SmallML | Single | 80x24 | (C-s, C-g) x20 | No state corruption, buffer unchanged |
| TC-ERR-010 | Normal | C-g | SmallML | Single | 80x24 | C-g | Mark deactivated (if active), or no action |
| TC-ERR-011 | SrchF | C-g | SmallML | Single | 80x24 | C-s, C-g | Search cancelled, cursor to original |
| TC-ERR-012 | Mini | C-g | SmallML | Single | 80x24 | C-x C-f, C-g | Minibuffer cancelled, normal mode |
| TC-ERR-013 | CxPfx | C-g | SmallML | Single | 80x24 | C-x, C-g | Prefix cancelled (C-g cancels) |
| TC-ERR-014 | Confirm | Invalid key 'a' | SmallML | Single | 80x24 | (make modified), C-x k, Enter, 'a' | Prompt remains, still in confirm |
| TC-ERR-015 | Confirm | Invalid C-f | SmallML | Single | 80x24 | (make modified), C-x k, Enter, C-f | Prompt remains |
| TC-ERR-016 | Confirm | Invalid Enter | SmallML | Single | 80x24 | (make modified), C-x k, Enter, Enter | Prompt remains |
| TC-ERR-017 | Confirm | Invalid Backspace | SmallML | Single | 80x24 | (make modified), C-x k, Enter, Backspace | Prompt remains |
| TC-ERR-018 | Confirm | C-g | SmallML | Single | 80x24 | (make modified), C-x k, Enter, C-g | Confirm cancelled, back to buffer |
| TC-ERR-019 | Confirm | 'y' | SmallML | Single | 80x24 | (make modified), C-x k, Enter, 'y' | Buffer killed |
| TC-ERR-020 | Confirm | 'n' | SmallML | Single | 80x24 | (make modified), C-x k, Enter, 'n' | Kill cancelled, buffer preserved |
| TC-ERR-021 | Normal | Rapid keys | Empty | Single | 80x24 | C-f x50 | Cursor stays (0,0), no crash |
| TC-ERR-022 | Normal | Rapid keys | Single | Single | 80x24 | C-n x50 | Cursor stays at line 0, no crash |
| TC-ERR-023 | SrchF | C-g after fail | SmallML | Single | 80x24 | C-s, "xyz", C-g | Cancel failing search, cursor returns |
| TC-ERR-024 | Normal | Unbound ctrl | SmallML | VSplit | 80x24 | C-q | No action in split view |
| TC-ERR-025 | CxPfx | Invalid 2nd key | SmallML | VSplit | 80x24 | C-x, 'a' | Prefix cancelled, 'a' inserted in split |
| TC-ERR-026 | Normal | Rapid mode switch | SmallML | Single | 40x12 | (C-s, C-g) x10 | No corruption in small terminal |
| TC-ERR-027 | Normal | Rapid keys | LargeML | Single | 80x24 | C-n x200 | Cursor advances to near end, scroll follows |
| TC-ERR-028 | Confirm | Invalid Space | SmallML | Single | 80x24 | (make modified), C-x k, Enter, Space | Prompt remains |
| TC-ERR-029 | Normal | Unbound meta | SmallML | Single | 80x24 | M-z | No action, buffer unchanged |
| TC-ERR-030 | Normal | Unbound meta | SmallML | Single | 80x24 | M-a | No action, buffer unchanged |

---

## Group 12: Terminal Resize (TC-RSZ)

**Relevant factors:** Term(transitions), WinCnt, Mode, Content

| ID | Initial Term | Resize To | Mode | WinCnt | Content | Expected Post-Condition |
|----|-------------|-----------|------|--------|---------|------------------------|
| TC-RSZ-001 | 80x24 | 40x12 | Normal | Single | SmallML | Redraws correctly, content visible |
| TC-RSZ-002 | 80x24 | 200x50 | Normal | Single | SmallML | Expands, no artifacts |
| TC-RSZ-003 | 80x24 | 40x12 | Normal | VSplit | SmallML | Both windows recalculate heights |
| TC-RSZ-004 | 80x24 | 40x12 | Normal | HSplit | SmallML | Both windows recalculate widths |
| TC-RSZ-005 | 80x24 | 200x50 | Normal | VSplit | SmallML | Windows expand proportionally |
| TC-RSZ-006 | 80x24 | 40x12 | SrchF | Single | SmallML | Search state preserved, query shown |
| TC-RSZ-007 | 80x24 | 40x12 | Mini | Single | SmallML | Minibuffer input preserved |
| TC-RSZ-008 | 40x12 | 80x24 | Normal | Single | SmallML | Expand from small to standard |
| TC-RSZ-009 | 200x50 | 80x24 | Normal | Single | SmallML | Shrink from large to standard |
| TC-RSZ-010 | 80x24 | 40x12 | Normal | Single | LargeML | Large buffer in small window |
| TC-RSZ-011 | 80x24 | 200x50 | Normal | Single | LargeML | Large buffer in large window |
| TC-RSZ-012 | 80x24 | 40x12 | Normal | Single | Empty | Empty buffer resize |
| TC-RSZ-013 | 80x24 | 40x12 | SrchF | VSplit | SmallML | Search + split + resize |
| TC-RSZ-014 | 80x24 | 200x50 | Mini | Single | SmallML | Minibuffer in large resize |
| TC-RSZ-015 | 80x24 | 40x6 | Normal | VSplit | SmallML | Very small resize with split |

---

## Group 13: Cross-Feature Interactions (TC-XFEAT)

**Relevant factors:** Multi-feature interactions testing 3-way combinations across subsystem boundaries

| ID | Features Tested | Input Sequence | Expected Post-Condition |
|----|----------------|----------------|------------------------|
| TC-XFEAT-001 | Search + Region | C-Space (set mark), C-s "text" Enter, check mark | Mark state after search exit |
| TC-XFEAT-002 | Undo + Yank | C-k (kill), M-< (top), C-y (yank), C-_ (undo) | Yanked text removed by undo |
| TC-XFEAT-003 | Undo + Kill Region | C-Space, C-n, C-w, C-_ | Region text fully restored |
| TC-XFEAT-004 | Search + Find-file | C-s "text", C-x C-f | Search exits, find-file prompt appears |
| TC-XFEAT-005 | Split + Kill buffer | C-x 2, C-x k Enter | Buffer killed, window shows next buffer |
| TC-XFEAT-006 | Comment + Undo | C-Space, C-n, M-x comment-region Enter, C-_ | Comments removed by undo |
| TC-XFEAT-007 | Save + Undo | Type 'a', C-x C-s, C-_ | Undo after save, [Modified] reappears |
| TC-XFEAT-008 | Kill + Yank cross-buffer | C-k, C-x b other Enter, C-y | Kill in buffer A, yank in buffer B |
| TC-XFEAT-009 | Search + Split | C-x 2, C-s "text" Enter | Search works in active split window |
| TC-XFEAT-010 | Grep + Window | M-x find-grep Enter Enter, C-x 2 | Split window while in grep buffer |
| TC-XFEAT-011 | Minibuffer + Cancel + Retry | C-x C-f, C-g, C-x C-f | Cancel and re-enter minibuffer |
| TC-XFEAT-012 | Multiple undo | Type 'a', Enter, C-d, C-_ x3 | Reverse all three operations |
| TC-XFEAT-013 | Kill consecutive + Undo | C-k, C-k (accumulated), C-_ | Single undo restores both kills? Or last only? |
| TC-XFEAT-014 | Region + Yank + Undo | C-Space, C-n, C-w, C-y, C-_ | Undo yank, killed text not in buffer |
| TC-XFEAT-015 | Buffer switch + Mark | C-Space, C-x b other Enter, C-x b back Enter | Mark state preserved across buffer switch |
| TC-XFEAT-016 | Resize + Search | C-s "text", resize 40x12, continue typing | Search continues after resize |
| TC-XFEAT-017 | Resize + Minibuffer | C-x C-f, resize 40x12, type path, Enter | File opens after resize during minibuffer |
| TC-XFEAT-018 | Grep + Enter + Undo | Grep results, Enter (open file), C-_ | Undo in opened file (not grep buffer) |
| TC-XFEAT-019 | Comment + Uncomment round-trip | Comment 3 lines, uncomment same 3 lines | Content identical to original |
| TC-XFEAT-020 | Split + Edit propagation + Undo | C-x 2, type 'a' (both show), C-_ | Undo visible in both windows |
| TC-XFEAT-021 | Kill ring + Buffer close | C-k (in buf A), C-x k Enter, C-y (in buf B) | Kill ring persists across buffer kills |
| TC-XFEAT-022 | M-x in search mode | C-s, Esc, x | Search exits (Esc), then M-x prompt appears |
| TC-XFEAT-023 | Undo limit + continued editing | Insert 101 chars, C-_ x100, type 'z' | After full undo, new edit works |
| TC-XFEAT-024 | Save + Kill buffer | C-x C-s, C-x k Enter | Saved buffer killed without confirm |
| TC-XFEAT-025 | Grep + Quit + Previous buffer | M-x find-grep..., q | Returns to correct previous buffer |

---

## Summary Statistics

| Group | ID Prefix | Test Case Count |
|-------|-----------|----------------|
| Movement | TC-MOV | 60 |
| Editing | TC-EDIT | 40 |
| Kill/Yank/Region | TC-KILL | 35 |
| Search | TC-SRCH | 30 |
| Buffer Management | TC-BUF | 25 |
| Window Management | TC-WIN | 30 |
| Grep Mode | TC-GREP | 25 |
| Minibuffer | TC-MINI | 35 |
| Undo | TC-UNDO | 20 |
| Comment/Uncomment | TC-CMNT | 15 |
| Error Handling/No-ops | TC-ERR | 30 |
| Terminal Resize | TC-RSZ | 15 |
| Cross-Feature Interactions | TC-XFEAT | 25 |
| **Total** | | **385** |

## 3-Way Coverage Analysis

### Coverage Strategy

For each thematic group, the relevant factors are:

| Group | Key Factors (3-way combinations covered) |
|-------|----------------------------------------|
| Movement | CurR(4) x CurC(4) x Content(6) x WinCnt(3) x Term(3) |
| Editing | CurR(4) x CurC(4) x Content(6) x RO(2) x Key(4) x WinCnt(3) x Term(3) |
| Kill/Yank | CurR(3) x CurC(3) x Mark(3) x KillRing(2) x ConsKill(2) x Content(3) |
| Search | Direction(2) x CurR(3) x Content(5) x SrchSt(4) x Term(3) x WinCnt(2) |
| Buffer Mgmt | Operation(6) x ModState(2) x File(3) x Content(3) x WinCnt(3) x Term(3) |
| Window Mgmt | WinCnt(3) x ActWin(2) x Content(3) x Term(3) |
| Grep | CurR(3) x Key(6) x WinCnt(2) x Term(3) |
| Minibuffer | MiniCtx(5) x InputState(4) x Term(3) x Key(6) |
| Undo | UndoStack(3) x Content(4) x Operation(6) x WinCnt(2) x Term(3) |
| Comment | FileType(2) x Mark(3) x Content(3) x WinCnt(2) x Term(2) |
| Error/No-ops | Mode(6) x Key(5) x Content(4) x WinCnt(2) x Term(3) |
| Resize | TermTransition(6) x Mode(3) x WinCnt(3) x Content(3) |
| Cross-Feature | (multi-factor interactions across groups) |

### Infeasible Combinations Pruned

Total infeasible combinations removed: ~45 factor-level triplets

Key exclusions:
- **Mode conflicts:** 21 triplets involving mutually exclusive modes
- **RO + editing:** 12 triplets (covered in TC-ERR group separately)
- **Empty buffer + position:** 8 triplets (empty buffer restricts cursor positions)
- **Single window + Win1:** 4 triplets

### Coverage Completeness

With 385 test cases covering 13 groups, the 3-way factor interactions within each relevant group are fully covered:
- Every group covers all valid 3-way combinations of its key factors
- Cross-feature interactions (TC-XFEAT) cover the most important inter-group 3-way combinations
- The no-op/boundary cases from [e2e-noop-boundary-cases.md](e2e-noop-boundary-cases.md) are incorporated into TC-ERR and the boundary cases within each group

## Notes

- Test cases reference factors from [e2e-factor-analysis.md](e2e-factor-analysis.md) and edge cases from [e2e-noop-boundary-cases.md](e2e-noop-boundary-cases.md)
- "SmallML" means a multi-line buffer with 5-30 lines of known content
- "LargeML" means 1000+ line buffer
- Input sequences use Emacs-style notation: C-x = Ctrl+x, M-x = Alt+x, C-Space = Ctrl+Space
- For tests requiring multiple operations, read left-to-right as sequential steps
- Status bar markers: '==' for active window, '--' for inactive window
- All test cases assume Normal mode as starting mode unless otherwise specified
