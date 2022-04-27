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
	"strconv"
	"strings"
	"time"
)

type middleView int

var (
	VIEWPORT  middleView = 1
	HISTORY   middleView = 2
	BOOKMARKS middleView = 3
)

type region int

var (
	URL region = 1
	VP  region = 2
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
	activeRegion region
	commandMode  bool
	command      string
	enabled      bool
	height       int
	middleView
}

func (m model) annotate(argsString string) (model, tea.Cmd) {
	args := strings.Split(argsString, " ")
	log.Println("Annotating", args)
	id := args[0][1:] // without the $
	annotation := args[1]

	intId, _ := strconv.ParseInt(id, 10, 64)
	m.persistence.AnnotateHistory(intId, annotation)

	updateHistory := func() tea.Msg {
		return lib.UpdateHistory
	}

	return m, updateHistory
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
	if m.activeRegion == URL {
		return lib.Url
	}

	return lib.Detail
}

func (m *model) resetCommandInput() {
	m.command = ""
	m.commandMode = false

	m.statusBar, _ = m.statusBar.Update(statusbar.UpdateFreetextCommand(""))
}

func (m model) performResize(msg tea.WindowSizeMsg) (model, tea.Cmd) {
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

func (m model) handleHistory(history lib.History) (model, tea.Cmd) {
	m.viewport, _ = m.viewport.Update(history)
	m.addressBar, _ = m.addressBar.Update(lib.Address{Url: history.Url, Method: history.Method})

	m.middleView = VIEWPORT

	return m, nil
}

func (m model) handleResponse(msg lib.Response) (model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.viewport, cmd = m.viewport.Update(m.viewport.SetResponse(msg))
	cmds = append(cmds, cmd)

	m.statusBar, _ = m.statusBar.Update(statusbar.UpdateStatus(lib.IDLE))
	m.statusBar, _ = m.statusBar.Update(statusbar.UpdateStatus(lib.Status(msg.Status)))
	m.statusBar, _ = m.statusBar.Update(
		// TODO: humanize the time
		lib.ShortMessage(fmt.Sprintf("%dms %s", msg.Time, humanize.Bytes(uint64(len(msg.Body))))),
	)

	m.enabled = true

	var headers = make(map[string][]string)

	for k, v := range msg.Request.Headers {
		headers[k] = v
	}

	m.persistence.SaveHistory(lib.History{
		Url:     msg.Request.Address.Url,
		Method:  msg.Request.Address.Method,
		Date:    time.Now(),
		Status:  msg.Status,
		Headers: headers,
		Body:    msg.Request.Body,
	})

	newHistory := func() tea.Msg { return lib.UpdateHistory }
	cmds = append(cmds, newHistory)

	return m, tea.Batch(cmds...)
}

func (m model) handleTrigger(msg lib.Trigger) (model, tea.Cmd, bool) {
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
			log.Println("Error parsing address", err)
			return m, nil, true
		}

		cmds = append(cmds, m.submitRequest(lib.Request{
			Address: address,
			Headers: headers,
			Body:    body,
		}))

		m.enabled = false

		return m, tea.Batch(cmds...), true

	case lib.UpdateHistory:
		m.historyList, _ = m.historyList.Update(msg)
		return m, nil, true

	default:
		m.viewport, _ = m.viewport.Update(msg)
		return m, nil, true
	}
}

func (m model) handeCommandInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	gainFocus := func() tea.Msg {
		return lib.GainFocus
	}

	if key.Matches(msg, m.keybinds.ToggleCommandMode) {
		m.resetCommandInput()
		m.statusBar, _ = m.statusBar.Update(m.getMode())
		return m, gainFocus
	}

	if key.Matches(msg, m.keybinds.Quit) {
		return m, tea.Quit
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
		command := m.command
		cmd := func() tea.Msg { return lib.FreeTextCommand(command) }
		m.resetCommandInput()

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

func (m model) Init() tea.Cmd {
	return m.addressBar.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

			if m.activeRegion == URL {
				m.activeRegion = VP
			} else {
				m.activeRegion = URL
			}

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
		parts := strings.SplitN(string(msg), " ", 2)
		for i, v := range parts {
			parts[i] = strings.Trim(v, " ")
		}

		switch parts[0] {
		case "/annotate":
			return m.annotate(parts[1])
		}

		m.resetCommandInput()

		return m, nil

	case lib.Trigger:
		t, cmd, done := m.handleTrigger(msg)
		if done {
			return t, cmd
		}

	case lib.Response:
		return m.handleResponse(msg)

	case lib.History:
		return m.handleHistory(msg)
	}

	var cmds []tea.Cmd

	// Forward the unhandled/trickled command to the active region
	switch m.activeRegion {
	case URL:
		{
			var cmd tea.Cmd
			m.addressBar, cmd = m.addressBar.Update(msg)

			cmds = append(cmds, cmd)

			switch msg := msg.(type) {
			case tea.KeyMsg:
				{
					if m.middleView != VIEWPORT {
						if msg.Type == tea.KeyUp || msg.Type == tea.KeyDown {
							switch m.middleView {
							case HISTORY:
								m.historyList, _ = m.historyList.Update(msg)
								return m, nil

							case BOOKMARKS:
								m.bookmarksList, _ = m.bookmarksList.Update(msg)
								return m, nil
							}
						}
					}

					addressEntry := m.addressBar.GetEntry()

					if (isHistoryEntry(addressEntry) || isBookmarkEntry(addressEntry)) &&
						msg.Type == tea.KeyEnter {
						// If still typing a bookmark or address, then  remove the Trigger (New Request)
						cmds = cmds[:len(cmds)-1]

						if m.middleView == HISTORY {
							selected := m.historyList.GetSelected()

							prefillWithHistory := func() tea.Msg { return selected }

							cmds = append(cmds, prefillWithHistory)
						}
					}

					if isHistoryEntry(addressEntry) {
						m.middleView = HISTORY
						m.historyList, _ = m.historyList.Update(history.Filter(addressEntry[1:]))
					} else if isBookmarkEntry(addressEntry) {
						m.middleView = BOOKMARKS
						m.bookmarksList, _ = m.bookmarksList.Update(bookmarks.Filter(addressEntry[1:]))
					} else {
						m.middleView = VIEWPORT
					}
				}
			}
		}

	case VP:
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
		addressBar:   addressbar.NewAddressBar(),
		client:       lib.NewHttpClient(),
		keybinds:     lib.DefaultKeyBinds(),
		statusBar:    statusbar.NewStatusBar(),
		viewport:     viewport.NewViewport(),
		persistence:  &db,
		historyList:  history.NewHistory(&db),
		activeRegion: URL,
		middleView:   VIEWPORT,

		enabled: true,
	}, nil
}

func isHistoryEntry(s string) bool {
	return strings.HasPrefix(s, "$")
}

func isBookmarkEntry(s string) bool {
	return strings.HasPrefix(s, "@")
}
