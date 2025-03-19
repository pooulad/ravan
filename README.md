# 🏄 Ravan

minimalist, dependency-free progress bar library for Go.

the package name(**Ravan**) means smooth or flowing in persian. Exactly the same as the progress bar.

![GitHub Release](https://img.shields.io/github/v/release/pooulad/ravan)

## Installation 🚀

```bash
go get github.com/pooulad/ravan@latest
```

## Usage 📝

```go
package main

import (
    "errors"
    "fmt"
    "time"

    "github.com/pooulad/ravan"
)

func main() {
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

```

See [examples](/examples/) directory that added examples of all versions.also you can read the [documentation](https://pkg.go.dev/github.com/pooulad/ravan) that is in next section.

## Bar characters 🔤

| Name        | Sign |
| ----------- | ---- |
| Hash        | "#"  |
| Asterisk    | "\*" |
| Equal       | "="  |
| Plus        | "+"  |
| Dash        | "-"  |
| GreaterThan | ">"  |
| LessThan    | "<"  |
| Colon       | ":"  |
| Exclamation | "!"  |
| DollarSign  | "$"  |
| AtSign      | "@"  |
| Percent     | "%"  |
| CircumFlex  | "^"  |
| And         | "&"  |
| Empty       | " "  |

📌Empty only available in WithIncompleteChar func.it is obvious that completed bar section can't fill with Empty char.

## Custom message 💭

in v0.0.3 and above you can set custom Message struct for handling in error and successful situation with WithMessage func.

```go
bar, _ := ravan.New(ravan.WithMessage(&ravan.Message{Failed: "error custom message", Success: "success custom message"}))
```

See and run this [example](/examples/v0.0.3/main.go)

## Documentation 📋

[![Go Reference](https://pkg.go.dev/badge/github.com/pooulad/ravan.svg)](https://pkg.go.dev/github.com/pooulad/ravan)

Check latest version please.

## Demo 🖼

![ravan_sample_screenshot](/assets/demo.gif)

## License 📏

Distributed under the MIT License. See [LICENSE](LICENSE) for more information.

## 👤Author

Created with ❤️ by pooulad.
