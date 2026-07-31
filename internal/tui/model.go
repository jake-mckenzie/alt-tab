// Package tui owns Alt-Tab's interactive terminal presentation.
package tui

import tea "charm.land/bubbletea/v2"

// Model holds the terminal UI state.
type Model struct{}

// New returns the initial terminal UI model.
func New() Model {
	return Model{}
}

// Init performs no startup I/O while the backend adapter is being introduced.
func (Model) Init() tea.Cmd {
	return nil
}

// Update handles application-level key events.
func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.KeyPressMsg:
		switch message.String() {
		case "ctrl+c", "q":
			return model, tea.Quit
		}
	}

	return model, nil
}

// View renders the initial full-screen application shell.
func (Model) View() tea.View {
	view := tea.NewView(
		"ALT-TAB\n" +
			"Guitar Chord Viewer\n\n" +
			"Bubble Tea interface scaffold\n\n" +
			"q  quit\n",
	)
	view.AltScreen = true
	view.WindowTitle = "Alt-Tab"
	return view
}
