package ravan

import (
	"bytes"
	"io"
	"os"
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
	r := New(expected)

	got := r.Width
	if expected != got {
		t.Errorf("expected %d. got %d", expected, got)
	}

}

// TestDraw verifies the Draw func works correctly.
func TestDraw(t *testing.T) {
	// Initialize the Ravan struct with a specific width
	r := New(10)

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
