package main

import (
	"fmt"
	"github.com/blackmann/gurl/lib"
	"github.com/blackmann/gurl/ui/addressbar"
	"github.com/blackmann/gurl/ui/bookmarks"
	"github.com/blackmann/gurl/ui/history"
	"github.com/blackmann/gurl/ui/statusbar"
	"github.com/blackmann/gurl/ui/viewport"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"log"
	"strings"
	"time"
)

type middleView int

var (
	VIEWPORT  middleView = 1
	HISTORY   middleView = 2
	BOOKMARKS middleView = 3
)

type model struct {
	// Config
	client      lib.Client
	keybinds    lib.Keymap
	persistence lib.Persistence

	// Views
	addressBar    addressbar.Model
	viewport      viewport.Model
	statusBar     statusbar.Model
	historyList   history.Model
	bookmarksList bookmarks.Model

	// State
	activeRegion int
	commandMode  bool
	command      string
	enabled      bool
	height       int
	middleView
}

func (m model) Init() tea.Cmd {
	return m.addressBar.Init()
}

func (m model) submitRequest(request lib.Request) tea.Cmd {
	return func() tea.Msg {
		if response, err := m.client.MakeRequest(request); err == nil {
			return response
		} else {
			// TODO: Return with good messaging
			return err
		}
	}
}

func (m model) getMode() lib.Mode {
	if m.activeRegion == 0 {
		return lib.Url
	}

	return lib.Detail
}

func (m model) performResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.height = msg.Height
	statusBarHeight := lipgloss.Height(m.statusBar.View())
	addressBarHeight := lipgloss.Height(m.addressBar.View())

	viewportHeight := msg.Height - (statusBarHeight + addressBarHeight)

	middleViewSize := tea.WindowSizeMsg{
		Height: viewportHeight,
		Width:  msg.Width,
	}

	m.viewport, _ = m.viewport.Update(middleViewSize)

	m.historyList, _ = m.historyList.Update(middleViewSize)
	m.bookmarksList, _ = m.bookmarksList.Update(middleViewSize)

	m.statusBar, _ = m.statusBar.Update(tea.WindowSizeMsg{Width: msg.Width})

	return m, nil
}

func (m model) handleResponse(msg lib.Response) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.viewport, cmd = m.viewport.Update(m.viewport.SetResponse(msg))
	cmds = append(cmds, cmd)

	m.statusBar, _ = m.statusBar.Update(statusbar.UpdateStatus(lib.IDLE))
	m.statusBar, _ = m.statusBar.Update(statusbar.UpdateStatus(lib.Status(msg.Status)))
	m.statusBar, _ = m.statusBar.Update(
		lib.ShortMessage(fmt.Sprintf("%dms %s", msg.Time, humanize.Bytes(uint64(len(msg.Body))))),
	)

	m.enabled = true

	go m.persistence.SaveHistory(lib.History{
		Url:    msg.Request.Address.Url,
		Method: msg.Request.Address.Method,
		Date:   time.Now(),
	})

	return m, tea.Batch(cmds...)
}

func (m model) handleTrigger(msg lib.Trigger) (tea.Model, tea.Cmd, bool) {
	switch msg {
	case lib.NewRequest:
		var cmds []tea.Cmd
		var cmd tea.Cmd

		// Status bar returns tick cmd when processing
		m.statusBar, _ = m.statusBar.Update(lib.ShortMessage(""))
		m.statusBar, cmd = m.statusBar.Update(statusbar.UpdateStatus(lib.PROCESSING))
		cmds = append(cmds, cmd)

		headers := m.viewport.GetHeaders()
		body := m.viewport.GetBody()

		address, err := m.addressBar.GetAddress()

		if err != nil {
			// TODO: Return error message (msg)
			log.Panicln("Error parsing address")
		}

		cmds = append(cmds, m.submitRequest(lib.Request{
			Address: address,
			Headers: headers,
			Body:    strings.NewReader(body),
		}))

		m.enabled = false

		return m, tea.Batch(cmds...), true

	default:
		m.viewport, _ = m.viewport.Update(msg)
		return m, nil, true
	}

	return nil, nil, false
}

