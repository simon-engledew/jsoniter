package jsoniter_test

import (
	"encoding/json"
	"github.com/simon-engledew/jsoniter"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestIterate(t *testing.T) {
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

	matcher := jsoniter.Matcher("some", jsoniter.Array, "nested", "structure")

	var pass bool

	fn := func(path []json.Token) error {
		if matcher(path) {
			pass = true
		}
		return nil
	}

	require.NoError(t, jsoniter.Iterate(d, fn))
	require.True(t, pass)
}
