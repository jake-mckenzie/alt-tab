#ifndef CHORDS_H
#define CHORDS_H

#include <stddef.h>
#include "matrix.h"

#define MAX_FRETS 9
#define MAX_CHORDS 64

typedef struct {
    int row;
    int col;
    char value;
} Fret;

typedef struct {
    const char *name;
    Fret frets[MAX_FRETS];
    size_t fret_count;
} Chord;

typedef struct {
    Chord chords[MAX_CHORDS];
    size_t count;
} ChordLibrary;

/* Public API */
void initChordLibrary(ChordLibrary *lib);
int  addChord(ChordLibrary *lib, Chord chord);
const Chord *findChord(const ChordLibrary *lib, const char *name);
void drawChord(Matrix *board, const Chord *chord);

#endif
