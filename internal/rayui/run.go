//go:build raylib

// Package rayui renders the alternate Raylib interface.
package rayui

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/jake-mckenzie/alt-tab/internal/app"
	"github.com/jake-mckenzie/alt-tab/internal/chords"
	"github.com/jake-mckenzie/alt-tab/internal/signal"
)

const (
	windowWidth         = 1280
	windowHeight        = 900
	minimumWindowWidth  = 960
	minimumWindowHeight = 700
	panelGap            = 12
	panelPadding        = 18
	audioSampleRate     = 48000
	audioBufferFrames   = 2048
	waveformSeconds     = 0.025
	fullNeckLastFret    = 27
)

// palette defines one graphical equivalent of the terminal themes.
type palette struct {
	name       string
	background rl.Color
	panel      rl.Color
	border     rl.Color
	text       rl.Color
	muted      rl.Color
	accent     rl.Color
	secondary  rl.Color
	error      rl.Color
}

// palettes begin with an SNES-inspired gray shell and purple controls.
var palettes = [...]palette{
	{
		name: "Super 16", background: rl.NewColor(166, 161, 178, 255),
		panel: rl.NewColor(43, 36, 56, 255), border: rl.NewColor(113, 86, 164, 255),
		text: rl.NewColor(242, 239, 247, 255), muted: rl.NewColor(187, 176, 202, 255),
		accent: rl.NewColor(203, 169, 247, 255), secondary: rl.NewColor(165, 105, 215, 255),
		error: rl.NewColor(176, 48, 73, 255),
	},
	{
		name: "Synthwave", background: rl.NewColor(12, 10, 28, 255),
		panel: rl.NewColor(25, 20, 48, 255), border: rl.NewColor(106, 77, 148, 255),
		text: rl.NewColor(241, 238, 255, 255), muted: rl.NewColor(153, 145, 181, 255),
		accent: rl.NewColor(255, 74, 181, 255), secondary: rl.NewColor(86, 224, 255, 255),
		error: rl.NewColor(255, 99, 115, 255),
	},
	{
		name: "Tidal", background: rl.NewColor(7, 22, 31, 255),
		panel: rl.NewColor(11, 39, 51, 255), border: rl.NewColor(39, 99, 120, 255),
		text: rl.NewColor(226, 247, 250, 255), muted: rl.NewColor(129, 170, 181, 255),
		accent: rl.NewColor(65, 220, 214, 255), secondary: rl.NewColor(72, 153, 255, 255),
		error: rl.NewColor(255, 112, 112, 255),
	},
	{
		name: "Ember", background: rl.NewColor(28, 17, 11, 255),
		panel: rl.NewColor(51, 29, 17, 255), border: rl.NewColor(125, 73, 37, 255),
		text: rl.NewColor(255, 242, 222, 255), muted: rl.NewColor(190, 154, 116, 255),
		accent: rl.NewColor(255, 145, 48, 255), secondary: rl.NewColor(255, 208, 88, 255),
		error: rl.NewColor(255, 91, 74, 255),
	},
	{
		name: "Evergreen", background: rl.NewColor(8, 24, 19, 255),
		panel: rl.NewColor(14, 43, 33, 255), border: rl.NewColor(47, 105, 81, 255),
		text: rl.NewColor(230, 250, 239, 255), muted: rl.NewColor(135, 176, 155, 255),
		accent: rl.NewColor(86, 224, 151, 255), secondary: rl.NewColor(157, 231, 111, 255),
		error: rl.NewColor(255, 108, 108, 255),
	},
}

// screenLayout contains the responsive rectangles used for input and drawing.
type screenLayout struct {
	header   rl.Rectangle
	dial     rl.Rectangle
	diagram  rl.Rectangle
	waveform rl.Rectangle
	spectrum rl.Rectangle
}

