package snapshot

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// Options for snapshot match
type Options struct {
	skip        int
	key         string
	ext         string
	transform   Transform
	marshal     Marshal
	clean       Clean
	assertEqual AssertEqual
	update      bool
}

// Option mutate SnapshotOptions
type Option func(*Options)

// OptionSkip add caller skip
func OptionSkip(skip int) Option {
	return func(so *Options) {
		so.skip += skip
	}
}

// OptionKey used in filename
func OptionKey(key string) Option {
	return func(so *Options) {
		so.key = key
	}
}

// OptionExt used as file extention
func OptionExt(ext string) Option {
	return func(so *Options) {
		so.ext = ext
	}
}

// OptionAssertEqual do assert
func OptionAssertEqual(assertEqual AssertEqual) Option {
	return func(so *Options) {
		so.assertEqual = assertEqual
	}
}

// OptionTransform object before marshal.
func OptionTransform(transform func(interface{}) interface{}) Option {
	return func(so *Options) {
		so.transform = transform
	}
}

// OptionMarshal object to bytes.
func OptionMarshal(marshal func(interface{}) ([]byte, error)) Option {
	return func(so *Options) {
		so.marshal = marshal
	}
}

// OptionClean dynamic data to make result deterministic.
// multiple clean option take affect in order.
func OptionClean(clean ...Clean) Option {
	return func(so *Options) {
		var orig = so.clean
		so.clean = func(v []byte) (ret []byte) {
			ret = v
			if orig != nil {
				ret = orig(ret)
			}
			for _, i := range clean {
				ret = i(ret)
			}
			return
		}
	}
}

// OptionCleanRegex replace all patten match by clean function,
// panic if any pattern is invalid.
func OptionCleanRegex(clean Clean, patterns ...string) Option {
	var c = make([]Clean, 0, len(patterns))
	for _, i := range patterns {
		c = append(c, CleanRegex(*regexp.MustCompile(i), clean))
	}
	return OptionClean(c...)
}

// OptionCleanRegexMask mask word matched by patterns to '*',
// panic if any pattern is invalid.
func OptionCleanRegexMask(patterns ...string) Option {
	return OptionCleanRegex(func(v []byte) []byte {
		return []byte(MaskString(string(v), '*', IsNonWord))
	}, patterns...)
}

// OptionUpdate is whether ignore existed file.
func OptionUpdate(v bool) Option {
	return func(so *Options) {
		so.update = v
	}
}

// Match compare object with file store under __snapshots__ folder
func Match(t *testing.T, actual interface{}, opts ...Option) {
	var args = new(Options)
	args.update = DefaultUpdate
	for _, i := range opts {
		i(args)
	}
	if args.transform == nil {
		args.transform = DefaultTransform
	}
	if args.marshal == nil {
		args.marshal = DefaultMarshal
	}
	if args.assertEqual == nil {
		args.assertEqual = DefaultAssertEqual
	}
	if args.key == "" {
		args.key = t.Name()
	}

	_, file, _, _ := runtime.Caller(args.skip + 1)
	p := filepath.Join(filepath.Dir(file), "__snapshots__", args.key+args.ext)
	require.NoError(t, os.MkdirAll(filepath.Dir(p), 0755))
	actualSnapshot, err := args.marshal(args.transform(actual))
	require.NoError(t, err)
	if args.clean != nil {
		actualSnapshot = args.clean(actualSnapshot)
	}

	update := func() error {
		return ioutil.WriteFile(p, actualSnapshot, 0644)
	}
	if args.update {
		require.NoError(t, update())
		return
	}
	expectedSnapshot, err := ioutil.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		require.NoError(t, update())
		return
	}
	require.NoError(t, err)
	args.assertEqual(t, expectedSnapshot, actualSnapshot)
}

// MatchJSON compare snapshot in json format
func MatchJSON(t *testing.T, actual interface{}, opts ...Option) {
	Match(t, actual,
		append(
			[]Option{
				OptionSkip(1),
				OptionExt(".json"),
				OptionAssertEqual(AssertEqualJSON),
			},
			opts...,
		)...,
	)
}
