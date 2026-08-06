//go:build raylib

// Package rayui renders the alternate Raylib interface.
package rayui

import (
	_ "embed"
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
	windowWidth               = 1280
	windowHeight              = 900
	minimumWindowWidth        = 960
	minimumWindowHeight       = 800
	panelGap                  = 12
	panelPadding              = 18
	diagramLeftPadding        = 88
	diagramRightPadding       = 28
	graphTopPadding           = 60
	waveLeftPadding           = 50
	waveRightPadding          = 32
	waveBottomPadding         = 42
	spectrumSidePadding       = 46
	spectrumBottomPad         = 62
	bodyTextSpacing           = 1.7
	audioSampleRate           = 44100
	audioBufferFrames         = 2048
	waveformSeconds           = 0.025
	fullNeckLastFret          = 27
	idleFrameRate       int32 = 30
	activeFrameRate           = 60
	uiFontSize                = 64
)

// boldFontData embeds the UI font so release builds remain self-contained.
//
//go:embed assets/SpaceMono-Bold.ttf
var boldFontData []byte

// uiFont is loaded after Raylib creates its graphics context.
var uiFont rl.Font

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
	signal     rl.Color
}

// palettes contains the approved graphical themes in runtime cycling order.
var palettes = [...]palette{
	{
		name: "Super Famicom", background: rl.NewColor(201, 197, 189, 255),
		panel: rl.NewColor(41, 37, 48, 255), border: rl.NewColor(98, 90, 112, 255),
		text: rl.NewColor(244, 240, 233, 255), muted: rl.NewColor(183, 173, 185, 255),
		accent: rl.NewColor(122, 74, 160, 255), secondary: rl.NewColor(208, 90, 146, 255),
		signal: rl.NewColor(232, 184, 63, 255),
	},
	{
		name: "Atomic Grape", background: rl.NewColor(22, 16, 31, 255),
		panel: rl.NewColor(39, 19, 56, 255), border: rl.NewColor(113, 67, 154, 255),
		text: rl.NewColor(251, 244, 255, 255), muted: rl.NewColor(184, 153, 204, 255),
		accent: rl.NewColor(194, 108, 255, 255), secondary: rl.NewColor(255, 79, 184, 255),
		signal: rl.NewColor(128, 243, 255, 255),
	},
	{
		name: "Paper Terminal", background: rl.NewColor(222, 211, 184, 255),
		panel: rl.NewColor(243, 236, 217, 255), border: rl.NewColor(39, 38, 41, 255),
		text: rl.NewColor(31, 29, 27, 255), muted: rl.NewColor(109, 102, 90, 255),
		accent: rl.NewColor(191, 63, 50, 255), secondary: rl.NewColor(37, 92, 122, 255),
		signal: rl.NewColor(32, 32, 32, 255),
	},
	{
		name: "Glacier Circuit", background: rl.NewColor(203, 219, 227, 255),
		panel: rl.NewColor(20, 40, 53, 255), border: rl.NewColor(68, 115, 138, 255),
		text: rl.NewColor(238, 251, 255, 255), muted: rl.NewColor(155, 187, 200, 255),
		accent: rl.NewColor(139, 131, 255, 255), secondary: rl.NewColor(85, 217, 236, 255),
		signal: rl.NewColor(197, 244, 255, 255),
	},
	{
		name: "Haunted Cartridge", background: rl.NewColor(23, 27, 24, 255),
		panel: rl.NewColor(35, 40, 37, 255), border: rl.NewColor(83, 107, 85, 255),
		text: rl.NewColor(232, 242, 223, 255), muted: rl.NewColor(148, 165, 139, 255),
		accent: rl.NewColor(181, 107, 228, 255), secondary: rl.NewColor(112, 219, 121, 255),
		signal: rl.NewColor(211, 255, 118, 255),
	},
	{
		name: "Oxide Industrial", background: rl.NewColor(74, 75, 64, 255),
		panel: rl.NewColor(25, 29, 28, 255), border: rl.NewColor(122, 117, 91, 255),
		text: rl.NewColor(235, 229, 207, 255), muted: rl.NewColor(170, 165, 143, 255),
		accent: rl.NewColor(210, 105, 53, 255), secondary: rl.NewColor(224, 180, 63, 255),
		signal: rl.NewColor(126, 177, 162, 255),
	},
	{
		name: "Cassette Future", background: rl.NewColor(23, 38, 59, 255),
		panel: rl.NewColor(11, 20, 34, 255), border: rl.NewColor(53, 90, 119, 255),
		text: rl.NewColor(232, 242, 245, 255), muted: rl.NewColor(138, 166, 179, 255),
		accent: rl.NewColor(89, 213, 216, 255), secondary: rl.NewColor(255, 113, 133, 255),
		signal: rl.NewColor(181, 156, 255, 255),
	},
	{
		name: "Royal Terminal", background: rl.NewColor(28, 25, 33, 255),
		panel: rl.NewColor(14, 13, 17, 255), border: rl.NewColor(84, 70, 95, 255),
		text: rl.NewColor(245, 237, 219, 255), muted: rl.NewColor(183, 169, 140, 255),
		accent: rl.NewColor(145, 102, 197, 255), secondary: rl.NewColor(213, 170, 81, 255),
		signal: rl.NewColor(244, 216, 134, 255),
	},
	{
		name: "Sakura Console", background: rl.NewColor(234, 219, 213, 255),
		panel: rl.NewColor(67, 44, 61, 255), border: rl.NewColor(137, 90, 119, 255),
		text: rl.NewColor(255, 245, 241, 255), muted: rl.NewColor(216, 185, 197, 255),
		accent: rl.NewColor(240, 144, 186, 255), secondary: rl.NewColor(159, 103, 192, 255),
		signal: rl.NewColor(255, 199, 217, 255),
	},
	{
		name: "CRT Amber", background: rl.NewColor(53, 44, 36, 255),
		panel: rl.NewColor(16, 14, 11, 255), border: rl.NewColor(106, 81, 52, 255),
		text: rl.NewColor(255, 217, 138, 255), muted: rl.NewColor(169, 131, 76, 255),
		accent: rl.NewColor(255, 181, 47, 255), secondary: rl.NewColor(239, 111, 46, 255),
		signal: rl.NewColor(255, 225, 161, 255),
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
	signal     signalRenderCache
}

// signalRenderCache retains chord-derived graph data until selection or size changes.
type signalRenderCache struct {
	name       string
	number     int
	notes      []signal.Note
	peaks      []spectrumPeak
	legend     string
	waveBounds rl.Rectangle
	wavePoints []rl.Vector2
}

// Run initializes Raylib and serves frames until the window closes.
func Run(catalog chords.Catalog) error {
	controller := app.NewController(catalog)
	if controller.Err() != nil {
		return controller.Err()
	}

	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagVsyncHint)
	rl.InitWindow(windowWidth, windowHeight, "Alt-Tab · Raylib")
	if !rl.IsWindowReady() {
		return errors.New("raylib window initialization failed")
	}
	defer rl.CloseWindow()
	rl.SetWindowMinSize(minimumWindowWidth, minimumWindowHeight)
	rl.SetExitKey(rl.KeyNull)
	if err := loadUIFont(); err != nil {
		return err
	}
	defer rl.UnloadFont(uiFont)
	rl.SetTargetFPS(idleFrameRate)

	rl.InitAudioDevice()
	if rl.IsAudioDeviceReady() {
		defer rl.CloseAudioDevice()
	}

	gui := viewer{controller: controller, families: controller.Families()}
	gui.audio = newAudioPlayer(controller.Voicing())
	defer gui.audio.close()
	frameRate := idleFrameRate

	for !gui.quit && !rl.WindowShouldClose() {
		layout := computeLayout(rl.GetScreenWidth(), rl.GetScreenHeight(), gui.fullNeck)
		gui.update(layout)
		gui.audio.update()
		desiredFrameRate := idleFrameRate
		if gui.audio.audible() {
			desiredFrameRate = activeFrameRate
		}
		if desiredFrameRate != frameRate {
			rl.SetTargetFPS(desiredFrameRate)
			frameRate = desiredFrameRate
		}
		gui.draw(layout)
	}
	return nil
}

