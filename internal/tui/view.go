package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

const (
	minimumTerminalWidth  = 40
	twoColumnMinimumWidth = 64
	sectionGapWidth       = 2
)

// render composes the titled selector, controls, diagram, and waveform sections.
func (model Model) render() string {
	width := model.width
	if width < minimumTerminalWidth {
		width = minimumTerminalWidth
	}
	// The outer style contributes two columns of padding on both sides.
	contentWidth := width - 4

	header := lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		model.styles.title.Render("ALT-TAB"),
		"  ",
		model.styles.subtitle.Render(
			fmt.Sprintf("Guitar Chord Viewer · %s", paletteAt(model.theme).name),
		),
	)
	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(header + "\n\n" + model.renderSections(contentWidth))
}

// renderSections keeps navigation above the diagram and waveform output.
func (model Model) renderSections(width int) string {
	var top string
	if width >= twoColumnMinimumWidth {
		leftWidth := (width - sectionGapWidth) / 2
		rightWidth := width - sectionGapWidth - leftWidth
		selector := model.renderPanel(model.renderChordSelector(), leftWidth, false)
		controls := model.renderPanel(model.renderControls(), rightWidth, false)
		top = lipgloss.JoinHorizontal(
			lipgloss.Top,
			selector,
			strings.Repeat(" ", sectionGapWidth),
			controls,
		)
	} else {
		selector := model.renderPanel(model.renderChordSelector(), width, true)
		controls := model.renderPanel(model.renderControls(), width, true)
		top = lipgloss.JoinVertical(lipgloss.Left, selector, "", controls)
	}

	if model.showHelp {
		help := model.renderPanel(model.renderHelp(), width, true)
		return lipgloss.JoinVertical(lipgloss.Left, top, "", help)
	}

	diagramContent := model.renderChordDiagram()
	if !model.fullNeck && width >= twoColumnMinimumWidth {
		diagram := model.renderPanel(diagramContent, width, true)
		waveformWidth := width - sectionGapWidth - lipgloss.Width(diagram)
		// Leave two spare cells because some terminals treat Braille width loosely.
		waveform := model.renderPanel(
			model.renderWaveformSection(waveformWidth-6),
			waveformWidth,
			false,
		)
		bottom := lipgloss.JoinHorizontal(
			lipgloss.Top,
			diagram,
			strings.Repeat(" ", sectionGapWidth),
			waveform,
		)
		return lipgloss.JoinVertical(lipgloss.Left, top, "", bottom)
	}

	diagram := model.renderPanel(diagramContent, width, !model.fullNeck)
	waveform := model.renderPanel(model.renderWaveformSection(width-6), width, false)
	return lipgloss.JoinVertical(
		lipgloss.Left,
		top,
		"",
		diagram,
		"",
		waveform,
	)
}

// renderPanel optionally shrinks a bordered section to its widest content line.
func (model Model) renderPanel(content string, maximumWidth int, shrink bool) string {
	// Add both border and padding cells when fitting raw content.
	styleWidth := maximumWidth
	if shrink {
		styleWidth = min(styleWidth, lipgloss.Width(content)+4)
	}
	return model.styles.panel.Width(styleWidth).Render(content)
}

// renderChordSelector draws every chord in one horizontal navigation row.
func (model Model) renderChordSelector() string {
	var output strings.Builder
	output.WriteString(model.styles.heading.Render("CHORD SELECTOR"))
	output.WriteString("\n\n")

	for index, name := range model.names {
		if index > 0 {
			output.WriteString("   ")
		}
		if index == model.selected {
			output.WriteString(model.styles.selected.Render("‹ " + name + " ›"))
		} else {
			output.WriteString(model.styles.normal.Render(name))
		}
	}

	return output.String()
}

// renderControls summarizes every available runtime command.
func (model Model) renderControls() string {
	return model.styles.heading.Render("CONTROLS") + "\n\n" +
		model.styles.muted.Render(
			"←/→ chord   ↑/↓ variation\n"+
				"f neck   t theme   w wave\n"+
				"? help   q quit",
		)
}

// renderChordDiagram draws the current chord heading, mode, and fretboard.
func (model Model) renderChordDiagram() string {
	if model.err != nil {
		return model.styles.err.Render(model.err.Error())
	}

	variationCount := model.catalog.VariationCount(model.voicing.Name)
	mode := "compact"
	if model.fullNeck {
		mode = fmt.Sprintf("full neck · frets 1–%d", fullNeckLastFret)
	}

	heading := fmt.Sprintf(
		"%s  ·  variation %d of %d",
		model.voicing.Name,
		model.voicing.Variation,
		variationCount,
	)
	header := model.styles.heading.Render("CHORD DIAGRAM") + "\n\n" +
		model.styles.heading.Render(heading) + "\n" +
		model.styles.accent.Render(mode)
	if model.fullNeck && model.width < fullNeckMinimumTerminalWidth {
		return header + "\n\n" + model.styles.err.Render(
			fmt.Sprintf(
				"Widen the terminal to at least %d columns (currently %d).",
				fullNeckMinimumTerminalWidth,
				model.width,
			),
		)
	}

	return header + "\n\n" +
		model.styles.normal.Render(renderFretboard(model.voicing, model.fullNeck))
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

// renderHelp lists every keyboard command in a dedicated panel.
func (model Model) renderHelp() string {
	return model.styles.heading.Render("KEYBOARD HELP") + "\n\n" +
		"← / h     Previous chord\n" +
		"→ / l     Next chord\n" +
		"↑ / k     Previous variation\n" +
		"↓ / j     Next variation\n" +
		"f         Toggle compact/full neck\n" +
		"t         Cycle color theme\n" +
		"w         Toggle detailed waveform\n" +
		"?         Open or close help\n" +
		"esc       Close help\n" +
		"q         Quit\n\n" +
		model.styles.muted.Render("Press ? or esc to return")
}
