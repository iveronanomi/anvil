package anvil

import (
	"reflect"
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
	Complex struct {
		Complex64  complex64
		Complex128 complex128
	}
	MyType struct {
		Embedded
		unexported string
		Pointer    *string
		JSON       int8 `json:"json_tag"`
		PointerStr *PointerStr
		Time       time.Time
		Face       IFace `json:"-,"`
		digits     Digits
	}
)

func TestAnvil_Notation_WithEmptySample(t *testing.T) {
	n, err := Notation(nil, SkipEmpty, ".")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(n) > 0 {
		t.Error("must be zero value for nil type")
		t.FailNow()
	}
}

func TestMapPrefix_WithIntKey(t *testing.T) {
	i := int(1)
	k := reflect.ValueOf(i)

	pref := mapPrefix("", k)

	if pref != "[1]" {
		t.Error("invalid value of map prefix for int key value")
	}
}

func TestMapPrefix_WithStringKey(t *testing.T) {
	i := "key"
	k := reflect.ValueOf(i)

	pref := mapPrefix("", k)

	if pref != "[key]" {
		t.Error("invalid value of map prefix for string key value")
	}
}

func TestMapPrefix_WithUintKey(t *testing.T) {
	i := uint(2)
	k := reflect.ValueOf(i)

	pref := mapPrefix("", k)

	if pref != "[2]" {
		t.Error("invalid value of map prefix for string key value")
	}
}

func TestMapPrefix_WithFloat32Key(t *testing.T) {
	i := float32(.1)
	k := reflect.ValueOf(i)

	pref := mapPrefix("", k)

	if pref != "[0.1]" {
		t.Error("invalid value of map prefix for string key value")
	}
}

func TestMapPrefix_WithFloat64Key(t *testing.T) {
	i := float64(.2)
	k := reflect.ValueOf(i)

	pref := mapPrefix("", k)

	if pref != "[0.2]" {
		t.Error("invalid value of map prefix for string key value")
	}
}

func TestMapPrefix_WithBool(t *testing.T) {
	k := reflect.ValueOf(true)

	pref := mapPrefix("", k)

	if pref != "[true]" {
		t.Error("invalid value of map prefix for string key value")
	}
}