// loadUIFont creates the single bold typeface used for all UI text.
func loadUIFont() error {
	uiFont = rl.LoadFontFromMemory(".ttf", boldFontData, uiFontSize, uiGlyphs())
	if !rl.IsFontValid(uiFont) {
		return errors.New("embedded UI font initialization failed")
	}
	rl.SetTextureFilter(uiFont.Texture, rl.FilterBilinear)
	return nil
}

// uiGlyphs limits the font atlas to the characters used by the interface.
func uiGlyphs() []rune {
	glyphs := make([]rune, 0, 102)
	for character := rune(32); character <= 126; character++ {
		glyphs = append(glyphs, character)
	}
	return append(glyphs, '·', '←', '→', '↑', '↓')
}

// computeLayout keeps compact graphs beside the diagram and full-neck panels stacked.
func computeLayout(width, height int, fullNeck bool) screenLayout {
	w := float32(width)
	h := float32(height)
	margin := float32(panelGap)
	headerHeight := clamp(h*0.105, 78, 102)
	dialHeight := clamp(h*0.19, 160, 180)
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

	diagramHeight := contentHeight * 0.41
	waveHeight := contentHeight * 0.28
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
	rl.DrawCircleV(rl.Vector2{X: bounds.X + 17, Y: bounds.Y + 27}, 4, theme.secondary)
	drawTextSpaced("ALT-TAB", bounds.X+38, bounds.Y+9, 34, 3.5, theme.accent)
	drawText("RAYLIB CHORD VIEWER - "+strings.ToUpper(theme.name),
		bounds.X+38, bounds.Y+53, 15, theme.muted)
	controls := "LEFT/RIGHT chord  UP/DOWN type  V voicing  F neck  N tab  T theme  SPACE play  F1 help  Q quit"
	size := measureText(controls, 15)
	controlWell := rl.Rectangle{
		X: bounds.X + bounds.Width - panelPadding - size.X - 10,
		Y: bounds.Y + 13, Width: size.X + 20, Height: 30,
	}
	rl.DrawRectangleRounded(controlWell, 0.35, 8, rl.Fade(theme.accent, 0.18))
	drawText(controls, bounds.X+bounds.Width-panelPadding-size.X,
		bounds.Y+19, 15, theme.text)
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
		Y:      bounds.Y + bounds.Height - 65,
		Width:  124,
		Height: 24,
	}
	current := gui.controller.Name()
	drawCenteredText("<  "+current+"  >", center, 22, theme.text)
	if family.Accidental != "" {
		drawCenteredText(family.Accidental, rl.Rectangle{
			X: center.X, Y: center.Y - 22, Width: center.Width, Height: 18,
		}, 17, theme.secondary)
	}
	if family.Minor != "" {
		drawCenteredText(family.Minor, rl.Rectangle{
			X: center.X, Y: center.Y + center.Height, Width: center.Width, Height: 18,
		}, 17, theme.secondary)
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
			Y:      bounds.Y + 40,
			Width:  cellWidth,
			Height: 28,
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
		X: bounds.X, Y: bounds.Y + 44, Width: bounds.Width, Height: 24,
	}, 20, theme.text)

	first, last := fretRange(voicing, gui.fullNeck)
	boardBottomPadding := float32(96)
	if !gui.fullNeck && !gui.tabNumbers {
		boardBottomPadding = 116
	}
	board := rl.Rectangle{
		X:      bounds.X + diagramLeftPadding,
		Y:      bounds.Y + 74,
		Width:  bounds.Width - diagramLeftPadding - diagramRightPadding,
		Height: bounds.Height - boardBottomPadding,
	}
	if board.Height < 75 {
		return
	}
	drawFretboard(board, voicing, first, last, gui.tabNumbers, theme)
	gui.drawDiagramLegend(bounds, theme)
}

