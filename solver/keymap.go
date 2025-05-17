package solver

import "github.com/charmbracelet/bubbles/v2/key"

type keyMap struct {
	Up 							key.Binding
	Down 						key.Binding
	Left 						key.Binding
	Right 					key.Binding
	Delete					key.Binding
	Quit 						key.Binding
	NextClue 				key.Binding
	PrevClue 				key.Binding
	ToggleDirection key.Binding
	ViewPreferences key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right"),
	),
	Delete: key.NewBinding(
		key.WithKeys("backspace", "delete"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	NextClue: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next clue"),
	),
	PrevClue: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "previous clue"),
	),
	ToggleDirection: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "change direction"),
	),
	ViewPreferences: key.NewBinding(
		key.WithKeys("ctrl+p"),
		key.WithHelp("ctrl+p", "change preferences"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NextClue, k.PrevClue, k.ToggleDirection, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextClue, k.PrevClue}, 
		{k.ToggleDirection},
		{k.ViewPreferences, k.Quit},
	}
}