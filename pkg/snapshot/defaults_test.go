package snapshot

import (
	"fmt"
	"testing"
)

func TestChangeDefaultTransform(t *testing.T) {
	t.Cleanup(ResetDefaults)
	DefaultTransform = func(value interface{}) interface{} {
		return value != nil
	}

	t.Run("transform to true", func(t *testing.T) {
		Match(t, "1111")
	})
	t.Run("transform to false", func(t *testing.T) {
		Match(t, nil)
	})

}

func TestChangeDefaultMarshall(t *testing.T) {
	t.Cleanup(ResetDefaults)
	DefaultMarshal = func(value interface{}) ([]byte, error) {
		return []byte(fmt.Sprintf("%#v", value)), nil
	}

	Match(t, nil)
}
