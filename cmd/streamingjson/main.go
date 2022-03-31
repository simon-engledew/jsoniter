package main

import (
	"encoding/json"
	"fmt"
	"github.com/simon-engledew/streamingjson"
	"os"
)

func callback(d *json.Decoder, path string) error {
	fmt.Println(path)
	return nil
}

func main() {
	err := streamingjson.Decode(json.NewDecoder(os.Stdin), callback)
	if err != nil {
		panic(err)
	}
}
