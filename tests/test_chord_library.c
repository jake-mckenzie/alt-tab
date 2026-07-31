#include <stdio.h>
#include <string.h>
#include "backend/chord_api.h"
#include "render/terminal_renderer.h"
#include "theory/chord_library.h"

#define ARRAY_COUNT( array ) ( sizeof( array ) / sizeof( ( array )[ 0 ] ) )

#define CHECK( condition )                                                   \
    do {                                                                     \
        if ( !( condition ) ) {                                              \
            fprintf( stderr, "check failed at line %d: %s\n",                \
                     __LINE__, #condition );                                  \
            return 0;                                                        \
        }                                                                    \
    } while ( 0 )

/* Mirrors the chord fields used by the regression data below. */
typedef struct {
    const char *name;
    int variation;
    int frets[ GUITAR_STRING_COUNT ];
    int fingers[ GUITAR_STRING_COUNT ];
} ExpectedChord;

/* Expected frets and one conventional fingering for every built-in chord. */
static const ExpectedChord EXPECTED_CHORDS[ ] = {
    { "A", 1, { 0, 2, 2, 2, 0, -1 }, { 0, 3, 2, 1, 0, 0 } },
    { "A", 2, { 5, 5, 6, 7, 7, 5 }, { 1, 1, 2, 4, 3, 1 } },
    { "Am", 1, { 0, 1, 2, 2, 0, -1 }, { 0, 1, 3, 2, 0, 0 } },
    { "Am", 2, { 5, 5, 5, 7, 7, 5 }, { 1, 1, 1, 4, 3, 1 } },
    { "B", 1, { 2, 4, 4, 4, 2, -1 }, { 1, 4, 3, 2, 1, 0 } },
    { "B", 2, { 7, 7, 8, 9, 9, 7 }, { 1, 1, 2, 4, 3, 1 } },
    { "Bb", 1, { 1, 3, 3, 3, 1, -1 }, { 1, 4, 3, 2, 1, 0 } },
    { "Bb", 2, { 6, 6, 7, 8, 8, 6 }, { 1, 1, 2, 4, 3, 1 } },
    { "C", 1, { 0, 1, 0, 2, 3, -1 }, { 0, 1, 0, 2, 3, 0 } },
    { "C", 2, { 3, 5, 5, 5, 3, -1 }, { 1, 4, 3, 2, 1, 0 } },
    { "Cm", 1, { 3, 4, 5, 5, 3, -1 }, { 1, 2, 4, 3, 1, 0 } },
    { "Cm", 2, { 8, 8, 8, 10, 10, 8 }, { 1, 1, 1, 4, 3, 1 } },
    { "D", 1, { 2, 3, 2, 0, -1, -1 }, { 2, 3, 1, 0, 0, 0 } },
    { "D", 2, { 5, 7, 7, 7, 5, -1 }, { 1, 4, 3, 2, 1, 0 } },
    { "Dm", 1, { 1, 3, 2, 0, -1, -1 }, { 1, 3, 2, 0, 0, 0 } },
    { "Dm", 2, { 5, 6, 7, 7, 5, -1 }, { 1, 2, 4, 3, 1, 0 } },
    { "E", 1, { 0, 0, 1, 2, 2, 0 }, { 0, 0, 1, 3, 2, 0 } },
    { "E", 2, { 7, 9, 9, 9, 7, -1 }, { 1, 4, 3, 2, 1, 0 } },
    { "Em", 1, { 0, 0, 0, 2, 2, 0 }, { 0, 0, 0, 3, 2, 0 } },
    { "Em", 2, { 7, 8, 9, 9, 7, -1 }, { 1, 2, 4, 3, 1, 0 } },
    { "F", 1, { 1, 1, 2, 3, 3, 1 }, { 1, 1, 2, 4, 3, 1 } },
    { "F", 2, { 8, 10, 10, 10, 8, -1 }, { 1, 4, 3, 2, 1, 0 } },
    { "F#", 1, { 2, 2, 3, 4, 4, 2 }, { 1, 1, 2, 4, 3, 1 } },
    { "F#", 2, { 9, 11, 11, 11, 9, -1 }, { 1, 4, 3, 2, 1, 0 } },
    { "G", 1, { 3, 0, 0, 0, 2, 3 }, { 3, 0, 0, 0, 1, 2 } },
    { "G", 2, { 3, 3, 4, 5, 5, 3 }, { 1, 1, 2, 4, 3, 1 } },
    { "Gm", 1, { 3, 3, 3, 5, 5, 3 }, { 1, 1, 1, 4, 3, 1 } },
    { "Gm", 2, { 10, 11, 12, 12, 10, -1 }, { 1, 2, 4, 3, 1, 0 } }
};

/* Verifies every stored fret and finger assignment. */
static int testDefaultLibrary( void )
{
    ChordLibrary library = chordLibraryDefault( );

    CHECK( library.count == ARRAY_COUNT( EXPECTED_CHORDS ) );

    for ( size_t chord_index = 0;
          chord_index < ARRAY_COUNT( EXPECTED_CHORDS );
          chord_index++ ) {
        const Chord *actual = &library.items[ chord_index ];
        const ExpectedChord *expected = &EXPECTED_CHORDS[ chord_index ];

        CHECK( strcmp( actual->name, expected->name ) == 0 );
        CHECK( actual->variation == expected->variation );
        for ( size_t string = 0; string < GUITAR_STRING_COUNT; string++ ) {
            CHECK(
                actual->strings[ string ].fret == expected->frets[ string ]
            );
            CHECK(
                actual->strings[ string ].finger ==
                expected->fingers[ string ]
            );
        }
    }

    return 1;
}

/* Verifies case-insensitive lookup and invalid arguments. */
static int testLookup( void )
{
    ChordLibrary library = chordLibraryDefault( );

    CHECK( chordLibraryFind( &library, "am" ) != NULL );
    CHECK( chordLibraryFind( &library, "BB" ) != NULL );
    CHECK( chordLibraryFind( &library, "f#" ) != NULL );
    CHECK( chordLibraryFind( &library, "missing" ) == NULL );
    CHECK( chordLibraryFind( NULL, "C" ) == NULL );
    CHECK( chordLibraryFind( &library, NULL ) == NULL );
    CHECK( chordLibraryVariationCount( &library, "c" ) == 2 );
    CHECK( chordLibraryVariationCount( &library, "missing" ) == 0 );
    CHECK( chordLibraryFindVariation( &library, "C", 2 ) != NULL );
    CHECK( chordLibraryFindVariation( &library, "C", 3 ) == NULL );

    return 1;
}

/* Verifies the UI-independent backend API and its ownership boundary. */
static int testChordApi( void )
{
    static const char *const EXPECTED_NAMES[ ] = {
        "A", "Am", "B", "Bb", "C", "Cm", "D",
        "Dm", "E", "Em", "F", "F#", "G", "Gm"
    };
    AltTabChordVoicing voicing;

    CHECK( altTabChordCount( ) == ARRAY_COUNT( EXPECTED_NAMES ) );
    for ( size_t index = 0; index < ARRAY_COUNT( EXPECTED_NAMES ); index++ ) {
        CHECK(
            strcmp( altTabChordNameAt( index ), EXPECTED_NAMES[ index ] ) == 0
        );
    }
    CHECK( altTabChordNameAt( ARRAY_COUNT( EXPECTED_NAMES ) ) == NULL );
    CHECK( altTabChordVariationCount( "c" ) == 2 );
    CHECK( altTabChordVariationCount( "missing" ) == 0 );

    for ( size_t index = 0; index < ARRAY_COUNT( EXPECTED_CHORDS ); index++ ) {
        const ExpectedChord *expected = &EXPECTED_CHORDS[ index ];

        CHECK(
            altTabChordLoad(
                expected->name,
                expected->variation,
                &voicing
            )
        );
        CHECK( voicing.variation == expected->variation );
        for ( size_t string = 0; string < ALT_TAB_STRING_COUNT; string++ ) {
            CHECK( voicing.strings[ string ].fret == expected->frets[ string ] );
            CHECK(
                voicing.strings[ string ].finger ==
                expected->fingers[ string ]
            );
        }
    }

    CHECK( !altTabChordLoad( "C", 3, &voicing ) );
    CHECK( !altTabChordLoad( "missing", 1, &voicing ) );
    CHECK( !altTabChordLoad( "C", 1, NULL ) );

    return 1;
}

/* Converts the root letter and accidental into a pitch class. */
static int rootPitchClass( const char *name )
{
    int root;

    switch ( name[ 0 ] ) {
        case 'A': root = 9; break;
        case 'B': root = 11; break;
        case 'C': root = 0; break;
        case 'D': root = 2; break;
        case 'E': root = 4; break;
        case 'F': root = 5; break;
        case 'G': root = 7; break;
        default: return -1;
    }

    if ( name[ 1 ] == '#' ) {
        root = ( root + 1 ) % 12;
    } else if ( name[ 1 ] == 'b' ) {
        root = ( root + 11 ) % 12;
    }

    return root;
}

/* Verifies that every voicing contains its complete named triad. */
static int testChordTones( void )
{
    static const int OPEN_STRING_PITCHES[ GUITAR_STRING_COUNT ] = {
        4, 11, 7, 2, 9, 4
    };
    ChordLibrary library = chordLibraryDefault( );

    for ( size_t index = 0; index < library.count; index++ ) {
        const Chord *chord = &library.items[ index ];
        int root = rootPitchClass( chord->name );
        int minor = strchr( chord->name, 'm' ) != NULL;
        int third = ( root + ( minor ? 3 : 4 ) ) % 12;
        int fifth = ( root + 7 ) % 12;
        int found_root = 0;
        int found_third = 0;
        int found_fifth = 0;

        CHECK( root >= 0 );

        for ( size_t string = 0; string < GUITAR_STRING_COUNT; string++ ) {
            int fret = chord->strings[ string ].fret;
            int pitch;

            if ( fret == GUITAR_MUTED_FRET ) {
                continue;
            }

            pitch = ( OPEN_STRING_PITCHES[ string ] + fret ) % 12;
            CHECK( pitch == root || pitch == third || pitch == fifth );
            found_root |= pitch == root;
            found_third |= pitch == third;
            found_fifth |= pitch == fifth;
        }

        CHECK( found_root && found_third && found_fifth );
    }

    return 1;
}

/* Verifies both renderer modes and their instructional output. */
static int testRenderer( void )
{
    ChordLibrary library = chordLibraryDefault( );
    const Chord *chord = chordLibraryFind( &library, "C" );
    Chord invalid_chord;
    FILE *stream = tmpfile( );
    char output[ 1024 ];
    size_t bytes_read;

    CHECK( chord != NULL );
    CHECK( stream != NULL );
    invalid_chord = *chord;
    invalid_chord.strings[ GUITAR_STRING_B ].finger = GUITAR_NO_FINGER;

    CHECK( !terminalRendererPrint(
        stream,
        &invalid_chord,
        5,
        TERMINAL_RENDER_COMPACT_TAB
    ) );
    CHECK( terminalRendererPrint(
        stream,
        chord,
        5,
        TERMINAL_RENDER_COMPACT_TAB
    ) );
    CHECK( terminalRendererPrint(
        stream,
        chord,
        5,
        TERMINAL_RENDER_FULL_NECK
    ) );
    CHECK( !terminalRendererPrint(
        stream,
        chord,
        0,
        TERMINAL_RENDER_COMPACT_TAB
    ) );
    CHECK( !terminalRendererPrint(
        stream,
        chord,
        5,
        ( TerminalRenderMode )99
    ) );

    rewind( stream );
    bytes_read = fread( output, 1, sizeof( output ) - 1, stream );
    output[ bytes_read ] = '\0';
    fclose( stream );

    CHECK( strstr( output, "Chord: C (variation 1)" ) != NULL );
    CHECK( strstr( output, "Fret numbers ->" ) != NULL );
    CHECK( strstr( output, "Fingers: 1 index" ) != NULL );
    CHECK( strstr( output, "e O|" ) != NULL );
    CHECK( strstr( output, "B  |--1--" ) != NULL );
    CHECK( strstr( output, "E X|" ) != NULL );
    CHECK( strstr( output, " 0 " ) == NULL );
    CHECK( strstr( output, " 5 " ) != NULL );

    return 1;
}

int main( void )
{
    if ( !testDefaultLibrary( ) ||
         !testLookup( ) ||
         !testChordApi( ) ||
         !testChordTones( ) ||
         !testRenderer( ) ) {
        return 1;
    }

    puts( "all tests passed" );
    return 0;
}
