package main

import (
	"fmt"
	"github.com/blackmann/gurl/handler"
	"github.com/blackmann/gurl/ui/addressbar"
	"github.com/blackmann/gurl/ui/statusbar"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"os"
)

type keymap struct {
	nextTab key.Binding
	quit    key.Binding
}

func getDefaultKeyBinds() keymap {
	return keymap{
		nextTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("^tab", "Toggle Regions")),

		quit: key.NewBinding(key.WithKeys("ctrl+c")),
	}
}

type model struct {
	// Config
	keybinds keymap
	handler  *handler.RequestHandler

	// Views
	addressBar addressbar.Model
	viewport   viewport.Model
	statusBar  statusbar.Model

	// State
	activeRegion int
}

func (m model) Init() tea.Cmd {
	return m.addressBar.Init()
}

func (m *model) resizeViewport(netHeight int) {
	statusBarHeight := lipgloss.Height(m.statusBar.View())
	addressBarHeight := lipgloss.Height(m.addressBar.View())

	m.viewport.Height = netHeight - (statusBarHeight + addressBarHeight)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keybinds.nextTab):
			m.activeRegion = (m.activeRegion + 1) % 2 // only two views
			return m, nil

		case key.Matches(msg, m.keybinds.quit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.resizeViewport(msg.Height)
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch m.activeRegion {
	// Request bar
	case 0:
		m.addressBar, cmd = m.addressBar.Update(msg)
		cmds = append(cmds, cmd)

	// Viewport
	case 1:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.statusBar, cmd = m.statusBar.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return fmt.Sprintf("%s\n%s\n%s",
		m.addressBar.View(),
		m.viewport.View(),
		m.statusBar.View())
}

func newAppModel() model {
	h := handler.NewRequestHandler()

	return model{
		addressBar: addressbar.NewAddressBar(&h),
		handler:    &h,
		keybinds:   getDefaultKeyBinds(),
		statusBar:  statusbar.NewStatusBar(&h),
	}
}

func main() {
	// Set up logger
	f, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {
		log.Panicln("Failed to open log file")
	}

	defer f.Close()

	log.SetOutput(f)

	// Initialize and start app
	app := tea.NewProgram(newAppModel(), tea.WithAltScreen())

	if err := app.Start(); err != nil {
		log.Panicln("Error occurred", err)
	}
}
