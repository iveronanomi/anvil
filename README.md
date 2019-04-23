# Anvil - Dot notation from Go type instance
- [Usage](#usage)
  - [As an anvil instance](#usage-as-an-anvil-instance)
  - [As package static call](#usage-as-package-static-call)
- [Modifier usage](#modifier-usage)
- [TODO List](#todo)

What is going on here?
Two main params used to make a notation from a type:
- `Glue` as glue for notation
- `Mode` as a mode for skipping empty values of a type.
	- `anvil.Skip` - skip empty values of type
	- `anvil.NoSkip` - do not skip empty values
	
In case if field of structure have a json tag name,
this tag used as a name for a field in notation

## Usage
### Usage: As an `anvil` instance
```go
do := anvil.Anvil{Mode:anvil.NoSkip, Glue:"."}
items, _ := do.Notation(v)
```
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

### Usage: As package static call
```go
items, _ := anvil.Notation(v, anvil.Skip, ".")
```
<details><summary>Full example</summary>
<p>

### Full example

```go
import (
	"fmt"
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
	PointerStr struct {
		F1 []string
		F2 []Sliced
	}
	Sliced struct {
		Key   string
		Value interface{}
		Bool  *bool
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

func main() {
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

	for i := range items {
		fmt.Printf("\n%#v", items[i])
	}
}
```

</p>
</details>

```go
Item{Key:"Test.Embedded.Boolean", Value:true}
Item{Key:"Test.unexported", Value:"string_val"}
Item{Key:"Test.json_tag", Value:1}
Item{Key:"Test.digits.Int8", Value:-1}
Item{Key:"Test.digits.Int16", Value:-16}
Item{Key:"Test.digits.Int32", Value:-32}
Item{Key:"Test.digits.Int64", Value:-64}
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

### TODO
- [ ] add modifiers executions for all types
- [ ] optimize types itteration
- [ ] support all built-in types