// drawDiagramLegend gives enlarged finger and symbol metadata enough room.
func (gui *viewer) drawDiagramLegend(bounds rl.Rectangle, theme palette) {
	const legendSize = 16
	symbols := "O open   X muted"
	if gui.tabNumbers {
		drawCenteredText(symbols, rl.Rectangle{
			X: bounds.X, Y: bounds.Y + bounds.Height - 30,
			Width: bounds.Width, Height: 20,
		}, legendSize, theme.muted)
		return
	}
	fingers := "Fingers: 1 index   2 middle   3 ring   4 little"
	if gui.fullNeck {
		drawCenteredText(fingers+"     -     "+symbols, rl.Rectangle{
			X: bounds.X, Y: bounds.Y + bounds.Height - 30,
			Width: bounds.Width, Height: 20,
		}, legendSize, theme.muted)
		return
	}
	drawCenteredText(fingers, rl.Rectangle{
		X: bounds.X, Y: bounds.Y + bounds.Height - 52,
		Width: bounds.Width, Height: 20,
	}, legendSize, theme.muted)
	drawCenteredText(symbols, rl.Rectangle{
		X: bounds.X, Y: bounds.Y + bounds.Height - 28,
		Width: bounds.Width, Height: 20,
	}, legendSize, theme.muted)
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
		drawHeavyCenteredText(label, rl.Rectangle{
			X: x, Y: bounds.Y - 2, Width: cellWidth, Height: stringGap,
		}, clamp(cellWidth*0.55, 17, 24), theme.text)
	}
	right := bounds.X + bounds.Width
	rl.DrawLineEx(rl.Vector2{X: right, Y: bounds.Y + stringGap},
		rl.Vector2{X: right, Y: bounds.Y + stringGap*6}, 1, theme.border)

	for index, placement := range voicing.Strings {
		y := bounds.Y + stringGap*float32(index+1)
		drawText(stringNames[index], bounds.X-58, y-9, 18, theme.muted)
		rl.DrawLineEx(rl.Vector2{X: bounds.X, Y: y}, rl.Vector2{X: right, Y: y}, 2, theme.text)
		if placement.Fret <= 0 {
			marker := "O"
			if placement.Fret < 0 {
				marker = "X"
			}
			drawText(marker, bounds.X-24, y-10, 19, theme.secondary)
			continue
		}
		if placement.Fret < first || placement.Fret > last {
			continue
		}
		x := bounds.X + (float32(placement.Fret-first)+0.5)*cellWidth
		radius := clamp(minFloat(cellWidth*0.30, stringGap*0.40), 8, 16)
		rl.DrawCircleV(rl.Vector2{X: x, Y: y}, radius+5, rl.Fade(theme.accent, 0.10))
		rl.DrawCircleV(rl.Vector2{X: x, Y: y}, radius, theme.accent)
		label := fmt.Sprintf("%d", placement.Finger)
		if tabNumbers {
			label = fmt.Sprintf("%d", placement.Fret)
		}
		drawHeavyCenteredText(label, rl.Rectangle{
			X: x - radius, Y: y - radius, Width: radius * 2, Height: radius * 2,
		}, clamp(radius*1.65, 15, 23), theme.panel)
	}
}

