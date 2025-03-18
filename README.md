# 🏄 Ravan

minimalist, dependency-free progress bar library for Go.

the package name(**Ravan**) means smooth or flowing in persian. Exactly the same as the progress bar.

## Installation 🚀

```bash
go get github.com/pooulad/ravan@latest
```

## Usage 📝

```go
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

```

See [examples](/examples/) directory that added examples of all versions.also you can read the documentation that is in next section.

## Documentation 📋

[![Go Reference](https://pkg.go.dev/badge/github.com/pooulad/ravan.svg)](https://pkg.go.dev/github.com/pooulad/ravan)

Check latest version please.

## License 📏

Distributed under the MIT License. See [LICENSE](LICENSE) for more information.

## 👤Author

Created with ❤️ by pooulad.
