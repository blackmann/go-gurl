package statusbar

import (
	"github.com/blackmann/go-gurl/lib"
	tea "github.com/charmbracelet/bubbletea"
)

// A command in this case is an entry from the keyboard
// that is mapped to an action. For example, ":q" to quit, lib.
// The action is not necessarily performed by statusbar. We're
// only using this type as a state update message type.
//
//  statusbar.Update(commandInput(":q"))
//
// This is a tea.Msg type
type commandInput string

func UpdateFreetextCommand(command string) tea.Msg {
	return commandInput(command)
}

func UpdateStatus(status lib.Status) tea.Msg {
	return statusUpdate(status)
}
