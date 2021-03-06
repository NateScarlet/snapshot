package snapshot

import (
	"testing"
	"time"
)

type EmptyStringer struct{}

func (EmptyStringer) String() string {
	return "test"
}

func TestMatch(t *testing.T) {

	type Object struct {
		A string
		B int
		C bool
	}
	t.Run("simple", func(t *testing.T) {
		Match(t, "text", OptionExt(".txt"))
	})
	t.Run("bytes", func(t *testing.T) {
		Match(t, []byte{0x01, 0x02, 0x03})
	})
	t.Run("custom bytes type", func(t *testing.T) {
		type CustomBytes []byte
		MatchJSON(t, CustomBytes{0x01, 0x02, 0x03})
	})
	t.Run("custom bytes array", func(t *testing.T) {
		type CustomByteArray [3]byte
		MatchJSON(t, CustomByteArray{0x01, 0x02, 0x03})
	})
	t.Run("object", func(t *testing.T) {
		type Object struct {
			A string
			B int
			C bool
		}
		MatchJSON(t, Object{})
	})
	t.Run("array of object", func(t *testing.T) {

		MatchJSON(t, []Object{{A: "1"}, {B: 2}, {C: true}})
	})
	t.Run("empty object", func(t *testing.T) {
		type EmptyObject struct{}
		MatchJSON(t, EmptyObject{})
	})
	t.Run("empty stringer", func(t *testing.T) {
		MatchJSON(t, EmptyStringer{})
	})
	t.Run("time", func(t *testing.T) {
		MatchJSON(t, time.Date(2020, 11, 26, 0, 0, 0, 0, time.UTC))
	})
	t.Run("map", func(t *testing.T) {
		MatchJSON(t, map[string]interface{}{
			"string": "1",
			"int":    2,
			"bool":   false,
			"time":   time.Date(2020, 11, 26, 0, 0, 0, 0, time.UTC),
		})
	})
	t.Run("int map", func(t *testing.T) {
		MatchJSON(t, map[int]interface{}{1: 2})
	})
	t.Run("object map", func(t *testing.T) {
		MatchJSON(t, map[Object]interface{}{{A: "1"}: 2})
	})
}
