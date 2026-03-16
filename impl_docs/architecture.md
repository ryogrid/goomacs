# Architecture Overview

goomacs is a lightweight, Emacs-like terminal text editor written in pure Go. It uses the [Chroma](https://github.com/alecthomas/chroma) library for syntax highlighting, [go-enry](https://github.com/go-enry/go-enry) for language detection, [go-runewidth](https://github.com/mattn/go-runewidth) for Unicode character width handling, and a custom ANSI/VT100 terminal backend. The architecture follows a clean four-layer design with an additional command/grep layer.

## Layer Diagram

```mermaid
graph TB
    subgraph "UI Layer (main.go)"
        EL[Event Loop]
        SM[Search Mode]
        MM[Minibuffer Mode]
        CM[Confirm Mode]
        CX[C-x Prefix Handler]
        WM[Window Manager]
        RD[Rendering Pipeline]
    end

    subgraph "Command Layer (command.go, grep.go)"
        CMD[Command Registry]
        CMT[Comment / Uncomment]
        GRP[Find-Grep]
        GP[Grep Output Parser]
        GM[Grep Mode Handler]
    end

    subgraph "Highlighting Layer (highlight.go)"
        HL[Highlighter]
        LX[Chroma Lexer]
        TH[Chroma Theme<br/>monokai]
        CC[RGB to 256-color<br/>Conversion]
    end

    subgraph "Data Layer (buffer.go)"
        BUF[Buffer Struct]
        ED[Editing Operations]
        KR[Kill Ring]
        US[Undo Stack]
        MR[Mark / Region]
        SR[Search Engine]
    end

    subgraph "Terminal Layer (term/)"
        SC[Screen Interface]
        TM[Terminal Implementation]
        IP[Input Parser]
        RN[ANSI 256-color Renderer]
        EV[Event System]
    end

    EL -->|key dispatch| SM
    EL -->|key dispatch| MM
    EL -->|key dispatch| CM
    EL -->|key dispatch| CX
    EL -->|manages| WM
    EL -->|calls| RD
    EL -->|editing commands| BUF
    EL -->|M-x dispatch| CMD
    EL -->|mode handler| GM
    CMD --- CMT
    CMD --- GRP
    GRP --- GP
    GRP -->|async exec| GM
    RD -->|read state| BUF
    RD -->|get styles| HL
    RD -->|SetContent, Show| SC
    EL -->|PollEvent| SC
    BUF -->|owns| HL
    HL --- LX
    HL --- TH
    HL --- CC
    SC --- TM
    TM --- IP
    TM --- RN
    TM --- EV
```

## Module Structure

```
goomacs/                      (module: goomacs)
├── main.go                  UI layer: event loop, window management, rendering, keybinding dispatch
├── buffer.go                Data layer: Buffer struct, text operations, file I/O
├── highlight.go             Highlighting layer: Chroma-based syntax coloring
├── command.go               Command registry, M-x dispatch, comment/uncomment operations
├── grep.go                  Find-grep command, grep output parsing, grep mode handler
├── buffer_test.go           Buffer unit tests (69 tests)
├── main_test.go             Main package tests
├── grep_test.go             Grep output parser tests
├── go.mod                   Module definition (Go 1.24, chroma/v2, go-enry/v2, go-runewidth)
├── go.sum                   Dependency checksums
├── term/                    Terminal layer (internal package)
│   ├── screen.go            Screen interface, Event types, Style, Color, KeyCode constants
│   ├── terminal.go          Terminal struct, raw mode, 256-color ANSI rendering, input parsing
│   └── terminal_test.go     Terminal backend tests (26 tests)
└── e2e/                     End-to-end tests (tmux-based integration tests)
    ├── harness.go           tmux session management and test lifecycle
    ├── assertions.go        Screen capture and golden file assertions
    ├── e2e_test.go          Test main and smoke tests
    ├── basic_test.go        Basic editing tests
    ├── navigation_test.go   Movement and scrolling tests
    ├── search_test.go       Incremental search tests
    ├── killyank_test.go     Kill, yank, undo tests
    ├── buffer_test.go       Multi-buffer management tests
    ├── window_test.go       Window splitting tests
    ├── grep_test.go         Find-grep E2E tests
    ├── highlight_test.go    Syntax highlighting smoke tests
    └── testdata/            Golden files for snapshot assertions
```

## Data Flow

### Input Processing

```mermaid
sequenceDiagram
    participant Stdin
    participant InputParser as Input Parser<br/>(goroutine)
    participant EventChan as Event Channel
    participant EventLoop as Event Loop<br/>(main.go)
    participant Buffer as Buffer<br/>(buffer.go)

    Stdin->>InputParser: raw bytes
    InputParser->>InputParser: parse control chars,<br/>ANSI sequences, UTF-8
    InputParser->>EventChan: KeyEvent / ResizeEvent
    EventChan->>EventLoop: PollEvent()
    EventLoop->>EventLoop: dispatch by mode
    EventLoop->>Buffer: editing method<br/>(InsertChar, KillLine, etc.)
    Buffer->>Buffer: modify Lines,<br/>update cursor,<br/>set HighlightDirty
```

### Rendering Pipeline

```mermaid
sequenceDiagram
    participant Loop as Event Loop
    participant Win as Window
    participant Buf as Buffer
    participant HL as Highlighter
    participant Screen as Screen<br/>(term.Terminal)
    participant Stdout

    Loop->>Win: AdjustScroll() (active window only)
    Loop->>Screen: Clear()

    loop each window
        Loop->>Buf: check HighlightDirty
        alt HighlightDirty and Highlighter != nil
            Loop->>HL: Highlight(buf.Lines)
            Loop->>Buf: HighlightDirty = false
        end
        loop each visible cell
            Loop->>HL: StyleAt(row, col)
            HL-->>Loop: base style (colors + bold)
            Loop->>Buf: InRegion(row, col)?
            Note over Loop: overlay reverse video<br/>for region/search
            Loop->>Screen: SetContent(x, y, ch, style)
        end
        Loop->>Screen: drawWindowStatusLine()
    end

    Loop->>Screen: ShowCursor(x, y)
    Loop->>Screen: Show()
    Screen->>Screen: diff cells vs prev
    Screen->>Stdout: ANSI 256-color<br/>escape sequences
```

## Multi-Buffer and Window Architecture

```mermaid
graph TB
    subgraph "Buffer Pool"
        B0["buffers[0]<br/>*scratch*"]
        B1["buffers[1]<br/>main.go"]
        B2["buffers[2]<br/>buffer.go"]
    end

    subgraph "Window Layout (vertical split)"
        W0["windows[0]<br/>StartRow=0, Height=12<br/>StartCol=0, Width=80"]
        W1["windows[1]<br/>StartRow=12, Height=12<br/>StartCol=0, Width=80"]
    end

    subgraph "Window Layout (horizontal split)"
        W2["windows[0]<br/>StartRow=0, Height=24<br/>StartCol=0, Width=40"]
        W3["windows[1]<br/>StartRow=0, Height=24<br/>StartCol=41, Width=39"]
    end

    W0 -->|"Buffer pointer"| B1
    W1 -->|"Buffer pointer"| B2

    ABI["activeBufferIdx = 1"]
    AWI["activeWindowIdx = 0"]
    SPM["splitMode = 'vertical' | 'horizontal'"]
```

- **Buffers** are stored in a flat slice (`buffers []*Buffer`). Each buffer holds its own content, cursor, kill ring, undo stack, optional `Highlighter`, and buffer-local `Mode` (e.g., `"grep"` for find-grep results).
- **Windows** are viewports into buffers. Each `Window` has its own `ScrollOffset`, `StartRow`, `Height`, `StartCol`, and `Width`. Multiple windows can reference the same buffer.
- `recalcWindows()` distributes screen space evenly among all windows. In vertical split mode (`C-x 2`), windows are stacked top-to-bottom. In horizontal split mode (`C-x 3`), windows are placed side-by-side with a `│` separator.
- `AdjustScroll()` is only called for the **active** window to prevent scroll bleeding when multiple windows share a buffer.

## File Responsibilities

| File | Responsibility |
|------|----------------|
| `main.go` | Event loop, multi-mode key dispatch (search, minibuffer, confirm, C-x prefix, buffer mode, normal), window management (vertical/horizontal split, close, switch), rendering pipeline (`drawWindowContent`, `drawWindowStatusLine`, `drawMessageLine`), tab expansion, minibuffer with cursor movement and tab completion |
| `buffer.go` | `Buffer` struct, cursor movement, editing (insert, delete, kill, yank), mark/region, incremental search, undo/redo, file I/O, highlight dirty flag, buffer-local `Mode` and `ReadOnly` flags |
| `command.go` | `Command` registry, `RegisterCommand` / `FindCommand` / `FindCommandsByPrefix` for M-x dispatch, `CommentStyle` types, language-aware comment/uncomment region using go-enry |
| `grep.go` | `findGrepCommand` (M-x find-grep), async grep execution, `ParseGrepOutput` / `ParseGrepLine` parsers, grep mode handler (Enter to jump, n/p navigation, M-n/M-p file navigation, g refresh, q quit) |
| `highlight.go` | `Highlighter` struct, Chroma lexer/theme integration, tokenization, RGB-to-256-color conversion, per-cell style caching |
| `term/screen.go` | `Screen` interface, `Event`/`KeyEvent`/`ResizeEvent` types, `Style` struct (fg/bg/reverse/bold), `Color` type, `KeyCode`/`ModMask` constants |
| `term/terminal.go` | `Terminal` struct implementing `Screen`, raw mode via termios syscalls, 256-color ANSI rendering with cell diffing, `writeStyledCell()` helper, keyboard input parsing, SIGWINCH resize handling |
| `e2e/` | tmux-based end-to-end test harness, screen capture assertions, golden file snapshot testing covering all major features |

## Design Principles

- **Separation of concerns** -- Buffer knows nothing about the terminal; the terminal knows nothing about text editing; main.go bridges them. The command layer registers handlers independently.
- **Extensible mode system** -- Buffer-local modes (`modeHandlers` map) allow special buffers like `*grep*` to override key dispatch without modifying the core event loop.
- **Minimal abstraction** -- No framework, no plugin system. Six source files plus one package.
- **Lazy re-highlighting** -- Syntax highlighting only re-runs when content changes (`HighlightDirty` flag), not on every redraw.
- **Snapshot-based undo** -- Full buffer state is saved before each edit. Simple but bounded (max 100 entries).
- **Diff-based rendering** -- Only changed screen cells are written to stdout, minimizing I/O.
- **Per-window scroll isolation** -- Each window maintains its own `ScrollOffset`; `AdjustScroll()` is only called for the active window.
- **Async command execution** -- Long-running commands (e.g., grep) run in goroutines, posting results via channels and waking the event loop with `PostEvent()`.
