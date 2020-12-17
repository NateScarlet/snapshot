package snapshot

import (
	"bytes"
	"encoding/json"
)

// Marshal go object to bytes
type Marshal = func(value interface{}) ([]byte, error)

// MarshalTextOrJSON marshal json if value not text.
func MarshalTextOrJSON(v interface{}) ([]byte, error) {
	switch o := v.(type) {
	case string:
		return []byte(o), nil
	case []byte:
		return o, nil
	}

	var b = new(bytes.Buffer)
	var enc = json.NewEncoder(b)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	var err = enc.Encode(v)
	return b.Bytes(), err
}
