package goro

import (
	"os"
)

var globalBackend Backend

// Init initializes a given Backend interface. Built-in options are EbitenBackend and TCellBackend.
func Init(backend Backend) error {
	globalBackend = backend
	return backend.Init()
}

// Quit tells the current Backend to Quit before calling os.Exit.
func Quit() {
	globalBackend.Quit()
	os.Exit(0)
}

// Run runs the Backend with the provided logic callback.
func Run(cb func(*Screen)) error {
	return globalBackend.Run(cb)
}
