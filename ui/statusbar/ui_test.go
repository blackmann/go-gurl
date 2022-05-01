package statusbar

import (
	"fmt"
	"github.com/blackmann/go-gurl/lib"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// Mostly view testing

func TestModel_Update_CommandInput(t *testing.T) {
	statusbar := NewStatusBar()
	command := "/annotate $1 ab"
	statusbar, _ = statusbar.Update(CommandInput(command))

	assert.True(t, strings.HasSuffix(statusbar.View(), fmt.Sprintf("%s : ", command)))
}

func TestModel_Update_Status_IDLE(t *testing.T) {
	statusbar := NewStatusBar()
	statusbar, _ = statusbar.Update(lib.IDLE)

	assert.Contains(t, statusbar.View(), "Idle")
}

func TestModel_Update_Status_PROCESSING(t *testing.T) {
	statusbar := NewStatusBar()
	statusbar, _ = statusbar.Update(lib.PROCESSING)

	assert.Contains(t, statusbar.View(), "Processing")
}

func TestModel_Update_CommandInput_Fit(t *testing.T) {
	statusbar := NewStatusBar()
	statusbar, _ = statusbar.Update(tea.WindowSizeMsg{Height: 30, Width: 70})

	statusbar, _ = statusbar.Update(lib.IDLE)
	cmd := "> @eg example.com"
	statusbar, _ = statusbar.Update(CommandInput(cmd))

	assert.Contains(t, statusbar.View(), cmd)
}

func TestModel_Update_CommandInput_Overflow(t *testing.T) {
	statusbar := NewStatusBar()
	statusbar, _ = statusbar.Update(tea.WindowSizeMsg{Height: 30, Width: 70})

	statusbar, _ = statusbar.Update(lib.IDLE)
	cmd := "> …mple https://api.example.com"
	statusbar, _ = statusbar.Update(CommandInput(cmd))

	assert.Contains(t, statusbar.View(), cmd)
}

func TestModel_Update_Message(t *testing.T) {
	statusbar := NewStatusBar()
	statusbar, _ = statusbar.Update(lib.IDLE)

	message := "140ms 34 B"
	statusbar, _ = statusbar.Update(lib.ShortMessage(message))

	assert.Contains(t, statusbar.View(), message)
}
