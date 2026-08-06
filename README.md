# Alt-Tab

[![CI](https://github.com/jake-mckenzie/alt-tab/actions/workflows/ci.yml/badge.svg)](https://github.com/jake-mckenzie/alt-tab/actions/workflows/ci.yml)

Alt-Tab is an interactive guitar-chord viewer for the terminal. A responsive
Bubble Tea interface displays chord voicings from a C11 chord library connected
to Go through cgo.

Diagrams follow standard tablature orientation: high e is at the top, low E is
at the bottom, and fret numbers increase from left to right. Compact mode shows
the relevant neck position; full-neck mode shows frets 1–27.

## TUI Preview

The default layout appears as follows in an 80-column terminal. Terminal colors
are omitted here; Alt-Tab adapts them to light and dark backgrounds.

<!-- BEGIN VERIFIED TUI OUTPUT -->
```text
  ╭──────────────────────────────────────────────────────────────────────────╮
  │                 ALT-TAB  Guitar Chord Viewer · Synthwave                 │
  │                                                                          │
  │ KEYS ←→ Base ↑↓ Type v Voicing f Neck n Tab t Theme w Wave ? Help q Quit │
  ╰──────────────────────────────────────────────────────────────────────────╯

  ╭──────────────────────────────────────────────────────────────────────────╮
  │                                CHORD DIAL                                │
  │                                                                          │
  │                     BASE CHORDS  A  B  C  D  E  F  G                     │
  │                                                                          │
  │                                                                          │
  │                                                                          │
  │                               ╭──────────╮                               │
  │                       G       │  ‹ A ›   │      B                        │
  │                               │          │                               │
  │                               │    Am    │                               │
  │                               ╰──────────╯                               │
  ╰──────────────────────────────────────────────────────────────────────────╯

  ╭────────────────────────────╮  ╭──────────────────────────────────────────╮
  │       CHORD DIAGRAM        │  │              NOTE WAVEFORM               │
  │          COMPACT           │  │                                          │
  │                            │  │  +1.0 |                               |  │
  │             A              │  │       | ⣷                     ⡼⡀      |  │
  │       VOICING 1 OF 2       │  │       |⢰⠉⡆                    ⡇⡇      |  │
  │                            │  │       |⢸ ⡇                   ⢸ ⢇      |  │
  │       1    2    3    4     │  │       |⢸ ⢱                   ⢸ ⢸      |  │
  │     ────────────────────   │  │       |⡇ ⢸  ⡀      ⡤⡀⢀⠎⡆ ⡰⡄  ⢸ ⠸⡀ ⡀   |  │
  │ e  O--------------------|  │  │   0.0 +⠧⠤⠬⡦⡴⠽⡤⠤⡴⢦⠤⡼⠤⠵⠮⠤⠼⠴⠥⢵⠤⠤⡧⠤⠤⡧⡴⠽⡤⠤⡼+  │
  │ B  |-------3------------|  │  │       |   ⠱⠃ ⢇⢰⠁ ⠓⠁       ⢸  ⡇  ⠱⠃ ⢣⢠⠃|  │
  │ G  |-------2------------|  │  │       |      ⠈⠁           ⠈⡆ ⡇     ⠈⠊ |  │
  │ D  |-------1------------|  │  │       |                    ⡇⢰⠁        |  │
  │ A  O--------------------|  │  │       |                    ⢱⢸         |  │
  │ E  X--------------------|  │  │       |                    ⠸⡜         |  │
  │                            │  │  -1.0 |                     ⠁         |  │
  │ Fingers: 1 index  2 middle │  │        0 ms                       25 ms  │
  │          3 ring   4 little │  │      Normalized amplitude over time      │
  │                            │  │                                          │
  │ Symbols: O open  X muted   │  │        │          │        │     │    │  │
  │                            │  │        ├──────────┼────────┼─────┼────┤  │
  │                            │  │        A2        E3       A3    C#4  E4  │
  │                            │  │        110 Hz                    330 Hz  │
  │                            │  │          Pitch range by semitone         │
  │                            │  │                                          │
  │                            │  │  Notes: E4  C#4  A3  E3  A2              │
  ╰────────────────────────────╯  ╰──────────────────────────────────────────╯
```
<!-- END VERIFIED TUI OUTPUT -->

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
| `↑` / `k` | Select an available accidental, or return from minor to base |
| `↓` / `j` | Select an available minor, or return from accidental to base |
| `v` | Cycle through voicings of the selected chord |
| `f` | Toggle compact or full-neck view |
| `n` | Toggle the fingered fretboard or fret-number tab view |
| `t` | Cycle color theme |
| `w` | Toggle high-resolution note waveform |
| `?` | Open or close help |
| `Esc` | Close help |
| Mouse wheel | Scroll the TUI vertically |
| `q` / `Ctrl+C` | Quit |

The base-chord row lists A through G above a three-position horizontal dial.
When the selected family has a flat, sharp, or minor chord, a nested vertical
dial shows only those real choices around the base chord. Families without
variants remain a single center cell, and unavailable directions do nothing.
Full-neck mode requires a terminal width of at least 98 columns. The fret-number
tab view uses one vertically aligned number per string, including `0` for open
strings and `X` for muted strings; selecting it exits full-neck mode.

The title and clearly labeled key controls share the first centered banner,
with the chord dial directly beneath it. Compact mode places the chord diagram
beside the note waveform at the same height and centers both output blocks.
Full-neck mode stacks those lower sections so the complete neck retains its
required width. The layout redraws when the terminal is resized and shows a
width notice below 80 columns instead of wrapping controls or fretboard cells.
Above 100 columns, the interface keeps its readable width instead of stretching.

Alt-Tab renders in a full-window terminal buffer and clips output to the current
terminal height. The scroll wheel moves the viewport without changing the
selected chord or voicing.

The stationary Braille waveform at the bottom plots 25 milliseconds of an
ideal equal-amplitude signal combining every sounding string in the selected
voicing. Braille subcells increase detail without increasing its character
width. A caption below the plot identifies its normalized amplitude over time.
The semitone-spaced pitch scale has a vertical marker for every sounding note,
frequency endpoints, its own caption, and a separated note legend; press `w` to
hide or restore the waveform section.

### Themes

Press `t` to cycle through four palettes. Each palette automatically selects
colors suited to the terminal's light or dark background.

| Theme | Character |
| --- | --- |
| `Synthwave` | Vivid magenta and violet; the default |
| `Tidal` | Cool cyan and ocean blue |
| `Ember` | Warm orange and amber |
| `Evergreen` | Calm green and mint |

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
classes, the C catalog API, cgo conversion, TUI navigation, compact and full-neck
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

Each chord has two conventional voicings. Finger numbers use `1` for index,
`2` for middle, `3` for ring, and `4` for little finger. `O` marks an open
string and `X` marks a muted string.

## Architecture

- `cmd/alt-tab` contains the application entry point.
- `internal/tui` owns Bubble Tea state, navigation, layout, styling, and
  fretboard rendering.
- `internal/chords` owns the Go chord model, cgo adapter, and immutable C chord
  table behind one catalog interface.
- `tests` contains the native chord regression suite.
- `.github/workflows` runs builds and quality checks for pushes and pull
  requests.

## Project Layout

```text
.
├── .github/workflows/
├── cmd/alt-tab/
├── internal/chords/
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
