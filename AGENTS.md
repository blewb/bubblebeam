# AGENTS.md

## Project Overview

**BubbleBeam** is a Go TUI (Terminal User Interface) application for logging hours to the Streamtime time-tracking system. It parses a simple text-based time log format, then provides an interactive interface to map entries to Streamtime jobs and job items via their API.

This is an actively developed personal tool — incomplete in places, with some TODO/commented-out code. See `PLAN.md` for the full product vision and planned features.

## Commands

### Build & Run

```sh
# Run the main TUI app (requires .env and temp/jobs.json)
go run ./cmd/app -s 1

# Run the API test harness
go run ./cmd/streamtest
```

### Flags (`cmd/app`)

| Flag | Short | Description |
|------|-------|-------------|
| `-sample N` | `-s N` | Use sample data file `sample/N.txt` |
| `-state N` | `-t N` | Launch into a specific app state (0–5) |

### Tests

```sh
# Run all tests
go test ./...

# Run span parser tests only
go test ./span/
```

There is no linter, formatter, or CI configuration in the project.

## Environment

The app requires a `.env` file in the project root (gitignored) with:

```
API_URL=<streamtime api base url>
API_TOKEN=<bearer token>
USER_ID=<numeric user id>
```

Both `cmd/app` and `cmd/streamtest` load this via `godotenv`.

## Development Shortcut

Job loading (`stream/jobs_search.go`) is currently hardcoded to read from `temp/jobs.json` instead of hitting the live API. The API call is commented out. The `temp/` directory is gitignored. To develop against the TUI you need a `temp/jobs.json` file containing a Streamtime search response.

## Project Structure

```
bubblebeam/
├── cmd/
│   ├── app/           # Main TUI application (BubbleTea)
│   │   ├── main.go    # Entry point, flag parsing, .env loading
│   │   ├── model.go   # BubbleTea model, state enum, Init()
│   │   ├── update.go  # Update() and per-state update handlers
│   │   ├── view.go    # View() and per-state view renderers
│   │   ├── boot.go    # Builder functions for BubbleTea components
│   │   ├── styles.go  # Lipgloss styles and header rendering
│   │   ├── jobs.go    # Job search logic (filtering loaded jobs)
│   │   └── title.txt  # ASCII art title (embedded via go:embed)
│   └── streamtest/    # Standalone script for testing the API layer
│       └── main.go
├── span/              # Library: text-based time log parser
│   ├── span.go        # Span type, file reading, day detection
│   ├── day.go         # Day type, validation, conflict detection
│   ├── entry.go       # Entry/Timestamp types, duration calc
│   ├── parser.go      # Line parsing (timestamps, ranges, tags)
│   ├── parser_test.go # Table-driven tests for parser
│   └── dates.go       # Datestamp generation for date selection UI
├── stream/            # Library: Streamtime API client
│   ├── api.go         # API struct, HTTP client, get/post/put helpers
│   ├── jobs.go        # Job/JobItem type definitions, ParsedJob structs
│   ├── jobs_search.go # LoadJobs() — fetches/loads active jobs list
│   ├── jobs_items.go  # GetJobItems(), GetJobItemUsers()
│   ├── jobs_search.json # Embedded JSON payload for search API call
│   └── todo.go        # Todo struct for submitting time entries
├── sample/            # Sample time-log text files for development
│   ├── 1.txt, 2.txt, 3.txt
└── PLAN.md            # Detailed product plan and feature spec
```

## Architecture & Key Patterns

### BubbleTea Model-Update-View (Elm Architecture)

The TUI in `cmd/app` follows the standard BubbleTea pattern:

- **Model** (`model.go`): Single `model` struct holds all state. A `modelState` enum (iota) drives which screen is shown.
- **Update** (`update.go`): Top-level `Update()` handles global keys (quit) and window resize, then dispatches to per-state update methods (`UpdateSelectDate`, `UpdateListEntries`, `UpdateSelectJob`, `UpdateSelectItem`). State update methods use pointer receivers (`*model`).
- **View** (`view.go`): Top-level `View()` switches on state and delegates to per-state view methods (`ViewSelectDate`, `ViewListEntries`, etc.). View methods use value receivers (`model`).

### App States (in order)

```
StateLoading → StateSelectDate → StateListEntries → StateSelectJob → StateSelectItem → StateConfirm
```

### Two Library Packages

1. **`span`** — Parses plain-text time logs. Format: lines starting with a weekday name begin a new day; entry lines are `HH:MM - HH:MM description #tag`. The parser validates timestamps (15-minute increments only), detects overlapping entries, and calculates durations.

2. **`stream`** — Streamtime REST API client. Wraps HTTP with bearer token auth. Types come in `Parsed*` (matching JSON) and flattened forms (`Job`, `JobItem`). The `API` struct manages its own state (`APIState` enum).

### Styling

- Uses **Lipgloss** for all styling (borders, colors, alignment).
- Uses **Bubbles** components: `table.Model`, `textinput.Model`, `paginator.Model`.
- Brand colors are blue gradient: `#00ccff`, `#00bbff`, `#00aaff`.
- ASCII title is embedded from `title.txt` via `//go:embed`.

## Code Conventions

- **Go 1.25** (uses modern features like range-over-int where applicable).
- **No external test framework** — standard `testing` package with table-driven tests using `map[string]struct{...}` pattern.
- **Constructors** follow `New*` or `Build*` naming: `NewSpan()`, `NewAPI()`, `BuildTable()`, `BuildPaginator()`.
- **Constants** use `UPPER_SNAKE_CASE` (e.g., `SELECTION_DAYS`, `REF_LENGTH`, `SLOTS_IN_DAY`).
- **Type enums** use `iota` with a descriptive type name (e.g., `modelState`, `APIState`).
- **JSON struct tags** on types that serialize/deserialize; `Parsed*` prefix for raw API response types.
- **Pointer receivers** for mutation methods, **value receivers** for read-only methods.
- **Blank lines** between logical sections within functions are used liberally.
- No explicit error wrapping — errors are returned or logged directly.

## Time Log Format (`sample/*.txt`)

```
Monday
------
 9:00 - 12:00 Did some work #job1
12:30 - 15:15 Another big chunk
15:15 - 17:00 Just a long meeting #job6
```

- Weekday name on its own line starts a new day.
- `------` separator lines and lines < 3 chars are skipped.
- Text before the first weekday is ignored.
- Time format: `HH:MM - HH:MM` (space-padded single-digit hours, e.g., ` 9:00`).
- Minutes must be in 15-minute increments (0, 15, 30, 45).
- Optional `#tag` at end of line (last `#` wins if multiple).

## Gotchas

- **`temp/` directory required**: The app will crash if `temp/jobs.json` doesn't exist, since `LoadJobs()` reads from disk rather than the API.
- **Must run from project root**: Sample files are referenced as `sample/N.txt` (relative path), and `.env` is loaded from CWD.
- **`Title.txt` is embedded**: Changing it requires recompilation.
- **`post()` and `put()` on API are defined but unused** — they exist for future time entry submission.
- **Job search is client-side**: All jobs are loaded upfront into memory; `SearchJobs()` in `cmd/app/jobs.go` filters the in-memory list, not via API.
- **State transitions use direct assignment**: `m.state = StateSelectJob` — no command/message pattern for state changes.
- **`jobsSearchJSON` is embedded but unused** in the current dev flow (API call is commented out in `jobs_search.go`).
