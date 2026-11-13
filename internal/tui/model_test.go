package tui

import (
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
)

// fakeCatalog isolates model behavior from the built-in chord table.
type fakeCatalog struct{}

// Names supplies the small ordered chord set used by model tests.
func (fakeCatalog) Names() []string {
	return []string{"C", "Cm", "G", "Gm"}
}

// VoicingCount gives each fake chord two navigable voicings.
func (fakeCatalog) VoicingCount(string) int {
	return 2
}

// Load constructs a predictable voicing without invoking the built-in catalog.
func (fakeCatalog) Load(name string, number int) (chords.Voicing, error) {
	if number < 1 || number > 2 {
		return chords.Voicing{}, errors.New("missing voicing")
	}

	return chords.Voicing{
		Name:   name,
		Number: number,
		Strings: [chords.StringCount]chords.StringPlacement{
			{Fret: 0},
			{Fret: 1, Finger: 1},
			{Fret: 0},
			{Fret: 2, Finger: 2},
			{Fret: 3, Finger: 3},
			{Fret: -1},
		},
	}, nil
}

// TestViewContainsChordAndFingering checks the essential rendered content.
func TestViewContainsChordAndFingering(t *testing.T) {
	view := New(fakeCatalog{}).View()

	for _, expected := range []string{"alt-tab", "C", "VOICING 1 OF 2", "--1--", "X"} {
		if !strings.Contains(view.Content, expected) {
			t.Fatalf("view does not contain %q", expected)
		}
	}
	if !view.AltScreen || view.MouseMode != tea.MouseModeCellMotion {
		t.Fatal("view does not use the resize-safe screen with wheel input")
	}
}

// TestChordListIsHorizontalAndAboveBothFretboardModes checks shared layout.
func TestChordListIsHorizontalAndAboveBothFretboardModes(t *testing.T) {
	for _, fullNeck := range []bool{false, true} {
		model := New(fakeCatalog{})
		model.fullNeck = fullNeck
		plain := ansiSequence.ReplaceAllString(model.View().Content, "")

		if !strings.Contains(plain, "BASE CHORDS  [C]  G") {
			t.Fatalf("fullNeck=%t does not show the base-chord row", fullNeck)
		}
		if strings.Index(plain, "CHORD SELECTOR") >
			strings.Index(plain, "VOICING 1 OF 2") {
			t.Fatalf("fullNeck=%t places the chord dial below the fretboard", fullNeck)
		}
	}
}

// TestSectionsAreClearlyLabeledAndOrdered checks the complete screen hierarchy.
func TestSectionsAreClearlyLabeledAndOrdered(t *testing.T) {
	plain := ansiSequence.ReplaceAllString(New(fakeCatalog{}).View().Content, "")
	positions := make(map[string]int)
	for _, title := range []string{
		"KEYS", "CHORD SELECTOR", "CHORD DIAGRAM", "WAVEFORM · AMPLITUDE / TIME", "FREQUENCY SPECTRUM",
	} {
		positions[title] = strings.Index(plain, title)
		if positions[title] < 0 {
			t.Fatalf("section %q is missing", title)
		}
	}
	if positions["KEYS"] > positions["CHORD SELECTOR"] ||
		positions["CHORD SELECTOR"] > positions["CHORD DIAGRAM"] ||
		positions["CHORD SELECTOR"] > positions["WAVEFORM · AMPLITUDE / TIME"] ||
		positions["CHORD SELECTOR"] > positions["FREQUENCY SPECTRUM"] {
		t.Fatal("title controls, dial, and output sections are out of order")
	}
}

// TestCompactPanelsShrinkToContent checks compact and full-neck width policy.
func TestCompactPanelsShrinkToContent(t *testing.T) {
	model := New(fakeCatalog{})
	const maximumWidth = 76

	selector := model.renderPanel(model.renderChordSelector(), maximumWidth, true)
	diagram := model.renderPanel(model.renderChordDiagram(), maximumWidth, true)
	if lipgloss.Width(selector) >= maximumWidth || lipgloss.Width(diagram) >= maximumWidth {
		t.Fatal("compact selector or diagram retains unnecessary trailing width")
	}

	model.fullNeck = true
	diagram = model.renderPanel(model.renderChordDiagram(), maximumWidth, false)
	if lipgloss.Width(diagram) != maximumWidth {
		t.Fatalf("full-neck diagram width = %d, want %d", lipgloss.Width(diagram), maximumWidth)
	}
}

