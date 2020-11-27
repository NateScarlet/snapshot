package snapshot

import (
	"testing"
	"time"
)

func TestTransformSchema(t *testing.T) {
	t.Cleanup(ResetDefaults)
	DefaultTransform = TransformSchema

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
		Match(t, CustomBytes{0x01, 0x02, 0x03})
	})
	t.Run("custom bytes array", func(t *testing.T) {
		type CustomByteArray [3]byte
		Match(t, CustomByteArray{0x01, 0x02, 0x03})
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
		Match(t, EmptyObject{})
	})
	t.Run("empty stringer", func(t *testing.T) {
		Match(t, EmptyStringer{})
	})
	t.Run("time", func(t *testing.T) {
		Match(t, time.Date(2020, 11, 26, 0, 0, 0, 0, time.UTC))
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

func TestDynamicData(t *testing.T) {
	type TODO struct {
		ID      int
		Created time.Time
		Name    string
	}
	MatchJSON(t,
		TODO{
			ID:      1,
			Created: time.Now(),
			Name:    "job1",
		},
		OptionTransform(func(i interface{}) interface{} {
			var m = MapFromStruct(i)
			if v, ok := m["Created"]; ok {
				m["Created"] = TransformSchema(v)
			}
			return m
		}),
	)
}
