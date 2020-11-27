package snapshot

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeDefaultTransform(t *testing.T) {
	t.Cleanup(ResetDefaults)
	var calledCount int
	DefaultTransform = func(value interface{}) interface{} {
		calledCount++
		return value != nil
	}

	t.Run("transform to true", func(t *testing.T) {
		Match(t, "1111")
	})
	assert.Equal(t, 1, calledCount)
	t.Run("transform to false", func(t *testing.T) {
		Match(t, nil)
	})
	assert.Equal(t, 2, calledCount)

}

func TestChangeDefaultMarshall(t *testing.T) {
	t.Cleanup(ResetDefaults)

	var calledCount int
	DefaultMarshal = func(value interface{}) ([]byte, error) {
		calledCount++
		return []byte(fmt.Sprintf("%#v", value)), nil
	}

	Match(t, nil)
	assert.Equal(t, 1, calledCount)
	Match(t, nil)
	assert.Equal(t, 2, calledCount)
}

func TestChangeDefaultAssertEqual(t *testing.T) {
	t.Cleanup(ResetDefaults)
	var calledCount int
	DefaultAssertEqual = func(t *testing.T, actual, expected []byte) {
		calledCount++
	}

	Match(t, nil, OptionUpdate(false))
	assert.Equal(t, 1, calledCount)
	Match(t, nil, OptionUpdate(false))
	assert.Equal(t, 2, calledCount)
}
