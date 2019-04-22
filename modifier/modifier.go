package modifier

import (
	"reflect"
)

var (
	// UUID as example of representation google uuid
	// uuid.UUID as a string value for notation
	UUID = func(v reflect.Value) (interface{}, bool, error) {
		var (
			value string
			err   error
		)
		// get uuid int val
		if m, ok := v.Type().MethodByName("ID"); ok && m.Func.Type().NumIn() == 1 {
			if m.Func.Call([]reflect.Value{v})[0].Uint() < 1 {
				return value, true, err
			}
		}
		// get uuid string value
		if m, ok := v.Type().MethodByName("String"); ok && m.Func.Type().NumIn() == 1 {
			value = m.Func.Call([]reflect.Value{v})[0].String()
			if len(value) > 0 {
				return value, false, err
			}
		}
		return value, true, err
	}

	// Time as example of representation build in time.Time type as string
	Time = func(v reflect.Value) (interface{}, bool, error) {
		var (
			value  string
			err    error
			layout = "2006-01-02T15:04:05.999999999Z07:00" // time.RFC3339Nano
		)
		if m, ok := v.Type().MethodByName("Format"); ok && m.Func.Type().NumIn() == 2 {
			value = m.Func.Call([]reflect.Value{v, reflect.ValueOf(layout)})[0].String()
			if len(value) > 0 {
				return value, false, err
			}
		}
		return value, true, err
	}
)
