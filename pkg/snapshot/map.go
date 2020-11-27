package snapshot

import "reflect"

// MapFromStruct convert struct to map.
// returns nil if input value is not a struct.
func MapFromStruct(v interface{}) (ret map[string]interface{}) {
	var rv = reflect.ValueOf(v)
	var rt = rv.Type()
	if rt.Kind() != reflect.Struct {
		return
	}
	ret = make(map[string]interface{})
	for i := 0; i < rv.NumField(); i++ {
		var f = rt.Field(i)
		// Skip un-exported fields.
		if f.PkgPath != "" {
			continue
		}
		ret[f.Name] = rv.Field(i).Interface()
	}
	return
}

// MapOmit delete given keys from map.
func MapOmit(v map[string]interface{}, keys ...string) map[string]interface{} {
	for _, i := range keys {
		delete(v, i)
	}
	return v
}
