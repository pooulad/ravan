package ravan

import (
	"fmt"
	"strings"
)

type BarCharacter string

const (
	Empty       BarCharacter = " "
	Hash        BarCharacter = "#"
	Asterisk    BarCharacter = "*"
	Equal       BarCharacter = "="
	Plus        BarCharacter = "+"
	Dash        BarCharacter = "-"
	GreaterThan BarCharacter = ">"
	Colon       BarCharacter = ":"
	Exclamation BarCharacter = "!"
	Dollar      BarCharacter = "$"
)

// Validation categories
const (
	CompleteType   = "complete"
	IncompleteType = "incomplete"
)

var (
	validCompleteChars = map[BarCharacter]bool{
		Hash: true, Asterisk: true, Equal: true, Dash: true,
		GreaterThan: true, Colon: true, Plus: true, Exclamation: true, Dollar: true,
	}

	validIncompleteChars = map[BarCharacter]bool{
		Empty: true, Hash: true, Asterisk: true, Equal: true,
		Dash: true, GreaterThan: true, Colon: true, Plus: true, Exclamation: true, Dollar: true,
	}
)

// Option pattern for configuration
type Option func(*Ravan) error

// Width option
func WithWidth(w int) Option {
	return func(r *Ravan) error {
		r.width = w
		return nil
	}
}

// Complete character option
func WithCompleteChar(c BarCharacter) Option {
	return func(r *Ravan) error {
		if !isValidCharacter(c, CompleteType) {
			return fmt.Errorf("invalid complete character: %s", c)
		}
		r.completeChar = c
		return nil
	}
}

// Incomplete character option
func WithIncompleteChar(c BarCharacter) Option {
	return func(r *Ravan) error {
		if !isValidCharacter(c, IncompleteType) {
			return fmt.Errorf("invalid incomplete character: %s", c)
		}
		r.incompleteChar = c
		return nil
	}
}

// Ravan struct
type Ravan struct {
	width          int
	completeChar   BarCharacter
	incompleteChar BarCharacter
}

// New creates a validated Ravan instance
func New(opts ...Option) (*Ravan, error) {
	r := &Ravan{
		width:          50,    // Default width
		completeChar:   Equal, // Default complete
		incompleteChar: Empty, // Default incomplete
	}

	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, err
		}
	}

	if r.completeChar == r.incompleteChar {
		return nil, fmt.Errorf("complete and incomplete characters must differ")
	}

	return r, nil
}

// Draw renders the progress bar on the terminal.
// progress should be a value between 0.0 and 1.0.
// When progress is 1.0 (100%), the bar is printed in green.
func (r *Ravan) Draw(progress float64) {
	complete := int(progress * float64(r.width))
	bar := strings.Repeat(string(r.completeChar), complete) +
		strings.Repeat(string(r.incompleteChar), r.width-complete)
	if progress >= 1.0 {
		// Green color: \033[32m ... \033[0m resets color.
		fmt.Printf("\r\033[32m[%s] %.0f%%\033[0m\n", bar, progress*100)
	} else {
		fmt.Printf("\r[%s] %.0f%%", bar, progress*100)
	}
}

// isValidCharacter checks validity using maps for O(1) lookups
func isValidCharacter(c BarCharacter, charType string) bool {
	switch charType {
	case CompleteType:
		return validCompleteChars[c]
	case IncompleteType:
		return validIncompleteChars[c]
	default:
		return false
	}
}