// viewer owns graphical-only state around the shared chord controller.
type viewer struct {
	controller *app.Controller
	families   []app.Family
	fullNeck   bool
	tabNumbers bool
	help       bool
	theme      int
	quit       bool
	audio      audioPlayer
}

// Run initializes Raylib and serves frames until the window closes.
func Run(catalog chords.Catalog) error {
	controller := app.NewController(catalog)
	if controller.Err() != nil {
		return controller.Err()
	}

	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagVsyncHint | rl.FlagMsaa4xHint)
	rl.InitWindow(windowWidth, windowHeight, "Alt-Tab · Raylib")
	if !rl.IsWindowReady() {
		return errors.New("raylib window initialization failed")
	}
	defer rl.CloseWindow()
	rl.SetWindowMinSize(minimumWindowWidth, minimumWindowHeight)
	rl.SetExitKey(rl.KeyNull)
	rl.SetTargetFPS(60)

	rl.InitAudioDevice()
	if rl.IsAudioDeviceReady() {
		defer rl.CloseAudioDevice()
	}

	gui := viewer{controller: controller, families: controller.Families()}
	gui.audio = newAudioPlayer(controller.Voicing())
	defer gui.audio.close()

	for !gui.quit && !rl.WindowShouldClose() {
		layout := computeLayout(rl.GetScreenWidth(), rl.GetScreenHeight(), gui.fullNeck)
		gui.update(layout)
		gui.audio.update()
		gui.draw(layout)
	}
	return nil
}

// computeLayout keeps compact graphs beside the diagram and full-neck panels stacked.
func computeLayout(width, height int, fullNeck bool) screenLayout {
	w := float32(width)
	h := float32(height)
	margin := float32(panelGap)
	headerHeight := clamp(h*0.105, 78, 102)
	dialHeight := clamp(h*0.145, 105, 140)
	header := rl.Rectangle{X: margin, Y: margin, Width: w - 2*margin, Height: headerHeight}
	dial := rl.Rectangle{
		X: margin, Y: header.Y + header.Height + margin,
		Width: w - 2*margin, Height: dialHeight,
	}
	contentY := dial.Y + dial.Height + margin
	contentHeight := maxFloat(120, h-contentY-margin)

	if !fullNeck {
		leftWidth := (w - 3*margin) * 0.54
		rightWidth := w - 3*margin - leftWidth
		waveHeight := (contentHeight - margin) * 0.58
		return screenLayout{
			header:  header,
			dial:    dial,
			diagram: rl.Rectangle{X: margin, Y: contentY, Width: leftWidth, Height: contentHeight},
			waveform: rl.Rectangle{
				X: margin + leftWidth + margin, Y: contentY,
				Width: rightWidth, Height: waveHeight,
			},
			spectrum: rl.Rectangle{
				X: margin + leftWidth + margin, Y: contentY + waveHeight + margin,
				Width: rightWidth, Height: contentHeight - waveHeight - margin,
			},
		}
	}

	diagramHeight := contentHeight * 0.42
	waveHeight := contentHeight * 0.30
	return screenLayout{
		header:  header,
		dial:    dial,
		diagram: rl.Rectangle{X: margin, Y: contentY, Width: w - 2*margin, Height: diagramHeight},
		waveform: rl.Rectangle{
			X: margin, Y: contentY + diagramHeight + margin,
			Width: w - 2*margin, Height: waveHeight - margin,
		},
		spectrum: rl.Rectangle{
			X: margin, Y: contentY + diagramHeight + waveHeight + margin,
			Width: w - 2*margin, Height: contentHeight - diagramHeight - waveHeight - margin,
		},
	}
}

