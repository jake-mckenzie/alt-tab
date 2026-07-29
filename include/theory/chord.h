#ifndef ALT_TAB_THEORY_CHORD_H
#define ALT_TAB_THEORY_CHORD_H

/* Identifies strings in the high-to-low order used by standard tablature. */
typedef enum {
    GUITAR_STRING_HIGH_E = 0,
    GUITAR_STRING_B,
    GUITAR_STRING_G,
    GUITAR_STRING_D,
    GUITAR_STRING_A,
    GUITAR_STRING_LOW_E,
    GUITAR_STRING_COUNT
} GuitarString;

#define GUITAR_MUTED_FRET ( -1 )
#define GUITAR_NO_FINGER 0

/* Describes the fret and fretting-hand finger used on one string. */
typedef struct {
    int fret;
    int finger;
} StringPlacement;

/* Stores one named six-string guitar-chord voicing. */
typedef struct {
    const char *name;
    int variation;
    StringPlacement strings[ GUITAR_STRING_COUNT ];
} Chord;

#endif
