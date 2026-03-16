# Find-Grep (grep.go)

The find-grep feature in `grep.go` provides project-wide text search with result navigation, modeled after Emacs `M-x grep`. Results are displayed in a special `*grep*` buffer with a dedicated mode handler for navigating and jumping to matches.

## Architecture

```mermaid
graph TB
    subgraph "User Interaction"
        MX["M-x find-grep"] --> MB["Minibuffer: edit grep command"]
        MB --> EXEC["executeGrepCommand(cmd)"]
    end

    subgraph "Async Execution"
        EXEC --> GR["goroutine: exec.Command('sh', '-c', cmd)"]
        GR --> CH["grepResultCh channel"]
        GR --> PE["PostEvent(KeyNUL) — wake loop"]
    end

    subgraph "Result Processing"
        CH --> LOOP["Event loop: non-blocking receive"]
        LOOP --> PARSE["ParseGrepOutput(stdout)"]
        PARSE --> BUF["Populate *grep* buffer"]
        BUF --> MODE["Set Mode = 'grep'"]
    end

    subgraph "Grep Mode"
        MODE --> GM["grepModeHandler"]
        GM --> JUMP["Enter: open file at line"]
        GM --> NAV["n/p: next/prev result"]
        GM --> FNAV["M-n/M-p: next/prev file"]
        GM --> REF["g: refresh (re-execute)"]
        GM --> QUIT["q: close grep buffer"]
    end
```

## Types

### GrepResult

```go
type GrepResult struct {
    File string  // filepath from grep output
    Line int     // line number
    Text string  // matched line text
}
```

### grepResultMsg

```go
type grepResultMsg struct {
    stdout string
    stderr string
    err    error
}
```

Internal message type for communicating async grep results via channel.

## Global State

| Variable | Type | Purpose |
|----------|------|---------|
| `lastGrepCommand` | `string` | Most recently executed grep command (used by `g` refresh) |
| `editorScreen` | `term.Screen` | Reference to terminal screen for `PostEvent()` from goroutines |
| `grepResultCh` | `chan grepResultMsg` | Channel for receiving async grep results |

## Grep Output Parsing

```mermaid
flowchart TD
    A["ParseGrepOutput(output)"] --> B["Split by newlines"]
    B --> C["For each non-empty line"]
    C --> D["ParseGrepLine(line)"]
    D --> E{"Matches 'filepath:linenum:text'?"}
    E -->|yes| F["GrepResult{File, Line, Text}"]
    E -->|no| G["Skip line"]
    F --> H["Append to results"]
```

`ParseGrepLine` uses `strings.SplitN(line, ":", 3)` and validates that the second field is a valid integer (line number).

## Async Execution Flow

```mermaid
sequenceDiagram
    participant User
    participant MB as Minibuffer
    participant Main as Event Loop
    participant GR as Goroutine
    participant Ch as grepResultCh
    participant Scr as Screen

    User->>MB: M-x find-grep
    MB->>MB: Show default command
    User->>MB: Edit and press Enter
    MB->>Main: callback with command string
    Main->>Main: lastGrepCommand = cmd
    Main->>GR: go executeGrepCommand(cmd)
    Main->>Main: message = "Running grep..."
    GR->>GR: exec.Command("sh", "-c", cmd)
    GR->>Ch: grepResultMsg{stdout, stderr, err}
    GR->>Scr: PostEvent(KeyNUL) — synthetic event
    Note over Main: Next PollEvent() returns
    Main->>Ch: select (non-blocking)
    Ch-->>Main: grepResultMsg

    alt err != nil
        Main->>Main: message = "Grep error: ..."
    else stdout empty
        Main->>Main: message = "Grep: no results"
    else success
        Main->>Main: Create/reuse *grep* buffer
        Main->>Main: Populate with formatted output
        Main->>Main: Set Mode="grep", ReadOnly=true
        Main->>Main: Switch to *grep* buffer
    end
```

The `PostEvent(KeyNUL)` trick is essential: it sends a no-op event to wake the event loop from its blocking `PollEvent()` call, allowing it to check `grepResultCh`.

## Default Grep Command

The default command presented in the minibuffer is:

```
find . -type f -exec grep -nH -e '' {} +
```

The user can edit this freely before execution. The cursor is pre-positioned at the empty search pattern (`''`).

## Grep Mode Handler

When a `*grep*` buffer is active, the `grepModeHandler` intercepts key events:

```mermaid
flowchart TD
    A[Key event in *grep* buffer] --> B{Key?}
    B -->|Enter| C["Parse line under cursor<br/>ParseGrepLine()"]
    C --> D{Valid grep result?}
    D -->|yes| E["Open File + goto Line"]
    D -->|no| F["Ignore"]

    B -->|n| G["Move cursor down<br/>to next result line"]
    B -->|p| H["Move cursor up<br/>to previous result line"]

    B -->|M-n| I["Jump to next file's<br/>first result"]
    B -->|M-p| J["Jump to previous file's<br/>first result"]

    B -->|g| K["Re-execute lastGrepCommand<br/>refresh results"]
    B -->|q| L["Switch to previous buffer<br/>close grep buffer"]

    B -->|other| M["Return false<br/>fall through to normal dispatch"]
```

### Enter: Jump to Source

```mermaid
sequenceDiagram
    participant GM as grepModeHandler
    participant Parse as ParseGrepLine
    participant Bufs as Buffer Pool
    participant Win as Active Window

    GM->>GM: Get current line text
    GM->>Parse: ParseGrepLine(lineText)
    Parse-->>GM: GrepResult{File, Line, Text}

    alt File already open
        GM->>Bufs: Find existing buffer
    else File not open
        GM->>Bufs: NewBufferFromFile(File)
    end

    GM->>Win: Set buffer to target
    GM->>Win: Move cursor to Line
    GM->>Win: AdjustScroll()
```

### M-n / M-p: File-Level Navigation

These commands skip to the first result of the next (or previous) file in the grep output. They parse lines to detect file boundaries, jumping past all results from the current file.

### g: Refresh

Re-executes `lastGrepCommand` via `executeGrepCommand()`, replacing the current grep buffer contents with fresh results. Useful when source files have been edited.

### q: Quit

Switches to `previousBufferIdx` and closes the `*grep*` buffer, restoring the editor to the state before grep was invoked.

## Grep Buffer Format

The `*grep*` buffer displays results in the standard grep format:

```
./main.go:42:    func main() {
./main.go:105:   screen.Clear()
./buffer.go:10:  type Buffer struct {
./buffer.go:55:  func (b *Buffer) InsertChar(ch rune) {
```

Each line is parseable by `ParseGrepLine()`, enabling the Enter-to-jump feature.

## Registration

Commands and mode handlers are registered in `grep.go`'s `init()` function:

```go
func init() {
    RegisterCommand("find-grep", findGrepCommand)
    modeHandlers["grep"] = grepModeHandler
}
```

This decoupled registration means `grep.go` depends on `command.go` (for `RegisterCommand`) and `main.go` (for `modeHandlers`), but neither depends on `grep.go`.
