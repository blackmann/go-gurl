package lib

import "github.com/charmbracelet/bubbletea"

func EnterString(m tea.Model, keys string) tea.Model {
	for _, k := range keys {
		hKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{k}}
		m, _ = m.Update(hKey)
	}

	return m
}
