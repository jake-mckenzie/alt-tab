package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

const (
	minimumTerminalWidth     = 80
	maximumLayoutWidth       = 100
	bottomColumnMinimumWidth = 64
	sectionGapWidth          = 2
	outerHorizontalPadding   = 2
	dialCellWidth            = 10
	verticalDialWidth        = 10
	verticalDialHeight       = 7
	verticalDialBaseRow      = verticalDialHeight / 2
	panelFrameHeight         = 2
)

// render composes navigation and the chord, waveform, and spectrum dashboard.
func (model Model) render() string {
	// Keep panels readable instead of expanding them across oversized windows.
	width := min(max(1, model.width), maximumLayoutWidth)
	// Preserve usable content in unusually narrow terminals.
	padding := min(outerHorizontalPadding, max(0, (width-1)/2))
	contentWidth := max(1, width-padding*2)
	header := model.renderHeader(contentWidth)
	// Avoid unbreakable fretboard cells when no readable diagram can fit.
	if width < minimumTerminalWidth {
		return lipgloss.NewStyle().Padding(1, padding).Render(
			header + "\n\n" + model.styles.muted.Render(
				fmt.Sprintf("Need %d+ columns", minimumTerminalWidth),
			),
		)
	}
	banner := model.renderCenteredPanel(
		header+"\n\n"+model.renderControls(),
		contentWidth,
	)

	return lipgloss.NewStyle().
		Padding(1, padding).
		Render(banner + "\n\n" + model.renderSections(contentWidth))
}

// renderHeader omits the subtitle when it cannot fit without terminal wrapping.
func (model Model) renderHeader(width int) string {
	title := model.styles.title.Render("ALT-TAB")
	subtitle := model.styles.subtitle.Render(
		fmt.Sprintf("Guitar Chord Viewer · %s", paletteAt(model.theme).name),
	)
	header := lipgloss.JoinHorizontal(lipgloss.Bottom, title, "  ", subtitle)
	if lipgloss.Width(header) <= width {
		return header
	}
	return title
}

// renderSections keeps navigation above the responsive dashboard modules.
func (model Model) renderSections(width int) string {
	selector := model.renderCenteredPanel(model.renderChordSelector(), width)

	if model.showHelp {
		help := model.renderPanel(model.renderHelp(), width, true)
		return lipgloss.JoinVertical(lipgloss.Left, selector, "", help)
	}

	diagramContent := model.renderChordDiagram()
	if !model.fullNeck && width >= bottomColumnMinimumWidth {
		diagramWidth := min(width, lipgloss.Width(diagramContent)+4)
		graphWidth := width - sectionGapWidth - diagramWidth
		graphContentWidth := graphWidth - 6
		waveformContent := model.renderWaveformSection(graphContentWidth)
		spectrumContent := model.renderSpectrumSection(graphContentWidth)
		waveform := model.renderCenteredPanel(waveformContent, graphWidth)
		spectrum := model.renderCenteredPanel(spectrumContent, graphWidth)
		graphs := lipgloss.JoinVertical(lipgloss.Left, waveform, "", spectrum)
		diagramContent = model.renderCompactChordDiagram(
			max(1, lipgloss.Height(graphs)-panelFrameHeight),
		)
		diagram := model.renderCenteredPanel(diagramContent, diagramWidth)
		bottom := lipgloss.JoinHorizontal(
			lipgloss.Top,
			diagram,
			strings.Repeat(" ", sectionGapWidth),
			graphs,
		)
		return lipgloss.JoinVertical(lipgloss.Left, selector, "", bottom)
	}

	diagram := model.renderCenteredPanel(diagramContent, width)
	graphGap := strings.Repeat(" ", sectionGapWidth)
	leftGraphWidth := (width - sectionGapWidth) / 2
	rightGraphWidth := width - sectionGapWidth - leftGraphWidth
	waveform, spectrum := model.renderMatchedCenteredPanels(
		model.renderWaveformSection(leftGraphWidth-6),
		leftGraphWidth,
		model.renderSpectrumSection(rightGraphWidth-6),
		rightGraphWidth,
	)
	graphs := lipgloss.JoinHorizontal(
		lipgloss.Top,
		waveform,
		graphGap,
		spectrum,
	)
	return lipgloss.JoinVertical(
		lipgloss.Left,
		selector,
		"",
		diagram,
		"",
		graphs,
	)
}

