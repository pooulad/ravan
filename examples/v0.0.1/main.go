package main

import (
	"fmt"
	"time"

	"github.com/pooulad/ravan"
)

func main() {
	total := 100
	r := ravan.New(50) // Create a progress bar of width 50 characters.

	for i := 0; i <= total; i++ {
		progress := float64(i) / float64(total)
		r.Draw(progress)
		time.Sleep(50 * time.Millisecond) // Simulate work
	}
	fmt.Println("progress bar completed")
}