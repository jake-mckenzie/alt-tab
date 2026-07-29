# Alt-Tab

Alt-Tab is a C11 terminal application that renders guitar chord voicings as
ASCII fretboards. Enter a supported chord name, and the program displays its
fingering across frets 0–27.

Chord lookup is case-insensitive, so inputs such as `Am`, `am`, and `AM`
are equivalent.

## Requirements

- A C11 compiler such as GCC or Clang
- Make
- AddressSanitizer and UndefinedBehaviorSanitizer support for the default
  debug build

## Example Output

The program prints frets 0–27; this example is shortened after fret 5.

```text
    Alt-Tab — Guitar Chord Viewer

    Enter chord name (or 'q' to quit): C
    C
           0  1  2  3  4  5
    e  |   O  -  -  -  -  -|
    B  |   -  *  -  -  -  -|
    G  |   O  -  -  -  -  -|
    D  |   -  -  *  -  -  -|
    A  |   -  -  -  *  -  -|
    E  |   X  -  -  -  -  -|
```

`O` is an open string, `*` is a fretted string, and `X` is a muted string.

## Project Structure

```text
.
├── include/
│   ├── render/
│   │   └── terminal_renderer.h
│   └── theory/
│       ├── chord.h
│       └── chord_library.h
├── src/
│   ├── app/
│   │   └── main.c
│   ├── render/
│   │   └── terminal_renderer.c
│   └── theory/
│       └── chord_library.c
├── tests/
│   └── test_chord_library.c
├── Makefile
└── README.md
```

## Building

The default debug build enables AddressSanitizer and
UndefinedBehaviorSanitizer:

```bash
make
```

Build an optimized release binary:

```bash
make clean
make BUILD=release
```

Remove generated objects and binaries:

```bash
make clean
```

## Running

Build and run:

```bash
make run
```

Or run an existing binary directly:

```bash
./alt-tab
```

Enter `q` or `Q` to quit. Unknown names produce a not-found message and return
to the prompt.

## Supported Chords

The immutable default library in `src/theory/chord_library.c` contains:

| Major | Minor | Accidental |
| --- | --- | --- |
| `A`, `B`, `C`, `D`, `E`, `F`, `G` | `Am`, `Cm`, `Dm`, `Em`, `Gm` | `Bb`, `F#` |

## Architecture

- `theory` defines six-string chord voicings and the immutable chord library.
  It performs no terminal input or output.
- `render` validates a chord and converts it into terminal output.
- `app` owns the input loop and coordinates lookup and rendering.

Fret arrays are ordered from the high E string to the low E string. A value of
`-1` means muted, `0` means open, and a positive value is the played fret.
Future audio modules can consume this model without depending on terminal
rendering.

## Tests

```bash
make test
```

The tests verify all 14 stored voicings, case-insensitive lookup, invalid lookup
arguments, and terminal-renderer output. Debug tests run with both configured
sanitizers.

## Roadmap

- Multiple voicings for one chord name
- Explicit barre visualization
- Audio input and chord detection
- Audio output and chord playback
- Piano mode