// lineWith returns the rendered line containing a section title.
func lineWith(output, title string) string {
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, title) {
			return line
		}
	}
	return ""
}

// TestDashboardLayoutGroupsOutputPanels checks responsive module placement.
func TestDashboardLayoutGroupsOutputPanels(t *testing.T) {
	model := New(fakeCatalog{})
	plain := ansiSequence.ReplaceAllString(model.View().Content, "")
	if !hasDoubleTopBorder(plain) {
		t.Fatal("compact diagram and graph stack are not side-by-side")
	}
	if strings.Index(plain, "WAVEFORM · AMPLITUDE / TIME") > strings.Index(plain, "FREQUENCY SPECTRUM") {
		t.Fatal("compact graph stack does not place the waveform above the spectrum")
	}

	model.fullNeck = true
	plain = ansiSequence.ReplaceAllString(model.View().Content, "")
	if hasDoubleTopBorder(plain) ||
		strings.Index(plain, "CHORD DIAGRAM") > strings.Index(plain, "WAVEFORM · AMPLITUDE / TIME") ||
		strings.Index(plain, "WAVEFORM · AMPLITUDE / TIME") > strings.Index(plain, "FREQUENCY SPECTRUM") {
		t.Fatal("full-neck dashboard does not use three stacked full-width modules")
	}
}

// hasDoubleTopBorder reports whether two panels begin on the same rendered row.
func hasDoubleTopBorder(output string) bool {
	for _, line := range strings.Split(output, "\n") {
		if strings.Count(line, "╭") == 2 {
			return true
		}
	}
	return false
}

// TestControlsStayOnOneClearlyLabeledLine checks the compact command legend.
func TestControlsStayOnOneClearlyLabeledLine(t *testing.T) {
	model := New(chords.NewCatalog())
	if strings.Contains(model.renderControls(), "\n") {
		t.Fatal("controls are not rendered on one line")
	}
	if !strings.HasPrefix(ansiSequence.ReplaceAllString(model.renderControls(), ""), "KEYS") {
		t.Fatal("controls do not have an explicit title")
	}
	for _, title := range []string{"Chord", "Voicing", "Neck", "Tab", "Theme", "Wave"} {
		if !strings.Contains(model.renderControls(), title) {
			t.Fatalf("controls omit the %s action title", title)
		}
	}
	panel := model.renderPanel(model.renderControls(), maximumLayoutWidth-4, true)
	if lipgloss.Width(panel) > maximumLayoutWidth-4 {
		t.Fatal("single-line controls exceed the documented terminal width")
	}
}

// TestHeaderAndDialFillSharedWidthAndCenterContent checks top section geometry.
func TestHeaderAndDialFillSharedWidthAndCenterContent(t *testing.T) {
	model := New(chords.NewCatalog())
	const width = maximumLayoutWidth - 4
	header := model.renderHeader(width)
	banner := ansiSequence.ReplaceAllString(
		model.renderCenteredPanel(header+"\n\n"+model.renderControls(), width),
		"",
	)
	dial := ansiSequence.ReplaceAllString(
		model.renderCenteredPanel(model.renderChordSelector(), width),
		"",
	)
	if lipgloss.Width(banner) != width || lipgloss.Width(dial) != width {
		t.Fatal("title banner and dial do not fill their shared section width")
	}
	for output, title := range map[string]string{
		banner: "alt-tab · Guitar Chord Viewer · Synthwave",
		dial:   "CHORD SELECTOR",
	} {
		line := lineWith(output, title)
		titleIndex := strings.Index(line, title)
		left := lipgloss.Width(line[:titleIndex])
		right := lipgloss.Width(line) - left - len(title)
		if absoluteInt(left-right) > 1 {
			t.Fatalf("%s is not centered: left=%d right=%d", title, left, right)
		}
	}
	if lineWith(banner, "Guitar Chord Viewer · Synthwave") == "" {
		t.Fatal("app subtitle is missing from the title line")
	}
}

