// Copyright (c) 2019, Ivan Eremin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modifier

import (
	"errors"
	"reflect"
)

var (
	// UUID as example of representation google uuid
	// uuid.UUID as a string value for notation
	UUID = func(v reflect.Value) (interface{}, bool, error) {
		var (
			value string
			empty bool
			err   error
		)
		// get uuid int val
		if m, ok := v.Type().MethodByName("ID"); ok && m.Func.Type().NumIn() == 1 {
			empty = m.Func.Call([]reflect.Value{v})[0].Uint() < 1
		} else {
			return value, true, errors.New("modifier:method ID not implemented")
		}
		// get uuid string value
		if m, ok := v.Type().MethodByName("String"); ok && m.Func.Type().NumIn() == 1 {
			value = m.Func.Call([]reflect.Value{v})[0].String()
		} else {
			err, empty = errors.New("modifier:method String not implemented"), true
		}
		return value, empty, err
	}

	// Time as example of representation build in time.Time type as string
	Time = func(v reflect.Value) (interface{}, bool, error) {
		var (
			value  string
			err    error
			empty  bool
			layout = "2006-01-02T15:04:05.999999999Z07:00" // time.RFC3339Nano
		)
		if m, ok := v.Type().MethodByName("IsZero"); ok && m.Func.Type().NumIn() == 1 {
			empty = m.Func.Call([]reflect.Value{v})[0].Bool()
		} else {
			return value, true, errors.New("modifier:method IsZero not implemented")
		}
		if m, ok := v.Type().MethodByName("Format"); ok && m.Func.Type().NumIn() == 2 {
			value = m.Func.Call([]reflect.Value{v, reflect.ValueOf(layout)})[0].String()
		} else {
			return value, true, errors.New("modifier:method Format not implemented")
		}
		return value, empty, err
	}
)
