#include "alt_tab.h"
#include "raylib.h"

#include <math.h>
#include <stdio.h>

static Color background = {201, 197, 189, 255};
static Color panel = {41, 37, 48, 255};
static Color border = {98, 90, 112, 255};
static Color text = {244, 240, 233, 255};
static Color muted = {183, 173, 185, 255};
static Color accent = {122, 74, 160, 255};
static Color secondary = {208, 90, 146, 255};
static Color signal = {232, 184, 63, 255};

static void panel_box(Rectangle box, const char *title) {
    DrawRectangleRec(box, panel);
    DrawRectangleLinesEx(box, 2, border);
    DrawText(title, (int)box.x + 28, (int)box.y + 20, 24, accent);
    DrawRectangle((int)box.x + 28, (int)box.y + 54, 180, 3, secondary);
}

static void draw_dial(Rectangle box, const AltTabController *controller) {
    panel_box(box, "CHORD DIAL");
    const int count = (int)alt_tab_family_count();
    float start = box.x + box.width / 2.0f - count * 40.0f;
    for (int i = 0; i < count; ++i) {
        Rectangle cell = {start + i * 80.0f, box.y + 62, 56, 32};
        if (i == controller->selected) DrawRectangleRounded(cell, 0.2f, 4, Fade(accent, 0.25f));
        DrawText(alt_tab_family_base(i), (int)cell.x + 20, (int)cell.y + 5, 22, i == controller->selected ? accent : muted);
    }
    char selected[24];
    snprintf(selected, sizeof selected, "<  %s  >", alt_tab_current_name(controller));
    int width = MeasureText(selected, 28);
    DrawText(selected, (int)(box.x + (box.width - width) / 2), (int)box.y + (int)box.height - 68, 28, text);
}

static void draw_fretboard(Rectangle box, const AltTabVoicing *voice) {
    static const char *names[] = {"e", "B", "G", "D", "A", "E"};
    float left = box.x + 95, top = box.y + 118, width = box.width - 135, height = box.height - 200;
    float cell = width / 4.0f, gap = height / 7.0f;
    for (int fret = 0; fret <= 4; ++fret) {
        float x = left + fret * cell;
        DrawLineEx((Vector2){x, top + gap}, (Vector2){x, top + gap * 6}, 1, border);
        char label[4]; snprintf(label, sizeof label, "%d", fret + 1);
        DrawText(label, (int)x + (int)cell / 2 - 8, (int)top + 2, 20, text);
    }
    DrawLineEx((Vector2){left + width, top + gap}, (Vector2){left + width, top + gap * 6}, 1, border);
    for (int i = 0; i < 6; ++i) {
        float y = top + gap * (i + 1);
        DrawText(names[i], (int)left - 58, (int)y - 10, 20, muted);
        DrawLineEx((Vector2){left, y}, (Vector2){left + width, y}, 2, text);
        int fret = voice->strings[i].fret;
        if (fret <= 0) { DrawText(fret < 0 ? "X" : "O", (int)left - 28, (int)y - 10, 20, secondary); continue; }
        if (fret > 4) continue;
        float x = left + (fret - 0.5f) * cell;
        DrawCircle((int)x, (int)y, 18, accent);
        char finger[4]; snprintf(finger, sizeof finger, "%d", voice->strings[i].finger);
        DrawText(finger, (int)x - 6, (int)y - 10, 20, panel);
    }
}

static void draw_wave(Rectangle box) {
    panel_box(box, "WAVEFORM · AMPLITUDE / TIME");
    Rectangle plot = {box.x + 70, box.y + 90, box.width - 110, box.height - 140};
    DrawLineEx((Vector2){plot.x, plot.y}, (Vector2){plot.x, plot.y + plot.height}, 2, border);
    DrawLineEx((Vector2){plot.x, plot.y + plot.height / 2}, (Vector2){plot.x + plot.width, plot.y + plot.height / 2}, 1, border);
    for (int x = 1; x < (int)plot.width; ++x) {
        float t0 = (x - 1) / plot.width * 25.0f;
        float t1 = x / plot.width * 25.0f;
        float y0 = sinf(t0 * 2.6f) * 0.35f + sinf(t0 * 5.0f) * 0.16f;
        float y1 = sinf(t1 * 2.6f) * 0.35f + sinf(t1 * 5.0f) * 0.16f;
        DrawLineEx((Vector2){plot.x + x - 1, plot.y + plot.height / 2 - y0 * plot.height}, (Vector2){plot.x + x, plot.y + plot.height / 2 - y1 * plot.height}, 2, signal);
    }
    DrawText("0 ms", (int)plot.x, (int)(plot.y + plot.height + 8), 16, muted);
    DrawText("25 ms", (int)(plot.x + plot.width - 45), (int)(plot.y + plot.height + 8), 16, muted);
}

