package streamingjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type PathExpr string

var Object PathExpr = "."
var Array PathExpr = "[]"

const (
	tokenArrayStart  = json.Delim('[')
	tokenArrayEnd    = json.Delim(']')
	tokenObjectStart = json.Delim('{')
	tokenObjectEnd   = json.Delim('}')
)

func PathString(path []json.Token) string {
	var b strings.Builder
	for _, t := range path {
		s := fmt.Sprintf("%v", t)
		if strings.ContainsAny(s, ".[]") {
			s = fmt.Sprintf("%q", s)
		}
		b.WriteString(s)
	}
	return b.String()
}

func value(d *json.Decoder, path []json.Token, fn func(path []json.Token) error) error {
	offset := d.InputOffset()

	if err := fn(path); err != nil {
		return err
	}
	// If the offset has moved on, do not consume the next token as the callback must have.
	if d.InputOffset() != offset {
		return nil
	}

	t, err := d.Token()
	if err != nil {
		return err
	}
	switch t {
	case tokenArrayStart:
		return array(d, append(path, Array), fn)
	case tokenObjectStart:
		return object(d, append(path, Object), fn)
	case tokenObjectEnd:
		return errors.New("unexpected delim }")
	case tokenArrayEnd:
		return errors.New("unexpected delim ]")
	}
	return nil
}

func object(d *json.Decoder, path []json.Token, fn func(path []json.Token) error) error {
	for d.More() {
		key, err := d.Token()
		if err != nil {
			return err
		}
		if err := value(d, append(path, key), fn); err != nil {
			return err
		}
	}
	t, err := d.Token()
	if err != nil {
		return err
	}
	if t != tokenObjectEnd {
		return errors.New("expected }")
	}
	return nil
}

func array(d *json.Decoder, path []json.Token, fn func(path []json.Token) error) error {
	for d.More() {
		err := value(d, path, fn)
		if err != nil {
			return err
		}
	}
	t, err := d.Token()
	if err != nil {
		return err
	}
	if t != tokenArrayEnd {
		return errors.New("expected ]")
	}
	return nil
}

// Iterate will call fn for every path in the JSON document.
func Iterate(d *json.Decoder, fn func(path []json.Token) error) error {
	for {
		err := value(d, make([]json.Token, 0, 64), fn)

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}
