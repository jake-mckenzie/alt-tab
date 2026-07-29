#define _CRT_SECURE_NO_WARNINGS

#include <stdio.h>
#include <string.h>

#include "render/terminal_renderer.h"
#include "theory/chord_library.h"

/* Recognizes the interactive quit command. */
static int shouldQuit( const char *input )
{
    return strcmp( input, "q" ) == 0 || strcmp( input, "Q" ) == 0;
}

/* Recognizes help aliases accepted at launch and in the prompt. */
static int isHelpCommand( const char *input )
{
    return strcmp( input, "?" ) == 0 ||
           strcmp( input, "help" ) == 0 ||
           strcmp( input, "--help" ) == 0 ||
           strcmp( input, "-h" ) == 0;
}

/* Recognizes aliases that enable the complete neck. */
static int isFullNeckCommand( const char *input )
{
    return strcmp( input, "--full-neck" ) == 0 ||
           strcmp( input, "-f" ) == 0;
}

/* Recognizes aliases that restore truncated tab output. */
static int isCompactCommand( const char *input )
{
    return strcmp( input, "--compact" ) == 0 ||
           strcmp( input, "-c" ) == 0;
}

/* Returns a human-readable label for the active renderer. */
static const char *renderModeName( TerminalRenderMode mode )
{
    return mode == TERMINAL_RENDER_FULL_NECK
        ? "full neck (frets 0-27)"
        : "compact horizontal tab";
}

/* Prints chord names directly from the immutable library. */
static void printSupportedChords( const ChordLibrary *library )
{
    for ( size_t index = 0; index < library->count; index++ ) {
        printf( "%s%s",
                index == 0 ? "" : "  ",
                library->items[ index ].name );
    }
    fputc( '\n', stdout );
}

/* Prints the accepted command-line syntax. */
static void printUsage( FILE *stream, const char *program )
{
    fprintf(
        stream,
        "Usage: %s [CHORD] [--full-neck|-f] [--compact|-c] [--help|-h]\n",
        program
    );
}

/* Shows interactive commands, chords, and launch examples. */
static void printHelp(
    const char *program,
    const ChordLibrary *library
)
{
    puts( "\nCOMMANDS" );
    puts( "  <chord>           Draw a compact horizontal tab" );
    puts( "  <chord> -f        Draw that chord on the complete neck" );
    puts( "  -f, --full-neck   Switch to the complete 0-27 fret neck" );
    puts( "  -c, --compact     Switch to compact horizontal tab" );
    puts( "  ?, help           Show this help" );
    puts( "  q                 Quit\n" );

    puts( "SUPPORTED CHORDS" );
    printf( "  " );
    printSupportedChords( library );

    puts( "\nLAUNCH EXAMPLES" );
    printf( "  %s A               Draw A in compact mode\n", program );
    printf( "  %s A -f            Draw A on the full neck\n", program );
    printf( "  %s -f              Start interactively in full-neck mode\n\n",
            program );
}

/* Presents concise startup guidance without overwhelming the prompt. */
static void printWelcome(
    const char *program,
    const ChordLibrary *library,
    TerminalRenderMode mode
)
{
    puts( "+----------------------------------------------------------+" );
    puts( "|  ALT-TAB                                                 |" );
    puts( "|  Guitar Chord Viewer                                     |" );
    puts( "+----------------------------------------------------------+" );
    printf( "  Mode    : %s\n", renderModeName( mode ) );
    printf( "  Chords  : " );
    printSupportedChords( library );
    puts( "  Commands: CHORD, CHORD -f, ?, or q" );
    printf( "  Example : %s A -f\n\n", program );
}

/* Parses one optional chord and display flags in any argument order. */
static int parseArguments(
    int argc,
    char *argv[ ],
    TerminalRenderMode *mode,
    const char **chord_name
)
{
    for ( int index = 1; index < argc; index++ ) {
        if ( isFullNeckCommand( argv[ index ] ) ) {
            *mode = TERMINAL_RENDER_FULL_NECK;
        } else if ( isCompactCommand( argv[ index ] ) ) {
            *mode = TERMINAL_RENDER_COMPACT_TAB;
        } else if ( isHelpCommand( argv[ index ] ) ) {
            return 0;
        } else if ( argv[ index ][ 0 ] == '-' ) {
            fprintf( stderr, "Unknown option: %s\n", argv[ index ] );
            printUsage( stderr, argv[ 0 ] );
            return -1;
        } else if ( *chord_name == NULL ) {
            *chord_name = argv[ index ];
        } else {
            fprintf( stderr, "Only one chord may be requested at a time.\n" );
            printUsage( stderr, argv[ 0 ] );
            return -1;
        }
    }

    return 1;
}

