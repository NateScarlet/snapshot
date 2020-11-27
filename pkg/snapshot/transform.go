package snapshot

import (
	"encoding/hex"
	"fmt"
	"reflect"
)

// Transform data before marshal
// transform may reuse input object.
type Transform = func(value interface{}) interface{}

// TransformJSON convert go object to json friendly data format.
// like mongo-compass used to edit document,
// types that json not support will be wrapped as { "$"+type: value  }.
func TransformJSON(v interface{}) interface{} {
	if v == nil {
		return v
	}

	switch o := v.(type) {
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
		return hex.EncodeToString(o)
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
		if rv.IsNil() {
			return nil
		}
		fallthrough
	case reflect.Array:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			var bytes = make([]byte, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				bytes[i] = uint8(rv.Index(i).Uint())
			}
			return map[string]interface{}{
				"$" + rv.Type().Name(): hex.EncodeToString(bytes),
			}
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
				"$" + rv.Type().Name(): fmt.Sprint(v),
			}
		}
		return map[string]interface{}{
			"$" + rv.Type().Name(): ret,
		}
	case reflect.Map:
		if rv.IsNil() {
			return nil
		}
		var ret = map[string]interface{}{}
		var it = rv.MapRange()
		for it.Next() {
			ret[fmt.Sprint(it.Key().Interface())] = DefaultTransform(it.Value().Interface())
		}
		return ret
	}

	return map[string]interface{}{
		"$" + rv.Type().Name(): v,
	}
}
