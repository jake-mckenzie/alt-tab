// Alt-Tab
// Takes in notes/chords via console, returns all possible variations of note/chord
// Guitar layout via console
// Version: 0.001
#define _CRT_SECURE_NO_WARNINGS

#include <stdio.h>
#include <ctype.h>
#include <string.h>

#include "matrix.h"
#include "chords.h"

int main(void)
{
    Matrix fretBoard = createMatrix();
    ChordLibrary library;

    initChordLibrary(&library);

    puts("Alt-Tab — Guitar Chord Viewer\n");

    for (;;) {
        char input[16];

        printf("Enter chord name (or 'q' to quit): ");
        if (!fgets(input, sizeof input, stdin)){
            break;
        }

        /* remove newline */
        input[strcspn(input, "\n")] = '\0';

        /* normalize to uppercase */
        for (char *p = input; *p; ++p) {
            *p = (char)toupper((unsigned char)*p);
        }

        if (strcmp(input, "Q") == 0){
            break;
        }

        const Chord *chord = findChord(&library, input);

        if (!chord) {
            printf("Chord '%s' not found.\n\n", input);
            continue;
        }

        drawChord(&fretBoard, chord);
        matrixPrinter(&fretBoard);
    }

    return 0;
}
