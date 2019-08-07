// Copyright (c) 2019, Ivan Eremin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package anvil

import (
	"testing"

	"github.com/iveronanomi/anvil/modifier"
)

// to avoid compiler optimisations eliminating the function under test
// and artificially lowering the run time of the benchmark.
var trash interface{}

func BenchmarkNotation_WithNoSkip(b *testing.B) {
	var r interface{}
	v := MyType{}
	a := Anvil{Mode: NoSkipEmpty, Glue: "."}
	a.RegisterModifierFunc(v, modifier.Time)

	for n := 0; n < b.N; n++ {
		r, _ = a.Notation(v)
	}

	trash = r
}

func BenchmarkNotation_WithForceSkip(b *testing.B) {
	var r interface{}
	v := MyType{}
	a := &Anvil{Mode: SkipEmpty, Glue: "."}
	a.RegisterModifierFunc(v, modifier.Time)

	for n := 0; n < b.N; n++ {
		r, _ = a.Notation(v)
	}

	trash = r
}