// TestBottomPanelsShareHeightAndCenterContent checks matched output rectangles.
func TestBottomPanelsShareHeightAndCenterContent(t *testing.T) {
	model := New(fakeCatalog{})
	diagram, waveform := model.renderMatchedCenteredPanels(
		model.renderChordDiagram(),
		28,
		model.renderWaveformSection(40),
		46,
	)
	if lipgloss.Height(diagram) != lipgloss.Height(waveform) {
		t.Fatalf(
			"bottom heights differ: diagram=%d waveform=%d",
			lipgloss.Height(diagram),
			lipgloss.Height(waveform),
		)
	}
	plain := ansiSequence.ReplaceAllString(diagram, "")
	titleLine := lineWith(plain, "CHORD DIAGRAM")
	titleIndex := strings.Index(titleLine, "CHORD DIAGRAM")
	left := lipgloss.Width(titleLine[:titleIndex])
	right := lipgloss.Width(titleLine) - left - len("CHORD DIAGRAM")
	if absoluteInt(left-right) > 1 {
		t.Fatal("diagram content is not horizontally centered")
	}
}

// TestCompactDiagramCentersBoardAndBottomAlignsLegend checks vertical layout.
func TestCompactDiagramCentersBoardAndBottomAlignsLegend(t *testing.T) {
	model := New(fakeCatalog{})
	const height = 25
	plain := ansiSequence.ReplaceAllString(model.renderCompactChordDiagram(height), "")
	lines := strings.Split(plain, "\n")
	firstString := strings.Index(plain, "e  O")
	lastString := strings.Index(plain, "E  X")
	if len(lines) != height || firstString < 0 || lastString < 0 {
		t.Fatal("compact diagram does not contain the expected fixed-height fretboard")
	}
	firstRow := strings.Count(plain[:firstString], "\n")
	lastRow := strings.Count(plain[:lastString], "\n")
	boardCenter := (firstRow + lastRow) / 2
	if absoluteInt(boardCenter-(height-1)/2) > 1 {
		t.Fatalf("fretboard center row = %d, want %d", boardCenter+1, (height-1)/2+1)
	}
	if strings.TrimSpace(lines[height-1]) != "Symbols: O open  X muted" {
		t.Fatalf("compact footer = %q, want symbol legend", lines[height-1])
	}
	if !strings.Contains(lines[height-4], "Fingers:") {
		t.Fatal("finger legend is not anchored with the compact footer")
	}
}

// TestFullNeckPanelsUseNaturalHeights avoids padding stacked output sections.
func TestFullNeckPanelsUseNaturalHeights(t *testing.T) {
	model := New(chords.NewCatalog())
	model.fullNeck = true
	diagram := model.renderCenteredPanel(model.renderChordDiagram(), maximumLayoutWidth-4)
	waveform := model.renderCenteredPanel(
		model.renderWaveformSection(maximumLayoutWidth-10),
		maximumLayoutWidth-4,
	)
	if lipgloss.Height(diagram) == lipgloss.Height(waveform) {
		t.Fatalf(
			"full-neck diagram and waveform were forced to the same height %d",
			lipgloss.Height(diagram),
		)
	}
}

// TestChordSelectorShowsCurrentQuality checks variant labels and stable height.
func TestChordSelectorShowsCurrentQuality(t *testing.T) {
	model := New(chords.NewCatalog())
	height := lipgloss.Height(model.renderChordSelector())
	for _, step := range []struct{ key, label string }{
		{",", "Am · minor"}, {",", "A · major"}, {"l", "B · major"},
		{",", "B · major"}, {".", "Bb · major (accidental)"}, {".", "B · major"},
	} {
		next, _ := model.Update(tea.KeyPressMsg{Code: rune(step.key[0]), Text: step.key})
		model = next.(Model)
		plain := ansiSequence.ReplaceAllString(model.renderChordSelector(), "")
		if !strings.Contains(plain, step.label) || lipgloss.Height(plain) != height {
			t.Fatalf("%s: unexpected selector: %s", step.key, plain)
		}
	}
}

// TestDiagramHeadingsAreStackedAndCentered checks the requested hierarchy.
func TestDiagramHeadingsAreStackedAndCentered(t *testing.T) {
	plain := ansiSequence.ReplaceAllString(New(fakeCatalog{}).renderChordDiagram(), "")
	lines := strings.Split(plain, "\n")
	want := []string{"CHORD DIAGRAM", "COMPACT", "", "C", "VOICING 1 OF 2"}
	if len(lines) < len(want) {
		t.Fatalf("diagram heading has %d lines, want at least %d", len(lines), len(want))
	}
	for index, expected := range want {
		if strings.TrimSpace(lines[index]) != expected {
			t.Fatalf("diagram heading line %d = %q, want %q", index+1, lines[index], expected)
		}
	}
}

