#define _CRT_SECURE_NO_WARNINGS

#include <stdio.h>
#include <string.h>

#include "render/terminal_renderer.h"
#include "theory/chord_library.h"

static int shouldQuit( const char *input )
{
    return strcmp( input, "q" ) == 0 || strcmp( input, "Q" ) == 0;
}

int main( void )
{
    ChordLibrary library = chordLibraryDefault( );

    puts( "Alt-Tab — Guitar Chord Viewer\n" );

    for ( ;; ) {
        char input[ 16 ];
        const Chord *chord;

        printf( "Enter chord name (or 'q' to quit): " );
        if ( fgets( input, sizeof( input ), stdin ) == NULL ) {
            break;
        }

        input[ strcspn( input, "\r\n" ) ] = '\0';
        if ( shouldQuit( input ) ) {
            break;
        }

        chord = chordLibraryFind( &library, input );
        if ( chord == NULL ) {
            printf( "Chord '%s' not found.\n\n", input );
            continue;
        }

        if ( !terminalRendererPrint( stdout, chord, TERMINAL_FRET_COUNT ) ) {
            fputs( "Unable to render chord.\n", stderr );
            return 1;
        }
    }

    return 0;
}
