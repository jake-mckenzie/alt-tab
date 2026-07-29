#include "theory/chord_library.h"

#include <ctype.h>

#define ARRAY_COUNT( array ) ( sizeof( array ) / sizeof( ( array )[ 0 ] ) )

static const Chord DEFAULT_CHORDS[ ] = {
    { "A",  { 0, 2, 2, 2, 0, GUITAR_MUTED_FRET } },
    { "Am", { 0, 1, 2, 2, 0, GUITAR_MUTED_FRET } },
    { "B",  { 2, 4, 4, 4, 2, GUITAR_MUTED_FRET } },
    { "Bb", { 1, 3, 3, 3, 1, GUITAR_MUTED_FRET } },
    { "C",  { 0, 1, 0, 2, 3, GUITAR_MUTED_FRET } },
    { "Cm", { 3, 4, 5, 5, 3, GUITAR_MUTED_FRET } },
    { "D",  { 2, 3, 2, 0, GUITAR_MUTED_FRET, GUITAR_MUTED_FRET } },
    { "Dm", { 1, 3, 2, 0, GUITAR_MUTED_FRET, GUITAR_MUTED_FRET } },
    { "E",  { 0, 0, 1, 2, 2, 0 } },
    { "Em", { 0, 0, 0, 2, 2, 0 } },
    { "F",  { 1, 1, 2, 3, 3, 1 } },
    { "F#", { 2, 2, 3, 4, 4, 2 } },
    { "G",  { 3, 0, 0, 0, 2, 3 } },
    { "Gm", { 3, 3, 3, 5, 5, 3 } }
};

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

ChordLibrary chordLibraryDefault( void )
{
    ChordLibrary library = {
        .items = DEFAULT_CHORDS,
        .count = ARRAY_COUNT( DEFAULT_CHORDS )
    };

    return library;
}

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
