package ravan

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

// Helper function to capture the output of a function
func captureOutput(f func()) string {
	// Create a pipe to read and write data
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	// Save the original os.Stdout
	originalStdout := os.Stdout
	// Set os.Stdout to the pipe's writer
	os.Stdout = w

	// Execute the function
	f()

	// Close the writer and restore os.Stdout
	w.Close()
	os.Stdout = originalStdout

	// Read the output from the pipe's reader
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// TestNew verifies the New func works correctly.
func TestNew(t *testing.T) {
	expected := 50
	r, err := New(WithWidth(expected))
	if err != nil {
		t.Error(err)
	}

	got := r.width
	if expected != got {
		t.Errorf("expected %d. got %d", expected, got)
	}

}

// TestDraw verifies the Draw func works correctly.
func TestDraw(t *testing.T) {
	// Create a Ravan instance with a configured width of 50.
	// In tests, getTerminalWidth will likely return 0 (non-terminal),
	// so fallback will be used: termWidth = r.width, effectiveWidth = 50 - overhead (7) = 43.
	r := &Ravan{
		width:          50,
		completeChar:   Equal, // "="
		incompleteChar: Empty, // " " (space)
	}

	// Test progress 0.5
	output := captureOutput(func() {
		r.Draw(0.5)
		// Wait a short time to allow printing to complete.
		time.Sleep(10 * time.Millisecond)
	})

	// Expected effectiveWidth is 43.
	expectedEffectiveWidth := 43
	completeCount := int(0.5 * float64(expectedEffectiveWidth)) // should be 21
	incompleteCount := expectedEffectiveWidth - completeCount   // 22

	expectedBar := strings.Repeat(string(r.completeChar), completeCount) +
		strings.Repeat(string(r.incompleteChar), incompleteCount)
	expectedOutput := fmt.Sprintf("\r[%s] %.0f%%", expectedBar, 50.0)
	if output != expectedOutput {
		t.Errorf("Draw(0.5) = %q; want %q", output, expectedOutput)
	}

	// Test progress 1.0 (complete)
	output = captureOutput(func() {
		r.Draw(1.0)
		// Wait a short time to allow printing to complete.
		time.Sleep(10 * time.Millisecond)
	})
	expectedBar = strings.Repeat(string(r.completeChar), expectedEffectiveWidth)
	// ANSI green color codes are added for complete progress.
	expectedOutput = fmt.Sprintf("\r\033[32m[%s] %.0f%%\033[0m\n", expectedBar, 100.0)
	if output != expectedOutput {
		t.Errorf("Draw(1.0) = %q; want %q", output, expectedOutput)
	}
}

func TestIsValidCharacter(t *testing.T) {
	tests := []struct {
		char     BarCharacter
		charType string
		expected bool
	}{
		// Valid cases
		{Hash, CompleteType, true},
		{Asterisk, IncompleteType, true},
		// Invalid cases
		{BarCharacter("q"), CompleteType, false},
		{BarCharacter("a"), IncompleteType, false},
		// Edge cases
		{Empty, CompleteType, false},
		{Empty, IncompleteType, true},
	}

	for _, test := range tests {
		t.Run(string(test.char)+test.charType, func(t *testing.T) {
			result := isValidCharacter(test.char, test.charType)
			if result != test.expected {
				t.Errorf("isValidCharacter(%v, %s) = %v; want %v", test.char, test.charType, result, test.expected)
			}
		})
	}
}

func TestWithWidth(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"positive width", 100, 100},
		{"zero width", 0, 0},
		{"negative width", -10, -10},
		{"large width", 500, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Ravan{}
			fn := WithWidth(tt.input)
			err := fn(r)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if r.width != tt.expected {
				t.Errorf("expected width %d, got %d", tt.expected, r.width)
			}
		})
	}
}

func TestWithCompleteChar(t *testing.T) {
	tests := []struct {
		name        string
		input       BarCharacter
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid complete character",
			input:   Hash,
			wantErr: false,
		},
		{
			name:        "invalid complete character",
			input:       Empty,
			wantErr:     true,
			errContains: "invalid complete character",
		},
		{
			name:    "another valid complete character",
			input:   Equal,
			wantErr: false,
		},
		{
			name:        "invalid type for complete",
			input:       "custom", // Test invalid string conversion
			wantErr:     true,
			errContains: "invalid complete character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Ravan{}
			fn := WithCompleteChar(tt.input)
			err := fn(r)

			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr: %v, got: %v", tt.wantErr, err)
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil || !contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q, got %v", tt.errContains, err)
				}
			}

			if !tt.wantErr && r.completeChar != tt.input {
				t.Errorf("expected %s, got %s", tt.input, r.completeChar)
			}
		})
	}
}

