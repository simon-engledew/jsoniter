package jsoniter

import (
	"encoding/json"
	"errors"
)

// Wildcard will match any json.Token when passed as a pattern argument to Matcher.
const Wildcard = iota

const (
	tokenArrayStart  = json.Delim('[')
	tokenArrayEnd    = json.Delim(']')
	tokenObjectStart = json.Delim('{')
	tokenObjectEnd   = json.Delim('}')
)

// Matcher returns a predicate which will match the series of tokens described by pattern.
func Matcher(pattern ...json.Token) func(path []json.Token) bool {
	return func(path []json.Token) bool {
		if len(pattern) != len(path) {
			return false
		}
		for i, v := range pattern {
			if v != Wildcard && v != path[i] {
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
	// If the offset has moved on then do not consume the value as the callback must have.
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
	index := 0
	for d.More() {
		if err := value(d, append(path, index), fn); err != nil {
			return err
		}
		index += 1
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

// Iterate triggers a callback for every object or array value in the JSON document being decoded.
//
// Rather than unmarshal the entire document at once, Iterate considers the document on a token by token basis.
// At each value callback will be passed the current path of the document, e.g:
//
// {"a": {"b": [1]}}
//
// will call back three times with the paths:
//
// * ["a"]
// * ["a", "b"]
// * ["a", "b", 0]
//
// To avoid allocations the path slice will be reused between callbacks and must not be retained or modified.
//
// It is safe to call `Decode` on the decoder during the callback to unmarshal a value from the document.
// Afterwards Iterate will continue to process the enclosing object or array:
//
//    var found any
//    jsoniter.Iterate(d, func(path []json.Token) error {
//        if matches(path) {
//            // decode the interesting section of the document:
//            return d.Decode(&found)
//        }
//        return nil
//    })
//
//
// Iterate will consume and discard the corresponding value if it has not been read by the callback.
func Iterate(d *json.Decoder, fn func(path []json.Token) error) error {
	return value(d, make([]json.Token, 0, 32), fn)
}
