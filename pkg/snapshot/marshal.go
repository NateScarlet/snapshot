package snapshot

import "encoding/json"

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
	return json.MarshalIndent(v, "", "  ")
}
