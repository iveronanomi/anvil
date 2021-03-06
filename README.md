[![Build Status](https://travis-ci.org/iveronanomi/anvil.svg?branch=master)](https://travis-ci.org/iveronanomi/anvil) [![Coverage Status](https://coveralls.io/repos/github/iveronanomi/anvil/badge.svg)](https://coveralls.io/github/iveronanomi/anvil) [![Go Report Card](https://goreportcard.com/badge/github.com/src-d/go-git)](https://goreportcard.com/report/github.com/src-d/go-git) [![GoDoc](https://godoc.org/github.com/iveronanomi/anvil?status.svg)](https://godoc.org/github.com/iveronanomi/anvil)

# Anvil - Dot notation from Go type instance
- [What is going on here?](#what-is-going-on)
- [Modifier usage](#modifier-usage)
- [TODO features](#todo-features)

What is going on here?
Two main params used to make a notation from a type:
- `Glue` as a glue for notation
- `Mode` as a mode for skipping empty values of a type.
	- `anvil.SkipEmpty` - skip empty values of type
	- `anvil.NoSkipEmpty` - do not skip empty values of type

In case of structure field have a `json` tag name - tag used as a name for a field in notation

## What is going on
```go
v := Test{
	Embedded: Embedded{
		Boolean: true,
	},
	unexported: "string_val", // todo: check `json` tag behaviors
	Json:       1,
	PointerStr: &PointerStr{
		F2: []Sliced{},
	},
	digits: Digits{
		Int:     0,
		Int8:    -1,
		Int16:   -16,
		Int32:   -32,
		Int64:   -64,
		Uint:    0,
		Uint8:   8,
		Uint16:  16,
		Uint32:  32,
		Uint64:  64,
		Float32: .32,
		Float64: -.64,
	},
}
items, _ := anvil.Notation(v, anvil.Skip, ".")
```
Result:
```go
Item{Key:"Test.Embedded.Boolean", Value:true}
Item{Key:"Test.unexported", Value:"string_val"}
Item{Key:"Test.json_tag", Value:1}
Item{Key:"Test.PointerStr.F1", Value:interface {}(nil)}
Item{Key:"Test.PointerStr.F2", Value:interface {}(nil)}
Item{Key:"Test.Face", Value:interface {}(nil)}
Item{Key:"Test.digits.Int", Value:0}
Item{Key:"Test.digits.Int8", Value:-1}
Item{Key:"Test.digits.Int16", Value:-16}
Item{Key:"Test.digits.Int32", Value:-32}
Item{Key:"Test.digits.Int64", Value:-64}
Item{Key:"Test.digits.zero", Value:0x0}
Item{Key:"Test.digits.Uint8", Value:0x8}
Item{Key:"Test.digits.Uint16", Value:0x10}
Item{Key:"Test.digits.Uint32", Value:0x20}
Item{Key:"Test.digits.Uint64", Value:0x40}
Item{Key:"Test.digits.Float32", Value:0.32}
Item{Key:"Test.digits.Float64", Value:-0.64}
```

## Modifier usage
It's a way how to represent a type in dot notation.
With complicated types (as `time.Time` in the example)
we interested in real value, not a full notation of type `time.Time`, modifiers will help with that.

```go
package main

import (
	"fmt"
	"time"

	"github.com/iveronanomi/anvil"
	"github.com/iveronanomi/anvil/modifier"
)

func main() {
	v := time.Now()

	do := anvil.Anvil{Mode:anvil.NoSkip, Glue:"."}
	do.Modifier(time.Now(), modifier.Time)
	items, _ := do.Notation(v)

	for i := range items {
		fmt.Printf("\n%#v", items[i])
	}
}
```
```go
Item{Key:"Time", Value:"2019-04-23T10:44:56.534221+03:00"}
```

### Features availability
|Type|Supported|Modifiers call|
|---:|:---:|:---:|
|`Array`|+|+|
|`Slice`|+|+|
|`Struct`|+|+|
|`Int`|+|-|
|`Int8`|+|-|
|`Int16`|+|-|
|`Int32`|+|-|
|`Int64`|+|-|
|`Float32`|+|-|
|`Float64`|+|-|
|`Uint`|+|-|
|`Uint8`|+|-|
|`Uint16`|+|-|
|`Uint32`|+|-|
|`Uint64`|+|-|
|`Bool`|+|-|
|`String`|+|-|
|`Interface`|+|-|
|`Complex64`|+|-|
|`Complex128`|+|-|
|`Map`|keys supported: Ints, Uints, Floats, Bool|-|
|`Uintptr`|-|-|
|`Ptr`|-|-|
|`UnsafePointer`|-|-|


### TODO Features
- unmarshal notation to type interface
- marshal type to notation interface
- codegeneration for unmarshal
