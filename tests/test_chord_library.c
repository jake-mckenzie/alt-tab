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

static const Chord EXPECTED_CHORDS[ ] = {
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

static int testDefaultLibrary( void )
{
    ChordLibrary library = chordLibraryDefault( );

    CHECK( library.count == ARRAY_COUNT( EXPECTED_CHORDS ) );

    for ( size_t chord_index = 0;
          chord_index < ARRAY_COUNT( EXPECTED_CHORDS );
          chord_index++ ) {
        const Chord *actual = &library.items[ chord_index ];
        const Chord *expected = &EXPECTED_CHORDS[ chord_index ];

        CHECK( strcmp( actual->name, expected->name ) == 0 );
        for ( size_t string = 0; string < GUITAR_STRING_COUNT; string++ ) {
            CHECK( actual->frets[ string ] == expected->frets[ string ] );
        }
    }

    return 1;
}

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

static int testRenderer( void )
{
    ChordLibrary library = chordLibraryDefault( );
    const Chord *chord = chordLibraryFind( &library, "C" );
    FILE *stream = tmpfile( );
    char output[ 1024 ];
    size_t bytes_read;

    CHECK( chord != NULL );
    CHECK( stream != NULL );
    CHECK( terminalRendererPrint( stream, chord, 6 ) );
    CHECK( !terminalRendererPrint( stream, chord, 0 ) );

    rewind( stream );
    bytes_read = fread( output, 1, sizeof( output ) - 1, stream );
    output[ bytes_read ] = '\0';
    fclose( stream );

    CHECK( strstr( output, "\nC\n" ) != NULL );
    CHECK( strstr( output, "e  |   O" ) != NULL );
    CHECK( strstr( output, "E  |   X" ) != NULL );

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
