// Copyright (c) 2019, Ivan Eremin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package modifier_test

import (
	"fmt"

	"github.com/iveronanomi/anvil"
	"github.com/iveronanomi/anvil/modifier"
)

type (
	// SolarSystem it's a set of planets
	SolarSystem []Planet
	// Planet structure used as example of stringer representation
	Planet struct {
		Name string
	}
)

// String representation of the planet instance
func (p Planet) String() string {
	return fmt.Sprintf("Planet %s is a good place to live.", p.Name)
}

// demonstrates how to use notation with a string modifier
func Example_with_string_modifier() {
	source := SolarSystem{
		{
			Name: "Mercury",
		},
		{
			Name: "Venus",
		},
	}
	do := &anvil.Anvil{Mode: anvil.SkipEmpty, Glue: "."}
	do.RegisterModifierFunc(Planet{}, modifier.String)

	items, _ := do.Notation(source)

	for i := range items {
		fmt.Printf("%#v\n", items[i])
	}
	// Output:
	// anvil.Item{Key:"SolarSystem[0]", Value:"Planet Mercury is a good place to live."}
	// anvil.Item{Key:"SolarSystem[1]", Value:"Planet Venus is a good place to live."}
}