func (m model) handeCommandInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	gainFocus := func() tea.Msg {
		return lib.GainFocus
	}

	if key.Matches(msg, m.keybinds.ToggleCommandMode) {
		m.commandMode = false
		m.command = ""
		m.statusBar, _ = m.statusBar.Update(statusbar.UpdateFreetextCommand(""))
		m.statusBar, _ = m.statusBar.Update(m.getMode())

		return m, gainFocus
	}

	switch msg.Type {
	case tea.KeyRunes:
		m.command += string(msg.Runes[0])

	case tea.KeyBackspace:
		if len(m.command) > 0 {
			m.command = m.command[:len(m.command)-1]
		}

	case tea.KeyRight:
		if len(m.command) == 0 {
			return m, lib.NavigateRight
		}

	case tea.KeyLeft:
		if len(m.command) == 0 {
			return m, lib.NavigateLeft
		}

	case tea.KeyEnter:
		cmd := func() tea.Msg { return lib.FreeTextCommand(m.command) }
		m.command = ""
		m.commandMode = false

		m.statusBar, _ = m.statusBar.Update(statusbar.UpdateFreetextCommand(""))

		return m, tea.Batch(cmd, gainFocus)
	}

	var prefix string
	if m.commandMode {
		prefix = "> "
	} else {
		prefix = ""
	}

	m.statusBar, _ = m.statusBar.Update(statusbar.
		UpdateFreetextCommand(fmt.Sprintf("%s%s", prefix, m.command)))

	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: Quick successive `esc` input should yield a reset state
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.commandMode {
			return m.handeCommandInput(msg)
		}

		switch {
		case key.Matches(msg, m.keybinds.NextTab):
			if m.middleView != VIEWPORT {
				// Showing history so no alternating between views
				return m, nil
			}

			m.activeRegion = (m.activeRegion + 1) % 2 // only two views
			m.statusBar, _ = m.statusBar.Update(m.getMode())

			// So that inputs receive "focus"
			return m, textinput.Blink

		case key.Matches(msg, m.keybinds.ToggleCommandMode):
			if m.middleView != VIEWPORT {
				m.middleView = VIEWPORT

				return m, nil
			}

			m.commandMode = true
			m.statusBar, _ = m.statusBar.Update(statusbar.UpdateFreetextCommand(">"))
			m.statusBar, _ = m.statusBar.Update(lib.Cmd)

			return m, func() tea.Msg {
				return lib.LostFocus
			}

		case key.Matches(msg, m.keybinds.Quit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		return m.performResize(msg)

	case lib.FreeTextCommand:
		m.viewport, _ = m.viewport.Update(msg)
		return m, nil

	case lib.Trigger:
		t, cmd, done := m.handleTrigger(msg)
		if done {
			return t, cmd
		}

	case lib.Response:
		return m.handleResponse(msg)
	}

	var cmds []tea.Cmd

	// Forward the unhandled/trickled command to the active region
	switch m.activeRegion {
	case 0:
		var cmd tea.Cmd
		m.addressBar, cmd = m.addressBar.Update(msg)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			if m.middleView != VIEWPORT {
				if msg.Type == tea.KeyUp || msg.Type == tea.KeyDown {
					switch m.middleView {
					case HISTORY:
						m.historyList, _ = m.historyList.Update(msg)
					case BOOKMARKS:
						m.bookmarksList, _ = m.bookmarksList.Update(msg)
					}
				}
			}
		default:
			{
				addressEntry := m.addressBar.GetEntry()
				if strings.HasPrefix(addressEntry, "$") {
					m.middleView = HISTORY
					m.historyList, _ = m.historyList.Update(history.Filter(addressEntry[1:]))
				} else if strings.HasPrefix(addressEntry, "@") {
					m.middleView = BOOKMARKS
					m.bookmarksList, _ = m.bookmarksList.Update(bookmarks.Filter(addressEntry[1:]))
				} else {
					m.middleView = VIEWPORT
				}
			}
		}

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
	var middleView string

	switch m.middleView {
	case HISTORY:
		middleView = m.historyList.View()
	case BOOKMARKS:
		middleView = m.bookmarksList.View()
	default:
		middleView = m.viewport.View()
	}

	return fmt.Sprintf("%s\n%s\n%s",
		m.addressBar.View(),
		middleView,
		m.statusBar.View())
}

func newAppModel() (model, error) {
	db, err := lib.NewDbPersistence()
	if err != nil {
		return model{}, err
	}

	return model{
		addressBar:  addressbar.NewAddressBar(),
		client:      lib.NewHttpClient(),
		keybinds:    lib.DefaultKeyBinds(),
		statusBar:   statusbar.NewStatusBar(),
		viewport:    viewport.NewViewport(),
		persistence: &db,
		historyList: history.NewHistory(&db),

		enabled: true,
	}, nil
}