// drawWaveform plots 25 milliseconds of the exact sounding-string signal.
func (gui *viewer) drawWaveform(bounds rl.Rectangle, theme palette) {
	drawPanel(bounds, theme)
	drawSectionTitle(bounds, "WAVEFORM · AMPLITUDE / TIME", theme)
	plot := inset(bounds, waveLeftPadding, graphTopPadding,
		waveRightPadding, waveBottomPadding)
	if plot.Width < 20 || plot.Height < 20 {
		return
	}
	data := gui.signalData()
	center := plot.Y + plot.Height/2
	// Matching vertical endpoints make the waveform's sampled time range explicit.
	rl.DrawLineEx(rl.Vector2{X: plot.X, Y: plot.Y},
		rl.Vector2{X: plot.X, Y: plot.Y + plot.Height}, 2, theme.border)
	rl.DrawLineEx(rl.Vector2{X: plot.X + plot.Width, Y: plot.Y},
		rl.Vector2{X: plot.X + plot.Width, Y: plot.Y + plot.Height}, 2, theme.border)
	rl.DrawLineEx(rl.Vector2{X: plot.X, Y: center},
		rl.Vector2{X: plot.X + plot.Width, Y: center}, 1, theme.border)
	gui.cacheWaveform(plot, data.notes)
	for index := 1; index < len(gui.signal.wavePoints); index++ {
		rl.DrawLineEx(gui.signal.wavePoints[index-1], gui.signal.wavePoints[index], 2, theme.signal)
	}
	drawAxisLabel("+1", plot.X, plot.Y, theme)
	drawAxisLabel("0", plot.X, center, theme)
	drawAxisLabel("-1", plot.X, plot.Y+plot.Height, theme)
	drawText("0 ms", plot.X, plot.Y+plot.Height+8, 14, theme.muted)
	right := "25 ms"
	drawText(right, plot.X+plot.Width-measureText(right, 14).X,
		plot.Y+plot.Height+8, 14, theme.muted)
}

// signalData rebuilds sorted notes, peaks, and labels only after a voicing change.
func (gui *viewer) signalData() *signalRenderCache {
	voicing := gui.controller.Voicing()
	if gui.signal.name == voicing.Name && gui.signal.number == voicing.Number {
		return &gui.signal
	}
	notes := signal.Notes(voicing)
	sort.Slice(notes, func(left, right int) bool { return notes[left].MIDI < notes[right].MIDI })
	legendNames := make([]string, len(notes))
	for index, note := range notes {
		legendNames[index] = note.Name
	}
	gui.signal = signalRenderCache{
		name: voicing.Name, number: voicing.Number, notes: notes,
		peaks: aggregateNotes(notes), legend: "Notes: " + strings.Join(legendNames, "  "),
	}
	return &gui.signal
}

// cacheWaveform samples the signal only when its voicing or plot geometry changes.
func (gui *viewer) cacheWaveform(plot rl.Rectangle, notes []signal.Note) {
	if gui.signal.waveBounds == plot && len(gui.signal.wavePoints) > 0 {
		return
	}
	center := plot.Y + plot.Height/2
	points := make([]rl.Vector2, maxInt(2, int(plot.Width)))
	for index := range points {
		t := float64(index) / float64(len(points)-1)
		points[index] = rl.Vector2{
			X: plot.X + float32(t)*plot.Width,
			Y: center - float32(signal.CompositeSample(t*waveformSeconds, notes))*plot.Height*0.43,
		}
	}
	gui.signal.waveBounds = plot
	gui.signal.wavePoints = points
}

