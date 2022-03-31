package main

import (
	"encoding/json"
	"fmt"
	"github.com/simon-engledew/streamingjson"
	"os"
)

func main() {
	d := json.NewDecoder(os.Stdin)

	callback := func(path []json.Token) error {
		if streamingjson.PathString(path) == `.some[].nested.structure` {
			var v any
			err := d.Decode(&v)
			if err != nil {
				return err
			}
			fmt.Println("found!", v)
		}
		return nil
	}

	err := streamingjson.Iterate(d, callback)
	if err != nil {
		panic(err)
	}
}
