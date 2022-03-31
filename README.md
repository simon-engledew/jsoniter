# jsoniter

Avoid loading large JSON documents into memory by only deserializing the parts you are interested in.

A basic streaming JSON wrapper for the Go encoding/json Decoder.

## Usage

Use `Iterate` to fire a callback for every object / array value in the JSON document:

```golang
d := json.NewDecoder(os.Stdin)

fn := func(path []json.Token) error {
    fmt.Println(path)
		return nil
}

jsoniter.Iterate(d, fn)
```

Path will contain a slice of string keys or int array indexes that describe the location of the value in the document.

To match values that you are interested in, `iterjson` provides a basic matcher:

```golang
d := json.NewDecoder(os.Stdin)

matcher := jsoniter.Matcher("some", jsoniter.Wildcard, "nested", "structure")

var hits int

fn := func(path []json.Token) error {
  if matcher(path) {
    hits += 1
  }
  return nil
}

jsoniter.Iterate(d, fn)
```

By closing over the decoder it is also possible to decode values during iteration:

```golang
d := json.NewDecoder(os.Stdin)

matcher := jsoniter.Matcher("some", 0, "nested", "structure")

var found any

fn := func(path []json.Token) error {
  if matcher(path) {
    return d.Decode(&found)
  }
  return nil
}

jsoniter.Iterate(d, fn)
  ```