// renderMatchedCenteredPanels gives paired sections equal top-aligned rectangles.
func (model Model) renderMatchedCenteredPanels(
	leftContent string,
	leftWidth int,
	rightContent string,
	rightWidth int,
) (string, string) {
	left := model.renderCenteredPanel(leftContent, leftWidth)
	right := model.renderCenteredPanel(rightContent, rightWidth)
	height := max(lipgloss.Height(left), lipgloss.Height(right))
	return model.renderCenteredPanelAtHeight(leftContent, leftWidth, height),
		model.renderCenteredPanelAtHeight(rightContent, rightWidth, height)
}

// renderCenteredPanel fills the shared section width and centers each content line.
func (model Model) renderCenteredPanel(content string, width int) string {
	return model.renderCenteredPanelAtHeight(content, width, 0)
}

// renderCenteredPanelAtHeight centers lines within an exact panel rectangle.
func (model Model) renderCenteredPanelAtHeight(content string, width, height int) string {
	panelWidth := max(1, width)
	innerWidth := max(1, panelWidth-4)
	centered := lipgloss.NewStyle().
		Width(innerWidth).
		Align(lipgloss.Center).
		Render(content)
	style := model.styles.panel.Width(panelWidth)
	if height > 0 {
		style = style.Height(height)
	}
	return style.Render(centered)
}

// renderPanel optionally shrinks a bordered section to its widest content line.
func (model Model) renderPanel(content string, maximumWidth int, shrink bool) string {
	// Add both border and padding cells when fitting raw content.
	styleWidth := max(1, maximumWidth)
	if shrink {
		styleWidth = min(styleWidth, lipgloss.Width(content)+4)
	}
	return model.styles.panel.Width(styleWidth).Render(content)
}

// renderChordSelector draws the base row and three-position chord-family dial.
func (model Model) renderChordSelector() string {
	var output strings.Builder
	output.WriteString(model.styles.heading.Render("CHORD DIAL"))
	output.WriteString("\n\n")
	output.WriteString(model.styles.muted.Render("BASE CHORDS  "))
	for index, family := range model.families {
		if index > 0 {
			output.WriteString("  ")
		}
		if index == model.selected {
			output.WriteString(model.styles.accent.Bold(true).Render(family.base))
		} else {
			output.WriteString(model.styles.normal.Render(family.base))
		}
	}
	output.WriteString("\n\n")

	if len(model.families) == 0 {
		return output.String()
	}
	family := model.families[model.selected]
	previous := model.families[wrap(model.selected-1, len(model.families))]
	next := model.families[wrap(model.selected+1, len(model.families))]
	output.WriteString(model.renderNestedDial(previous.base, family, next.base))

	return output.String()
}

