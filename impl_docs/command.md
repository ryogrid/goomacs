# Command System (command.go)

The command system in `command.go` provides a registry for named commands accessible via `M-x`, along with language-aware comment/uncomment operations using [go-enry](https://github.com/go-enry/go-enry) for language detection.

## Architecture

```mermaid
graph TB
    subgraph "M-x Flow"
        MX["M-x key press"] --> MB["Minibuffer prompt<br/>'M-x '"]
        MB --> TAB["Tab completion<br/>FindCommandsByPrefix()"]
        MB --> ENTER["Enter"]
        ENTER --> FC["FindCommand(name)"]
        FC --> EXEC["command.Fn(buf, &message)"]
    end

    subgraph "Command Registry"
        REG["RegisterCommand(name, fn)"]
        CMDS["commands []Command"]
        REG --> CMDS
        FC --> CMDS
        TAB --> CMDS
    end

    subgraph "Registered Commands"
        CR["comment-region"]
        UR["uncomment-region"]
        FG["find-grep"]
    end

    CMDS --- CR
    CMDS --- UR
    CMDS --- FG
```

## Types

### Command

```go
type Command struct {
    Name string
    Fn   func(*Buffer, *string)
}
```

- `Name` -- Command name as typed in the M-x prompt (e.g., `"comment-region"`)
- `Fn` -- Function receiving the active buffer and a pointer to the message string for status feedback

### CommentStyle

```go
type CommentStyle struct {
    LinePrefix string  // e.g., "//"
    BlockStart string  // e.g., "/*"
    BlockEnd   string  // e.g., "*/"
}
```

Defines comment delimiters for a programming language. `LinePrefix` is used for single-line comments; `BlockStart`/`BlockEnd` are used for block comments.

## Registry API

| Function | Purpose |
|----------|---------|
| `RegisterCommand(name, fn)` | Adds a command to the global registry |
| `FindCommand(name) *Command` | Returns exact match by name, or nil |
| `FindCommandsByPrefix(prefix) []Command` | Returns all commands with matching prefix (for tab completion) |

Commands are registered via `init()` functions, allowing `grep.go` and `command.go` to independently register their commands without coupling.

## Comment / Uncomment Operations

### Language Detection

```mermaid
flowchart TD
    A["detectCommentStyle(filename, content)"] --> B["go-enry: GetLanguage(filename, content)"]
    B --> C{Language known?}
    C -->|yes| D["Look up commentStyles map"]
    C -->|no| E["Default: LinePrefix = '#'"]
    D --> F{Found in map?}
    F -->|yes| G["Return CommentStyle"]
    F -->|no| E
    E --> G
```

### Supported Languages

The `commentStyles` map covers 30+ languages:

| Language Group | Languages | Line Prefix | Block Delimiters |
|---------------|-----------|-------------|-----------------|
| C-family | Go, C, C++, Java, JavaScript, TypeScript, Rust, Kotlin, Swift, PHP, SCSS | `//` | `/* */` |
| Scripting | Python, Ruby, Shell, Perl, R, YAML, TOML, Makefile | `#` | -- |
| Lisp-family | Lisp, Clojure, Scheme | `;;` | -- |
| Markup | HTML, XML | -- | `<!-- -->` |
| Functional | Haskell | `--` | `{- -}` |
| Other | Lua (`--`), SQL (`--`), Erlang (`%`), Elixir (`#`), CSS (`/* */`) | varies | varies |

### comment-region

```mermaid
flowchart TD
    A["comment-region called"] --> B{MarkActive?}
    B -->|no| C["message: 'No region selected'"]
    B -->|yes| D["detectCommentStyle(filename, content)"]
    D --> E["Get region bounds (startR, endR)"]
    E --> F["SaveUndo()"]
    F --> G["Prepend LinePrefix + space to each line"]
    G --> H["HighlightDirty = true"]
```

### uncomment-region

```mermaid
flowchart TD
    A["uncomment-region called"] --> B{MarkActive?}
    B -->|no| C["message: 'No region selected'"]
    B -->|yes| D["detectCommentStyle(filename, content)"]
    D --> E["Get region bounds (startR, endR)"]
    E --> F["SaveUndo()"]
    F --> G["For each line in region"]
    G --> H{Has LinePrefix?}
    H -->|yes| I["Strip LinePrefix (+ optional space)"]
    H -->|no| J{Has BlockStart/BlockEnd?}
    J -->|yes| K["Strip block delimiters"]
    J -->|no| L["Skip line"]
    I --> M["HighlightDirty = true"]
    K --> M
```

Uncomment tolerates leading whitespace before comment markers.