// update maps keyboard and chord-row clicks onto shared navigation state.
func (gui *viewer) update(layout screenLayout) {
	selectionChanged := false
	if keyPressed(rl.KeyLeft, rl.KeyH) {
		gui.controller.MoveChord(-1)
		selectionChanged = true
	}
	if keyPressed(rl.KeyRight, rl.KeyL) {
		gui.controller.MoveChord(1)
		selectionChanged = true
	}
	if keyPressed(rl.KeyUp, rl.KeyK) {
		before := gui.controller.Name()
		gui.controller.MoveKind(-1)
		selectionChanged = selectionChanged || before != gui.controller.Name()
	}
	if keyPressed(rl.KeyDown, rl.KeyJ) {
		before := gui.controller.Name()
		gui.controller.MoveKind(1)
		selectionChanged = selectionChanged || before != gui.controller.Name()
	}
	if rl.IsKeyPressed(rl.KeyV) {
		gui.controller.CycleVoicing()
		selectionChanged = true
	}
	if rl.IsKeyPressed(rl.KeyF) {
		gui.fullNeck = !gui.fullNeck
		if gui.fullNeck {
			gui.tabNumbers = false
		}
	}
	if rl.IsKeyPressed(rl.KeyN) {
		gui.tabNumbers = !gui.tabNumbers
		if gui.tabNumbers {
			gui.fullNeck = false
		}
	}
	if rl.IsKeyPressed(rl.KeyT) {
		gui.theme = app.Wrap(gui.theme+1, len(palettes))
	}
	if rl.IsKeyPressed(rl.KeySpace) || rl.IsKeyPressed(rl.KeyP) {
		gui.audio.toggle()
	}
	if rl.IsKeyPressed(rl.KeySlash) || rl.IsKeyPressed(rl.KeyF1) {
		gui.help = !gui.help
	}
	if rl.IsKeyPressed(rl.KeyEscape) && gui.help {
		gui.help = false
	} else if rl.IsKeyPressed(rl.KeyQ) || rl.IsKeyPressed(rl.KeyEscape) {
		gui.quit = true
	}

	if gui.handleDialClick(layout.dial) {
		selectionChanged = true
	}
	if selectionChanged {
		gui.audio.setVoicing(gui.controller.Voicing())
	}
}

// keyPressed accepts both an arrow and its Vim-style alternative.
func keyPressed(primary, alternate int32) bool {
	return rl.IsKeyPressed(primary) || rl.IsKeyPressed(alternate)
}

// handleDialClick selects a base chord from the visible row.
func (gui *viewer) handleDialClick(bounds rl.Rectangle) bool {
	if !rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		return false
	}
	row := baseChordRow(bounds, len(gui.families))
	mouse := rl.GetMousePosition()
	for index, cell := range row {
		if rl.CheckCollisionPointRec(mouse, cell) {
			gui.controller.MoveChord(index - gui.controller.SelectedIndex())
			return true
		}
	}
	return false
}

// draw renders one complete immediate-mode frame.
func (gui *viewer) draw(layout screenLayout) {
	theme := palettes[gui.theme]
	rl.BeginDrawing()
	defer rl.EndDrawing()
	rl.ClearBackground(theme.background)

	gui.drawHeader(layout.header, theme)
	gui.drawDial(layout.dial, theme)
	gui.drawDiagram(layout.diagram, theme)
	gui.drawWaveform(layout.waveform, theme)
	gui.drawSpectrum(layout.spectrum, theme)
	if gui.help {
		gui.drawHelp(theme)
	}
}

// drawHeader combines application identity, mode, and controls.
func (gui *viewer) drawHeader(bounds rl.Rectangle, theme palette) {
	drawPanel(bounds, theme)
	// The small status lamp and recessed controls echo the original console shell.
	rl.DrawCircleV(rl.Vector2{X: bounds.X + 17, Y: bounds.Y + 25}, 4, theme.secondary)
	drawText("ALT-TAB", bounds.X+38, bounds.Y+13, 28, theme.accent)
	drawText("RAYLIB CHORD VIEWER - "+strings.ToUpper(theme.name),
		bounds.X+38, bounds.Y+46, 14, theme.muted)
	controls := "LEFT/RIGHT chord  UP/DOWN type  V voicing  F neck  N tab  T theme  SPACE play  F1 help  Q quit"
	size := measureText(controls, 14)
	controlWell := rl.Rectangle{
		X: bounds.X + bounds.Width - panelPadding - size.X - 10,
		Y: bounds.Y + 13, Width: size.X + 20, Height: 30,
	}
	rl.DrawRectangleRounded(controlWell, 0.35, 8, rl.Fade(theme.accent, 0.18))
	drawText(controls, bounds.X+bounds.Width-panelPadding-size.X,
		bounds.Y+20, 14, theme.text)
}