// TestChordFamiliesMatchCatalog checks every accidental and minor dial position.
func TestChordFamiliesMatchCatalog(t *testing.T) {
	want := []chordFamily{
		{base: "A", minor: "Am"},
		{base: "B", accidental: "Bb"},
		{base: "C", minor: "Cm"},
		{base: "D", minor: "Dm"},
		{base: "E", minor: "Em"},
		{base: "F", accidental: "F#"},
		{base: "G", minor: "Gm"},
	}
	got := buildChordFamilies(chords.NewCatalog().Names())
	if len(got) != len(want) {
		t.Fatalf("family count = %d, want %d", len(got), len(want))
	}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("family %d = %+v, want %+v", index, got[index], want[index])
		}
	}
}

// TestResizeRedrawsWithinTerminalWidth prevents wrapping outside Bubble Tea's frame.
func TestResizeRedrawsWithinTerminalWidth(t *testing.T) {
	for _, fullNeck := range []bool{false, true} {
		model := New(chords.NewCatalog())
		model.fullNeck = fullNeck
		for _, width := range []int{120, 80, 40, 20, 98, 30, 64, 39, 97, 79} {
			updated, command := model.Update(tea.WindowSizeMsg{Width: width, Height: 40})
			model = updated.(Model)
			if command != nil {
				t.Fatalf("width %d returned an unnecessary manual redraw", width)
			}

			plain := ansiSequence.ReplaceAllString(model.View().Content, "")
			if viewHeight := lipgloss.Height(plain); viewHeight > model.height {
				t.Fatalf("width %d rendered %d rows into %d", width, viewHeight, model.height)
			}
			for lineNumber, line := range strings.Split(plain, "\n") {
				if lineWidth := lipgloss.Width(line); lineWidth > width {
					t.Fatalf(
						"fullNeck=%t width %d line %d renders %d columns: %q\n%s",
						fullNeck,
						width,
						lineNumber+1,
						lineWidth,
						line,
						plain,
					)
				}
			}
		}
	}
}

// TestWideTerminalDoesNotStretchLayout keeps unused columns outside the UI.
func TestWideTerminalDoesNotStretchLayout(t *testing.T) {
	for _, fullNeck := range []bool{false, true} {
		model := New(chords.NewCatalog())
		model.width = 180
		model.fullNeck = fullNeck
		plain := ansiSequence.ReplaceAllString(model.View().Content, "")
		for lineNumber, line := range strings.Split(plain, "\n") {
			if width := lipgloss.Width(line); width > maximumLayoutWidth {
				t.Fatalf(
					"fullNeck=%t line %d stretched to %d columns",
					fullNeck,
					lineNumber+1,
					width,
				)
			}
		}
	}
}

// TestMouseWheelScrollsViewportWithoutChangingVoicing checks wheel isolation.
func TestMouseWheelScrollsViewportWithoutChangingVoicing(t *testing.T) {
	model := New(fakeCatalog{})
	model.width = 80
	model.height = 12
	before := model.View().Content
	updated, _ := model.Update(tea.MouseWheelMsg{Button: tea.MouseWheelDown})
	model = updated.(Model)

	if model.scroll == 0 || model.View().Content == before {
		t.Fatal("mouse wheel did not scroll the viewport")
	}
	if model.voicing.Name != "C" || model.voicing.Number != 1 {
		t.Fatal("mouse wheel changed the selected voicing")
	}
}

// TestNavigationUpdatesSelection checks dial, alternate, and voicing controls.
func TestNavigationUpdatesSelection(t *testing.T) {
	model := New(fakeCatalog{})
	updated, _ := model.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	model = updated.(Model)

	if model.voicing.Name != "G" || model.voicing.Number != 1 {
		t.Fatalf("right selected %s:%d, want G:1", model.voicing.Name, model.voicing.Number)
	}

	updated, _ = model.Update(tea.KeyPressMsg{Code: ',', Text: ","})
	model = updated.(Model)
	if model.voicing.Name != "Gm" || model.voicing.Number != 1 {
		t.Fatalf("down selected %s:%d, want Gm:1", model.voicing.Name, model.voicing.Number)
	}

	updated, _ = model.Update(tea.KeyPressMsg{Code: '2', Text: "2"})
	model = updated.(Model)
	if model.voicing.Name != "Gm" || model.voicing.Number != 2 {
		t.Fatalf("v selected %s:%d, want Gm:2", model.voicing.Name, model.voicing.Number)
	}
}

