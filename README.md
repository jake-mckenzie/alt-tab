# Alt-Tab

[![CI](https://github.com/jake-mckenzie/alt-tab/actions/workflows/ci.yml/badge.svg)](https://github.com/jake-mckenzie/alt-tab/actions/workflows/ci.yml)

Alt-Tab is an interactive guitar-chord viewer for the terminal. Its responsive
Bubble Tea interface displays explicit open chords and generated movable
voicings from a compact, pure-Go catalog.

Diagrams follow standard tablature orientation: high e is at the top, low E is
at the bottom, and fret numbers increase from left to right. Compact mode shows
the relevant neck position; full-neck mode shows frets 1–27.

## TUI Preview

The default layout appears as follows in an 80-column terminal. Terminal colors
are omitted here; Alt-Tab adapts them to light and dark backgrounds.

<!-- BEGIN VERIFIED TUI OUTPUT -->
```text
  ╭──────────────────────────────────────────────────────────────────────────╮
  │                                 ALT-TAB                                  │
  │                     Guitar Chord Viewer · Synthwave                      │
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
  │       CHORD DIAGRAM        │  │       WAVEFORM · AMPLITUDE / TIME        │
  │          COMPACT           │  │                                          │
  │                            │  │  +1.0 |                               |  │
  │             A              │  │       | ⣷                     ⡼⡀      |  │
  │       VOICING 1 OF 3       │  │       |⢰⠉⡆                    ⡇⡇      |  │
  │                            │  │       |⢸ ⡇                   ⢸ ⢇      |  │
  │                            │  │       |⢸ ⢱                   ⢸ ⢸      |  │
  │                            │  │       |⡇ ⢸         ⡤⡀⢀⠎⡆ ⡰⡄  ⢸ ⠸⡀ ⡀   |  │
  │                            │  │   0.0 +⠧⠤⠬⡦⡴⠵⡤⠤⡴⢦⠤⡼⠤⠵⠮⠤⠼⠴⠥⢵⠤⠤⡧⠤⠤⡧⡴⠽⡤⠤⡼+  │
  │                            │  │       |   ⠱⠃ ⢇⢰⠁ ⠓⠁       ⢸  ⡇  ⠱⠃ ⢣⢠⠃|  │
  │                            │  │       |      ⠈⠁           ⠈⡆ ⡇     ⠈⠊ |  │
  │       1    2    3    4     │  │       |                    ⡇⢰⠁        |  │
  │     ────────────────────   │  │       |                    ⢱⢸         |  │
  │ e  O--------------------|  │  │       |                    ⠸⡜         |  │
  │ B  |-------3------------|  │  │  -1.0 |                     ⠁         |  │
  │ G  |-------2------------|  │  │        0 ms                       25 ms  │
  │ D  |-------1------------|  │  ╰──────────────────────────────────────────╯
  │ A  O--------------------|  │
  │ E  X--------------------|  │  ╭──────────────────────────────────────────╮
  │                            │  │            FREQUENCY SPECTRUM            │
  │                            │  │                                          │
  │                            │  │  1.0 |█          █       █     █    █|   │
  │                            │  │      |█          █       █     █    █|   │
  │                            │  │  0.5 |█          █       █     █    █|   │
  │                            │  │      |█          █       █     █    █|   │
  │                            │  │      |█          █       █     █    █|   │
  │                            │  │  0.0 +┼──────────┼───────┼─────┼────┼+   │
  │ Fingers: 1 index  2 middle │  │        A2        E3      A3    C#4  E4   │
  │          3 ring   4 little │  │        110 Hz                   330 Hz   │
  │                            │  │                                          │
  │ Symbols: O open  X muted   │  │        Notes: E4  C#4  A3  E3  A2        │
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

Before building or running, Make automatically checks for Go, a resolvable
dependency graph, and valid dependency checksums.

Alt-Tab is an interactive TUI and must be run in a terminal. It does not accept
the command-line chord and display flags used by the former interface.

## Requirements

- macOS, Linux, or another Unix-like environment
- Go 1.25 or newer
- Make

Confirm the required tools are available:

```bash
go version
make --version
```

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
tab view preserves the compact neck but places each fret number directly on its
string instead of showing the separate fret scale or finger numbers. `O` and
`X` still mark open and muted strings; selecting tab view exits full-neck mode.

The title and clearly labeled key controls share the first centered banner,
with the chord dial directly beneath it. Inspired by htop++, compact mode uses
a primary chord panel on the left and stacked waveform and spectrum modules on
the right. Full-neck mode stacks all three modules at the same full width.
The layout redraws when the terminal is resized and shows a
width notice below 80 columns instead of wrapping controls or fretboard cells.
Above 100 columns, the interface keeps its readable width instead of stretching.

Alt-Tab renders in a full-window terminal buffer and clips output to the current
terminal height. The scroll wheel moves the viewport without changing the
selected chord or voicing.

The stationary Braille waveform plots 25 milliseconds of an ideal
equal-amplitude signal combining every sounding string in the selected voicing.
Braille subcells increase detail without increasing its character width. Its
section heading also identifies the amplitude-over-time axes without a redundant
graph title. The separate frequency spectrum shows a normalized peak for each
sounding note, with pitch names, frequency endpoints, and a centered note
legend. Press `w` to hide or restore the time-domain waveform; the spectrum
remains visible.

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

Run all tests:

```bash
make test
```

Run individual quality checks:

```bash
make fmt-check
make test-go
make race
make vet
```

The optional race check, and therefore `make check`, requires a platform and
toolchain supported by Go's race detector. Normal builds and `make test` do not
require cgo or a C compiler.

Run every formatting, test, race, and vet check used by CI:

```bash
make check
```

The suite verifies every generated fret and finger assignment, chord pitch
class, waveform frequency and note label, spectrum legend, catalog lookup, TUI
navigation, compact and full-neck diagrams, help behavior, and the README output
snapshot.

## Cleaning

Remove the application binary and generated build artifacts:

```bash
make clean
```

## Build Troubleshooting

If dependency checksums or downloads fail, refresh the module cache and retry:

```bash
go mod download
make
```

## Supported Chords

| Major | Minor | Accidental |
| --- | --- | --- |
| `A`, `B`, `C`, `D`, `E`, `F`, `G` | `Am`, `Cm`, `Dm`, `Em`, `Gm` | `Bb`, `F#` |

The catalog contains 14 chord names and three conventional voicings per chord,
for 42 voicings total. Open positions are stored explicitly, while movable E-,
A-, and D-shape chords are generated at the required fret. Finger numbers use
`1` for index, `2` for middle, `3` for ring, and `4` for little finger. `O`
marks an open string and `X` marks a muted string.

## Architecture

- `cmd/alt-tab` contains the application entry point.
- `internal/tui` owns Bubble Tea state, navigation, layout, styling, and
  fretboard rendering.
- `internal/chords` owns the pure-Go chord model, compact shape definitions,
  generated voicings, and immutable catalog interface.
- `.github/workflows` runs builds and quality checks for pushes and pull
  requests.

## Project Layout

```text
.
├── .github/workflows/
├── cmd/alt-tab/
├── internal/chords/
├── internal/tui/
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
