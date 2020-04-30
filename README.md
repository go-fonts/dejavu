# dejavu

[![GitHub release](https://img.shields.io/github/release/go-fonts/dejavu.svg)](https://github.com/go-fonts/dejavu/releases)
[![GoDoc](https://godoc.org/github.com/go-fonts/dejavu?status.svg)](https://godoc.org/github.com/go-fonts/dejavu)
[![License](https://img.shields.io/badge/License-BSD--3-blue.svg)](https://github.com/go-fonts/dejavu/raw/master/LICENSE)

`dejavu` provides the [DejaVu](https://dejavu-fonts.github.io/) fonts as importable Go packages.

The fonts are released under the [DejaVu font](https://github.com/go-fonts/dejavu/raw/master/LICENSE-DejaVu) license.
The Go packages under the [BSD-3](https://github.com/go-fonts/dejavu/raw/master/LICENSE) license.

## Example

```go
import (
	"fmt"
	"log"

	"github.com/go-fonts/dejavu/dejavumathtexgyre"
	"golang.org/x/image/font/sfnt"
)

func Example() {
	ttf, err := sfnt.Parse(dejavumathtexgyre.TTF)
	if err != nil {
		log.Fatalf("could not parse DejaVu Math font: %+v", err)
	}
	fmt.Printf("num glyphs: %d\n", ttf.NumGlyphs())

	// Output:
	// num glyphs: 4282
}
```
