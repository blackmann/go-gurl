package main

import (
	"fmt"
	"github.com/blackmann/gurl/lib"
	"github.com/blackmann/gurl/ui/addressbar"
	"github.com/blackmann/gurl/ui/statusbar"
	"github.com/blackmann/gurl/ui/viewport"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"log"
	"net/http"
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
	height   int

	// Views
	addressBar addressbar.Model
	viewport   viewport.Model
	statusBar  statusbar.Model

	// State
	activeRegion int
	commandMode  bool
	command      string
	enabled      bool
}

func (m model) Init() tea.Cmd {
	return m.addressBar.Init()
}

func (m model) submitRequest(address lib.Address, headers http.Header, body io.Reader) tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{}
		req, err := http.NewRequest(address.Method, address.Url, body)
		req.Header = headers

		res, err := client.Do(req)
		defer res.Body.Close()

		if err != nil {
			log.Println("Error occurred", err)
			return nil
		}

		body, err := io.ReadAll(res.Body)

		if err != nil {
			log.Println("Error occurred while reading response body", err)
			return nil
		}

		return lib.Response{Body: body, Headers: res.Header}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.commandMode {
			gainFocus := func() tea.Msg {
				return lib.GainFocus
			}

			if key.Matches(msg, m.keybinds.toggleCommandMode) {
				m.commandMode = false
				m.command = ""
				m.statusBar, _ = m.statusBar.Update(statusbar.FreeTextCommandCmd(""))

				return m, gainFocus
			}

			switch msg.Type {
			case tea.KeyRunes:
				m.command += string(msg.Runes[0])

			case tea.KeyBackspace:
				if len(m.command) > 0 {
					m.command = m.command[:len(m.command)-1]
				}

			case tea.KeyEnter:
				cmd := getFreeTextCommand(m.command)
				m.command = ""
				m.commandMode = false

				m.statusBar, _ = m.statusBar.Update(statusbar.FreeTextCommandCmd(""))

				return m, tea.Batch(cmd, gainFocus)
			}

			var prefix string
			if m.commandMode {
				prefix = "> "
			} else {
				prefix = ""
			}

			m.statusBar, _ = m.statusBar.Update(statusbar.FreeTextCommandCmd(fmt.Sprintf("%s%s", prefix, m.command)))

			return m, nil
		}

		switch {
		case key.Matches(msg, m.keybinds.nextTab):
			m.activeRegion = (m.activeRegion + 1) % 2 // only two views
			return m, nil

		case key.Matches(msg, m.keybinds.toggleCommandMode):
			m.commandMode = true
			m.statusBar, _ = m.statusBar.Update(statusbar.FreeTextCommandCmd(">"))

			return m, func() tea.Msg {
				return lib.LostFocus
			}

		case key.Matches(msg, m.keybinds.quit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.height = msg.Height
		statusBarHeight := lipgloss.Height(m.statusBar.View())
		addressBarHeight := lipgloss.Height(m.addressBar.View())

		viewportHeight := msg.Height - (statusBarHeight + addressBarHeight)

		m.viewport, _ = m.viewport.Update(tea.WindowSizeMsg{Height: viewportHeight, Width: msg.Width})

		return m, nil

	case lib.FreeText:
		m.viewport, _ = m.viewport.Update(msg)
		return m, nil

	case lib.Trigger:
		switch msg {
		case lib.NewRequest:
			var cmds []tea.Cmd
			var cmd tea.Cmd

			m.statusBar, cmd = m.statusBar.Update(statusbar.StatusCmd(lib.PROCESSING))
			cmds = append(cmds, cmd)

			headers := m.viewport.GetHeaders()
			body := m.viewport.GetBody()

			address, err := m.addressBar.GetAddress()

			if err != nil {
				// TODO: Return error message (msg)
				log.Panicln("Error parsing address")
			}

			cmds = append(cmds, m.submitRequest(address, headers, body))

			m.enabled = false

			return m, tea.Batch(cmds...)
		}

	case lib.Response:
		var cmds []tea.Cmd
		var cmd tea.Cmd

		m.viewport, cmd = m.viewport.Update(m.viewport.SetResponse(msg))
		cmds = append(cmds, cmd)

		m.statusBar, _ = m.statusBar.Update(statusbar.StatusCmd(lib.IDLE))
		m.enabled = true

		return m, tea.Batch(cmds...)
	}

	var cmds []tea.Cmd

	// Forward the unhandled command to the active region
	switch m.activeRegion {
	case 0:
		var cmd tea.Cmd
		m.addressBar, cmd = m.addressBar.Update(msg)

		cmds = append(cmds, cmd)

	case 1:
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)

		cmds = append(cmds, cmd)
	}

	// For ticks
	var cmd tea.Cmd
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
	return model{
		addressBar: addressbar.NewAddressBar(),
		keybinds:   getDefaultKeyBinds(),
		statusBar:  statusbar.NewStatusBar(),
		viewport:   viewport.NewViewport(),

		enabled: true,
	}
}

func getFreeTextCommand(cmd string) tea.Cmd {
	return func() tea.Msg {
		return lib.FreeText(cmd)
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
