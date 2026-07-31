#ifndef ALT_TAB_NATIVE_CHORD_API_H
#define ALT_TAB_NATIVE_CHORD_API_H

#include <stddef.h>

#define ALT_TAB_STRING_COUNT 6
#define ALT_TAB_MUTED_FRET ( -1 )

/* Describes one string without exposing internal chord-library types. */
typedef struct {
    int fret;
    int finger;
} AltTabStringPlacement;

/* Contains one chord voicing in standard high-e-to-low-E tab order. */
typedef struct {
    int variation;
    AltTabStringPlacement strings[ ALT_TAB_STRING_COUNT ];
} AltTabChordVoicing;

/* Returns the number of distinct chord names in the built-in library. */
size_t altTabChordCount( void );

/* Returns a read-only chord name, or NULL when the index is out of range. */
const char *altTabChordNameAt( size_t index );

/* Returns the number of available voicings for a chord name. */
size_t altTabChordVariationCount( const char *name );

/* Copies one voicing into caller-owned memory; returns zero when unavailable. */
int altTabChordLoad(
    const char *name,
    int variation,
    AltTabChordVoicing *output
);

#endif
