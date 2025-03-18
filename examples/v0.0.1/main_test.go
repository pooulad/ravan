package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMainFunction(t *testing.T) {
	// Save the original stdout
	originalStdout := os.Stdout

	// Create a pipe to capture stdout
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	// Redirect stdout to the write end of the pipe
	os.Stdout = writePipe

	// Channel to capture output
	outputChan := make(chan string)
	// Channel to capture errors
	errorChan := make(chan error)

	// Goroutine to read from the read end of the pipe
	go func() {
		var buf bytes.Buffer
		_, err := buf.ReadFrom(readPipe)
		if err != nil {
			errorChan <- err
			return
		}
		outputChan <- buf.String()
	}()

	// Run the main function
	main()

	// Close the write end of the pipe after main completes
	writePipe.Close()

	// Wait for the output or an error
	select {
	case output := <-outputChan:
		// Restore original stdout
		os.Stdout = originalStdout

		// Close the read end of the pipe
		readPipe.Close()

		// Check if the output contains the expected completion message
		if !strings.Contains(output, "progress bar completed") {
			t.Errorf("Expected completion message not found in output")
		}
	case err := <-errorChan:
		t.Fatalf("Error reading from pipe: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatalf("Test timed out")
	}
}
