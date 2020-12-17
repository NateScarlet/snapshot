package snapshot

import "testing"

func TestMarshalTextOrJSON(t *testing.T) {
	type M = map[string]interface{}
	t.Run("string", func(t *testing.T) {
		Match(t, "test")
	})
	t.Run("bytes", func(t *testing.T) {
		Match(t, []byte("test"))
	})
	t.Run("object", func(t *testing.T) {
		Match(t, M{"test": M{}})
	})
	t.Run("non ascii", func(t *testing.T) {
		Match(t, M{"测试": M{}})
	})
	t.Run("should not escape html", func(t *testing.T) {
		Match(t, M{"<tag>": M{}})
	})
}
