package snapshot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// AssertEqual fails test if snapshot not equal with actual value.
type AssertEqual func(t *testing.T, expected, actual []byte)

// AssertEqualBytes compare each byte
func AssertEqualBytes(t *testing.T, expected, actual []byte) {
	assert.Equal(t, expected, actual)
}

// AssertEqualJSON compare value as json, output more human friendly error for json.
func AssertEqualJSON(t *testing.T, expected, actual []byte) {
	assert.JSONEq(t, string(expected), string(actual))
}