// renderNestedDial embeds real accidental and minor choices in the center cell.
func (model Model) renderNestedDial(left string, family chordFamily, right string) string {
	if family.accidental == "" && family.minor == "" {
		center := model.styles.selected.Render("‹ " + family.base + " ›")
		lines := make([]string, verticalDialHeight)
		for row := range lines {
			if row == verticalDialBaseRow {
				lines[row] = model.renderHorizontalDialLine(left, center, right)
			} else {
				lines[row] = model.renderHorizontalDialLine("", "", "")
			}
		}
		return strings.Join(lines, "\n")
	}

	border := model.styles.accent.Render
	entry := func(name string, kind int) string {
		entry := model.styles.normal.Render(name)
		if model.chordKind == kind {
			entry = model.styles.selected.Render("‹ " + name + " ›")
		}
		return border("│") + centerDisplayText(entry, verticalDialWidth) + border("│")
	}

	centers := make([]string, verticalDialHeight)
	top := verticalDialBaseRow - 1
	bottom := verticalDialBaseRow + 1
	if family.accidental != "" {
		top = 0
		centers[1] = entry(family.accidental, accidentalChord)
	}
	if family.minor != "" {
		bottom = verticalDialHeight - 1
		centers[verticalDialHeight-2] = entry(family.minor, minorChord)
	}
	centers[top] = border("╭" + strings.Repeat("─", verticalDialWidth) + "╮")
	centers[verticalDialBaseRow] = entry(family.base, baseChord)
	centers[bottom] = border("╰" + strings.Repeat("─", verticalDialWidth) + "╯")
	for row := top + 1; row < bottom; row++ {
		if centers[row] == "" {
			centers[row] = border("│") +
				strings.Repeat(" ", verticalDialWidth) + border("│")
		}
	}

	lines := make([]string, verticalDialHeight)
	for row, center := range centers {
		if row == verticalDialBaseRow {
			lines[row] = model.renderHorizontalDialLine(left, center, right)
		} else {
			lines[row] = model.renderHorizontalDialLine("", center, "")
		}
	}
	return strings.Join(lines, "\n")
}

// renderHorizontalDialLine aligns neighboring bases with the nested center cell.
func (model Model) renderHorizontalDialLine(left, center, right string) string {
	return centerDisplayText(model.styles.normal.Render(left), dialCellWidth) +
		strings.Repeat(" ", sectionGapWidth) +
		centerDisplayText(center, verticalDialWidth+2) +
		strings.Repeat(" ", sectionGapWidth) +
		centerDisplayText(model.styles.normal.Render(right), dialCellWidth)
}

// centerDisplayText pads styled terminal text without counting ANSI bytes.
func centerDisplayText(value string, width int) string {
	padding := max(0, width-lipgloss.Width(value))
	left := padding / 2
	return strings.Repeat(" ", left) + value + strings.Repeat(" ", padding-left)
}

// renderControls summarizes every available runtime command.
func (model Model) renderControls() string {
	hints := [][2]string{
		{"←→", "Base"},
		{"↑↓", "Type"},
		{"v", "Voicing"},
		{"f", "Neck"},
		{"n", "Tab"},
		{"t", "Theme"},
		{"w", "Wave"},
		{"?", "Help"},
		{"q", "Quit"},
	}
	parts := make([]string, len(hints))
	for index, hint := range hints {
		key := model.styles.accent.Bold(true).Render(hint[0])
		parts[index] = key + " " + model.styles.muted.Render(hint[1])
	}
	return model.styles.heading.Render("KEYS") + " " +
		strings.Join(parts, " ")
}

// renderChordDiagram draws the current chord heading, mode, and fretboard.
func (model Model) renderChordDiagram() string {
	if model.err != nil {
		return model.styles.err.Render(model.err.Error())
	}
	header := model.renderChordDiagramHeader()
	if model.fullNeck && model.width < fullNeckMinimumTerminalWidth {
		return header + "\n\n" + model.styles.err.Render(
			fmt.Sprintf(
				"Widen the terminal to at least %d columns (currently %d).",
				fullNeckMinimumTerminalWidth,
				model.width,
			),
		)
	}

	board, legend := renderFretboardParts(
		model.voicing,
		model.fullNeck,
		model.tabNotation,
	)
	return header + "\n\n" + model.styles.normal.Render(board) +
		"\n\n" + model.styles.normal.Render(legend)
}

