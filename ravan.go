package ravan

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
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
	LessThan    BarCharacter = "<"
	Colon       BarCharacter = ":"
	Exclamation BarCharacter = "!"
	DollarSign  BarCharacter = "$"
	AtSign      BarCharacter = "@"
	Percent     BarCharacter = "%"
	CircumFlex  BarCharacter = "^"
	And         BarCharacter = "&"
)

// winsize handles window size.
type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// Validation categories
const (
	CompleteType   = "complete"
	IncompleteType = "incomplete"
)

var (
	validCompleteChars = map[BarCharacter]bool{
		Hash: true, Asterisk: true, Equal: true, Dash: true,
		GreaterThan: true, Colon: true, Plus: true, Exclamation: true,
		DollarSign: true, LessThan: true,
		And: true, AtSign: true, Percent: true, CircumFlex: true,
	}

	validIncompleteChars = map[BarCharacter]bool{
		Empty: true, Hash: true, Asterisk: true, Equal: true,
		Dash: true, GreaterThan: true, Colon: true, Plus: true,
		Exclamation: true, DollarSign: true, LessThan: true,
		And: true, AtSign: true, Percent: true, CircumFlex: true,
	}

	// colors for messagees
	successColor = "\033[32m"
	errorColor   = "\033[31m"
	resetColor   = "\033[0m"
)

// Option pattern for configuration
// You can use:
//
//	WithWidth
//	WithCompleteChar
//	WithIncompleteChar
//	WithMessage
type Option func(*Ravan) error

// Message struct for showing with Ravan progress bar
type Message struct {
	Failed  string
	Success string
}

// Ravan struct
type Ravan struct {
	width          int
	completeChar   BarCharacter
	incompleteChar BarCharacter
	message        Message
}

// New creates a validated Ravan instance
func New(opts ...Option) (*Ravan, error) {
	r := &Ravan{
		width:          50,    // Default width
		completeChar:   Equal, // Default complete
		incompleteChar: Empty, // Default incomplete
		message: Message{
			Failed:  "Operation failed",
			Success: "Operation successful",
		}, // Default message
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
	termWidth := getTerminalWidth()
	if termWidth == 0 {
		termWidth = r.width // fallback if terminal width cannot be determined
	}

	// Overhead accounts for extra characters like "[", "]", " 100%"
	overhead := 7
	effectiveWidth := r.width
	if termWidth-overhead < effectiveWidth {
		effectiveWidth = termWidth - overhead
		if effectiveWidth < 1 {
			effectiveWidth = 1
		}
	}

	complete := int(progress * float64(effectiveWidth))
	bar := strings.Repeat(string(r.completeChar), complete) +
		strings.Repeat(string(r.incompleteChar), effectiveWidth-complete)

	if progress >= 1.0 {
		// Print in green when complete
		fmt.Printf("\r\033[32m[%s] %.0f%%\033[0m\n", bar, progress*100)
	} else {
		fmt.Printf("\r[%s] %.0f%%", bar, progress*100)
	}
}

// FailMsg shows error (if provided) and/or custom failure message
// Usage:
// r.FailMsg()                      // Shows only custom message if set
// r.FailMsg(err)                   // Shows error + custom message
// r.FailMsg(err, customMessage)    // Optional: Override default custom message
func (r *Ravan) FailMsg(err ...interface{}) {
	var e error
	customMsg := r.message.Failed // Default to initialized message

	// Parse arguments
	for _, arg := range err {
		switch v := arg.(type) {
		case error:
			e = v
		case string:
			customMsg = v // Allow message override
		}
	}

	msg := strings.Builder{}
	msg.WriteString("\n" + errorColor)

	if e != nil {
		msg.WriteString(fmt.Sprintf("Error: %v. ", e))
	}

	if customMsg != "" {
		msg.WriteString(customMsg)
	}

	msg.WriteString(resetColor + "\n")
	fmt.Print(msg.String())
}

// SuccessMsg prints a success message with successColor.
func (r *Ravan) SuccessMsg() {
	fmt.Printf("%sSuccess: %s%s\n", successColor, r.message.Success, resetColor)
}

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

// WithMessage sets the custom messages.
func WithMessage(msg *Message) Option {
	return func(r *Ravan) error {
		if msg.Failed != "" {
            r.message.Failed = msg.Failed
        }
        if msg.Success != "" {
            r.message.Success = msg.Success
        }
		return nil
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

// getTerminalWidth returns the number of columns in the terminal.
// If an error occurs, it returns 0.
func getTerminalWidth() int {
	ws := &winsize{}
	// Use syscall.Syscall to call the TIOCGWINSZ ioctl on stdout.
	retCode, _, err := syscall.Syscall(syscall.SYS_IOCTL, os.Stdout.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(ws)))
	if int(retCode) == -1 || err != 0 {
		return 0 // error: fallback will be used
	}
	return int(ws.Col)
}
