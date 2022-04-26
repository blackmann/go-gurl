package lib

import "github.com/charmbracelet/bubbles/key"

type Keymap struct {
	NextTab           key.Binding
	Quit              key.Binding
	ToggleCommandMode key.Binding
}

func DefaultKeyBinds() Keymap {
	return Keymap{
		NextTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("^tab", "Toggle Regions")),

		Quit: key.NewBinding(key.WithKeys("ctrl+c")),

		ToggleCommandMode: key.NewBinding(key.WithKeys("esc")),
	}
}
