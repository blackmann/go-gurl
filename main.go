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
	app := tea.NewProgram(newAppModel(), tea.WithAltScreen() /*tea.WithMouseCellMotion()*/)

	if err := app.Start(); err != nil {
		log.Panicln("Error occurred", err)
	}
}
