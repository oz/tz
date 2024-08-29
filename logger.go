package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var logger *log.Logger

type NoOpWriter struct{}

func (w *NoOpWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func SetupLogger() {

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Fatalf("failed to open log file: %v", err)
		}
		logger = log.New(f, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		logger = log.New(&NoOpWriter{}, "", 0)
	}
}
