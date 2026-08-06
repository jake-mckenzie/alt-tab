#ifndef ALT_TAB_CHORDS_H
#define ALT_TAB_CHORDS_H

#include <stddef.h>

/* A standard guitar voicing always contains six strings. */
enum { GUITAR_STRING_COUNT = 6 };

/* Negative fret values distinguish muted strings from open strings. */
#define GUITAR_MUTED_FRET ( -1 )

/* Describes the fret and fretting-hand finger used on one string. */
typedef struct {
    int fret;
    int finger;
} StringPlacement;

/* Stores one voicing in standard high-e-to-low-E tablature order. */
typedef struct {
    const char *name;
    int variation;
    StringPlacement strings[ GUITAR_STRING_COUNT ];
} Chord;

/* Returns the total number of stored voicings. */
size_t altTabChordVoicingCount( void );

/* Returns a stored voicing, or NULL when the index is out of range. */
const Chord *altTabChordVoicingAt( size_t index );

/* Returns the number of distinct chord names. */
size_t altTabChordNameCount( void );

/* Returns a stored chord name, or NULL when the index is out of range. */
const char *altTabChordNameAt( size_t index );

/* Returns the number of voicings stored for a case-insensitive chord name. */
size_t altTabChordVariationCount( const char *name );

/* Returns a numbered voicing by case-insensitive name, or NULL if absent. */
const Chord *altTabChordFind( const char *name, int variation );

#endif
