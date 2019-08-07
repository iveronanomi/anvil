// Copyright (c) 2019, Ivan Eremin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package anvil

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type (
	// mode of (non-)skipping empty values
	mode int
	// Anvil executor structure
	Anvil struct {
		//Mode behavior for skipping empty values
		Mode mode
		//Glue string to glue fields
		Glue string
		// modifier it's a list of functions used as a rule
		// to find out empty or not empty value of a field with given type and
		// type representation
		// exported type key as a key and list of functions to execute.
		modifier map[string]func(f reflect.Value) (interface{}, bool, error)
		// collection of []{key => value}
		items []Item
		deep  int // reserved for a future features
	}
	// Item field with typed value as a result of notation
	Item struct {
		Key   string
		Value interface{}
	}
)

const (
	// NoSkipEmpty fields with empty values
	NoSkipEmpty mode = iota
	// SkipEmpty fields with empty values
	SkipEmpty
)

// RegisterModifierFunc - assign a modifier function
// to extract value of given type as
// result of callback function used (value, isEmpty, error) where
// value is an interface{} value, isEmpty - valuable for
// behaviour Mode, and error if error occurred, used to stop execution
func (s *Anvil) RegisterModifierFunc(t interface{}, mod func(f reflect.Value) (interface{}, bool, error)) *Anvil {
	if s.modifier == nil {
		s.modifier = make(map[string]func(f reflect.Value) (interface{}, bool, error))
	}
	s.modifier[reflect.TypeOf(t).String()] = mod
	return s
}

// Notation of go type as a list of []Item
// where key is a string and value is a typed interface value
func Notation(source interface{}, behaviour mode, glue string) ([]Item, error) {
	if source == nil {
		return nil, nil
	}
	s := &Anvil{
		Glue:     glue,
		Mode:     behaviour,
		modifier: make(map[string]func(f reflect.Value) (interface{}, bool, error)),
	}
	return s.notation("", reflect.ValueOf(source), false)
}

// Notation of go type as a list of []Item
// where key is a string and value is a typed interface value
func (s *Anvil) Notation(sample interface{}) ([]Item, error) {
	if sample == nil {
		return nil, nil
	}
	return s.notation("", reflect.ValueOf(sample), false)
}