// drawDial renders clickable base chords and the current family variants.
func (gui *viewer) drawDial(bounds rl.Rectangle, theme palette) {
	drawPanel(bounds, theme)
	drawSectionTitle(bounds, "CHORD DIAL", theme)
	for index, cell := range baseChordRow(bounds, len(gui.families)) {
		color := theme.muted
		if index == gui.controller.SelectedIndex() {
			color = theme.accent
			rl.DrawRectangleRounded(cell, 0.35, 8, rl.Fade(theme.accent, 0.12))
		}
		drawCenteredText(gui.families[index].Base, cell, 20, color)
	}

	if len(gui.families) == 0 {
		return
	}
	family := gui.families[gui.controller.SelectedIndex()]
	center := rl.Rectangle{
		X:      bounds.X + bounds.Width/2 - 62,
		Y:      bounds.Y + bounds.Height - 40,
		Width:  124,
		Height: 24,
	}
	current := gui.controller.Name()
	drawCenteredText("<  "+current+"  >", center, 22, theme.text)
	if family.Accidental != "" {
		drawCenteredText(family.Accidental, rl.Rectangle{
			X: center.X, Y: center.Y - 20, Width: center.Width, Height: 18,
		}, 15, theme.secondary)
	}
	if family.Minor != "" {
		drawCenteredText(family.Minor, rl.Rectangle{
			X: center.X, Y: center.Y + center.Height, Width: center.Width, Height: 18,
		}, 15, theme.secondary)
	}
}

// baseChordRow returns stable mouse and drawing cells for all natural chords.
func baseChordRow(bounds rl.Rectangle, count int) []rl.Rectangle {
	if count == 0 {
		return nil
	}
	cellWidth := float32(54)
	total := cellWidth * float32(count)
	start := bounds.X + (bounds.Width-total)/2
	row := make([]rl.Rectangle, count)
	for index := range row {
		row[index] = rl.Rectangle{
			X:      start + float32(index)*cellWidth,
			Y:      bounds.Y + 42,
			Width:  cellWidth,
			Height: 32,
		}
	}
	return row
}

// drawDiagram renders a scalable horizontal guitar fretboard.
func (gui *viewer) drawDiagram(bounds rl.Rectangle, theme palette) {
	drawPanel(bounds, theme)
	mode := "COMPACT"
	if gui.fullNeck {
		mode = "FULL NECK"
	} else if gui.tabNumbers {
		mode = "TAB NUMBERS"
	}
	drawSectionTitle(bounds, "CHORD DIAGRAM · "+mode, theme)
	voicing := gui.controller.Voicing()
	heading := fmt.Sprintf("%s   -   VOICING %d OF %d",
		voicing.Name, voicing.Number, gui.controller.VoicingCount())
	drawCenteredText(heading, rl.Rectangle{
		X: bounds.X, Y: bounds.Y + 35, Width: bounds.Width, Height: 28,
	}, 20, theme.text)

	first, last := fretRange(voicing, gui.fullNeck)
	board := rl.Rectangle{
		X:      bounds.X + 54,
		Y:      bounds.Y + 82,
		Width:  bounds.Width - 82,
		Height: bounds.Height - 122,
	}
	if board.Height < 75 {
		return
	}
	drawFretboard(board, voicing, first, last, gui.tabNumbers, theme)
	legend := "O open   X muted"
	if !gui.tabNumbers {
		legend = "Fingers: 1 index   2 middle   3 ring   4 little     -     " + legend
	}
	drawCenteredText(legend, rl.Rectangle{
		X: bounds.X, Y: bounds.Y + bounds.Height - 29,
		Width: bounds.Width, Height: 20,
	}, 14, theme.muted)
}

