#ifndef ALT_TAB_H
#define ALT_TAB_H

#include <stddef.h>

enum { ALT_TAB_STRING_COUNT = 6, ALT_TAB_CHORD_COUNT = 14 };

typedef struct {
    int fret;
    int finger;
} AltTabPlacement;

typedef struct {
    const char *name;
    int number;
    AltTabPlacement strings[ALT_TAB_STRING_COUNT];
} AltTabVoicing;

typedef struct {
    int selected;
    int kind;
    int voicing_number;
} AltTabController;

void alt_tab_controller_init(AltTabController *controller);
void alt_tab_move_chord(AltTabController *controller, int delta);
void alt_tab_move_kind(AltTabController *controller, int delta);
void alt_tab_cycle_voicing(AltTabController *controller);
const AltTabVoicing *alt_tab_current_voicing(const AltTabController *controller);
const char *alt_tab_current_name(const AltTabController *controller);
const char *alt_tab_family_base(int index);
size_t alt_tab_family_count(void);

#endif