func TestWithIncompleteChar(t *testing.T) {
	tests := []struct {
		name        string
		input       BarCharacter
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid incomplete character",
			input:   Empty,
			wantErr: false,
		},
		{
			name:    "valid shared character",
			input:   Hash,
			wantErr: false,
		},
		{
			name:        "invalid incomplete character",
			input:       "invalid",
			wantErr:     true,
			errContains: "invalid incomplete character",
		},
		{
			name:    "another valid incomplete character",
			input:   Dash,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Ravan{}
			fn := WithIncompleteChar(tt.input)
			err := fn(r)

			if (err != nil) != tt.wantErr {
				t.Fatalf("wantErr: %v, got: %v", tt.wantErr, err)
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil || !contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q, got %v", tt.errContains, err)
				}
			}

			if !tt.wantErr && r.incompleteChar != tt.input {
				t.Errorf("expected %s, got %s", tt.input, r.incompleteChar)
			}
		})
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// TestGetTerminalWidthNonTerminal simulates a non-terminal stdout and expects a width of 0.
func TestGetTerminalWidthNonTerminal(t *testing.T) {
	// Create a temporary file which is not a terminal.
	tmp, err := os.CreateTemp("", "nonterminal")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmp.Name())

	// Save the original os.Stdout and override it with our temp file.
	originalStdout := os.Stdout
	os.Stdout = tmp

	// Call getTerminalWidth; since tmp is not a terminal, we expect 0.
	width := getTerminalWidth()
	if width != 0 {
		t.Errorf("Expected getTerminalWidth to return 0 for non-terminal, got %d", width)
	}

	// Restore the original os.Stdout.
	os.Stdout = originalStdout
}

// TestGetTerminalWidth_Terminal checks that getTerminalWidth returns a positive value when os.Stdout is a terminal.
// If os.Stdout is not a terminal, the test is skipped.
func TestGetTerminalWidth_Terminal(t *testing.T) {
	// Check if os.Stdout is a terminal by inspecting its file mode.
	fi, err := os.Stdout.Stat()
	if err != nil {
		t.Fatalf("failed to stat os.Stdout: %v", err)
	}
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		t.Skip("os.Stdout is not a terminal; skipping test")
	}

	// Now call getTerminalWidth. It should return a positive value.
	width := getTerminalWidth()
	if width <= 0 {
		t.Errorf("Expected positive terminal width, got %d", width)
	}
}

// TestWithMessageOptionDefaults verifies WithMessage works correctlly with default messages.
func TestWithMessageOptionDefaults(t *testing.T) {
	// Create a new Ravan instance with no message override.
	r, err := New()
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	// Check that default messages are set.
	expectedFailed := "Operation failed"
	expectedSuccess := "Operation successful"

	if r.message.Failed != expectedFailed {
		t.Errorf("Expected default Failed message %q, got %q", expectedFailed, r.message.Failed)
	}
	if r.message.Success != expectedSuccess {
		t.Errorf("Expected default Success message %q, got %q", expectedSuccess, r.message.Success)
	}
}

// TestWithMessageOptionCustom verifies WithMessage works correctlly with custom messages.
func TestWithMessageOptionCustom(t *testing.T) {
	// Create a custom message.
	customMsg := &Message{
		Failed:  "Custom failed message",
		Success: "Custom success message",
	}

	// Create a new Ravan instance with custom messages.
	r, err := New(WithMessage(customMsg))
	if err != nil {
		t.Fatalf("New(WithMessage) error: %v", err)
	}

	// Verify that the custom messages are set.
	if r.message.Failed != customMsg.Failed {
		t.Errorf("Expected Failed message %q, got %q", customMsg.Failed, r.message.Failed)
	}
	if r.message.Success != customMsg.Success {
		t.Errorf("Expected Success message %q, got %q", customMsg.Success, r.message.Success)
	}
}

// TestWithMessageOptionPartial verifies WithMessage works correctlly with custom  messages with only Failed message provided.
func TestWithMessageOptionPartial(t *testing.T) {
	// Create a custom message with only Failed message provided.
	customMsg := &Message{
		Failed:  "Only custom failed message",
		Success: "",
	}

	// Create a new Ravan instance with the partial custom message.
	r, err := New(WithMessage(customMsg))
	if err != nil {
		t.Fatalf("New(WithMessage) error: %v", err)
	}

	// Expected: custom Failed should override, while Success remains default.
	if r.message.Failed != customMsg.Failed {
		t.Errorf("Expected Failed message %q, got %q", customMsg.Failed, r.message.Failed)
	}

	expectedSuccess := "Operation successful"
	if r.message.Success != expectedSuccess {
		t.Errorf("Expected default Success message %q, got %q", expectedSuccess, r.message.Success)
	}
}
