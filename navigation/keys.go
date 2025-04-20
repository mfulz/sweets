package navigation

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines keybindings. It satisfies to the help.KeyMap interface, which
// is used to render the menu.
type KeyMap struct {
	// Keybindings used when selecting models.
	Next   key.Binding
	Prev   key.Binding
	Select key.Binding
	Leave  key.Binding

	// Quitting
	Quit key.Binding

	// Help toggle keybindings.
	ShowFullHelp  key.Binding
	CloseFullHelp key.Binding
}

var keyMap *KeyMap = &KeyMap{
	// Browsing.
	Next: key.NewBinding(
		key.WithKeys("right", "l", "n"),
		key.WithHelp("→/l/n", "next model"),
	),
	Prev: key.NewBinding(
		key.WithKeys("left", "h", "p"),
		key.WithHelp("←/h/p", "prev model"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select model"),
	),
	Leave: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "leave model"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q", "Quit"),
	),
	ShowFullHelp: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "more"),
	),
	CloseFullHelp: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "close help"),
	),
}

// DefaultKeyMap returns a default set of keybindings.
func DefaultKeyMap() *KeyMap {
	return keyMap
}