// fretRange keeps common positions readable while preserving exact labels.
func fretRange(voicing chords.Voicing, fullNeck bool) (int, int) {
	if fullNeck {
		return 1, fullNeckLastFret
	}
	lowest := 0
	highest := 0
	for _, placement := range voicing.Strings {
		if placement.Fret > 0 && (lowest == 0 || placement.Fret < lowest) {
			lowest = placement.Fret
		}
		if placement.Fret > highest {
			highest = placement.Fret
		}
	}
	if highest <= 4 {
		return 1, 4
	}
	return maxInt(1, lowest), maxInt(lowest+3, highest)
}

// drawFretboard draws strings, frets, labels, and finger placements.
func drawFretboard(
	bounds rl.Rectangle,
	voicing chords.Voicing,
	first, last int,
	tabNumbers bool,
	theme palette,
) {
	stringNames := [...]string{"e", "B", "G", "D", "A", "E"}
	fretCount := last - first + 1
	cellWidth := bounds.Width / float32(fretCount)
	stringGap := bounds.Height / float32(chords.StringCount+1)

	for fret := first; fret <= last; fret++ {
		x := bounds.X + float32(fret-first)*cellWidth
		rl.DrawLineEx(rl.Vector2{X: x, Y: bounds.Y + stringGap},
			rl.Vector2{X: x, Y: bounds.Y + stringGap*6}, 1, theme.border)
		label := fmt.Sprintf("%d", fret)
		drawCenteredText(label, rl.Rectangle{
			X: x, Y: bounds.Y - 2, Width: cellWidth, Height: stringGap,
		}, clamp(cellWidth*0.36, 10, 15), theme.muted)
	}
	right := bounds.X + bounds.Width
	rl.DrawLineEx(rl.Vector2{X: right, Y: bounds.Y + stringGap},
		rl.Vector2{X: right, Y: bounds.Y + stringGap*6}, 1, theme.border)

	for index, placement := range voicing.Strings {
		y := bounds.Y + stringGap*float32(index+1)
		drawText(stringNames[index], bounds.X-32, y-8, 16, theme.muted)
		rl.DrawLineEx(rl.Vector2{X: bounds.X, Y: y}, rl.Vector2{X: right, Y: y}, 2, theme.text)
		if placement.Fret <= 0 {
			marker := "O"
			if placement.Fret < 0 {
				marker = "X"
			}
			drawText(marker, bounds.X-15, y-9, 17, theme.secondary)
			continue
		}
		if placement.Fret < first || placement.Fret > last {
			continue
		}
		x := bounds.X + (float32(placement.Fret-first)+0.5)*cellWidth
		radius := clamp(minFloat(cellWidth*0.27, stringGap*0.30), 7, 14)
		rl.DrawCircleV(rl.Vector2{X: x, Y: y}, radius+5, rl.Fade(theme.accent, 0.10))
		rl.DrawCircleV(rl.Vector2{X: x, Y: y}, radius, theme.accent)
		label := fmt.Sprintf("%d", placement.Finger)
		if tabNumbers {
			label = fmt.Sprintf("%d", placement.Fret)
		}
		drawCenteredText(label, rl.Rectangle{
			X: x - radius, Y: y - radius, Width: radius * 2, Height: radius * 2,
		}, clamp(radius*1.15, 10, 16), theme.panel)
	}
}

