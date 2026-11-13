// Package tui owns Alt-Tab's interactive terminal presentation.
package tui

import (
	"errors"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/jake-mckenzie/alt-tab/internal/chords"
)

// Model holds the terminal UI state.
type Model struct {
	player      AudioPlayer
	playing     bool
	held        string
	keyReleases bool
	catalog     chords.Catalog
	families    []chordFamily
	selected    int
	chordKind   int
	voicing     chords.Voicing
	fullNeck    bool
	tabNotation bool
	showHelp    bool
	dark        bool
	theme       int
	waveform    bool
	width       int
	height      int
	scroll      int
	err         error
	styles      styles
}

const (
	accidentalChord = -1
	baseChord       = 0
	minorChord      = 1
)

// chordFamily groups one natural chord with its available accidental and minor.
type chordFamily struct {
	base       string
	accidental string
	minor      string
}

// New returns the initial terminal UI model.
func New(catalog chords.Catalog) Model {
	model := Model{
		catalog:  catalog,
		dark:     true,
		waveform: true,
		width:    100,
		height:   100,
		styles:   newStyles(true, 0),
	}

	if catalog == nil {
		model.err = errors.New("chord catalog is unavailable")
		return model
	}

	model.families = buildChordFamilies(catalog.Names())
	model.loadSelection(1)
	return model
}

// AudioPlayer is the platform-independent output contract used by the UI.
type AudioPlayer interface{ Set(chords.Voicing, bool) }

// WithAudio enables playback without coupling navigation to a platform backend.
func (model Model) WithAudio(player AudioPlayer) Model {
	model.player = player
	return model
}

// setPlaying updates playback only when audio and a valid voicing are available.
func (model *Model) setPlaying(active bool) {
	model.playing = active && model.player != nil && model.err == nil
	if model.player != nil {
		model.player.Set(model.voicing, model.playing)
	}
}

// Init asks Bubble Tea to report the terminal background color.
func (Model) Init() tea.Cmd {
	return tea.RequestBackgroundColor
}

// Update handles application-level key events.
func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.KeyboardEnhancementsMsg:
		model.keyReleases = message.SupportsEventTypes()
	case tea.BlurMsg:
		model.held = ""
		model.setPlaying(false)
	case tea.KeyReleaseMsg:
		if message.String() == model.held {
			model.held = ""
			model.setPlaying(false)
		}
	case tea.BackgroundColorMsg:
		model.dark = message.IsDark()
		model.styles = newStyles(model.dark, model.theme)
	case tea.WindowSizeMsg:
		model.width = max(1, message.Width)
		model.height = max(1, message.Height)
		model.scroll = min(model.scroll, model.maximumScroll())
	case tea.MouseWheelMsg:
		switch message.Button {
		case tea.MouseWheelUp:
			model.moveViewport(-3)
		case tea.MouseWheelDown:
			model.moveViewport(3)
		}
	case tea.KeyPressMsg:
		key := message.String()
		if message.IsRepeat && (key == "," || key == "." || key == "p" || key == "space" || key == " " || (len(key) == 1 && key[0] >= 'a' && key[0] <= 'g')) {
			return model, nil
		}
		if key == "ctrl+c" || key == "q" {
			model.setPlaying(false)
			return model, tea.Quit
		}
		if key == "t" {
			model.theme = wrap(model.theme+1, len(palettes))
			model.styles = newStyles(model.dark, model.theme)
			return model, nil
		}
		if key == "w" {
			model.waveform = !model.waveform
			return model, nil
		}
		if model.showHelp {
			if key == "esc" || key == "?" {
				model.showHelp = false
			}
			return model, nil
		}

		if (key == "space" || key == " " || key == "p") && model.err == nil {
			model.held = ""
			model.setPlaying(!model.playing)
			return model, nil
		}
		if len(key) == 1 && key[0] >= 'a' && key[0] <= 'g' {
			if model.held == key {
				return model, nil
			}
			wasPlaying := model.playing && model.currentChordName() == strings.ToUpper(key)
			for i, family := range model.families {
				if family.base == strings.ToUpper(key) {
					model.selected, model.chordKind = i, baseChord
					model.loadSelection(max(1, model.voicing.Number))
					if model.keyReleases {
						model.held = key
						model.setPlaying(true)
					} else {
						model.setPlaying(!wasPlaying)
					}
					break
				}
			}
			return model, nil
		}
		if len(key) == 1 && key[0] >= '1' && key[0] <= '9' {
			number := int(key[0] - '0')
			if model.catalog != nil && number <= model.catalog.VoicingCount(model.currentChordName()) {
				model.loadSelection(number)
			}
			return model, nil
		}
		switch key {
		case ",":
			model.toggleChordKind(minorChord)
		case ".":
			model.toggleChordKind(accidentalChord)
		case "left", "h":
			model.moveChord(-1)
		case "right", "l":
			model.moveChord(1)
		case "x":
			model.fullNeck = !model.fullNeck
			if model.fullNeck {
				model.tabNotation = false
			}
		case "n":
			model.tabNotation = !model.tabNotation
			if model.tabNotation {
				model.fullNeck = false
			}
		case "?":
			model.showHelp = true
			model.held = ""
			model.setPlaying(false)
		}
	}

	return model, nil
}

