package snapshot

import "testing"

func TestMatch(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		Match(t, "text", OptionExt(".txt"))
	})
}
