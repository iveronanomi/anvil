package anvil

import (
	"testing"
	"time"

	"github.com/iveronanomi/anvil/modifier"
)

// MyType structures
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
	MyType struct {
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

func TestAnvil_Notation_TimeModifier_ExpectedStringValue(t *testing.T) {
	v := time.Now()
	expected := []Item{
		{Key: "Time", Value: v.Format(time.RFC3339Nano)},
	}
	a := &Anvil{Mode: NoSkip, Glue: "."}
	a.Modifier(time.Time{}, modifier.Time)

	r, err := a.Notation(v)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	check(t, expected, r)
}

func TestAnvil_Notation_NoSkip(t *testing.T) {
	s := "string_val"
	f1 := []string{"one", "two", "three"}
	clock := time.Now()
	v := MyType{
		Embedded: Embedded{
			Boolean: true,
		},
		unexported: s, // todo: check `json` tag behaviors
		Pointer:    &s,
		Json:       1,
		PointerStr: &PointerStr{
			F1: f1,
			F2: []Sliced{},
		},
		Time: clock,
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
	expected := []Item{
		{Key: "MyType.Embedded.Boolean", Value: v.Embedded.Boolean},
		{Key: "MyType.unexported", Value: v.unexported},
		{Key: "MyType.Pointer", Value: *v.Pointer},
		{Key: "MyType.json_tag", Value: int8(1)},
		{Key: "MyType.PointerStr.F1[0]", Value: v.PointerStr.F1[0]},
		{Key: "MyType.PointerStr.F1[1]", Value: v.PointerStr.F1[1]},
		{Key: "MyType.PointerStr.F1[2]", Value: v.PointerStr.F1[2]},
		{Key: "MyType.PointerStr.F2", Value: nil},
		{Key: "MyType.Time", Value: clock.Format(time.RFC3339Nano)},
		{Key: "MyType.Face", Value: nil},
		{Key: "MyType.digits.Int", Value: v.digits.Int},
		{Key: "MyType.digits.Int8", Value: v.digits.Int8},
		{Key: "MyType.digits.Int16", Value: v.digits.Int16},
		{Key: "MyType.digits.Int32", Value: v.digits.Int32},
		{Key: "MyType.digits.Int64", Value: v.digits.Int64},
		{Key: "MyType.digits.zero", Value: v.digits.Uint},
		{Key: "MyType.digits.Uint8", Value: v.digits.Uint8},
		{Key: "MyType.digits.Uint16", Value: v.digits.Uint16},
		{Key: "MyType.digits.Uint32", Value: v.digits.Uint32},
		{Key: "MyType.digits.Uint64", Value: v.digits.Uint64},
		{Key: "MyType.digits.Float32", Value: v.digits.Float32},
		{Key: "MyType.digits.Float64", Value: v.digits.Float64},
	}
	a := &Anvil{Mode: NoSkip, Glue: "."}
	a.Modifier(time.Time{}, modifier.Time)

	r, err := a.Notation(v)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	check(t, expected, r)
}

func TestAnvil_Notation_Skip(t *testing.T) {
	s := "string_val"
	clock := time.Now()
	//id := uuid.MustParse("18bc60b8-17a1-4548-8471-73d30d240c99")
	tr := true
	fa := false
	v := MyType{
		Embedded: Embedded{
			Boolean: true,
		},
		unexported: s,
		Pointer:    &s,
		Json:       1,
		PointerStr: &PointerStr{
			F1: []string{"", "two", " "},
			F2: []Sliced{
				{Key: "1", Value: 1, Bool: &tr},
				{Key: "2", Value: "2", Bool: &fa},
				{Key: "0", Value: 0, Bool: &fa},
			},
		},
		Time: clock,
		//UUID: &id,
	}
	expected := []Item{
		{Key: "MyType.Embedded.Boolean", Value: true},
		{Key: "MyType.unexported", Value: v.unexported},
		{Key: "MyType.Pointer", Value: *v.Pointer},
		{Key: "MyType.json_tag", Value: v.Json},
		{Key: "MyType.PointerStr.F1[1]", Value: v.PointerStr.F1[1]},
		{Key: "MyType.PointerStr.F1[2]", Value: v.PointerStr.F1[2]},
		{Key: "MyType.PointerStr.F2[0].Key", Value: v.PointerStr.F2[0].Key},
		{Key: "MyType.PointerStr.F2[0].Value", Value: v.PointerStr.F2[0].Value},
		{Key: "MyType.PointerStr.F2[0].Bool", Value: *v.PointerStr.F2[0].Bool},
		{Key: "MyType.PointerStr.F2[1].Key", Value: v.PointerStr.F2[1].Key},
		{Key: "MyType.PointerStr.F2[1].Value", Value: v.PointerStr.F2[1].Value},
		//todo: bool pointer is it empty with a `false` value?
		//{Key: "MyType.PointerStr.F2[1].Bool", Value: *v.PointerStr.F2[1].Bool},
		{Key: "MyType.PointerStr.F2[2].Key", Value: v.PointerStr.F2[2].Key},
		//{Key: "MyType.PointerStr.F2[2].Bool", Value: *v.PointerStr.F2[2].Bool},
		{Key: "MyType.Time", Value: clock.Format(time.RFC3339Nano)},
		//{Key: "MyType.uuid", Value: v.UUID.String()},
	}
	a := &Anvil{Mode: Skip, Glue: "."}
	a.Modifier(time.Time{}, modifier.Time)

	r, err := a.Notation(v)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	check(t, expected, r)
}

func check(t *testing.T, expected, occurred []Item) {
	t.Helper()
	var (
		a, b   = expected, occurred
		at, bt = "expected", "occurred"
		fail   = len(a) != len(b)
		err    = "% 3d| %s:\033[00;31m%v\033[00m %s:\033[00;31m%v\033[00m"
		info   = "% 3d| %s:\033[00;30m%v\033[00m %s:\033[00;30m%v\033[00m"
		skip   = "% 3d| %s:\033[00;31m%v\033[00m %s: -"
		length = "%s \033[00;31m%d\033[00m fields %s:\033[00;31m%d\033[00m"
	)
	if len(a) != len(b) {
		t.Errorf(length, at, len(a), bt, len(b))
	}
	if len(a) < len(b) {
		a, b, at, bt = b, a, bt, at
	}
	var i int
	for ; i < len(b); i++ {
		if a[i].Key == b[i].Key && a[i].Value == b[i].Value {
			t.Logf(info, i, at, a[i], bt, b[i])
			continue
		}
		fail = true
		t.Logf(err, i, at, a[i], bt, b[i])
	}
	for ; i < len(a); i++ {
		t.Logf(skip, i, at, a[i], bt)
	}

	if fail {
		t.FailNow()
	}
}
