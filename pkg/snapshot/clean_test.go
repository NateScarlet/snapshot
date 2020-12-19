package snapshot

import (
	"testing"
	"time"
)

func TestCleanRegexp(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		Match(t, "text", OptionCleanRegex(MaskNonWordAsAsterisk, "ex"))
	})
	t.Run("group", func(t *testing.T) {
		Match(t, "text", OptionCleanRegex(MaskNonWordAsAsterisk, "t(e)x"))
	})
	t.Run("multiple group", func(t *testing.T) {
		Match(t, "text", OptionCleanRegex(MaskNonWordAsAsterisk, "t(e)(x)"))
	})
	t.Run("object", func(t *testing.T) {
		type Object struct {
			A string
			B int
			C bool
			D time.Time
		}
		MatchJSON(t, Object{
			A: "a",
			B: 1,
			C: true,
			D: time.Now(),
		}, OptionCleanRegex(CleanAs("*time*"), `"D": {\s+"\$Time": "(.+)"\s+}`))
	})
	t.Run("multiple match", func(t *testing.T) {
		type Object struct {
			A      int
			B      int
			ACount int
			BCount int
		}
		MatchJSON(t, Object{
			A:      1,
			B:      2,
			ACount: 1,
			BCount: 2,
		}, OptionCleanRegex(CleanAs(`"*count*"`), `(?m)^\s*".+Count": (\d+),?$`))
	})

}