// drawWaveform plots 25 milliseconds of the exact sounding-string signal.
func (gui *viewer) drawWaveform(bounds rl.Rectangle, theme palette) {
	drawPanel(bounds, theme)
	drawSectionTitle(bounds, "WAVEFORM · AMPLITUDE / TIME", theme)
	plot := inset(bounds, panelPadding, 48, panelPadding, 38)
	if plot.Width < 20 || plot.Height < 20 {
		return
	}
	notes := signal.Notes(gui.controller.Voicing())
	center := plot.Y + plot.Height/2
	rl.DrawLineEx(rl.Vector2{X: plot.X, Y: center},
		rl.Vector2{X: plot.X + plot.Width, Y: center}, 1, theme.border)
	points := make([]rl.Vector2, maxInt(2, int(plot.Width)))
	for index := range points {
		t := float64(index) / float64(len(points)-1)
		sample := signal.CompositeSample(t*waveformSeconds, notes)
		points[index] = rl.Vector2{
			X: plot.X + float32(t)*plot.Width,
			Y: center - float32(sample)*plot.Height*0.43,
		}
	}
	rl.DrawLineStrip(points, rl.Fade(theme.accent, 0.22))
	for index := 1; index < len(points); index++ {
		rl.DrawLineEx(points[index-1], points[index], 2, theme.accent)
	}
	drawText("+1", plot.X-2, plot.Y-14, 12, theme.muted)
	drawText("-1", plot.X-2, plot.Y+plot.Height-10, 12, theme.muted)
	drawText("0 ms", plot.X, plot.Y+plot.Height+3, 12, theme.muted)
	right := "25 ms"
	drawText(right, plot.X+plot.Width-measureText(right, 12).X,
		plot.Y+plot.Height+3, 12, theme.muted)
}

// drawSpectrum places one exact peak per sounding pitch on a logarithmic scale.
func (gui *viewer) drawSpectrum(bounds rl.Rectangle, theme palette) {
	drawPanel(bounds, theme)
	drawSectionTitle(bounds, "FREQUENCY SPECTRUM · LOG Hz", theme)
	plot := inset(bounds, panelPadding, 48, panelPadding, 44)
	if plot.Width < 20 || plot.Height < 20 {
		return
	}
	notes := append([]signal.Note(nil), signal.Notes(gui.controller.Voicing())...)
	sort.Slice(notes, func(left, right int) bool { return notes[left].MIDI < notes[right].MIDI })
	if len(notes) == 0 {
		return
	}
	peaks := aggregateNotes(notes)
	low := peaks[0].note.MIDI
	high := peaks[len(peaks)-1].note.MIDI
	span := maxInt(1, high-low)
	axisY := plot.Y + plot.Height
	rl.DrawLineEx(rl.Vector2{X: plot.X, Y: axisY},
		rl.Vector2{X: plot.X + plot.Width, Y: axisY}, 1, theme.border)
	maxCount := 1
	for _, peak := range peaks {
		maxCount = maxInt(maxCount, peak.count)
	}
	for _, peak := range peaks {
		t := float32(peak.note.MIDI-low) / float32(span)
		x := plot.X + t*plot.Width
		height := plot.Height * float32(peak.count) / float32(maxCount)
		top := axisY - height
		rl.DrawLineEx(rl.Vector2{X: x, Y: top}, rl.Vector2{X: x, Y: axisY}, 4, theme.secondary)
		rl.DrawCircleV(rl.Vector2{X: x, Y: top}, 7, theme.accent)
		drawCenteredText(peak.note.Name, rl.Rectangle{
			X: x - 25, Y: axisY + 4, Width: 50, Height: 18,
		}, 12, theme.text)
	}
	legendNames := make([]string, len(notes))
	for index, note := range notes {
		legendNames[index] = note.Name
	}
	legend := "Notes: " + strings.Join(legendNames, "  ")
	drawCenteredText(legend, rl.Rectangle{
		X: bounds.X, Y: bounds.Y + bounds.Height - 22,
		Width: bounds.Width, Height: 16,
	}, 12, theme.muted)
}

