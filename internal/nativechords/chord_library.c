#include "chord_library.h"
#include <ctype.h>

#define ARRAY_COUNT( array ) ( sizeof( array ) / sizeof( ( array )[ 0 ] ) )

/* Built-in voicings use conventional fretting-hand finger assignments. */
static const Chord DEFAULT_CHORDS[ ] = {
    { "A", 1, {
        { 0, 0 }, { 2, 3 }, { 2, 2 },
        { 2, 1 }, { 0, 0 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "A", 2, {
        { 5, 1 }, { 5, 1 }, { 6, 2 },
        { 7, 4 }, { 7, 3 }, { 5, 1 }
    } },
    { "Am", 1, {
        { 0, 0 }, { 1, 1 }, { 2, 3 },
        { 2, 2 }, { 0, 0 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "Am", 2, {
        { 5, 1 }, { 5, 1 }, { 5, 1 },
        { 7, 4 }, { 7, 3 }, { 5, 1 }
    } },
    { "B", 1, {
        { 2, 1 }, { 4, 4 }, { 4, 3 },
        { 4, 2 }, { 2, 1 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "B", 2, {
        { 7, 1 }, { 7, 1 }, { 8, 2 },
        { 9, 4 }, { 9, 3 }, { 7, 1 }
    } },
    { "Bb", 1, {
        { 1, 1 }, { 3, 4 }, { 3, 3 },
        { 3, 2 }, { 1, 1 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "Bb", 2, {
        { 6, 1 }, { 6, 1 }, { 7, 2 },
        { 8, 4 }, { 8, 3 }, { 6, 1 }
    } },
    { "C", 1, {
        { 0, 0 }, { 1, 1 }, { 0, 0 },
        { 2, 2 }, { 3, 3 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "C", 2, {
        { 3, 1 }, { 5, 4 }, { 5, 3 },
        { 5, 2 }, { 3, 1 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "Cm", 1, {
        { 3, 1 }, { 4, 2 }, { 5, 4 },
        { 5, 3 }, { 3, 1 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "Cm", 2, {
        { 8, 1 }, { 8, 1 }, { 8, 1 },
        { 10, 4 }, { 10, 3 }, { 8, 1 }
    } },
    { "D", 1, {
        { 2, 2 }, { 3, 3 }, { 2, 1 },
        { 0, 0 }, { GUITAR_MUTED_FRET, 0 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "D", 2, {
        { 5, 1 }, { 7, 4 }, { 7, 3 },
        { 7, 2 }, { 5, 1 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "Dm", 1, {
        { 1, 1 }, { 3, 3 }, { 2, 2 },
        { 0, 0 }, { GUITAR_MUTED_FRET, 0 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "Dm", 2, {
        { 5, 1 }, { 6, 2 }, { 7, 4 },
        { 7, 3 }, { 5, 1 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "E", 1, {
        { 0, 0 }, { 0, 0 }, { 1, 1 },
        { 2, 3 }, { 2, 2 }, { 0, 0 }
    } },
    { "E", 2, {
        { 7, 1 }, { 9, 4 }, { 9, 3 },
        { 9, 2 }, { 7, 1 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "Em", 1, {
        { 0, 0 }, { 0, 0 }, { 0, 0 },
        { 2, 3 }, { 2, 2 }, { 0, 0 }
    } },
    { "Em", 2, {
        { 7, 1 }, { 8, 2 }, { 9, 4 },
        { 9, 3 }, { 7, 1 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "F", 1, {
        { 1, 1 }, { 1, 1 }, { 2, 2 },
        { 3, 4 }, { 3, 3 }, { 1, 1 }
    } },
    { "F", 2, {
        { 8, 1 }, { 10, 4 }, { 10, 3 },
        { 10, 2 }, { 8, 1 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "F#", 1, {
        { 2, 1 }, { 2, 1 }, { 3, 2 },
        { 4, 4 }, { 4, 3 }, { 2, 1 }
    } },
    { "F#", 2, {
        { 9, 1 }, { 11, 4 }, { 11, 3 },
        { 11, 2 }, { 9, 1 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "G", 1, {
        { 3, 3 }, { 0, 0 }, { 0, 0 },
        { 0, 0 }, { 2, 1 }, { 3, 2 }
    } },
    { "G", 2, {
        { 3, 1 }, { 3, 1 }, { 4, 2 },
        { 5, 4 }, { 5, 3 }, { 3, 1 }
    } },
    { "Gm", 1, {
        { 3, 1 }, { 3, 1 }, { 3, 1 },
        { 5, 4 }, { 5, 3 }, { 3, 1 }
    } },
    { "Gm", 2, {
        { 10, 1 }, { 11, 2 }, { 12, 4 },
        { 12, 3 }, { 10, 1 }, { GUITAR_MUTED_FRET, 0 }
    } }
};

/* Compares two chord names using ASCII-compatible case folding. */
static int namesEqual( const char *left, const char *right )
{
    while ( *left != '\0' && *right != '\0' ) {
        if ( tolower( ( unsigned char )*left ) !=
             tolower( ( unsigned char )*right ) ) {
            return 0;
        }

        left++;
        right++;
    }

    return *left == *right;
}

/* Exposes the static chord table without copying its entries. */
ChordLibrary chordLibraryDefault( void )
{
    ChordLibrary library = {
        .items = DEFAULT_CHORDS,
        .count = ARRAY_COUNT( DEFAULT_CHORDS )
    };

    return library;
}

/* Performs a case-insensitive linear lookup in a chord library. */
const Chord *chordLibraryFind( const ChordLibrary *library, const char *name )
{
    return chordLibraryFindVariation( library, name, 1 );
}

/* Finds a specific numbered voicing in a chord library. */
const Chord *chordLibraryFindVariation(
    const ChordLibrary *library,
    const char *name,
    int variation
)
{
    if ( library == NULL || name == NULL ) {
        return NULL;
    }

    for ( size_t index = 0; index < library->count; index++ ) {
        if ( namesEqual( library->items[ index ].name, name ) &&
             library->items[ index ].variation == variation ) {
            return &library->items[ index ];
        }
    }

    return NULL;
}

/* Counts every voicing that shares the requested chord name. */
size_t chordLibraryVariationCount(
    const ChordLibrary *library,
    const char *name
)
{
    size_t count = 0;

    if ( library == NULL || name == NULL ) {
        return 0;
    }

    for ( size_t index = 0; index < library->count; index++ ) {
        if ( namesEqual( library->items[ index ].name, name ) ) {
            count++;
        }
    }

    return count;
}