/* Runs one-shot chord output or the interactive chord prompt. */
int main( int argc, char *argv[ ] )
{
    ChordLibrary library = chordLibraryDefault( );
    TerminalRenderMode render_mode = TERMINAL_RENDER_COMPACT_TAB;
    const char *requested_chord = NULL;
    int argument_result = parseArguments(
        argc,
        argv,
        &render_mode,
        &requested_chord
    );

    if ( argument_result < 0 ) {
        return 1;
    }

    if ( argument_result == 0 ) {
        printUsage( stdout, argv[ 0 ] );
        printHelp( argv[ 0 ], &library );
        return 0;
    }

    if ( requested_chord != NULL ) {
        const Chord *chord = chordLibraryFind( &library, requested_chord );

        if ( chord == NULL ) {
            fprintf( stderr, "Unknown chord: %s\n", requested_chord );
            return 1;
        }

        return terminalRendererPrint(
            stdout,
            chord,
            TERMINAL_FRET_COUNT,
            render_mode
        ) ? 0 : 1;
    }

    printWelcome( argv[ 0 ], &library, render_mode );

    for ( ;; ) {
        char input[ 64 ];
        char *tokens[ 3 ];
        size_t token_count = 0;
        char *token;
        const char *chord_name = NULL;
        int mode_changed = 0;
        TerminalRenderMode command_mode = render_mode;
        const Chord *chord;

        printf( "chord> " );
        if ( fgets( input, sizeof( input ), stdin ) == NULL ) {
            break;
        }

        input[ strcspn( input, "\r\n" ) ] = '\0';
        token = strtok( input, " \t" );
        while ( token != NULL && token_count < 3 ) {
            tokens[ token_count++ ] = token;
            token = strtok( NULL, " \t" );
        }

        if ( token != NULL ) {
            puts( "Too many values. Type ? for usage.\n" );
            continue;
        }

        if ( token_count == 0 ) {
            continue;
        }

        if ( token_count == 1 && shouldQuit( tokens[ 0 ] ) ) {
            break;
        }

        if ( token_count == 1 && isHelpCommand( tokens[ 0 ] ) ) {
            printHelp( argv[ 0 ], &library );
            continue;
        }

        for ( size_t index = 0; index < token_count; index++ ) {
            if ( isFullNeckCommand( tokens[ index ] ) ) {
                command_mode = TERMINAL_RENDER_FULL_NECK;
                mode_changed = 1;
            } else if ( isCompactCommand( tokens[ index ] ) ) {
                command_mode = TERMINAL_RENDER_COMPACT_TAB;
                mode_changed = 1;
            } else if ( tokens[ index ][ 0 ] == '-' ) {
                printf( "Unknown option: %s\n\n", tokens[ index ] );
                chord_name = NULL;
                mode_changed = 0;
                break;
            } else if ( chord_name == NULL ) {
                chord_name = tokens[ index ];
            } else {
                puts( "Only one chord may be requested at a time.\n" );
                chord_name = NULL;
                mode_changed = 0;
                break;
            }
        }

        if ( chord_name == NULL && mode_changed ) {
            render_mode = command_mode;
            printf( "Display changed to %s.\n\n", renderModeName( render_mode ) );
            continue;
        }

        if ( chord_name == NULL ) {
            continue;
        }

        chord = chordLibraryFind( &library, chord_name );
        if ( chord == NULL ) {
            printf( "Unknown chord: %s\n", chord_name );
            puts( "Type ? to see supported chords and commands.\n" );
            continue;
        }

        if ( !terminalRendererPrint(
                 stdout,
                 chord,
                 TERMINAL_FRET_COUNT,
                 command_mode
             ) ) {
            fputs( "Unable to render chord.\n", stderr );
            return 1;
        }
    }

    return 0;
}
