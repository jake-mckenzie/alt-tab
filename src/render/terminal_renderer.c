#include "render/terminal_renderer.h"

static const char *const STRING_NAMES[ GUITAR_STRING_COUNT ] = {
    "e", "B", "G", "D", "A", "E"
};

static int chordIsRenderable( const Chord *chord, size_t fret_count )
{
    if ( chord == NULL || chord->name == NULL || fret_count == 0 ) {
        return 0;
    }

    for ( size_t string = 0; string < GUITAR_STRING_COUNT; string++ ) {
        int fret = chord->frets[ string ];

        if ( fret < GUITAR_MUTED_FRET ||
             ( fret >= 0 && ( size_t )fret >= fret_count ) ) {
            return 0;
        }
    }

    return 1;
}

static char markerAt(
    const Chord *chord,
    size_t string,
    size_t fret
)
{
    int played_fret = chord->frets[ string ];

    if ( played_fret == GUITAR_MUTED_FRET && fret == 0 ) {
        return 'X';
    }

    if ( played_fret == 0 && fret == 0 ) {
        return 'O';
    }

    if ( played_fret > 0 && ( size_t )played_fret == fret ) {
        return '*';
    }

    return '-';
}

int terminalRendererPrint(
    FILE *stream,
    const Chord *chord,
    size_t fret_count
)
{
    if ( stream == NULL || !chordIsRenderable( chord, fret_count ) ) {
        return 0;
    }

    fprintf( stream, "\n%s\n     ", chord->name );
    for ( size_t fret = 0; fret < fret_count; fret++ ) {
        fprintf( stream, "%3zu", fret );
    }
    fputc( '\n', stream );

    for ( size_t string = 0; string < GUITAR_STRING_COUNT; string++ ) {
        fprintf( stream, "%s  | ", STRING_NAMES[ string ] );
        for ( size_t fret = 0; fret < fret_count; fret++ ) {
            fprintf( stream, "%3c", markerAt( chord, string, fret ) );
        }
        fputs( "|\n", stream );
    }

    fputc( '\n', stream );
    return ferror( stream ) == 0;
}