// spectrumPeak combines duplicate pitches so its height matches oscillator amplitude.
type spectrumPeak struct {
	note  signal.Note
	count int
}

// aggregateNotes counts unison strings in an already pitch-sorted note slice.
func aggregateNotes(notes []signal.Note) []spectrumPeak {
	peaks := make([]spectrumPeak, 0, len(notes))
	for _, note := range notes {
		last := len(peaks) - 1
		if last >= 0 && peaks[last].note.MIDI == note.MIDI {
			peaks[last].count++
			continue
		}
		peaks = append(peaks, spectrumPeak{note: note, count: 1})
	}
	return peaks
}

// drawHelp overlays the complete graphical control reference.
func (gui *viewer) drawHelp(theme palette) {
	width := float32(520)
	height := float32(390)
	bounds := rl.Rectangle{
		X:      (float32(rl.GetScreenWidth()) - width) / 2,
		Y:      (float32(rl.GetScreenHeight()) - height) / 2,
		Width:  width,
		Height: height,
	}
	rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()),
		rl.Fade(rl.Black, 0.55))
	drawPanel(bounds, theme)
	drawCenteredText("KEYBOARD HELP", rl.Rectangle{
		X: bounds.X, Y: bounds.Y + 22, Width: bounds.Width, Height: 30,
	}, 25, theme.accent)
	lines := []string{
		"← / H    Previous chord", "→ / L    Next chord",
		"↑ / K    Accidental or base", "↓ / J    Minor or base",
		"V        Next voicing", "F        Compact / full neck",
		"N        Finger / tab numbers", "T        Cycle theme",
		"SPACE    Play / stop chord", "F1 / ?   Close this help",
		"Q / Esc  Quit",
	}
	y := bounds.Y + 75
	for _, line := range lines {
		drawText(line, bounds.X+70, y, 18, theme.text)
		y += 26
	}
}

// drawPanel paints the common rounded module background and border.
func drawPanel(bounds rl.Rectangle, theme palette) {
	shadow := bounds
	shadow.X += 3
	shadow.Y += 3
	rl.DrawRectangleRounded(shadow, 0.045, 12, rl.Fade(rl.Black, 0.22))
	rl.DrawRectangleRounded(bounds, 0.045, 12, theme.panel)
	rl.DrawRectangleRoundedLinesEx(bounds, 0.045, 12, 2, theme.border)
	// A subtle upper highlight makes each module read like molded plastic.
	rl.DrawLineEx(
		rl.Vector2{X: bounds.X + 14, Y: bounds.Y + 5},
		rl.Vector2{X: bounds.X + bounds.Width - 14, Y: bounds.Y + 5},
		1,
		rl.Fade(rl.White, 0.32),
	)
}

// drawSectionTitle positions one concise label at a panel's upper-left edge.
func drawSectionTitle(bounds rl.Rectangle, title string, theme palette) {
	drawText(title, bounds.X+panelPadding, bounds.Y+14, 16, theme.accent)
	underlineWidth := minFloat(measureText(title, 16).X, 48)
	rl.DrawLineEx(
		rl.Vector2{X: bounds.X + panelPadding, Y: bounds.Y + 34},
		rl.Vector2{X: bounds.X + panelPadding + underlineWidth, Y: bounds.Y + 34},
		3,
		theme.secondary,
	)
}

// drawText adds a tight second pass for a heavier console-style default font.
func drawText(text string, x, y, size float32, color rl.Color) {
	font := rl.GetFontDefault()
	position := rl.Vector2{X: x, Y: y}
	rl.DrawTextEx(font, text, position, size, 1, color)
	position.X += clamp(size*0.045, 0.5, 0.8)
	rl.DrawTextEx(font, text, position, size, 1, color)
}

// drawCenteredText centers one label within a rectangular region.
func drawCenteredText(text string, bounds rl.Rectangle, size float32, color rl.Color) {
	measured := measureText(text, size)
	drawText(text,
		bounds.X+(bounds.Width-measured.X)/2,
		bounds.Y+(bounds.Height-measured.Y)/2,
		size,
		color,
	)
}