// notation structure nested
func (s *Anvil) notation(key string, v reflect.Value, title bool) (items []Item, err error) {
	var (
		value interface{}
		empty = true
		skip  = s.Mode == SkipEmpty
	)
	// get value by pointer if it is
	v = reflect.Indirect(v)

	// set default prefix for a field
	if len(key) < 1 {
		key = v.Type().Name()
	}
	switch v.Kind() {
	case reflect.Invalid:
		return nil, errors.New("anvil:invalid value of " + v.Type().Name())
	case reflect.Array:
		if v.Len() < 1 {
			break
		}
		if value, empty, err = s.modify(v); err != nil {
			break
		}
		if !empty {
			break
		}
		for i := 0; i < v.Len(); i++ {
			n, err := s.notation(arrayPrefix(key, i), v.Index(i), true)
			if err != nil {
				return nil, err
			}
			if len(n) < 1 {
				continue
			}
			items = append(items, n...)
		}
	case reflect.Slice:
		if v.IsNil() {
			break
		}
		if value, empty, err = s.modify(v); err != nil {
			break
		}
		if v.Len() < 1 {
			break
		}
		for i := 0; i < v.Len(); i++ {
			if v.Index(i).CanAddr() {
				n, err := s.notation(arrayPrefix(key, i), reflect.Indirect(v.Index(i).Addr()), true)
				if err != nil {
					return nil, err
				}
				if len(n) < 1 {
					continue
				}
				items = append(items, n...)
			}
		}
	case reflect.Struct:
		if value, empty, err = s.modify(v); err != nil {
			return nil, err
		}
		if !empty {
			break
		}
		l := v.NumField()
		for i := 0; i < l; i++ {
			f := reflect.Indirect(v.Field(i))
			// skip invalid field
			if f.Kind() == reflect.Invalid {
				continue
			}
			n, err := s.notation(s.key(key, v.Type().Field(i), false), f, true)
			if err != nil {
				return nil, err
			}
			if len(n) < 1 {
				continue
			}
			items = append(items, n...)
		}
	case reflect.Interface:
		if !v.Elem().IsValid() {
			break
		}
		n, err := s.notation(key, v.Elem(), true)
		if err != nil {
			return nil, err
		}
		if len(n) < 1 {
			break
		}
		items = n
	case reflect.Int:
		value, empty = int(v.Int()), v.Int() == 0
	case reflect.Int8:
		value, empty = int8(v.Int()), v.Int() == 0
	case reflect.Int16:
		value, empty = int16(v.Int()), v.Int() == 0
	case reflect.Int32:
		value, empty = int32(v.Int()), v.Int() == 0
	case reflect.Int64:
		value, empty = v.Int(), v.Int() == 0
	case reflect.Float32:
		value, empty = float32(v.Float()), v.Float() == .0
	case reflect.Float64:
		value, empty = v.Float(), v.Float() == 0
	case reflect.Uint:
		value, empty = uint(v.Uint()), v.Uint() == 0
	case reflect.Uint8:
		value, empty = uint8(v.Uint()), v.Uint() == 0
	case reflect.Uint16:
		value, empty = uint16(v.Uint()), v.Uint() == 0
	case reflect.Uint32:
		value, empty = uint32(v.Uint()), v.Uint() == 0
	case reflect.Uint64:
		value, empty = v.Uint(), v.Uint() == 0
	case reflect.Bool:
		value, empty = v.Bool(), !v.Bool()
	case reflect.String:
		value, empty = v.String(), len(v.String()) < 1
	case reflect.Map:
		if v.IsNil() || v.Len() < 1 {
			break
		}
		keys := v.MapKeys()
		for i := range keys {
			n, err := s.notation(mapPrefix(key, keys[i]), v.MapIndex(keys[i]), true)
			if err != nil {
				return nil, err
			}
			if len(n) < 1 {
				continue
			}
			items = append(items, n...)
		}
	case reflect.Complex64:
		value = complex64(v.Complex())
		empty = complex64(reflect.Zero(v.Type()).Complex()) == value
	case reflect.Complex128:
		value = v.Complex()
		empty = reflect.Zero(v.Type()).Complex() == value
	case reflect.Uintptr, reflect.Ptr, reflect.UnsafePointer:
		fallthrough
	default:
		return nil, errors.New("anvil:not implemented for " + v.Kind().String())
	}
	if len(items) > 0 {
		return items, err
	}
	if empty && skip {
		return nil, err
	}
	return append(s.items, Item{Key: key, Value: value}), err
}

// modify - call modifier function if presented for a given type
func (s *Anvil) modify(v reflect.Value) (interface{}, bool, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("anvil: %v on appendix call", r)
		}
	}()
	if fn, ok := s.modifier[v.Type().String()]; ok {
		return fn(v)
	}
	return nil, true, err
}

// key of field
func (s *Anvil) key(pref string, v reflect.StructField, omit bool) string {
	if omit {
		return pref
	}
	var title string
	json, ok := v.Tag.Lookup("json")
	if !ok || len(json) < 1 {
		title = v.Name
	} else {
		tags := strings.Split(json, ",")
		if len(tags[0]) > 1 || tags[0] != "-" {
			title = tags[0]
		}
	}
	if len(title) < 1 {
		title = v.Name
	}
	return pref + s.Glue + title
}

// arrayPrefix - make a notation prefix for a slice/array fields
func arrayPrefix(pref string, idx int) string {
	return pref + "[" + strconv.Itoa(idx) + "]"
}

// mapPrefix - make a notation prefix for a map fields
func mapPrefix(pref string, idx reflect.Value) string {
	var val string
	switch idx.Kind() {
	case reflect.String:
		val = idx.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val = strconv.FormatInt(idx.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val = strconv.FormatUint(idx.Uint(), 10)
	case reflect.Float32:
		val = strconv.FormatFloat(idx.Float(), 'f', -1, 32)
	case reflect.Float64:
		val = strconv.FormatFloat(idx.Float(), 'f', -1, 64)
	case reflect.Bool:
		val = strconv.FormatBool(idx.Bool())
	}
	return pref + "[" + val + "]"
}
