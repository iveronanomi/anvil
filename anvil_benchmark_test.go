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
	v := Test{}
	a := Anvil{Mode: NoSkip, Glue: "."}
	a.Modifier(v, modifier.Time)

	for n := 0; n < b.N; n++ {
		r, _ = a.Notation(v)
	}

	trash = r
}

func BenchmarkNotation_WithForceSkip(b *testing.B) {
	var r interface{}
	v := Test{}
	a := &Anvil{Mode: NoSkip, Glue: "."}
	a.Modifier(v, modifier.Time)

	for n := 0; n < b.N; n++ {
		r, _ = a.Notation(v)
	}

	trash = r
}