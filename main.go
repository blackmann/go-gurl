package main

import (
	"fmt"
	"github.com/blackmann/gurl/common/commands"
	"github.com/blackmann/gurl/common/status"
	"github.com/blackmann/gurl/handler"
	"github.com/blackmann/gurl/ui/addressbar"
	"github.com/blackmann/gurl/ui/statusbar"
	"github.com/blackmann/gurl/ui/viewport"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"os"
)

type keymap struct {
	nextTab           key.Binding
	quit              key.Binding
	toggleCommandMode key.Binding
}

func getDefaultKeyBinds() keymap {
	return keymap{
		nextTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("^tab", "Toggle Regions")),

		quit: key.NewBinding(key.WithKeys("ctrl+c")),

		toggleCommandMode: key.NewBinding(key.WithKeys("esc")),
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
	commandMode  bool
	command      string
}

func (m model) Init() tea.Cmd {
	return m.addressBar.Init()
}

func (m *model) resizeViewport(netHeight int, netWidth int) {
	statusBarHeight := lipgloss.Height(m.statusBar.View())
	addressBarHeight := lipgloss.Height(m.addressBar.View())

	m.viewport.Height = netHeight - (statusBarHeight + addressBarHeight)
	m.statusBar, _ = m.statusBar.Update(statusbar.Width(netWidth))
}

func (m model) submitRequest() tea.Msg {
	m.handler.MakeRequest()

	return commands.CompleteRequest
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.commandMode {
			if key.Matches(msg, m.keybinds.toggleCommandMode) {
				m.commandMode = false
				m.command = ""
				m.statusBar, _ = m.statusBar.Update(statusbar.Command(""))

				return m, nil
			}

			var cmd tea.Cmd

			// Only accept letters
			switch msg.Type {
			case tea.KeyRunes:
				m.command += string(msg.Runes[0])

			case tea.KeyBackspace:
				if len(m.command) > 0 {
					m.command = m.command[:len(m.command)-1]
				}

			case tea.KeyEnter:
				cmd = commands.CreateFreeCommand(m.command)
				m.command = ""
				m.commandMode = false
			}

			var prefix string
			if m.commandMode {
				prefix = "> "
			} else {
				prefix = ""
			}

			m.statusBar, _ = m.statusBar.Update(statusbar.Command(fmt.Sprintf("%s%s", prefix, m.command)))

			return m, cmd
		}

		switch {
		case key.Matches(msg, m.keybinds.nextTab):
			m.activeRegion = (m.activeRegion + 1) % 2 // only two views
			return m, nil

		case key.Matches(msg, m.keybinds.toggleCommandMode):
			m.commandMode = true
			m.statusBar, _ = m.statusBar.Update(statusbar.Command(">"))

			return m, nil

		case key.Matches(msg, m.keybinds.quit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.resizeViewport(msg.Height, msg.Width)
		return m, nil

	case commands.FreeCommand:
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)

		return m, cmd

	case commands.AppCommand:
		switch msg {
		case commands.NewRequest:
			m.statusBar = m.statusBar.SetStatus(status.PROCESSING)
			cmds = append(cmds, m.submitRequest)

		case commands.CompleteRequest:
			m.statusBar = m.statusBar.SetStatus(status.IDLE)
		}
	}

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
		addressBar: addressbar.NewAddressBar(),
		handler:    &h,
		keybinds:   getDefaultKeyBinds(),
		statusBar:  statusbar.NewStatusBar(),
		viewport:   viewport.NewViewport(),
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