func TestAnvil_Notation_TimeModifier_ExpectedStringValue(t *testing.T) {
	v := time.Now()
	expected := []Item{
		{Key: "Time", Value: v.Format(time.RFC3339Nano)},
	}
	a := &Anvil{Mode: NoSkipEmpty, Glue: "."}
	a.RegisterModifierFunc(time.Time{}, modifier.Time)

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
		JSON:       1,
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
	a := &Anvil{Mode: NoSkipEmpty, Glue: "."}
	a.RegisterModifierFunc(time.Time{}, modifier.Time)

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
	tr := true
	fa := false
	v := MyType{
		Embedded: Embedded{
			Boolean: true,
		},
		unexported: s,
		Pointer:    &s,
		JSON:       1,
		PointerStr: &PointerStr{
			F1: []string{"", "two", " "},
			F2: []Sliced{
				{Key: "1", Value: 1, Bool: &tr},
				{Key: "2", Value: "2", Bool: &fa},
				{Key: "0", Value: 0, Bool: &fa},
			},
		},
		Time: clock,
	}
	expected := []Item{
		{Key: "MyType.Embedded.Boolean", Value: true},
		{Key: "MyType.unexported", Value: v.unexported},
		{Key: "MyType.Pointer", Value: *v.Pointer},
		{Key: "MyType.json_tag", Value: v.JSON},
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
	a := &Anvil{Mode: SkipEmpty, Glue: "."}
	a.RegisterModifierFunc(time.Time{}, modifier.Time)

	r, err := a.Notation(v)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	check(t, expected, r)
}

func TestAnvil_Notation_Map_WithInt16Keys(t *testing.T) {
	t.Skip()
	type Str struct {
		Map map[int16]string
	}
	m := map[int16]string{
		-1: "One",
		2:  "Two",
	}
	expected := []Item{
		{Key: "Str.Map[-1]", Value: "One"},
		{Key: "Str.Map[2]", Value: "Two"},
	}
	v := Str{Map: m}

	a := &Anvil{Mode: SkipEmpty, Glue: "."}

	r, err := a.Notation(v)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	check(t, expected, r)
}

func TestAnvil_Notation_Map_WithUint8Keys(t *testing.T) {
	t.Skip()
	type Str struct {
		Map map[uint8]string
	}
	m := map[uint8]string{
		1: "One",
		2: "Two",
	}
	expected := []Item{
		{Key: "Str.Map[1]", Value: "One"},
		{Key: "Str.Map[2]", Value: "Two"},
	}
	v := Str{Map: m}

	a := &Anvil{Mode: SkipEmpty, Glue: "."}

	r, err := a.Notation(v)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	check(t, expected, r)
}

func TestAnvil_Notation_Map_WithFloat64Keys(t *testing.T) {
	t.Skip()
	type Str struct {
		Map map[float64]string
	}
	m := map[float64]string{
		.12345678901: "One",
		-23456789.01: "Two",
	}
	expected := []Item{
		{Key: "Str.Map[0.12345678901]", Value: "One"},
		{Key: "Str.Map[-23456789.01]", Value: "Two"},
	}
	v := Str{Map: m}

	a := &Anvil{Mode: SkipEmpty, Glue: "."}

	r, err := a.Notation(v)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	check(t, expected, r)
}

func TestAnvil_Notation_Map_WithBoolKeys(t *testing.T) {
	t.Skip()
	type Str struct {
		MapBool map[bool]string
	}
	expected := []Item{
		{Key: "Str.MapBool[true]", Value: "Uno"},
		{Key: "Str.MapBool[false]", Value: "Dos"},
	}
	m := map[bool]string{
		true:  "Uno",
		false: "Dos",
	}
	v := Str{MapBool: m}

	a := &Anvil{Mode: SkipEmpty, Glue: "."}

	r, err := a.Notation(v)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	check(t, expected, r)
}

// in case of not implemented types
func TestNotation_Map_WithNotImplementedKeysTypes_ExpectedEmptyKey(t *testing.T) {
	t.Skip()
	type Str struct {
		MapBool map[struct{ T string }]string
	}
	expected := []Item{
		{Key: "Str.MapBool[]", Value: "One"},
		{Key: "Str.MapBool[]", Value: "Two"},
	}
	m := map[struct{ T string }]string{
		struct{ T string }{T: "Uno"}: "One",
		struct{ T string }{T: "Dos"}: "Two",
	}
	v := Str{MapBool: m}

	r, err := Notation(v, SkipEmpty, ".")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	check(t, expected, r)
}

func TestNotation_Complex(t *testing.T) {
	v := Complex{
		Complex64:  complex64(complex(.1, .0)),
		Complex128: complex128(complex(.0, .1)), // must be skipped
	}
	expected := []Item{
		{Key: "Complex.Complex64", Value: v.Complex64},
		{Key: "Complex.Complex128", Value: v.Complex128},
	}

	r, err := Notation(v, SkipEmpty, ".")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	check(t, expected, r)
}

func TestNotation_Interface(t *testing.T) {
	type Interface struct {
		v interface{}
	}
	v := Interface{v: 1}
	expected := []Item{
		{Key: "Interface.v", Value: 1},
	}
	r, err := Notation(v, SkipEmpty, ".")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	check(t, expected, r)
}

func TestNotation_Array(t *testing.T) {
	type Array struct {
		Val [8]uint
	}
	v := Array{Val: [8]uint{1, 2}}
	expected := []Item{
		{Key: "Array.Val[0]", Value: uint(1)},
		{Key: "Array.Val[1]", Value: uint(2)},
	}
	r, err := Notation(v, SkipEmpty, ".")

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
