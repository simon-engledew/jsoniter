package jsoniter_test

import (
	"encoding/json"
	"errors"
	"github.com/simon-engledew/jsoniter"
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	doc := `{
	  "some": [{
		"nested": {
		  "structure": {
			"a": 1
		  }
		}
	  }]
	}`

	d := json.NewDecoder(strings.NewReader(doc))

	matcher := jsoniter.Matcher("some", 0, "nested", "structure")

	var found map[string]int

	err := jsoniter.Iterate(d, func(path []json.Token) error {
		if matcher(path) {
			return d.Decode(&found)
		}
		return nil
	})

	require.NoError(t, err)
	require.Equal(t, map[string]int{"a": 1}, found)
}

func TestInvalid(t *testing.T) {
	doc := `{
	  "some": [}`

	d := json.NewDecoder(strings.NewReader(doc))

	err := jsoniter.Iterate(d, func(path []json.Token) error {
		return nil
	})

	require.ErrorContains(t, err, `invalid character '}' looking for beginning of value`)
}

func TestEOF(t *testing.T) {
	doc := `{
	  "some": [{`

	d := json.NewDecoder(strings.NewReader(doc))

	err := jsoniter.Iterate(d, func(path []json.Token) error {
		return nil
	})

	require.ErrorIs(t, err, io.EOF)
}

func count(d *json.Decoder, matcher func(path []json.Token) bool) (int, error) {
	var hits int

	err := jsoniter.Iterate(d, func(path []json.Token) error {
		if matcher(path) {
			hits += 1
		}
		return nil
	})

	return hits, err
}

func TestIterate(t *testing.T) {
	doc := `{
	  "some": [{
		"nested": {
		  "structure": {
			"a": 1
		  }
		}
	  }, {
		"nested": {
		  "structure": {
			"b": 2
		  }
		}
      }]
	}`

	d := json.NewDecoder(strings.NewReader(doc))

	matcher := jsoniter.Matcher("some", jsoniter.Wildcard, "nested", "structure")

	hits, err := count(d, matcher)

	require.NoError(t, err)
	require.Equal(t, 2, hits)
	require.Equal(t, d.InputOffset(), int64(len(doc)))
}

func TestStop(t *testing.T) {
	doc := `{
	  "some": [{
		"nested": {
		  "structure": {
			"a": 1
		  }
		}
	  }]
	}`

	d := json.NewDecoder(strings.NewReader(doc))

	stopErr := errors.New("stop")

	err := jsoniter.Iterate(d, func(path []json.Token) error {
		return stopErr
	})

	require.ErrorIs(t, err, stopErr)
	require.Less(t, d.InputOffset(), int64(len(doc)))
}

func TestMiss(t *testing.T) {
	doc := `{
	  "some": [{
		"nested": {
		  "structure": {
			"a": 1
		  }
		}
	  }]
	}`

	d := json.NewDecoder(strings.NewReader(doc))

	matcher := jsoniter.Matcher("some", jsoniter.Wildcard, "nested", "structure", "c")

	hits, err := count(d, matcher)

	require.NoError(t, err)
	require.Zero(t, hits)
}
