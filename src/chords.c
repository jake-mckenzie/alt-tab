#include "alt_tab.h"

#include <string.h>

enum { MUTED = -1, ACCIDENTAL = -1, BASE = 0, MINOR = 1 };

typedef struct { int frets[6]; int fingers[6]; int movable; } Shape;
typedef struct { int shape; int base_fret; } Spec;
typedef struct { const char *name; Spec voices[3]; } Definition;

static const Shape shapes[] = {
    {{0,2,2,2,0,-1},{0,3,2,1,0,0},0}, {{0,1,2,2,0,-1},{0,1,3,2,0,0},0},
    {{0,1,0,2,3,-1},{0,1,0,2,3,0},0}, {{2,3,2,0,-1,-1},{2,3,1,0,0,0},0},
    {{1,3,2,0,-1,-1},{1,3,2,0,0,0},0}, {{0,0,1,2,2,0},{0,0,1,3,2,0},0},
    {{0,0,0,2,2,0},{0,0,0,3,2,0},0}, {{3,0,0,0,2,3},{3,0,0,0,1,2},0},
    {{0,0,1,2,2,0},{1,1,2,4,3,1},1}, {{0,2,2,2,0,-1},{1,4,3,2,1,0},1},
    {{0,0,0,2,2,0},{1,1,1,4,3,1},1}, {{0,1,2,2,0,-1},{1,2,4,3,1,0},1},
    {{2,3,2,0,-1,-1},{3,4,2,1,0,0},1}, {{1,3,2,0,-1,-1},{2,4,3,1,0,0},1}
};

static const Definition definitions[] = {
    {"A",{{0,0},{8,5},{12,7}}}, {"Am",{{1,0},{10,5},{13,7}}},
    {"B",{{9,2},{8,7},{12,9}}}, {"Bb",{{9,1},{8,6},{12,8}}},
    {"C",{{2,0},{9,3},{12,10}}}, {"Cm",{{11,3},{10,8},{13,10}}},
    {"D",{{3,0},{9,5},{12,12}}}, {"Dm",{{4,0},{11,5},{13,12}}},
    {"E",{{5,0},{9,7},{12,2}}}, {"Em",{{6,0},{11,7},{13,2}}},
    {"F",{{8,1},{9,8},{12,3}}}, {"F#",{{8,2},{9,9},{12,4}}},
    {"G",{{7,0},{8,3},{12,5}}}, {"Gm",{{10,3},{11,10},{13,5}}}
};

static const char *families[] = {"A","B","C","D","E","F","G"};

static int wrap(int value, int size) { value %= size; return value < 0 ? value + size : value; }
static int index_for_name(const char *name) {
    for (int i = 0; i < ALT_TAB_CHORD_COUNT; ++i) if (strcmp(definitions[i].name, name) == 0) return i;
    return 0;
}
static const char *name_for(const AltTabController *c) {
    static const char *accidentals[] = {NULL,"Bb",NULL,NULL,NULL,"F#",NULL};
    static const char *minors[] = {"Am",NULL,"Cm","Dm","Em",NULL,"Gm"};
    if (c->kind == ACCIDENTAL && accidentals[c->selected]) return accidentals[c->selected];
    if (c->kind == MINOR && minors[c->selected]) return minors[c->selected];
    return families[c->selected];
}

void alt_tab_controller_init(AltTabController *c) { c->selected = 0; c->kind = BASE; c->voicing_number = 1; }
void alt_tab_move_chord(AltTabController *c, int delta) { c->selected = wrap(c->selected + delta, 7); c->kind = BASE; c->voicing_number = 1; }
void alt_tab_move_kind(AltTabController *c, int delta) {
    static const int has_accidental[] = {0,1,0,0,0,1,0};
    static const int has_minor[] = {1,0,1,1,1,0,1};
    int next = c->kind;
    if (delta < 0 && c->kind == MINOR) next = BASE;
    else if (delta < 0 && c->kind == BASE && has_accidental[c->selected]) next = ACCIDENTAL;
    else if (delta > 0 && c->kind == ACCIDENTAL) next = BASE;
    else if (delta > 0 && c->kind == BASE && has_minor[c->selected]) next = MINOR;
    if (next != c->kind) { c->kind = next; c->voicing_number = 1; }
}
void alt_tab_cycle_voicing(AltTabController *c) { c->voicing_number = c->voicing_number % 3 + 1; }
const char *alt_tab_current_name(const AltTabController *c) { return name_for(c); }
const AltTabVoicing *alt_tab_current_voicing(const AltTabController *c) {
    static AltTabVoicing voice;
    const Definition *definition = &definitions[index_for_name(name_for(c))];
    const Spec spec = definition->voices[c->voicing_number - 1];
    const Shape *shape = &shapes[spec.shape];
    voice.name = definition->name; voice.number = c->voicing_number;
    for (int i = 0; i < 6; ++i) {
        voice.strings[i].fret = shape->frets[i];
        if (shape->movable && voice.strings[i].fret != MUTED) voice.strings[i].fret += spec.base_fret;
        voice.strings[i].finger = shape->fingers[i];
    }
    return &voice;
}
const char *alt_tab_family_base(int index) { return families[index]; }
size_t alt_tab_family_count(void) { return sizeof(families) / sizeof(families[0]); }