// drawAxisLabel right-aligns a value inside the waveform's left gutter.
func drawAxisLabel(label string, axisX, centerY float32, theme palette) {
	const labelSize = 14
	size := measureText(label, labelSize)
	drawText(label, axisX-size.X-8, centerY-size.Y/2, labelSize, theme.muted)
}

// drawSpectrum places one exact peak per sounding pitch on a logarithmic scale.
func (gui *viewer) drawSpectrum(bounds rl.Rectangle, theme palette) {
	drawPanel(bounds, theme)
	drawSectionTitle(bounds, "FREQUENCY SPECTRUM · LOG Hz", theme)
	plot := inset(bounds, spectrumSidePadding, graphTopPadding,
		spectrumSidePadding, spectrumBottomPad)
	if plot.Width < 20 || plot.Height < 20 {
		return
	}
	data := gui.signalData()
	if len(data.notes) == 0 {
		return
	}
	peaks := data.peaks
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
		rl.DrawCircleV(rl.Vector2{X: x, Y: top}, 7, theme.signal)
		drawCenteredText(peak.note.Name, rl.Rectangle{
			X: x - 28, Y: axisY + 4, Width: 56, Height: 20,
		}, 14, theme.text)
	}
	drawCenteredText(data.legend, rl.Rectangle{
		X: bounds.X, Y: bounds.Y + bounds.Height - 25,
		Width: bounds.Width, Height: 18,
	}, 14, theme.muted)
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
	const (
		titleSize    = 20
		titleSpacing = 2.5
	)
	drawTextSpaced(title, bounds.X+panelPadding, bounds.Y+12,
		titleSize, titleSpacing, theme.accent)
	underlineWidth := measureTextSpaced(title, titleSize, titleSpacing).X + 1
	rl.DrawLineEx(
		rl.Vector2{X: bounds.X + panelPadding, Y: bounds.Y + 36},
		rl.Vector2{X: bounds.X + panelPadding + underlineWidth, Y: bounds.Y + 36},
		3,
		theme.secondary,
	)
}

// drawText applies consistent tracking and the heavier console treatment.
func drawText(text string, x, y, size float32, color rl.Color) {
	drawTextSpaced(text, x, y, size, bodyTextSpacing, color)
}

// drawTextSpaced renders the bundled bold font with configurable tracking.
func drawTextSpaced(text string, x, y, size, spacing float32, color rl.Color) {
	position := rl.Vector2{X: x, Y: y}
	rl.DrawTextEx(uiFont, text, position, size, spacing, color)
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

// drawHeavyCenteredText centers prominent diagram numbers using the bold font.
func drawHeavyCenteredText(text string, bounds rl.Rectangle, size float32, color rl.Color) {
	measured := measureText(text, size)
	drawTextSpaced(text,
		bounds.X+(bounds.Width-measured.X)/2,
		bounds.Y+(bounds.Height-measured.Y)/2,
		size, bodyTextSpacing, color)
}

// measureText returns default-font dimensions matching drawText.
func measureText(text string, size float32) rl.Vector2 {
	return measureTextSpaced(text, size, bodyTextSpacing)
}

// measureTextSpaced mirrors the tracking used by drawTextSpaced.
func measureTextSpaced(text string, size, spacing float32) rl.Vector2 {
	return rl.MeasureTextEx(uiFont, text, size, spacing)
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
	started    bool
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
		if player.started {
			rl.ResumeAudioStream(player.stream)
		} else {
			rl.PlayAudioStream(player.stream)
			player.started = true
		}
	}
}

// update refills every processed half-buffer from the main thread.
func (player *audioPlayer) update() {
	if !player.ready || !player.started || len(player.notes) == 0 {
		return
	}
	if !player.active && player.gain < 0.0001 {
		player.gain = 0
		rl.PauseAudioStream(player.stream)
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

// audible reports whether smooth playback needs the active frame rate.
func (player *audioPlayer) audible() bool {
	return player.ready && player.started && (player.active || player.gain >= 0.0001)
}

// close releases the native stream before the audio device shuts down.
func (player *audioPlayer) close() {
	if !player.ready {
		return
	}
	if player.started {
		rl.StopAudioStream(player.stream)
	}
	rl.UnloadAudioStream(player.stream)
}
