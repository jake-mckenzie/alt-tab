#ifndef ALT_TAB_RENDER_TERMINAL_RENDERER_H
#define ALT_TAB_RENDER_TERMINAL_RENDERER_H

#include <stddef.h>
#include <stdio.h>

#include "theory/chord.h"

#define TERMINAL_FRET_COUNT 28

int terminalRendererPrint(
    FILE *stream,
    const Chord *chord,
    size_t fret_count
);

#endif
