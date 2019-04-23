package anvil

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type (
	// mode of skipping empty values
	mode int
	// Anvil executor structure
	Anvil struct {
		//Mode behavior for skipping empty values
		Mode mode
		//Glue string to glue fields
		Glue string
		//sample of structure to make a notation
		sample interface{}
		// modifier it's a list of functions to use as a rule
		// to find out empty or not empty value of a field with given type
		// exported type key as a key and list of functions to execute.
		// functions executing in fifo order
		modifier map[string]func(f reflect.Value) (interface{}, bool, error)
		// collection of []{key => value}
		items []Item
		deep  int
	}
	// Item field with typed value
	// as result of notation
	Item struct {
		Key   string
		Value interface{}
		//Kind  reflect.Kind
		//Empty bool
	}
)

const (
	// NoSkip fields with empty values
	NoSkip mode = iota
	// Skip fields with empty values
	Skip
)

// Modifier add function executor to extract value given type as
// result of callback function used (value, isEmpty, error) where
// value is interface{} value to assign, isEmpty valuable for
// behaviour Mode, and error if error occurred, used to stop execution
func (s *Anvil) Modifier(t interface{}, call func(f reflect.Value) (interface{}, bool, error)) *Anvil {
	if s.modifier == nil {
		s.modifier = make(map[string]func(f reflect.Value) (interface{}, bool, error))
	}
	s.modifier[reflect.TypeOf(t).String()] = call
	return s
}

// Notation structure as a plain list of [][]interface{}
// where k(string) = v(typed interface)
func Notation(sample interface{}, behaviour mode, glue string) ([]Item, error) {
	s := &Anvil{
		sample:   sample,
		Glue:     glue,
		Mode:     behaviour,
		modifier: make(map[string]func(f reflect.Value) (interface{}, bool, error)),
	}
	return s.notation("", reflect.ValueOf(sample), false)
}

// Notation make a notation of sample structure
// with fields glued by separator in one row
func (s *Anvil) Notation(sample interface{}) ([]Item, error) {
	return s.notation("", reflect.ValueOf(sample), false)
}

// slicePrefix - make a notation prefix before slice fields
func slicePrefix(key string, idx int) string {
	return key + "[" + strconv.Itoa(idx) + "]"
}

// notation structure nested
func (s *Anvil) notation(key string, v reflect.Value, title bool) (items []Item, err error) {
	var (
		value interface{}
		empty = true
		skip  = s.Mode == Skip
	)
	// get value by pointer if it is
	v = reflect.Indirect(v)

	// set default prefix for a field
	if len(key) < 1 {
		key = v.Type().Name()
	}
	//fmt.Printf("switch: %d, %s\n", j, v.Kind())
	switch v.Kind() {
	case reflect.Invalid:
		return nil, errors.New("squeezer:invalid value of " + v.Type().Name())
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
			if v.Index(i).CanAddr() {
				n, err := s.notation(slicePrefix(key, i), v.Index(i).Addr(), true)
				if err != nil {
					return nil, err
				}
				if len(n) < 1 {
					continue
				}
				items = append(items, n...)
			}
		}
	case reflect.Slice:
		if v.IsNil() {
			break
		}
		if v.Len() < 1 {
			break
		}
		for i := 0; i < v.Len(); i++ {
			if v.Index(i).CanAddr() {
				n, err := s.notation(slicePrefix(key, i), v.Index(i).Addr(), true)
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
			//log.Printf("field:v.Type().Field(%d).Name=%s key=%s\n",i, v.Type().Field(i).Name, s.key(key, v.Type().Field(i),false))
			n, err := s.notation(s.key(key, v.Type().Field(i), false), v.Field(i), true)
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
		var n []Item
		if n, err = s.notation(key, v.Elem(), true); err != nil {
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
	case reflect.Complex64, reflect.Complex128, reflect.Uintptr, reflect.Map:
		return nil, errors.New("squeezer:not implemented for " + v.Kind().String())
	case reflect.Ptr, reflect.UnsafePointer:
		fallthrough
	default:
		return nil, errors.New("squeezer:not implemented for " + v.Kind().String())
	}
	//fmt.Printf("{idx}: %d:\n{key}: %s\n{value}: %v\n{items}: %v\n{skip}: %v\n{empty}: %v\n",
	//	j, key, value, s.items, skip, empty)
	if len(items) > 0 {
		return items, err
	}
	if empty && skip {
		return nil, err
	}
	return append(s.items, Item{Key: key, Value: value}), err
}

// call modifier fn if it possible
func (s *Anvil) modify(v reflect.Value) (interface{}, bool, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("squeezer: %v on appendix call for type", r)
		}
	}()
	if fn, ok := s.modifier[v.Type().String()]; ok {
		return fn(v)
	}
	return nil, true, err
}

// key of
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