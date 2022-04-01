# jsoniter

Avoid loading large JSON documents into memory by only deserializing the parts you are interested in.

A basic streaming JSON wrapper for the Go `encoding/json` `Decoder`.

## Usage

Use `jsoniter.Iterate` to fire a callback for every object and array value in the JSON document:

```golang
package main

import (
	"encoding/json"
	"fmt"
	"github.com/simon-engledew/jsoniter"
	"os"
)

func main() {
	d := json.NewDecoder(os.Stdin)

	err := jsoniter.Iterate(d, func(path []json.Token) error {
		fmt.Println(path)
		return nil
	})
	if err != nil {
		panic(err)
    }
}
```

`path` will contain a slice of `string` keys or `int` array indexes that describe the location of the value in the document. Do not modify or retain it during the callback.

To match values that you are interested in, `jsoniter` provides a basic matcher:

```golang
package main

import (
	"encoding/json"
	"fmt"
	"github.com/simon-engledew/jsoniter"
	"os"
)

func main() {
	d := json.NewDecoder(os.Stdin)

	// jsoniter.Wildcard will match any token
	// .some[*].nested.structure
	matcher := jsoniter.Matcher("some", jsoniter.Wildcard, "nested", "structure")

	var hits int

	err := jsoniter.Iterate(d, func(path []json.Token) error {
		if matcher(path) {
			hits += 1
		}
		return nil
	})
	if err != nil {
		panic(err)
    }
	fmt.Printf("found %d items matching path\n", hits)
}
```

By closing over the decoder it is also possible to decode values during iteration:

```golang
package main

import (
	"encoding/json"
	"fmt"
	"github.com/simon-engledew/jsoniter"
	"os"
)

func main() {
	d := json.NewDecoder(os.Stdin)

	// .some[0].nested.structure
	matcher := jsoniter.Matcher("some", 0, "nested", "structure")

	var found any

	err := jsoniter.Iterate(d, func(path []json.Token) error {
		if matcher(path) {
			return d.Decode(&found)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("found: ", found)
}
```

If a value is consumed by the callback `Iterate` will continue on with the rest of the document.
