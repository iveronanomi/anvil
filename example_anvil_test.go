// Copyright (c) 2019, Ivan Eremin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package anvil_test

import (
	"fmt"
	"github.com/iveronanomi/anvil/modifier"
	"time"

	"github.com/iveronanomi/anvil"
)

type (
	SolarSystem []Planet
	Planet      struct {
		Name  string `json:"name"`
		mass  float32
		rings bool
		moons int8
	}
)

// ExampleNotation_noSkipEmpty demonstrates a technique
// for make notation from structure without empty values
func ExampleNotation_noSkipEmpty() {
	source := SolarSystem{
		{
			Name: "Mercury",
			mass: .0553,
		},
		{
			Name: "Venus",
			mass: .815,
		},
		{
			Name:  "Earth",
			mass:  1,
			moons: 1,
		},
		{
			Name:  "Mars",
			mass:  .11,
			moons: 2,
		},
		{
			Name:  "Jupiter",
			mass:  317.8,
			rings: true,
			moons: 79,
		},
		{
			Name:  "Saturn",
			mass:  95.2,
			rings: true,
			moons: 62,
		},
		{
			Name:  "Uranus",
			mass:  14.6,
			rings: true,
			moons: 27,
		},
		{
			Name:  "Neptune",
			mass:  17.2,
			rings: true,
			moons: 14,
		},
	}
	items, _ := anvil.Notation(source, anvil.NoSkipEmpty, ".")

	for i := range items {
		fmt.Printf("%#v\n", items[i])
	}
	// Output:
	// anvil.Item{Key:"SolarSystem[0].name", Value:"Mercury"}
	// anvil.Item{Key:"SolarSystem[0].mass", Value:0.0553}
	// anvil.Item{Key:"SolarSystem[0].rings", Value:false}
	// anvil.Item{Key:"SolarSystem[0].moons", Value:0}
	// anvil.Item{Key:"SolarSystem[1].name", Value:"Venus"}
	// anvil.Item{Key:"SolarSystem[1].mass", Value:0.815}
	// anvil.Item{Key:"SolarSystem[1].rings", Value:false}
	// anvil.Item{Key:"SolarSystem[1].moons", Value:0}
	// anvil.Item{Key:"SolarSystem[2].name", Value:"Earth"}
	// anvil.Item{Key:"SolarSystem[2].mass", Value:1}
	// anvil.Item{Key:"SolarSystem[2].rings", Value:false}
	// anvil.Item{Key:"SolarSystem[2].moons", Value:1}
	// anvil.Item{Key:"SolarSystem[3].name", Value:"Mars"}
	// anvil.Item{Key:"SolarSystem[3].mass", Value:0.11}
	// anvil.Item{Key:"SolarSystem[3].rings", Value:false}
	// anvil.Item{Key:"SolarSystem[3].moons", Value:2}
	// anvil.Item{Key:"SolarSystem[4].name", Value:"Jupiter"}
	// anvil.Item{Key:"SolarSystem[4].mass", Value:317.8}
	// anvil.Item{Key:"SolarSystem[4].rings", Value:true}
	// anvil.Item{Key:"SolarSystem[4].moons", Value:79}
	// anvil.Item{Key:"SolarSystem[5].name", Value:"Saturn"}
	// anvil.Item{Key:"SolarSystem[5].mass", Value:95.2}
	// anvil.Item{Key:"SolarSystem[5].rings", Value:true}
	// anvil.Item{Key:"SolarSystem[5].moons", Value:62}
	// anvil.Item{Key:"SolarSystem[6].name", Value:"Uranus"}
	// anvil.Item{Key:"SolarSystem[6].mass", Value:14.6}
	// anvil.Item{Key:"SolarSystem[6].rings", Value:true}
	// anvil.Item{Key:"SolarSystem[6].moons", Value:27}
	// anvil.Item{Key:"SolarSystem[7].name", Value:"Neptune"}
	// anvil.Item{Key:"SolarSystem[7].mass", Value:17.2}
	// anvil.Item{Key:"SolarSystem[7].rings", Value:true}
	// anvil.Item{Key:"SolarSystem[7].moons", Value:14}
}

// ExampleAnvil_Notation_map demonstrates a technique
// for make notation from map
func ExampleAnvil_Notation_map() {
	source := map[string]string{
		"One": "Uno",
		"Two": "Dos",
	}
	squeezer := &anvil.Anvil{Mode: anvil.SkipEmpty, Glue: "."}

	items, _ := squeezer.Notation(source)

	for i := range items {
		fmt.Printf("%#v\n", items[i])
	}
	// Output:
	// anvil.Item{Key:"[One]", Value:"Uno"}
	// anvil.Item{Key:"[Two]", Value:"Dos"}
}

type DateRange struct {
	Range map[string]time.Time
}

// ExampleAnvil_RegisterModifierFunc describes how to use modifiers to change behavior of types representation
func ExampleAnvil_RegisterModifierFunc() {
	t1, _ := time.Parse(time.Kitchen, "1:01AM")
	source := DateRange{
		Range: map[string]time.Time{
			"date_1": t1,
			"date_2": t1.Add(-time.Minute),
		},
	}

	squeezer := &anvil.Anvil{Mode: anvil.SkipEmpty, Glue: "."}

	items, _ := squeezer.
		RegisterModifierFunc(time.Time{}, modifier.Time).
		Notation(source)

	for i := range items {
		fmt.Printf("%#v\n", items[i])
	}
	// Output:
	// anvil.Item{Key:"DateRange.Range[date_1]", Value:"0000-01-01T01:01:00Z"}
	// anvil.Item{Key:"DateRange.Range[date_2]", Value:"0000-01-01T01:00:00Z"}
}
