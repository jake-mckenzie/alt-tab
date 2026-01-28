#include "chords.h"
#include <string.h>

/* ---- Internal helpers (not in header) ---- */

static int chordNameEquals(const char *a, const char *b)
{
    return strcmp(a, b) == 0;
}

/* ---- Public functions ---- */

int addChord(ChordLibrary *lib, Chord chord)
{
    if (lib->count >= MAX_CHORDS)
        return 0;

    lib->chords[lib->count++] = chord;
    return 1;
}

const Chord *findChord(const ChordLibrary *lib, const char *name)
{
    for (size_t i = 0; i < lib->count; i++) {
        if (chordNameEquals(lib->chords[i].name, name))
            return &lib->chords[i];
    }
    return NULL;
}

void drawChord(Matrix *board, const Chord *chord)
{
    clearFretBoard(board);

    for (size_t i = 0; i < chord->fret_count; i++) {
        int r = chord->frets[i].row;
        int c = chord->frets[i].col;

        if (r >= 0 && r < ROWS && c >= 0 && c < COLS)
            board->matrix[r][c] = chord->frets[i].value;
    }
}

/* ---- Library initialization ---- */

void initChordLibrary(ChordLibrary *lib)
{
    lib->count = 0;

    addChord(lib, (Chord){
        .name = "C",
        .frets = {
            {0, 0, 'E'}, {1, 1, 'C'}, {2, 0, 'G'},
            {3, 2, 'E'}, {4, 3, 'C'}
        },
        .fret_count = 5
    });

    addChord(lib, (Chord){
        .name = "G",
        .frets = {
            {0, 3, 'G'}, {1, 0, 'B'}, {2, 0, 'G'},
            {3, 0, 'D'}, {4, 2, 'B'}, {5, 3, 'G'}
        },
        .fret_count = 6
    });
}
