# 🏄 Ravan

minimalist, dependency-free progress bar library for Go.

the package name(**Ravan**) means smooth or flowing in persian. Exactly the same as the progress bar.

## Installation 🚀

```bash
go get github.com/pooulad/ravan
```

## Usage 📝

```go
package main

import (
    "fmt"
    "time"

    "github.com/pooulad/ravan"
)

func main() {
    total := 100
    r := ravan.New(50)

    for i := 0; i <= total; i++ {
        progress := float64(i) / float64(total)
        r.Draw(progress)
        time.Sleep(50 * time.Millisecond) // Simulate work
    }
    fmt.Println("progress bar completed")
}
```

See [examples](/examples/) directory that added with test for you.

## Documentation 📋

[![Go Reference](https://pkg.go.dev/badge/github.com/pooulad/ravan.svg)](https://pkg.go.dev/github.com/pooulad/ravan)

## License 📏

Distributed under the MIT License. See [LICENSE](LICENSE) for more information.

## 👤Author

Created with ❤️ by pooulad.
