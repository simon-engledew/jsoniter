package streamingjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type Callback func(d *json.Decoder, path string) error

const (
	tokenArrayStart  = json.Delim('[')
	tokenArrayEnd    = json.Delim(']')
	tokenObjectStart = json.Delim('{')
	tokenObjectEnd   = json.Delim('}')
)

func value(d *json.Decoder, path string, fn Callback) error {
	offset := d.InputOffset()
	if err := fn(d, path); err != nil {
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
		return array(d, path, fn)
	case tokenObjectStart:
		return object(d, path, fn)
	case tokenObjectEnd:
		return errors.New("unexpected delim }")
	case tokenArrayEnd:
		return errors.New("unexpected delim ]")
	}
	return nil
}

func object(d *json.Decoder, path string, fn Callback) error {
	for d.More() {
		key, err := d.Token()
		if err != nil {
			return err
		}
		if err := value(d, fmt.Sprintf("%s.%q", path, key), fn); err != nil {
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

func array(d *json.Decoder, path string, fn Callback) error {
	for d.More() {
		err := value(d, fmt.Sprintf("%s[]", path), fn)
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

// Decode will call fn for every path in the JSON document.
func Decode(d *json.Decoder, fn Callback) error {
	for {
		err := value(d, "", fn)

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}
