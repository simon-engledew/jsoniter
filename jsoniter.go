package jsoniter

import (
	"encoding/json"
	"errors"
	"io"
)

type PathExpr string

var Array PathExpr = "[]"

const (
	tokenArrayStart  = json.Delim('[')
	tokenArrayEnd    = json.Delim(']')
	tokenObjectStart = json.Delim('{')
	tokenObjectEnd   = json.Delim('}')
)

func Matcher(pattern ...json.Token) func(path []json.Token) bool {
	return func(path []json.Token) bool {
		if len(pattern) != len(path) {
			return false
		}
		for i := range pattern {
			if pattern[i] != path[i] {
				return false
			}
		}
		return true
	}
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
		return object(d, path, fn)
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
		err := value(d, make([]json.Token, 0, 32), fn)

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}