// TestNumberedVoicings checks direct selection and the removal of cycling.
func TestNumberedVoicings(t *testing.T) {
	model := New(chords.NewCatalog())
	for _, number := range []rune{'3', '1', '2'} {
		next, _ := model.Update(tea.KeyPressMsg{Code: number, Text: string(number)})
		model = next.(Model)
		if model.voicing.Number != int(number-'0') {
			t.Fatal("number did not select voicing")
		}
	}
	for _, key := range []rune{'v', '9'} {
		next, _ := model.Update(tea.KeyPressMsg{Code: key, Text: string(key)})
		model = next.(Model)
		if model.voicing.Number != 2 {
			t.Fatal("removed or unavailable key changed voicing")
		}
	}
}

// TestDiagramModesAndHelpToggle checks mutually exclusive diagram modes.
func TestDiagramModesAndHelpToggle(t *testing.T) {
	model := New(fakeCatalog{})
	updated, _ := model.Update(tea.KeyPressMsg{Code: 'x', Text: "x"})
	model = updated.(Model)
	if !model.fullNeck ||
		!strings.Contains(model.View().Content, " 1  2  3  4  5  6  7  8  9") ||
		!strings.Contains(model.View().Content, " 25 26 27") {
		t.Fatal("x did not enable the full-neck view")
	}

	updated, _ = model.Update(tea.KeyPressMsg{Code: 'n', Text: "n"})
	model = updated.(Model)
	plain := ansiSequence.ReplaceAllString(model.View().Content, "")
	if !model.tabNotation || model.fullNeck ||
		!strings.Contains(plain, "TAB NUMBERS") ||
		!strings.Contains(plain, "B  |--1-----------------|") ||
		strings.Contains(plain, "Fingers:") ||
		!strings.Contains(plain, "Symbols: O open  X muted") {
		t.Fatal("n did not enable the fret-number tab view")
	}

	updated, _ = model.Update(tea.KeyPressMsg{Code: 'x', Text: "x"})
	model = updated.(Model)
	if !model.fullNeck || model.tabNotation {
		t.Fatal("full-neck mode did not replace tab-number mode")
	}

	updated, _ = model.Update(tea.KeyPressMsg{Code: '?', Text: "?"})
	model = updated.(Model)
	if !model.showHelp || !strings.Contains(model.View().Content, "KEYBOARD HELP") {
		t.Fatal("? did not open help")
	}
}

// TestThemeToggleCyclesPalettes checks runtime switching and wraparound.
func TestThemeToggleCyclesPalettes(t *testing.T) {
	model := New(fakeCatalog{})
	if model.theme != 0 || !strings.Contains(model.View().Content, "Synthwave") {
		t.Fatal("new model does not use the Synthwave theme")
	}

	for index := 1; index <= len(palettes); index++ {
		updated, _ := model.Update(tea.KeyPressMsg{Code: 't', Text: "t"})
		model = updated.(Model)
		want := index % len(palettes)
		if model.theme != want ||
			!strings.Contains(model.View().Content, paletteAt(want).name) {
			t.Fatalf("theme index = %d, want %d", model.theme, want)
		}
	}
}

// TestFullNeckRequiresNinetyEightColumns checks both sides of the boundary.
func TestFullNeckRequiresNinetyEightColumns(t *testing.T) {
	model := New(fakeCatalog{})
	model.width = fullNeckMinimumTerminalWidth - 1
	model.fullNeck = true

	if !strings.Contains(model.View().Content, "Widen the terminal") {
		t.Fatal("narrow full-neck view does not request a wider terminal")
	}

	model.width = fullNeckMinimumTerminalWidth
	if strings.Contains(model.View().Content, "Widen the terminal") {
		t.Fatal("full-neck view rejects a 98-column terminal")
	}
	if !strings.Contains(model.View().Content, " 25 26 27") {
		t.Fatal("full-neck view does not render at 98 columns")
	}
}

// TestQuitKeyReturnsCommand checks that quitting is delegated to Bubble Tea.
func TestQuitKeyReturnsCommand(t *testing.T) {
	model := New(fakeCatalog{})
	_, command := model.Update(tea.KeyPressMsg{Code: 'q', Text: "q"})

	if command == nil {
		t.Fatal("q did not return a quit command")
	}
}
