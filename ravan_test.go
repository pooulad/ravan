package ravan

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
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
	// Initialize the Ravan struct with a specific width
	expected := 10
	r, err := New(WithWidth(expected))
	if err != nil {
		t.Error(err)
	}

	// Define test cases with different progress values
	testCases := []struct {
		progress       float64
		expectedOutput string
	}{
		{0.0, "\r[          ] 0%"},
		{0.5, "\r[=====     ] 50%"},
		{1.0, "\r\033[32m[==========] 100%\033[0m\n"},
	}

	for _, tc := range testCases {
		// Capture the output of the Draw method
		output := captureOutput(func() {
			r.Draw(tc.progress)
		})

		// Compare the captured output with the expected output
		if output != tc.expectedOutput {
			t.Errorf("For progress %.1f, expected output %q, but got %q", tc.progress, tc.expectedOutput, output)
		}
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
