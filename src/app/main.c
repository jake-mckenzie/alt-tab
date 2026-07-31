#define _CRT_SECURE_NO_WARNINGS

#include <limits.h>
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

/* Recognizes aliases that request every variation of one chord. */
static int isAllVariationsCommand( const char *input )
{
    return strcmp( input, "--all-variations" ) == 0 ||
           strcmp( input, "-a" ) == 0;
}

/* Returns a human-readable label for the active renderer. */
static const char *renderModeName( TerminalRenderMode mode )
{
    return mode == TERMINAL_RENDER_FULL_NECK
        ? "full neck (frets 1-27)"
        : "compact horizontal tab";
}

/* Prints chord names directly from the immutable library. */
static void printSupportedChords( const ChordLibrary *library )
{
    const char *previous_name = NULL;
    int first = 1;

    for ( size_t index = 0; index < library->count; index++ ) {
        const char *name = library->items[ index ].name;

        if ( previous_name != NULL && strcmp( previous_name, name ) == 0 ) {
            continue;
        }

        printf( "%s%s",
                first ? "" : "  ",
                name );
        previous_name = name;
        first = 0;
    }
    fputc( '\n', stdout );
}

/* Prints a small sample instead of crowding the startup screen. */
static void printChordExamples(
    const ChordLibrary *library,
    size_t maximum
)
{
    const char *previous_name = NULL;
    size_t printed = 0;

    for ( size_t index = 0;
          index < library->count && printed < maximum;
          index++ ) {
        const char *name = library->items[ index ].name;

        if ( previous_name != NULL && strcmp( previous_name, name ) == 0 ) {
            continue;
        }

        printf( "%s%s", printed == 0 ? "" : "  ", name );
        previous_name = name;
        printed++;
    }
    fputc( '\n', stdout );
}

