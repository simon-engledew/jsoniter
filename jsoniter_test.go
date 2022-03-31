package jsoniter_test

import (
	"encoding/json"
	"github.com/simon-engledew/jsoniter"
	"github.com/stretchr/testify/require"
	"io"
	"strings"
	"testing"
)

func TestInvalid(t *testing.T) {
	doc := `{
	  "some": [}`

	d := json.NewDecoder(strings.NewReader(doc))

	fn := func(path []json.Token) error {
		return nil
	}

	require.ErrorContains(t, jsoniter.Iterate(d, fn), `invalid character '}' looking for beginning of value`)
}

func TestEOF(t *testing.T) {
	doc := `{
	  "some": [{`

	d := json.NewDecoder(strings.NewReader(doc))

	fn := func(path []json.Token) error {
		return nil
	}

	require.ErrorIs(t, jsoniter.Iterate(d, fn), io.EOF)
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

	matcher := jsoniter.Matcher("some", jsoniter.Array, "nested", "structure")

	var hits int

	fn := func(path []json.Token) error {
		if matcher(path) {
			hits += 1
		}
		return nil
	}

	require.NoError(t, jsoniter.Iterate(d, fn))
	require.Equal(t, 2, hits)
}