// measureText returns default-font dimensions matching drawText.
func measureText(text string, size float32) rl.Vector2 {
	return rl.MeasureTextEx(rl.GetFontDefault(), text, size, 1)
}

// inset removes independent padding from each rectangle side.
func inset(bounds rl.Rectangle, left, top, right, bottom float32) rl.Rectangle {
	return rl.Rectangle{
		X:      bounds.X + left,
		Y:      bounds.Y + top,
		Width:  maxFloat(1, bounds.Width-left-right),
		Height: maxFloat(1, bounds.Height-top-bottom),
	}
}

// clamp restricts one floating-point layout value to an inclusive range.
func clamp(value, minimum, maximum float32) float32 {
	return minFloat(maxFloat(value, minimum), maximum)
}

// minFloat returns the smaller floating-point value.
func minFloat(left, right float32) float32 {
	if left < right {
		return left
	}
	return right
}

// maxFloat returns the larger floating-point value.
func maxFloat(left, right float32) float32 {
	if left > right {
		return left
	}
	return right
}

// maxInt returns the larger integer value.
func maxInt(left, right int) int {
	if left > right {
		return left
	}
	return right
}

// audioPlayer continuously fills one Raylib stream with chord oscillators.
type audioPlayer struct {
	stream     rl.AudioStream
	ready      bool
	active     bool
	notes      []signal.Note
	phases     []float64
	gain       float64
	targetGain float64
	buffer     []float32
}

// newAudioPlayer initializes a silent stream when the audio device is ready.
func newAudioPlayer(voicing chords.Voicing) audioPlayer {
	player := audioPlayer{}
	if !rl.IsAudioDeviceReady() {
		return player
	}
	rl.SetAudioStreamBufferSizeDefault(audioBufferFrames)
	player.stream = rl.LoadAudioStream(audioSampleRate, 32, 1)
	if !rl.IsAudioStreamValid(player.stream) {
		return player
	}
	player.ready = true
	player.buffer = make([]float32, audioBufferFrames)
	player.setVoicing(voicing)
	rl.SetAudioStreamVolume(player.stream, 0.65)
	rl.PlayAudioStream(player.stream)
	return player
}

// setVoicing updates oscillator frequencies without reallocating audio buffers.
func (player *audioPlayer) setVoicing(voicing chords.Voicing) {
	player.notes = signal.Notes(voicing)
	player.phases = make([]float64, len(player.notes))
}

// toggle fades synthesized playback in or out to avoid clicks.
func (player *audioPlayer) toggle() {
	if !player.ready {
		return
	}
	player.active = !player.active
	player.targetGain = 0
	if player.active {
		player.targetGain = 0.28
	}
}

// update refills every processed half-buffer from the main thread.
func (player *audioPlayer) update() {
	if !player.ready || len(player.notes) == 0 {
		return
	}
	for rl.IsAudioStreamProcessed(player.stream) {
		for frame := range player.buffer {
			// The short per-sample ramp prevents discontinuities during toggles.
			player.gain += (player.targetGain - player.gain) * 0.0025
			var sample float64
			for index, note := range player.notes {
				sample += math.Sin(player.phases[index])
				player.phases[index] += 2 * math.Pi * note.Frequency / audioSampleRate
				if player.phases[index] >= 2*math.Pi {
					player.phases[index] -= 2 * math.Pi
				}
			}
			player.buffer[frame] = float32(sample / float64(len(player.notes)) * player.gain)
		}
		rl.UpdateAudioStream(player.stream, player.buffer)
	}
}

// close releases the native stream before the audio device shuts down.
func (player *audioPlayer) close() {
	if !player.ready {
		return
	}
	rl.StopAudioStream(player.stream)
	rl.UnloadAudioStream(player.stream)
}
