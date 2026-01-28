# alt-tab

**Console Guitar Chord Viewer written in C**

Alt-Tab is a terminal-based application that renders **guitar chords as
ASCII fretboards**.\
Given a chord name (e.g. `C`, `G`, `Dm7`), the program looks it up in a
chord library and prints a visual representation of the fret positions.

Written in **pure C (C11)** and focuses on:
  - data modeling with structs
  - pointer correctness
  - clean modular design
  - professional build tooling (Make, sanitizers)

------------------------------------------------------------------------

## Example Output

    Alt-Tab — Guitar Chord Viewer

    Enter chord name (or 'q' to quit): C
    |--@---@---@---------------|
    |--@---@---@---------------|
    |--------------------------|
    |--------------------------|
    |--------------------------|
    |--------------------------|

------------------------------------------------------------------------

## Project Structure

    .
    ├── include/
    │   ├── matrix.h
    │   └── chords.h
    │
    ├── src/
    │   ├── main.c
    │   ├── matrix.c
    │   └── chords.c
    │
    ├── build/
    ├── Makefile

------------------------------------------------------------------------

## Building

### Debug build (default)

``` bash
make
```

or:

``` bash
make BUILD=debug
```

### Release build

``` bash
make BUILD=release
```

### Clean

``` bash
make clean
```

## Running

``` bash
./alt-tab
```

Type `q` to quit.

------------------------------------------------------------------------

## Supported Chords

Defined in `chords.c` (examples: `C`, `G`, `D`, `Am`, `Em`).

------------------------------------------------------------------------

## Development Notes

-   C11
-   macOS / Linux
-   Sanitizers enabled in debug builds

------------------------------------------------------------------------

## Roadmap

-   Multiple chord shapes
-   Barre chords
-   Fret numbers
-   Piano mode
-   Tests

------------------------------------------------------------------------
