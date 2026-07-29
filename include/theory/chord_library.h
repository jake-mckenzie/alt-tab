#ifndef ALT_TAB_THEORY_CHORD_LIBRARY_H
#define ALT_TAB_THEORY_CHORD_LIBRARY_H

#include <stddef.h>

#include "theory/chord.h"

typedef struct {
    const Chord *items;
    size_t count;
} ChordLibrary;

ChordLibrary chordLibraryDefault( void );
const Chord *chordLibraryFind( const ChordLibrary *library, const char *name );

#endif
