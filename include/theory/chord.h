#ifndef ALT_TAB_THEORY_CHORD_H
#define ALT_TAB_THEORY_CHORD_H

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

typedef struct {
    const char *name;
    int frets[ GUITAR_STRING_COUNT ];
} Chord;

#endif