// buildChordFamilies preserves natural-chord order and attaches known variants.
func buildChordFamilies(names []string) []chordFamily {
	families := make([]chordFamily, 0, len(names))
	indexes := make(map[string]int)
	for _, name := range names {
		if len(name) != 1 || name[0] < 'A' || name[0] > 'G' {
			continue
		}
		indexes[name] = len(families)
		families = append(families, chordFamily{base: name})
	}

	for _, name := range names {
		if len(name) != 2 {
			continue
		}
		index, exists := indexes[name[:1]]
		if !exists {
			continue
		}
		switch name[1] {
		case 'b', '#':
			families[index].accidental = name
		case 'm':
			families[index].minor = name
		}
	}
	return families
}

// currentChordName resolves the highlighted dial position to a catalog name.
func (model Model) currentChordName() string {
	if len(model.families) == 0 {
		return ""
	}
	family := model.families[model.selected]
	switch model.chordKind {
	case accidentalChord:
		return family.accidental
	case minorChord:
		return family.minor
	default:
		return family.base
	}
}

// View renders the application in a resize-safe, scrollable full-screen buffer.
func (model Model) View() tea.View {
	view := tea.NewView(model.renderViewport())
	if background := paletteAt(model.theme).background; background != "" {
		view.BackgroundColor = lipgloss.Color(background)
	}
	view.WindowTitle = "alt-tab"
	view.AltScreen = true
	view.ReportFocus = true
	view.KeyboardEnhancements.ReportEventTypes = true
	view.MouseMode = tea.MouseModeCellMotion
	return view
}

// renderViewport clips the document to the terminal and applies wheel scrolling.
func (model Model) renderViewport() string {
	lines := strings.Split(model.render(), "\n")
	height := max(1, model.height)
	start := min(model.scroll, max(0, len(lines)-height))
	end := min(len(lines), start+height)
	return strings.Join(lines[start:end], "\n")
}

// maximumScroll returns the last offset that still fills the viewport.
func (model Model) maximumScroll() int {
	return max(0, lipgloss.Height(model.render())-max(1, model.height))
}

// moveViewport scrolls without changing the selected chord or voicing.
func (model *Model) moveViewport(delta int) {
	model.scroll = min(max(0, model.scroll+delta), model.maximumScroll())
}

// loadSelection refreshes the active voicing after navigation.
func (model *Model) loadSelection(number int) {
	if len(model.families) == 0 {
		model.err = errors.New("no chords are available")
		return
	}

	voicing, err := model.catalog.Load(model.currentChordName(), number)
	model.voicing = voicing
	model.err = err
	if model.player != nil {
		model.setPlaying(model.playing)
	}
}

// moveChord rotates the base-chord dial and selects its first voicing.
func (model *Model) moveChord(delta int) {
	if len(model.families) == 0 {
		return
	}

	model.selected = wrap(model.selected+delta, len(model.families))
	model.chordKind = baseChord
	model.loadSelection(1)
}

// toggleChordKind swaps an available variant with the natural major chord.
func (model *Model) toggleChordKind(kind int) {
	if len(model.families) == 0 {
		return
	}
	family := model.families[model.selected]
	if (kind == minorChord && family.minor == "") || (kind == accidentalChord && family.accidental == "") {
		return
	}
	if model.chordKind == kind {
		kind = baseChord
	}
	model.chordKind = kind
	model.loadSelection(1)
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
