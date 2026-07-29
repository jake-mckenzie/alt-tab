#include <stdio.h>
#include <string.h>

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
    int frets[ GUITAR_STRING_COUNT ];
    int fingers[ GUITAR_STRING_COUNT ];
} ExpectedChord;

/* Expected frets and one conventional fingering for every built-in chord. */
static const ExpectedChord EXPECTED_CHORDS[ ] = {
    { "A",  { 0, 2, 2, 2, 0, -1 }, { 0, 3, 2, 1, 0, 0 } },
    { "Am", { 0, 1, 2, 2, 0, -1 }, { 0, 1, 3, 2, 0, 0 } },
    { "B",  { 2, 4, 4, 4, 2, -1 }, { 1, 4, 3, 2, 1, 0 } },
    { "Bb", { 1, 3, 3, 3, 1, -1 }, { 1, 4, 3, 2, 1, 0 } },
    { "C",  { 0, 1, 0, 2, 3, -1 }, { 0, 1, 0, 2, 3, 0 } },
    { "Cm", { 3, 4, 5, 5, 3, -1 }, { 1, 2, 4, 3, 1, 0 } },
    { "D",  { 2, 3, 2, 0, -1, -1 }, { 2, 3, 1, 0, 0, 0 } },
    { "Dm", { 1, 3, 2, 0, -1, -1 }, { 1, 3, 2, 0, 0, 0 } },
    { "E",  { 0, 0, 1, 2, 2, 0 }, { 0, 0, 1, 3, 2, 0 } },
    { "Em", { 0, 0, 0, 2, 2, 0 }, { 0, 0, 0, 3, 2, 0 } },
    { "F",  { 1, 1, 2, 3, 3, 1 }, { 1, 1, 2, 4, 3, 1 } },
    { "F#", { 2, 2, 3, 4, 4, 2 }, { 1, 1, 2, 4, 3, 1 } },
    { "G",  { 3, 0, 0, 0, 2, 3 }, { 3, 0, 0, 0, 1, 2 } },
    { "Gm", { 3, 3, 3, 5, 5, 3 }, { 1, 1, 1, 4, 3, 1 } }
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
        6,
        TERMINAL_RENDER_COMPACT_TAB
    ) );
    CHECK( terminalRendererPrint(
        stream,
        chord,
        6,
        TERMINAL_RENDER_COMPACT_TAB
    ) );
    CHECK( terminalRendererPrint(
        stream,
        chord,
        6,
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
        6,
        ( TerminalRenderMode )99
    ) );

    rewind( stream );
    bytes_read = fread( output, 1, sizeof( output ) - 1, stream );
    output[ bytes_read ] = '\0';
    fclose( stream );

    CHECK( strstr( output, "Chord: C" ) != NULL );
    CHECK( strstr( output, "Fret numbers ->" ) != NULL );
    CHECK( strstr( output, "Fingers: 1 index" ) != NULL );
    CHECK( strstr( output, "e  |--O--" ) != NULL );
    CHECK( strstr( output, "B  |-------1-" ) != NULL );
    CHECK( strstr( output, "E  |--X--" ) != NULL );
    CHECK( strstr( output, " 5 " ) != NULL );

    return 1;
}

int main( void )
{
    if ( !testDefaultLibrary( ) || !testLookup( ) || !testRenderer( ) ) {
        return 1;
    }

    puts( "all tests passed" );
    return 0;
}
