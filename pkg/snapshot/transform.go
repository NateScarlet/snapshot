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
		return TransformJSON(rv.Elem().Interface())
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
			ret = append(ret, TransformJSON(rv.Index(i).Interface()))
		}
		return ret
	case reflect.Struct:
		var m = MapFromStruct(v)
		for k, v := range m {
			m[k] = TransformJSON(v)
		}
		if len(m) == 0 {
			return map[string]interface{}{
				"$" + rv.Type().Name(): fmt.Sprint(v),
			}
		}
		return map[string]interface{}{
			"$" + rv.Type().Name(): m,
		}
	case reflect.Map:
		if rv.IsNil() {
			return nil
		}
		var ret = map[string]interface{}{}
		var it = rv.MapRange()
		for it.Next() {
			ret[fmt.Sprint(it.Key().Interface())] = TransformJSON(it.Value().Interface())
		}
		return ret
	}

	return map[string]interface{}{
		"$" + rv.Type().Name(): v,
	}
}

// TransformSchema that only keep type info.
// Useful when handling dynamic data like timestamp.
//
// Returns:
//
//   nil
//     when value is nil
//
//   $[]byte
//     when value is byte slice
//
//   $[length]byte
//     when value is byte array
//
//   []interface{}
//     when value is slice or array
//
//   {"$pointer": value }
//     when value is potinter
//
//   {"$"+type: map[string]interace{}}
//     when value is struct and has any exported fields.
//
//   map[string]interace{}
//     when value is map
//
//   "$"+type
//     default
func TransformSchema(v interface{}) interface{} {
	if v == nil {
		return v
	}

	type M = map[string]interface{}
	type A = []interface{}

	var rv = reflect.ValueOf(v)
	var rt = rv.Type()
	switch rt.Kind() {
	case reflect.Interface:
		return TransformSchema(rv.Elem().Interface())
	case reflect.Ptr:
		return M{"$pointer": TransformSchema(rv.Elem().Interface())}
	case reflect.Slice:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return "$[]byte"
		}
		fallthrough
	case reflect.Array:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return fmt.Sprintf("$[%d]byte", rv.Len())
		}
		var ret = []interface{}{}
		for i := 0; i < rv.Len(); i++ {
			ret = append(ret, TransformSchema(rv.Index(i).Interface()))
		}
		return ret
	case reflect.Struct:
		var m = MapFromStruct(v)
		for k, v := range m {
			m[k] = TransformSchema(v)
		}
		if len(m) == 0 {
			return "$" + rv.Type().Name()
		}
		return map[string]interface{}{"$" + rt.Name(): m}
	case reflect.Map:
		var ret = map[string]interface{}{}
		var it = rv.MapRange()
		for it.Next() {
			ret[fmt.Sprint(it.Key().Interface())] = TransformSchema(it.Value().Interface())
		}
		return ret
	default:
		return "$" + rt.Name()
	}
}
