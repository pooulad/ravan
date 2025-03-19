package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/pooulad/ravan"
)

func main() {
	fmt.Println("Success example:")
	SuccessExample()
	fmt.Println("Error example:")
	ErrorExample()
}

// Simulated file processing function
func processFile(filename string) error {
	if filename == "file5.txt" {
		// Intentionally create an error for file3.txt
		return errors.New("failed to read file")
	}
	time.Sleep(500 * time.Millisecond) // Simulate processing time
	return nil
}

func SuccessExample() {
	files := []string{"file1.txt", "file2.txt", "file3.txt", "file4.txt", "file6.txt"}

	bar, _ := ravan.New(ravan.WithMessage(&ravan.Message{Failed: "error custom message", Success: "success custom message"}))
	for i, file := range files {
		err := processFile(file)
		if err != nil {
			bar.FailMsg(fmt.Errorf("error processing %s: %w", file, err))
			return
		}
		bar.Draw(float64(i+1) / float64(len(files)))
	}

	bar.SuccessMsg()
}

func ErrorExample() {
	files := []string{"file1.txt", "file2.txt", "file3.txt", "file4.txt", "file5.txt"}

	bar, _ := ravan.New(ravan.WithMessage(&ravan.Message{Failed: "error custom message", Success: "success custom message"}))
	for i, file := range files {
		err := processFile(file)
		if err != nil {
			bar.FailMsg(fmt.Errorf("error processing %s: %w", file, err))
			return
		}
		bar.Draw(float64(i+1) / float64(len(files)))
	}

	bar.SuccessMsg()
}
