// Copyright (c) 2019, Ivan Eremin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package anvil_test

import (
	"fmt"

	"github.com/iveronanomi/anvil"
	"github.com/iveronanomi/anvil/modifier"
)

type (
	// Drinks it's a set of drinks
	Drinks []Drink
	// Drink structure used as example of stringer representation
	Drink struct {
		Title string
	}
)

// String representation of the Drink instance
func (p Drink) String() string {
	return fmt.Sprintf("%s is a good drink.", p.Title)
}

// This example demonstrates how to transmits a structure with stringer interface to string with a String modifier
func Example_stringModifier() {
	// bunch of drinks as data example
	data := Drinks{
		{
			Title: "Tea",
		},
		{
			Title: "Coffee",
		},
	}

	// create anvil instance, and register string modifier for Drink structure
	do := &anvil.Anvil{Mode: anvil.SkipEmpty, Glue: "."}
	do.RegisterModifierFunc(Drink{}, modifier.String)

	// extract notation from Drinks list
	items, _ := do.Notation(data)

	for i := range items {
		fmt.Printf("%v\n", items[i])
	}

	// Output:
	// {Drinks[0] Tea is a good drink.}
	// {Drinks[1] Coffee is a good drink.}
}
