package snapshot

import (
	"regexp"
	"strings"
)

// Clean data before compare.
type Clean = func(v []byte) []byte

// MaskString replace all runes that not keep to given rune.
func MaskString(v string, as rune, keep func(c rune) bool) string {
	var b = strings.Builder{}
	b.Grow(len(v))
	for _, c := range v {
		if keep != nil && keep(c) {
			b.WriteRune(c)
		} else {
			b.WriteRune(as)
		}
	}

	return b.String()
}

// IsNonWord return true when rune in [0-9a-zA-Z].
func IsNonWord(c rune) bool {
	if '0' <= c && '9' >= c {
		return false
	}
	if 'a' <= c && 'z' >= c {
		return false
	}
	if 'A' <= c && 'Z' >= c {
		return false
	}
	return true
}

// MaskAsAsterisk treat data as string and replace all characters to `*`.
func MaskAsAsterisk(v []byte) []byte {
	return []byte(MaskString(string(v), '*', func(c rune) bool { return false }))
}

// MaskNonWordAsAsterisk treat data as string and replace all non word characters to `*`.
func MaskNonWordAsAsterisk(v []byte) []byte {
	return []byte(MaskString(string(v), '*', IsNonWord))
}

// CleanString by a function that takes string as argument
func CleanString(fn func(v string) string) Clean {
	return func(v []byte) (ret []byte) {
		return []byte(fn(string(v)))
	}
}

// CleanAs convert any data to given string.
func CleanAs(as string) Clean {
	return func([]byte) []byte {
		return []byte(as)
	}
}

// CleanRegex clean all data matched by given regexp.
// If regexp has group, clean function called on all submatch,
// otherwise, it called on match.
func CleanRegex(pattern regexp.Regexp, clean Clean) Clean {
	return func(v []byte) (ret []byte) {
		ret = v
		match := pattern.FindAllSubmatchIndex(v, -1)
		var offset = 0
		for _, submatch := range match {
			var i = 1
			var j = 2
			if len(submatch) > 2 {
				i = 3
				j = len(submatch)
			}
			for ; i < j; i += 2 {
				var (
					p = submatch[i-1] + offset
					q = submatch[i] + offset
				)
				var r = clean(ret[p:q])
				offset += len(r) - (q - p)
				var b = make([]byte, 0, len(v)+offset)
				ret = append(append(append(b, ret[:p]...), r...), ret[q:]...)
			}
		}
		return
	}
}
