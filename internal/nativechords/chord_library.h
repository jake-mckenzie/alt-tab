#ifndef ALT_TAB_NATIVE_CHORD_LIBRARY_H
#define ALT_TAB_NATIVE_CHORD_LIBRARY_H

#include <stddef.h>
#include "chord.h"

/* Provides a read-only view of a chord collection. */
typedef struct {
    const Chord *items;
    size_t count;
} ChordLibrary;

/* Returns the built-in set of supported chord voicings. */
ChordLibrary chordLibraryDefault( void );

/* Finds a chord by name without regard to letter case. */
const Chord *chordLibraryFind( const ChordLibrary *library, const char *name );

/* Finds one numbered variation of a chord name. */
const Chord *chordLibraryFindVariation(
    const ChordLibrary *library,
    const char *name,
    int variation
);

/* Counts the available variations for one chord name. */
size_t chordLibraryVariationCount(
    const ChordLibrary *library,
    const char *name
);

#endif