static void draw_spectrum(Rectangle box) {
    panel_box(box, "FREQUENCY SPECTRUM · LOG Hz");
    float baseline = box.y + box.height - 54;
    DrawLineEx((Vector2){box.x + 50, baseline}, (Vector2){box.x + box.width - 50, baseline}, 1, border);
    const char *notes[] = {"A2", "E3", "A3", "C#4", "E4"};
    for (int i = 0; i < 5; ++i) {
        float x = box.x + 70 + i * (box.width - 140) / 4.0f;
        DrawLineEx((Vector2){x, baseline}, (Vector2){x, box.y + 104}, 4, secondary);
        DrawCircle((int)x, (int)box.y + 104, 9, signal);
        DrawText(notes[i], (int)x - 12, (int)baseline + 12, 16, text);
    }
}

int main(void) {
    AltTabController controller;
    alt_tab_controller_init(&controller);
    SetConfigFlags(FLAG_WINDOW_RESIZABLE | FLAG_VSYNC_HINT);
    InitWindow(1280, 900, "Alt-Tab · Raylib");
    SetWindowMinSize(960, 800);
    SetTargetFPS(60);
    while (!WindowShouldClose()) {
        if (IsKeyPressed(KEY_LEFT) || IsKeyPressed(KEY_H)) alt_tab_move_chord(&controller, -1);
        if (IsKeyPressed(KEY_RIGHT) || IsKeyPressed(KEY_L)) alt_tab_move_chord(&controller, 1);
        if (IsKeyPressed(KEY_UP) || IsKeyPressed(KEY_K)) alt_tab_move_kind(&controller, -1);
        if (IsKeyPressed(KEY_DOWN) || IsKeyPressed(KEY_J)) alt_tab_move_kind(&controller, 1);
        if (IsKeyPressed(KEY_V)) alt_tab_cycle_voicing(&controller);
        float w = (float)GetScreenWidth(), h = (float)GetScreenHeight(), gap = 12;
        Rectangle header = {gap, gap, w - gap * 2, 105};
        Rectangle dial = {gap, 129, w - gap * 2, 180};
        float content_y = 321, content_h = h - content_y - gap, left_w = (w - gap * 3) * .54f;
        Rectangle diagram = {gap, content_y, left_w, content_h};
        Rectangle wave = {gap * 2 + left_w, content_y, w - left_w - gap * 3, (content_h - gap) * .58f};
        Rectangle spectrum = {wave.x, wave.y + wave.height + gap, wave.width, content_h - wave.height - gap};
        BeginDrawing(); ClearBackground(background);
        DrawRectangleRec(header, panel); DrawRectangleLinesEx(header, 2, border);
        DrawCircle((int)header.x + 26, (int)header.y + 38, 5, secondary);
        DrawText("ALT-TAB", (int)header.x + 58, (int)header.y + 18, 42, accent);
        DrawText("RAYLIB CHORD VIEWER - SUPER FAMICOM", (int)header.x + 58, (int)header.y + 72, 18, muted);
        DrawText("←/→ CHORD   ↑/↓ TYPE   V VOICING   F NECK   N TAB   T THEME   SPACE PLAY   F1 HELP   Q QUIT", (int)header.x + 560, (int)header.y + 33, 16, text);
        draw_dial(dial, &controller); panel_box(diagram, "CHORD DIAGRAM · COMPACT");
        const AltTabVoicing *voice = alt_tab_current_voicing(&controller); char label[64];
        snprintf(label, sizeof label, "%s   -   VOICING %d OF 3", voice->name, voice->number);
        DrawText(label, (int)(diagram.x + diagram.width / 2 - MeasureText(label, 24) / 2), (int)diagram.y + 70, 24, text);
        draw_fretboard(diagram, voice); draw_wave(wave); draw_spectrum(spectrum); EndDrawing();
    }
    CloseWindow(); return 0;
}
