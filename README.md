# Alt-Tab

Alt-Tab is a C11 terminal application that renders guitar chord voicings as
readable horizontal tablature. The default diagram shows frets 0–4 or extends
to the highest required fret, while full-neck mode displays frets 0–27.

Chord lookup is case-insensitive, so inputs such as `Am`, `am`, and `AM`
are equivalent. Diagrams follow standard tab orientation: high e is on top,
low E is on the bottom, and fret numbers increase from left to right.

## Requirements

- A C11 compiler such as GCC or Clang
- Make
- AddressSanitizer and UndefinedBehaviorSanitizer support for the default
  debug build

## Example Output

```text
+----------------------------------------------------------+
|  ALT-TAB                                                 |
|  Guitar Chord Viewer                                     |
+----------------------------------------------------------+
  Mode    : compact horizontal tab
  Chords  : A  Am  B  Bb  C  Cm  D  Dm  E  Em  F  F#  G  Gm
  Commands: CHORD, CHORD -f, ?, or q
  Example : ./alt-tab A -f

chord> C

Chord: C

    Fret numbers ->
      0    1    2    3    4
e  |--O----------------------|
B  |-------1-----------------|
G  |--O----------------------|
D  |------------2------------|
A  |-----------------3-------|
E  |--X----------------------|

    Fingers: 1 index  2 middle  3 ring  4 little
    Symbols: O open   X muted
```

Finger numbers show one conventional fretting-hand placement for the stored
voicing. `O` is an open string and `X` is a muted string.

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

Pass options through the Makefile with `ARGS`, or use the dedicated full-neck
target:

```bash
make run ARGS=--full-neck
make run ARGS="A -f"
make run-full
```

Do not use `make run -f`: `-f` is an option consumed by Make itself.

Start the interactive prompt:

```bash
./alt-tab
```

Print one chord and exit:

```bash
./alt-tab A
./alt-tab A -f
```

Start interactively with the complete 0–27 fret neck:

```bash
./alt-tab --full-neck
./alt-tab -f
```

The display mode can also be changed while the program is running:

```text
chord> A -f
chord> --full-neck
chord> --compact
```

`A -f` draws A on the full neck; a flag entered by itself changes the mode for
subsequent chords. Options and the chord name can appear in either order.
Use `?` or `help` at the `chord>` prompt for supported chords and commands.
Display command-line help without starting the prompt:

```bash
./alt-tab --help
```

Enter `q` or `Q` to quit. Unknown names suggest opening help and return to the
prompt. Unknown command-line options print usage and exit with an error.

## Supported Chords

The immutable default library in `src/theory/chord_library.c` contains:

| Major | Minor | Accidental |
| --- | --- | --- |
| `A`, `B`, `C`, `D`, `E`, `F`, `G` | `Am`, `Cm`, `Dm`, `Em`, `Gm` | `Bb`, `F#` |

## Architecture

- `theory` defines fret and finger placement for six-string chord voicings and
  owns the immutable chord library. It performs no terminal input or output.
- `render` validates a chord and provides compact and full-neck horizontal tab.
- `app` owns the input loop and coordinates lookup and rendering.

String placements are ordered from high e to low E. A fret value of `-1` means
muted, `0` means open, and a positive value is the played fret. Fretted strings
also store a finger number from `1` through `4`. Future audio modules can
consume this model without depending on terminal rendering.

## Tests

```bash
make test
```

The tests verify all 14 stored fret and finger assignments, case-insensitive
lookup, invalid lookup arguments, both renderer modes, chord-plus-flag commands,
interactive display switching, and help output. Debug tests run with both
configured sanitizers.

## Roadmap

- Multiple voicings for one chord name
- Explicit barre visualization
- Audio input and chord detection
- Audio output and chord playback
- Piano mode
