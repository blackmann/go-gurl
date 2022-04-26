package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
)

func main() {
	// Set up logger
	f, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {
		log.Panicln("Failed to open log file")
	}

	defer f.Close()

	log.SetOutput(f)

	// Initialize and start app
	// TODO/DECIDE: Alternate between cell motion
	model, err := newAppModel()

	if err != nil {
		log.Panicln("Error occurred while creating model", err)
	}

	app := tea.NewProgram(model, tea.WithAltScreen() /*tea.WithMouseCellMotion()*/)

	if err := app.Start(); err != nil {
		log.Panicln("Error occurred", err)
	}
}
