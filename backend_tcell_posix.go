// +build !windows,!disableTCell

package goro

import "fmt"

// SetTitle sets the backend window's title.
func (backend *BackendTCell) SetTitle(title string) {
	fmt.Printf("\033]0;" + title + "\007")
}
