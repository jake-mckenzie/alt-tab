#include "alt_tab.h"
#include <assert.h>
#include <string.h>

int main(void) {
    AltTabController controller;
    alt_tab_controller_init(&controller);
    assert(strcmp(alt_tab_current_name(&controller), "A") == 0);
    assert(alt_tab_current_voicing(&controller)->strings[1].fret == 2);
    alt_tab_move_kind(&controller, 1);
    assert(strcmp(alt_tab_current_name(&controller), "Am") == 0);
    alt_tab_move_chord(&controller, 1);
    assert(strcmp(alt_tab_current_name(&controller), "B") == 0);
    alt_tab_cycle_voicing(&controller);
    assert(alt_tab_current_voicing(&controller)->number == 2);
    return 0;
}