// renderChordDiagramHeader stacks the mode, active chord, and voicing number.
func (model Model) renderChordDiagramHeader() string {
	voicingCount := model.catalog.VoicingCount(model.voicing.Name)
	mode := "COMPACT"
	if model.fullNeck {
		mode = fmt.Sprintf("FULL NECK · FRETS 1–%d", fullNeckLastFret)
	} else if model.tabNotation {
		mode = "TAB NUMBERS"
	}

	return model.styles.heading.Render("CHORD DIAGRAM") + "\n" +
		model.styles.accent.Render(mode) + "\n\n" +
		model.styles.selected.Render("  "+model.voicing.Name+"  ") + "\n" +
		model.styles.muted.Render(fmt.Sprintf(
			"VOICING %d OF %d",
			model.voicing.Number,
			voicingCount,
		))
}

// renderCompactChordDiagram anchors legends and vertically centers the neck.
func (model Model) renderCompactChordDiagram(height int) string {
	if model.err != nil {
		return model.styles.err.Render(model.err.Error())
	}
	board, legend := renderFretboardParts(model.voicing, false, model.tabNotation)
	return arrangeCompactDiagram(
		model.renderChordDiagramHeader(),
		model.styles.normal.Render(board),
		model.styles.normal.Render(legend),
		height,
	)
}

// arrangeCompactDiagram fixes the header, board center, and legend footer.
func arrangeCompactDiagram(header, board, legend string, height int) string {
	headerLines := strings.Split(header, "\n")
	boardLines := strings.Split(board, "\n")
	legendLines := strings.Split(legend, "\n")
	minimumHeight := len(headerLines) + len(boardLines) + len(legendLines) + 2
	height = max(height, minimumHeight)
	boardTop := max(len(headerLines)+1, (height-len(boardLines))/2)
	legendTop := height - len(legendLines)
	if boardTop+len(boardLines) >= legendTop {
		boardTop = len(headerLines) + 1
		legendTop = boardTop + len(boardLines) + 1
		height = legendTop + len(legendLines)
	}

	lines := make([]string, height)
	copy(lines, headerLines)
	copy(lines[boardTop:], boardLines)
	copy(lines[legendTop:], legendLines)
	return strings.Join(lines, "\n")
}

// renderWaveformSection labels the note plot or explains how to restore it.
func (model Model) renderWaveformSection(width int) string {
	heading := model.styles.heading.Render("NOTE WAVEFORM")
	if !model.waveform {
		return heading + "\n\n" +
			model.styles.muted.Render("Hidden · press w to show")
	}
	if model.err != nil {
		return heading + "\n\n" + model.styles.err.Render(model.err.Error())
	}

	return heading + "\n\n" + model.styles.accent.Render(
		renderWaveform(width, model.voicing),
	)
}

// renderSpectrumSection owns the frequency-domain view and its note legend.
func (model Model) renderSpectrumSection(width int) string {
	heading := model.styles.heading.Render("FREQUENCY SPECTRUM")
	if model.err != nil {
		return heading + "\n\n" + model.styles.err.Render(model.err.Error())
	}

	return heading + "\n\n" + model.styles.accent.Render(
		renderSpectrum(width, model.voicing),
	)
}

// renderHelp lists every keyboard command in a dedicated panel.
func (model Model) renderHelp() string {
	return model.styles.heading.Render("KEYBOARD HELP") + "\n\n" +
		"← / h     Previous chord\n" +
		"→ / l     Next chord\n" +
		"↑ / k     Accidental chord or return to base\n" +
		"↓ / j     Minor chord or return to base\n" +
		"v         Next voicing\n" +
		"f         Toggle compact/full neck\n" +
		"n         Toggle fingered fretboard/tab numbers\n" +
		"t         Cycle color theme\n" +
		"w         Toggle detailed waveform\n" +
		"wheel     Scroll viewport\n" +
		"?         Open or close help\n" +
		"esc       Close help\n" +
		"q         Quit\n\n" +
		model.styles.muted.Render("Press ? or esc to return")
}
