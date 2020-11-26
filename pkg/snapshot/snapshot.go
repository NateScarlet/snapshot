package snapshot

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// DefaultTransform implementation
func DefaultTransform(v interface{}) interface{} {
	if v == nil {
		return v
	}

	switch v.(type) {
	case bool:
		return v
	case uint:
		return v
	case uint8:
		return v
	case uint16:
		return v
	case uint32:
		return v
	case uint64:
		return v
	case int:
		return v
	case int8:
		return v
	case int16:
		return v
	case int32:
		return v
	case int64:
		return v
	case float32:
		return v
	case float64:
		return v
	case string:
		return v
	case []byte:
		return v
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Interface:
		fallthrough
	case reflect.Ptr:
		if rv.IsNil() {
			return nil
		}
		return DefaultTransform(rv.Elem().Interface())
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		if rv.IsNil() {
			return nil
		}
		var ret = []interface{}{}
		for i := 0; i < rv.Len(); i++ {
			ret = append(ret, DefaultTransform(rv.Index(i).Interface()))
		}
		return ret
	case reflect.Struct:
		var ret = map[string]interface{}{}
		for i := 0; i < rv.NumField(); i++ {
			var f = rv.Type().Field(i)
			// Skip un-exported fields.
			if f.PkgPath != "" {
				continue
			}
			ret[f.Name] = DefaultTransform(rv.Field(i).Interface())
		}
		if len(ret) == 0 {
			return map[string]interface{}{
				"$" + rv.Type().Name(): v,
			}
		}
		return ret
	case reflect.Map:
		if rv.IsNil() {
			return nil
		}
		var ret = map[string]interface{}{}
		var it = rv.MapRange()
		for it.Next() {
			ret[it.Key().String()] = DefaultTransform(it.Value().Interface())
		}
		return ret
	}

	return map[string]interface{}{
		"$" + rv.Type().Name(): v,
	}
}

// DefaultMarshal implementation
func DefaultMarshal(v interface{}) ([]byte, error) {
	switch o := v.(type) {
	case string:
		return []byte(o), nil
	case []byte:
		return o, nil
	}
	return json.MarshalIndent(v, "", "  ")
}

// Options for snapshot match
type Options struct {
	skip      int
	key       string
	ext       string
	match     func(expected, actual []byte)
	transform func(interface{}) interface{}
	marshal   func(interface{}) ([]byte, error)
	update    bool
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

// OptionMatch do assert
func OptionMatch(match func(a, b []byte)) Option {
	return func(so *Options) {
		so.match = match
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

// OptionUpdate is whether ignore existed file.
func OptionUpdate(v bool) Option {
	return func(so *Options) {
		so.update = v
	}
}

// Match compare object with file store under __snapshots__ folder
func Match(t *testing.T, actual interface{}, opts ...Option) {
	var args = new(Options)
	args.update = os.Getenv("SNAPSHOT_UPDATE") == "true"
	for _, i := range opts {
		i(args)
	}
	if args.transform == nil {
		args.transform = DefaultTransform
	}
	if args.marshal == nil {
		args.marshal = DefaultMarshal
	}
	if args.match == nil {
		args.match = func(a, b []byte) {
			assert.Equal(t, string(a), string(b))
		}
	}
	if args.key == "" {
		args.key = t.Name()
	}

	_, file, _, _ := runtime.Caller(args.skip + 1)
	p := filepath.Join(filepath.Dir(file), "__snapshots__", args.key+args.ext)
	require.NoError(t, os.MkdirAll(filepath.Dir(p), 0755))
	actualSnapshot, err := args.marshal(args.transform(actual))
	require.NoError(t, err)
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
	args.match(expectedSnapshot, actualSnapshot)
}

// MatchJSON compare snapshot in json format
func MatchJSON(t *testing.T, actual interface{}, opts ...Option) {
	Match(t, actual,
		append(
			[]Option{
				OptionSkip(1),
				OptionExt(".json"),
				OptionMatch(func(a, b []byte) {
					assert.JSONEq(t, string(a), string(b))
				}),
			},
			opts...,
		)...,
	)
}
