package ravan

import (
	"fmt"
	"strings"
)

// Ravan holds configuration for the progress bar.
type Ravan struct {
	Width int
}

// New creates a new Ravan struct with the given width.
func New(width int) *Ravan {
	return &Ravan{Width: width}
}

// Draw renders the progress bar on the terminal.
// progress should be a value between 0.0 and 1.0.
// When progress is 1.0 (100%), the bar is printed in green.
func (r *Ravan) Draw(progress float64) {
	complete := int(progress * float64(r.Width))
	bar := strings.Repeat("=", complete) + strings.Repeat(" ", r.Width-complete)
	if progress >= 1.0 {
		// Green color: \033[32m ... \033[0m resets color.
		fmt.Printf("\r\033[32m[%s] %.0f%%\033[0m", bar, progress*100)
		fmt.Println()
	} else {
		fmt.Printf("\r[%s] %.0f%%", bar, progress*100)
	}
}
