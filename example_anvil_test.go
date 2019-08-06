// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package anvil_test

import (
	"fmt"
	"time"

	"github.com/iveronanomi/anvil"
)

type (
	IFace interface {
		Name() interface{}
		Complex128() complex128
	}
	Embedded struct {
		Boolean bool
	}
	Sliced struct {
		Key   string
		Value interface{}
		Bool  *bool
	}
	PointerStr struct {
		F1 []string
		F2 []Sliced
	}
	Nested struct {
		*Nested
	}
	Digits struct {
		Int     int
		Int8    int8
		Int16   int16
		Int32   int32
		Int64   int64
		Uint    uint `json:"zero"`
		Uint8   uint8
		Uint16  uint16
		Uint32  uint32
		Uint64  uint64
		Float32 float32
		Float64 float64
	}
	Test struct {
		Embedded
		unexported string
		Pointer    *string
		Json       int8 `json:"json_tag"`
		PointerStr *PointerStr
		Time       time.Time
		Face       IFace `json:"-,"`
		digits     Digits
	}
)

// ExampleNotation_with_SkipEmpty demonstrates a technique
// for make notation from structure without empty values
func ExampleNotation_with_SkipEmpty() {
	v := Test{
		Embedded: Embedded{
			Boolean: true,
		},
		unexported: "string_val",
		Json:       1,
		PointerStr: &PointerStr{
			F2: []Sliced{},
		},
		digits: Digits{
			Int8:    -1,
			Int16:   -16,
			Int32:   -32,
			Int64:   -64,
			Uint8:   8,
			Uint16:  16,
			Uint32:  32,
			Uint64:  64,
			Float32: .32,
			Float64: -.64,
		},
	}
	items, _ := anvil.Notation(v, anvil.SkipEmpty, ".")

	for i := range items {
		fmt.Printf("%v\n", items[i])
	}
	// Output:
	// {Test.Embedded.Boolean true}
	// {Test.unexported string_val}
	// {Test.json_tag 1}
	// {Test.digits.Int8 -1}
	// {Test.digits.Int16 -16}
	// {Test.digits.Int32 -32}
	// {Test.digits.Int64 -64}
	// {Test.digits.Uint8 8}
	// {Test.digits.Uint16 16}
	// {Test.digits.Uint32 32}
	// {Test.digits.Uint64 64}
	// {Test.digits.Float32 0.32}
	// {Test.digits.Float64 -0.64}
}

func ExampleAnvil_Notation_map_string_keys() {
	type Str struct {
		Map map[string]string
	}
	m := map[string]string{
		"One": "Uno",
		"Two": "Dos",
	}

	v := Str{Map: m}

	a := &anvil.Anvil{Mode: anvil.SkipEmpty, Glue: "."}

	r, _ := a.Notation(v)

	for i := range r {
		fmt.Println(r[i])
	}
	// Unordered output:
	// {Str.Map[One] Uno}
	// {Str.Map[Two] Dos}
}
