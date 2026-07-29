#include "theory/chord_library.h"

#include <ctype.h>

#define ARRAY_COUNT( array ) ( sizeof( array ) / sizeof( ( array )[ 0 ] ) )

/* Built-in voicings use conventional fretting-hand finger assignments. */
static const Chord DEFAULT_CHORDS[ ] = {
    { "A",  {
        { 0, 0 }, { 2, 3 }, { 2, 2 },
        { 2, 1 }, { 0, 0 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "Am", {
        { 0, 0 }, { 1, 1 }, { 2, 3 },
        { 2, 2 }, { 0, 0 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "B",  {
        { 2, 1 }, { 4, 4 }, { 4, 3 },
        { 4, 2 }, { 2, 1 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "Bb", {
        { 1, 1 }, { 3, 4 }, { 3, 3 },
        { 3, 2 }, { 1, 1 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "C",  {
        { 0, 0 }, { 1, 1 }, { 0, 0 },
        { 2, 2 }, { 3, 3 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "Cm", {
        { 3, 1 }, { 4, 2 }, { 5, 4 },
        { 5, 3 }, { 3, 1 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "D",  {
        { 2, 2 }, { 3, 3 }, { 2, 1 },
        { 0, 0 }, { GUITAR_MUTED_FRET, 0 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "Dm", {
        { 1, 1 }, { 3, 3 }, { 2, 2 },
        { 0, 0 }, { GUITAR_MUTED_FRET, 0 }, { GUITAR_MUTED_FRET, 0 }
    } },
    { "E",  {
        { 0, 0 }, { 0, 0 }, { 1, 1 },
        { 2, 3 }, { 2, 2 }, { 0, 0 }
    } },
    { "Em", {
        { 0, 0 }, { 0, 0 }, { 0, 0 },
        { 2, 3 }, { 2, 2 }, { 0, 0 }
    } },
    { "F",  {
        { 1, 1 }, { 1, 1 }, { 2, 2 },
        { 3, 4 }, { 3, 3 }, { 1, 1 }
    } },
    { "F#", {
        { 2, 1 }, { 2, 1 }, { 3, 2 },
        { 4, 4 }, { 4, 3 }, { 2, 1 }
    } },
    { "G",  {
        { 3, 3 }, { 0, 0 }, { 0, 0 },
        { 0, 0 }, { 2, 1 }, { 3, 2 }
    } },
    { "Gm", {
        { 3, 1 }, { 3, 1 }, { 3, 1 },
        { 5, 4 }, { 5, 3 }, { 3, 1 }
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
    if ( library == NULL || name == NULL ) {
        return NULL;
    }

    for ( size_t index = 0; index < library->count; index++ ) {
        if ( namesEqual( library->items[ index ].name, name ) ) {
            return &library->items[ index ];
        }
    }

    return NULL;
}
