#ifndef ALT_TAB_RENDER_TERMINAL_RENDERER_H
#define ALT_TAB_RENDER_TERMINAL_RENDERER_H

#include <stddef.h>
#include <stdio.h>
#include "theory/chord.h"

#define TERMINAL_LAST_FRET 27

/* Selects a truncated or complete horizontal fretboard. */
typedef enum {
    TERMINAL_RENDER_COMPACT_TAB = 0,
    TERMINAL_RENDER_FULL_NECK
} TerminalRenderMode;

/* Writes a validated chord diagram to the requested stream. */
int terminalRendererPrint(
    FILE *stream,
    const Chord *chord,
    size_t last_fret,
    TerminalRenderMode mode
);

#endif
