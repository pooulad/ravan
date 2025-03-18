package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pooulad/ravan"
)

func main() {
	total := 100
	// if you dont add any option the default values are set. -> ravan.New()
	r, err := ravan.New(ravan.WithWidth(50), ravan.WithCompleteChar("-"))
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i <= total; i++ {
		progress := float64(i) / float64(total)
		r.Draw(progress)
		time.Sleep(50 * time.Millisecond) // Simulate work
	}
	fmt.Println("progress bar completed")
}
