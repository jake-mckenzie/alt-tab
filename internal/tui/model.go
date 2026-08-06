// Package tui owns Alt-Tab's interactive terminal presentation.
package tui

import (
	"errors"

	tea "charm.land/bubbletea/v2"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
)

// Model holds the terminal UI state.
type Model struct {
	catalog   chords.Catalog
	names     []string
	selected  int
	variation int
	voicing   chords.Voicing
	fullNeck  bool
	showHelp  bool
	dark      bool
	theme     int
	width     int
	err       error
	styles    styles
}

// New returns the initial terminal UI model.
func New(catalog chords.Catalog) Model {
	model := Model{
		catalog:   catalog,
		variation: 1,
		dark:      true,
		width:     100,
		styles:    newStyles(true, 0),
	}

	if catalog == nil {
		model.err = errors.New("chord catalog is unavailable")
		return model
	}

	model.names = catalog.Names()
	model.loadSelection()
	return model
}

// Init asks Bubble Tea to report the terminal background color.
func (Model) Init() tea.Cmd {
	return tea.RequestBackgroundColor
}

// Update handles application-level key events.
func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.BackgroundColorMsg:
		model.dark = message.IsDark()
		model.styles = newStyles(model.dark, model.theme)
	case tea.WindowSizeMsg:
		model.width = message.Width
	case tea.KeyPressMsg:
		key := message.String()
		if key == "ctrl+c" || key == "q" {
			return model, tea.Quit
		}
		if key == "t" {
			model.theme = wrap(model.theme+1, len(palettes))
			model.styles = newStyles(model.dark, model.theme)
			return model, nil
		}
		if model.showHelp {
			if key == "esc" || key == "?" {
				model.showHelp = false
			}
			return model, nil
		}

		switch key {
		case "up", "k":
			model.moveVariation(-1)
		case "down", "j":
			model.moveVariation(1)
		case "left", "h":
			model.moveChord(-1)
		case "right", "l":
			model.moveChord(1)
		case "f":
			model.fullNeck = !model.fullNeck
		case "?":
			model.showHelp = true
		}
	}

	return model, nil
}

// View renders the full-screen application.
func (model Model) View() tea.View {
	view := tea.NewView(model.render())
	view.AltScreen = true
	view.WindowTitle = "Alt-Tab"
	return view
}

// loadSelection refreshes the active voicing after navigation.
func (model *Model) loadSelection() {
	if len(model.names) == 0 {
		model.err = errors.New("no chords are available")
		return
	}

	voicing, err := model.catalog.Load(model.names[model.selected], model.variation)
	model.voicing = voicing
	model.err = err
}

// moveChord wraps through chord names and selects their first variation.
func (model *Model) moveChord(delta int) {
	if len(model.names) == 0 {
		return
	}

	model.selected = wrap(model.selected+delta, len(model.names))
	model.variation = 1
	model.loadSelection()
}

// moveVariation wraps through the active chord's available voicings.
func (model *Model) moveVariation(delta int) {
	if len(model.names) == 0 {
		return
	}

	count := model.catalog.VariationCount(model.names[model.selected])
	if count == 0 {
		return
	}

	// Convert the public one-based variation number around zero-based wrapping.
	model.variation = wrap(model.variation-1+delta, count) + 1
	model.loadSelection()
}

// wrap confines a possibly negative index to a non-empty collection.
func wrap(value, size int) int {
	if size <= 0 {
		return 0
	}

	value %= size
	if value < 0 {
		value += size
	}
	return value
}
