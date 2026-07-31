# Alt-Tab

[![CI](https://github.com/jake-mckenzie/alt-tab/actions/workflows/ci.yml/badge.svg)](https://github.com/jake-mckenzie/alt-tab/actions/workflows/ci.yml)

Alt-Tab is an interactive guitar-chord viewer for the terminal. A responsive
Bubble Tea interface displays chord voicings from a C11 chord library connected
to Go through cgo.

Diagrams follow standard tablature orientation: high e is at the top, low E is
at the bottom, and fret numbers increase from left to right. Compact mode shows
the relevant neck position; full-neck mode shows frets 1–27.

## Quick Start

```bash
git clone https://github.com/jake-mckenzie/alt-tab.git
cd alt-tab
make
./bin/alt-tab
```

`make run` combines the build and run steps:

```bash
make run
```

Before building or running, Make automatically checks for Go, a C compiler,
cgo, a resolvable dependency graph, and valid dependency checksums.

Alt-Tab is an interactive TUI and must be run in a terminal. It does not accept
the command-line chord and display flags used by the former interface.

## Requirements

- macOS, Linux, or another Unix-like environment
- Go 1.25 or newer with cgo enabled
- A C11 compiler such as Clang or GCC
- Make

The default test build also requires AddressSanitizer and
UndefinedBehaviorSanitizer support from the C compiler.

Confirm the required tools are available:

```bash
go version
gcc --version
make --version
go env CGO_ENABLED
```

`go env CGO_ENABLED` must print `1`.

## Building

Run the dependency check by itself:

```bash
make check-deps
```

Build the development binary:

```bash
make
```

The resulting executable is:

```text
./bin/alt-tab
```

Build a smaller release binary with paths and debug symbols removed:

```bash
make BUILD=release
```

You can also build directly with Go:

```bash
mkdir -p bin
go build -o bin/alt-tab ./cmd/alt-tab
```

Missing dependencies listed in `go.mod` and `go.sum` are downloaded
automatically, while cached dependencies work offline. Modules are resolved in
read-only mode and verified against their recorded checksums. The check runs
automatically before `make`, `make run`, and `make test`.

## Running

Build and launch:

```bash
make run
```

Launch an existing build:

```bash
./bin/alt-tab
```

### Controls

| Key | Action |
| --- | --- |
| `←` / `h` | Previous chord |
| `→` / `l` | Next chord |
| `↑` / `k` | Previous variation |
| `↓` / `j` | Next variation |
| `f` | Toggle compact or full-neck view |
| `?` | Open or close help |
| `Esc` | Close help |
| `q` / `Ctrl+C` | Quit |

The complete horizontal chord list appears above the fretboard in both compact
and full-neck modes. Full-neck mode requires a terminal width of at least 98
columns.

## TUI Output

The default layout appears as follows in an 80-column terminal. Terminal colors
are omitted here; Alt-Tab adapts them to light and dark backgrounds.

<!-- BEGIN VERIFIED TUI OUTPUT -->
```text
  ALT-TAB  Guitar Chord Viewer

  ╭────────────────────────────────────────────────────────────────────────╮
  │ CHORDS                                                                 │
  │                                                                        │
  │ ‹ A ›   Am   B   Bb   C   Cm   D   Dm   E   Em   F   F#   G   Gm       │
  ╰────────────────────────────────────────────────────────────────────────╯

  ╭────────────────────────────────────────────────────────────────────────╮
  │ A  ·  variation 1 of 2                                                 │
  │ compact                                                                │
  │                                                                        │
  │         1    2    3    4                                               │
  │ e    O--------------------|                                            │
  │ B    |-------3------------|                                            │
  │ G    |-------2------------|                                            │
  │ D    |-------1------------|                                            │
  │ A    O--------------------|                                            │
  │ E    X--------------------|                                            │
  │                                                                        │
  │ Fingers: 1 index  2 middle                                             │
  │          3 ring   4 little                                             │
  │ Symbols: O open  X muted                                               │
  ╰────────────────────────────────────────────────────────────────────────╯

  ←/→ chord  ↑/↓ variation  f full neck  ? help  q quit
```
<!-- END VERIFIED TUI OUTPUT -->

## Testing

Run the C backend tests and all Go tests:

```bash
make test
```

The default test configuration enables AddressSanitizer and
UndefinedBehaviorSanitizer for the C backend. To test an optimized build without
sanitizers, clean first so every C object is rebuilt with release flags:

```bash
make clean
make BUILD=release test
```

Optional Go-only checks:

```bash
make fmt-check
make test-go
make race
make vet
```

Run every formatting, test, race, and vet check used by CI:

```bash
make check
```

The suite verifies every stored fret and finger assignment, chord pitch
classes, the public C API, cgo conversion, TUI navigation, compact and full-neck
diagrams, help behavior, and the README output snapshot.

## Cleaning

Remove the application binary and generated test objects:

```bash
make clean
```

## Build Troubleshooting

If cgo is disabled, enable it for the build:

```bash
CGO_ENABLED=1 make
```

Select a specific compiler when the default `cc` is unsuitable:

```bash
make CC=clang
make CC=gcc
```

If dependency checksums or downloads fail, refresh the module cache and retry:

```bash
go mod download
make
```

## Supported Chords

| Major | Minor | Accidental |
| --- | --- | --- |
| `A`, `B`, `C`, `D`, `E`, `F`, `G` | `Am`, `Cm`, `Dm`, `Em`, `Gm` | `Bb`, `F#` |

Each chord has two conventional variations. Finger numbers use `1` for index,
`2` for middle, `3` for ring, and `4` for little finger. `O` marks an open
string and `X` marks a muted string.

## Architecture

- `cmd/alt-tab` contains the application entry point.
- `internal/tui` owns Bubble Tea state, navigation, layout, styling, and
  fretboard rendering.
- `internal/chords` defines Go-owned chord types and the catalog interface.
- `internal/nativechords` owns the cgo adapter, stable C API, and immutable
  chord library.
- `tests` contains the native backend regression suite.
- `.github/workflows` runs builds and quality checks for pushes and pull
  requests.

Terminal concerns remain separate from the chord model so future audio input or
output packages can consume the same Go chord types.

## Project Layout

```text
.
├── .github/workflows/
├── cmd/alt-tab/
├── internal/chords/
├── internal/nativechords/
├── internal/tui/
├── tests/
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Roadmap

- Additional chord voicings
- Explicit barre visualization
- Audio input and chord detection
- Audio output and chord playback
- Piano mode
