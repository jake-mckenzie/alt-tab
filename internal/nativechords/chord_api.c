#include "chord_api.h"
#include <string.h>
#include "chord_library.h"

_Static_assert(
    ALT_TAB_STRING_COUNT == GUITAR_STRING_COUNT,
    "public and internal string counts must match"
);

/* Detects the first occurrence of a chord name in the ordered library. */
static int isFirstNamedChord( const ChordLibrary *library, size_t index )
{
    const char *name = library->items[ index ].name;

    for ( size_t previous = 0; previous < index; previous++ ) {
        if ( strcmp( library->items[ previous ].name, name ) == 0 ) {
            return 0;
        }
    }

    return 1;
}

/* Counts distinct chord names without exposing duplicate variations. */
size_t altTabChordCount( void )
{
    ChordLibrary library = chordLibraryDefault( );
    size_t count = 0;

    for ( size_t index = 0; index < library.count; index++ ) {
        count += isFirstNamedChord( &library, index );
    }

    return count;
}

/* Maps a public name index to the immutable library string. */
const char *altTabChordNameAt( size_t requested_index )
{
    ChordLibrary library = chordLibraryDefault( );
    size_t name_index = 0;

    for ( size_t index = 0; index < library.count; index++ ) {
        if ( !isFirstNamedChord( &library, index ) ) {
            continue;
        }
        if ( name_index == requested_index ) {
            return library.items[ index ].name;
        }
        name_index++;
    }

    return NULL;
}

/* Delegates variation counting through the stable backend boundary. */
size_t altTabChordVariationCount( const char *name )
{
    ChordLibrary library = chordLibraryDefault( );

    return chordLibraryVariationCount( &library, name );
}

/* Copies internal chord data so callers never own a library pointer. */
int altTabChordLoad(
    const char *name,
    int variation,
    AltTabChordVoicing *output
)
{
    ChordLibrary library = chordLibraryDefault( );
    const Chord *chord;

    if ( output == NULL ) {
        return 0;
    }

    chord = chordLibraryFindVariation( &library, name, variation );
    if ( chord == NULL ) {
        return 0;
    }

    output->variation = chord->variation;
    for ( size_t string = 0; string < ALT_TAB_STRING_COUNT; string++ ) {
        output->strings[ string ].fret = chord->strings[ string ].fret;
        output->strings[ string ].finger = chord->strings[ string ].finger;
    }

    return 1;
}
