#include "render/terminal_renderer.h"

/* Standard tab labels run from the highest string to the lowest. */
static const char *const STRING_NAMES[ GUITAR_STRING_COUNT ] = {
    "e", "B", "G", "D", "A", "E"
};

#define COMPACT_MINIMUM_LAST_FRET 4
#define COMPACT_CELL_WIDTH 5
#define FULL_NECK_CELL_WIDTH 3

/* Rejects invalid fret ranges and missing finger assignments. */
static int chordIsRenderable( const Chord *chord, size_t fret_count )
{
    if ( chord == NULL ||
         chord->name == NULL ||
         chord->variation < 1 ||
         fret_count == 0 ) {
        return 0;
    }

    for ( size_t string = 0; string < GUITAR_STRING_COUNT; string++ ) {
        StringPlacement placement = chord->strings[ string ];

        if ( placement.fret < GUITAR_MUTED_FRET ||
             ( placement.fret >= 0 &&
               ( size_t )placement.fret >= fret_count ) ) {
            return 0;
        }

        if ( placement.fret > 0 &&
             ( placement.finger < 1 || placement.finger > 4 ) ) {
            return 0;
        }

        if ( placement.fret <= 0 &&
             placement.finger != GUITAR_NO_FINGER ) {
            return 0;
        }
    }

    return 1;
}

/* Finds the final occupied fret for compact output. */
static size_t highestPlayedFret( const Chord *chord )
{
    size_t highest = 0;

    for ( size_t string = 0; string < GUITAR_STRING_COUNT; string++ ) {
        int fret = chord->strings[ string ].fret;

        if ( fret > 0 && ( size_t )fret > highest ) {
            highest = ( size_t )fret;
        }
    }

    return highest;
}

/* Keeps compact diagrams short while showing every occupied fret. */
static size_t compactLastFret(
    const Chord *chord,
    size_t fret_count
)
{
    size_t highest = highestPlayedFret( chord );
    size_t last_fret = highest > COMPACT_MINIMUM_LAST_FRET
        ? highest
        : COMPACT_MINIMUM_LAST_FRET;

    return last_fret < fret_count ? last_fret : fret_count - 1;
}

/* Selects the symbol displayed at one string-and-fret intersection. */
static char markerAt(
    const Chord *chord,
    size_t string,
    size_t fret
)
{
    StringPlacement placement = chord->strings[ string ];

    if ( fret == 0 && placement.fret == GUITAR_MUTED_FRET ) {
        return 'X';
    }

    if ( fret == 0 && placement.fret == 0 ) {
        return 'O';
    }

    if ( placement.fret > 0 && ( size_t )placement.fret == fret ) {
        return ( char )( '0' + placement.finger );
    }

    return '\0';
}

/* Draws one fixed-width segment of a horizontal string. */
static void printCell(
    FILE *stream,
    char marker,
    size_t width
)
{
    size_t marker_column = width / 2;

    for ( size_t column = 0; column < width; column++ ) {
        if ( marker != '\0' && column == marker_column ) {
            fputc( marker, stream );
        } else {
            fputc( '-', stream );
        }
    }
}

/* Counts decimal digits so fret labels can be centered. */
static size_t digitCount( size_t value )
{
    size_t digits = 1;

    while ( value >= 10 ) {
        value /= 10;
        digits++;
    }

    return digits;
}

/* Centers a fret number over its corresponding string segment. */
static void printFretNumber(
    FILE *stream,
    size_t fret,
    size_t width
)
{
    size_t digits = digitCount( fret );
    size_t left_padding = width > digits ? ( width - digits ) / 2 : 0;
    size_t right_padding = width > digits
        ? width - digits - left_padding
        : 0;

    for ( size_t column = 0; column < left_padding; column++ ) {
        fputc( ' ', stream );
    }
    fprintf( stream, "%zu", fret );
    for ( size_t column = 0; column < right_padding; column++ ) {
        fputc( ' ', stream );
    }
}

/* Prints six horizontal strings in standard high-e-to-low-E tab order. */
static void printTab(
    FILE *stream,
    const Chord *chord,
    size_t last_fret,
    size_t cell_width
)
{
    fputs( "    ", stream );
    for ( size_t fret = 0; fret <= last_fret; fret++ ) {
        printFretNumber( stream, fret, cell_width );
    }
    fputc( '\n', stream );

    for ( size_t string = 0; string < GUITAR_STRING_COUNT; string++ ) {
        fprintf( stream, "%s  |", STRING_NAMES[ string ] );
        for ( size_t fret = 0; fret <= last_fret; fret++ ) {
            printCell(
                stream,
                markerAt( chord, string, fret ),
                cell_width
            );
        }
        fputs( "|\n", stream );
    }
}

/* Renders either the compact or full-neck tab layout. */
int terminalRendererPrint(
    FILE *stream,
    const Chord *chord,
    size_t fret_count,
    TerminalRenderMode mode
)
{
    size_t last_fret;
    size_t cell_width;

    if ( stream == NULL ||
         !chordIsRenderable( chord, fret_count ) ||
         ( mode != TERMINAL_RENDER_COMPACT_TAB &&
           mode != TERMINAL_RENDER_FULL_NECK ) ) {
        return 0;
    }

    if ( mode == TERMINAL_RENDER_FULL_NECK ) {
        last_fret = fret_count - 1;
        cell_width = FULL_NECK_CELL_WIDTH;
    } else {
        last_fret = compactLastFret( chord, fret_count );
        cell_width = COMPACT_CELL_WIDTH;
    }

    fprintf(
        stream,
        "\nChord: %s (variation %d)\n\n",
        chord->name,
        chord->variation
    );
    fputs( "    Fret numbers ->\n", stream );
    printTab( stream, chord, last_fret, cell_width );
    fputs(
        "\n    Fingers: 1 index  2 middle  3 ring  4 little\n"
        "    Symbols: O open   X muted\n\n",
        stream
    );

    return ferror( stream ) == 0;
}