/* Prints the accepted command-line syntax. */
static void printUsage( FILE *stream, const char *program )
{
    fprintf(
        stream,
        "Usage: %s [CHORD[:VARIATION]] [OPTIONS]\n",
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
    puts( "  <chord>           Draw variation 1" );
    puts( "  <chord>:<number>  Draw one variation, such as C:2" );
    puts( "  <chord> -a        Draw every variation" );
    puts( "  <chord> -f        Draw on the complete neck" );
    puts( "  -a, --all-variations  Show all variations of a chord" );
    puts( "  -f, --full-neck   Switch to the complete 1-27 fret neck" );
    puts( "  -c, --compact     Switch to compact horizontal tab" );
    puts( "  ?, help           Show this help" );
    puts( "  q                 Quit\n" );

    puts( "SUPPORTED CHORDS" );
    printf( "  " );
    printSupportedChords( library );

    puts( "\nLAUNCH EXAMPLES" );
    printf( "  %s C:2             Draw C variation 2\n", program );
    printf( "  %s C -a            Draw every C variation\n", program );
    printf( "  %s A -f            Draw A on the full neck\n", program );
    printf( "  %s -f              Start interactively in full-neck mode\n\n",
            program );
}

/* Presents concise startup guidance without overwhelming the prompt. */
static void printWelcome(
    const ChordLibrary *library,
    TerminalRenderMode mode
)
{
    puts( "+----------------------------------------------------------+" );
    puts( "|  ALT-TAB                                                 |" );
    puts( "|  Guitar Chord Viewer                                     |" );
    puts( "+----------------------------------------------------------+" );
    printf( "  Mode    : %s\n", renderModeName( mode ) );
    printf( "  Examples: " );
    printChordExamples( library, 4 );
    puts( "  Variants: C:1  C:2  C -a" );
    puts( "  Commands: CHORD, CHORD:2, CHORD -a, ?, or q" );
    fputc( '\n', stdout );
}

/* Parses one optional chord and display flags in any argument order. */
static int parseArguments(
    int argc,
    char *argv[ ],
    TerminalRenderMode *mode,
    const char **chord_name,
    int *show_all
)
{
    for ( int index = 1; index < argc; index++ ) {
        if ( isFullNeckCommand( argv[ index ] ) ) {
            *mode = TERMINAL_RENDER_FULL_NECK;
        } else if ( isCompactCommand( argv[ index ] ) ) {
            *mode = TERMINAL_RENDER_COMPACT_TAB;
        } else if ( isAllVariationsCommand( argv[ index ] ) ) {
            *show_all = 1;
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

/* Parses a selector such as C or C:2 without modifying its source. */
static int parseChordSelector(
    const char *selector,
    char *name,
    size_t name_capacity,
    int *variation,
    int *has_variation
)
{
    const char *separator = strchr( selector, ':' );
    size_t name_length = separator == NULL
        ? strlen( selector )
        : ( size_t )( separator - selector );

    if ( name_length == 0 || name_length >= name_capacity ) {
        return 0;
    }

    memcpy( name, selector, name_length );
    name[ name_length ] = '\0';
    *variation = 1;
    *has_variation = separator != NULL;

    if ( separator != NULL ) {
        const char *digit = separator + 1;
        int parsed = 0;

        if ( *digit == '\0' ) {
            return 0;
        }

        while ( *digit != '\0' ) {
            int value;

            if ( *digit < '0' || *digit > '9' ) {
                return 0;
            }

            value = *digit - '0';
            if ( parsed > ( INT_MAX - value ) / 10 ) {
                return 0;
            }

            parsed = parsed * 10 + value;
            digit++;
        }

        if ( parsed < 1 ) {
            return 0;
        }
        *variation = parsed;
    }

    return 1;
}

/* Resolves a selector and renders one or every matching variation. */
static int renderChordSelection(
    const ChordLibrary *library,
    const char *selector,
    int show_all,
    TerminalRenderMode mode
)
{
    char name[ 16 ];
    int variation;
    int has_variation;
    size_t variation_count;

    if ( !parseChordSelector(
             selector,
             name,
             sizeof( name ),
             &variation,
             &has_variation
         ) ) {
        fprintf( stderr, "Invalid chord selector: %s\n", selector );
        return 0;
    }

    variation_count = chordLibraryVariationCount( library, name );
    if ( variation_count == 0 ) {
        fprintf( stderr, "Unknown chord: %s\n", name );
        return 0;
    }

    if ( show_all && has_variation ) {
        fputs( "Choose a variation or use -a, not both.\n", stderr );
        return 0;
    }

    if ( show_all ) {
        for ( size_t index = 1; index <= variation_count; index++ ) {
            const Chord *chord = chordLibraryFindVariation(
                library,
                name,
                ( int )index
            );

            if ( chord == NULL ||
                 !terminalRendererPrint(
                     stdout,
                     chord,
                     TERMINAL_LAST_FRET,
                     mode
                 ) ) {
                return 0;
            }
        }
        return 1;
    }

    {
        const Chord *chord = chordLibraryFindVariation(
            library,
            name,
            variation
        );

        if ( chord == NULL ) {
            fprintf(
                stderr,
                "%s has %zu variations; variation %d does not exist.\n",
                name,
                variation_count,
                variation
            );
            return 0;
        }

        return terminalRendererPrint(
            stdout,
            chord,
            TERMINAL_LAST_FRET,
            mode
        );
    }
}

/* Runs one-shot chord output or the interactive chord prompt. */
int main( int argc, char *argv[ ] )
{
    ChordLibrary library = chordLibraryDefault( );
    TerminalRenderMode render_mode = TERMINAL_RENDER_COMPACT_TAB;
    const char *requested_chord = NULL;
    int show_all = 0;
    int argument_result = parseArguments(
        argc,
        argv,
        &render_mode,
        &requested_chord,
        &show_all
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
        return renderChordSelection(
            &library,
            requested_chord,
            show_all,
            render_mode
        ) ? 0 : 1;
    }

    if ( show_all ) {
        fputs( "-a requires a chord name.\n", stderr );
        return 1;
    }

    printWelcome( &library, render_mode );

    for ( ;; ) {
        char input[ 64 ];
        char *tokens[ 3 ];
        size_t token_count = 0;
        char *token;
        const char *chord_name = NULL;
        int mode_changed = 0;
        int command_show_all = 0;
        TerminalRenderMode command_mode = render_mode;

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
            } else if ( isAllVariationsCommand( tokens[ index ] ) ) {
                command_show_all = 1;
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

        if ( chord_name == NULL && command_show_all ) {
            puts( "-a requires a chord name.\n" );
            continue;
        }

        if ( chord_name == NULL ) {
            continue;
        }

        if ( !renderChordSelection(
                 &library,
                 chord_name,
                 command_show_all,
                 command_mode
             ) ) {
            puts( "Type ? to see supported chords and commands.\n" );
        }
    }

    return 0;
}
