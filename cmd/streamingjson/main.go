package main

import (
	"encoding/json"
	"fmt"
	"github.com/simon-engledew/streamingjson"
	"os"
)

func callback(d *json.Decoder, path string) error {
	if path == `."some"[]."nested"."structure"` {
		var v any
		err := d.Decode(&v)
		if err != nil {
			return err
		}
		fmt.Println("found!", path, v)
	}
	return nil
}

func main() {
	err := streamingjson.Decode(json.NewDecoder(os.Stdin), callback)
	if err != nil {
		panic(err)
	}
}
