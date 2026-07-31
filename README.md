# Alt-Tab

Alt-Tab is an interactive guitar-chord viewer built with Bubble Tea. It uses a
C11 backend for chord data and a responsive Go terminal interface for browsing
chords, selecting voicings, and viewing accurate finger placement.

Diagrams follow standard tablature orientation: high e is on top, low E is on
the bottom, and fret numbers increase from left to right. Compact mode shows the
relevant neck position; full-neck mode shows frets 1–27.

## Requirements

- Go 1.25 or newer with cgo enabled
- A C11 compiler such as GCC or Clang
- Make
- AddressSanitizer and UndefinedBehaviorSanitizer support for debug backend
  tests

## Project Structure

```text
.
├── cmd/
│   └── alt-tab/
│       └── main.go
├── include/
│   ├── backend/
│   │   └── chord_api.h
│   └── theory/
│       ├── chord.h
│       └── chord_library.h
├── internal/
│   ├── chords/
│   │   ├── catalog.go
│   │   ├── native.go
│   │   ├── native_backend.c
│   │   └── native_test.go
│   └── tui/
│       ├── fretboard.go
│       ├── fretboard_test.go
│       ├── model.go
│       ├── model_test.go
│       ├── styles.go
│       └── view.go
├── src/
│   ├── backend/
│   │   └── chord_api.c
│   └── theory/
│       └── chord_library.c
├── tests/
│   └── test_chord_library.c
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Building

Build the `alt-tab` TUI:

```bash
make
```

Build a smaller release binary without debug symbols:

```bash
make BUILD=release
```

Remove generated binaries and test objects:

```bash
make clean
```

## Running

```bash
make run
```

You can also launch the built binary directly:

```bash
./alt-tab
```

### Controls

| Key | Action |
| --- | --- |
| `↑` / `k` | Previous chord |
| `↓` / `j` | Next chord |
| `←` / `h` | Previous variation |
| `→` / `l` | Next variation |
| `f` | Toggle compact or full-neck view |
| `?` | Open or close help |
| `Esc` | Close help |
| `q` / `Ctrl+C` | Quit |

## Supported Chords

| Major | Minor | Accidental |
| --- | --- | --- |
| `A`, `B`, `C`, `D`, `E`, `F`, `G` | `Am`, `Cm`, `Dm`, `Em`, `Gm` | `Bb`, `F#` |

Each chord currently has two conventional variations. Finger numbers use
`1` for index, `2` for middle, `3` for ring, and `4` for little finger. `O`
marks an open string and `X` marks a muted string.

## Architecture

- `theory` owns the immutable chord library and has no UI dependencies.
- `backend` exposes a stable C API for listing chords and copying voicings into
  caller-owned memory.
- `internal/chords` converts C data into Go-owned structs behind a catalog
  interface.
- `internal/tui` owns Bubble Tea state, navigation, responsive layout, styling,
  and fretboard rendering.
- `cmd/alt-tab` is the only application entry point.

This separation keeps terminal concerns out of the chord model and allows
future audio input or output packages to consume the same Go chord types.

## Tests

```bash
make test
```

The suite verifies every stored fret and finger assignment, chord pitch
classes, the public C API, cgo data conversion, navigation state, compact and
full-neck diagrams, and help behavior. Debug backend tests run with
AddressSanitizer and UndefinedBehaviorSanitizer.

## Roadmap

- Additional chord voicings
- Explicit barre visualization
- Audio input and chord detection
- Audio output and chord playback
- Piano mode
