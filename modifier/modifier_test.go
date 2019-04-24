package modifier

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func Test_TimeMod(t *testing.T) {
	expectedVal := "2019-04-22T15:49:32.556091+03:00"
	tm, _ := time.Parse(time.RFC3339Nano, expectedVal)
	v := reflect.ValueOf(tm)

	val, empty, err := Time(v)

	if err != nil {
		t.Error(err)
	}
	if empty {
		t.Error("`empty` must be `false`")
	}
	if val.(string) != expectedVal {
		t.Errorf("`time` must be equal with %v", val)
	}
	if t.Failed() {
		t.FailNow()
	}
}

func Test_TimeMod_WithZeroTime(t *testing.T) {
	expectedVal := "0001-01-01T00:00:00Z"
	v := reflect.ValueOf(time.Time{})

	val, empty, err := Time(v)

	if err != nil {
		t.Errorf("`err` must be a <nil>, occured `%v`", err)
	}
	if val != expectedVal {
		t.Errorf("`val` must be `%v` for zero time value", expectedVal)
	}
	if !empty {
		t.Error("`empty` must be `true` for zero time value")
	}
	if t.Failed() {
		t.FailNow()
	}
}

type hasIsZeroMethod struct{}

func (t hasIsZeroMethod) IsZero() bool {
	return false
}

func Test_TimeMod_WithNotImplementedIsZeroMethod(t *testing.T) {
	v := reflect.ValueOf(int(0))
	expectedVal := ""
	expectedErr := errors.New("modifier:method IsZero not implemented")

	val, empty, err := Time(v)

	if err == nil {
		t.Errorf("`err` must be `%v`", expectedErr)
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf("expected `%v` error", expectedErr)
	}
	if !empty {
		t.Errorf("`empty` must be `true`, not `%v`", empty)
	}
	if val != expectedVal {
		t.Errorf("`val` must be <nil>, not `%v`", val)
	}
	if t.Failed() {
		t.FailNow()
	}
}

func Test_TimeMod_WithNotImplementedFormatMethod(t *testing.T) {
	v := reflect.ValueOf(hasIsZeroMethod{})
	expectedVal := ""
	expectedErr := errors.New("modifier:method Format not implemented")

	val, empty, err := Time(v)

	if err == nil {
		t.Errorf("`err` must be `%v`", expectedErr)
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf("expected `%v` error", expectedErr)
	}
	if !empty {
		t.Errorf("`empty` must be `true`, not `%v`", empty)
	}
	if val != expectedVal {
		t.Errorf("`val` must be <nil>, not `%v`", val)
	}
	if t.Failed() {
		t.FailNow()
	}
}

type dummyUUIDType struct {
	str string
	id  uint
}

func (f dummyUUIDType) String() string { return f.str }
func (f dummyUUIDType) ID() uint       { return f.id }

func Test_UUIDMod(t *testing.T) {
	expectedVal := "6354e816-551d-11e9-92ee-acde48001122"
	uuid := dummyUUIDType{id: 1, str: expectedVal}
	v := reflect.ValueOf(uuid)

	val, empty, err := UUID(v)

	if err != nil {
		t.Errorf("`err` must be <nil>, not %v", err)
	}
	if empty {
		t.Errorf("`empty` must be `false`, not %v", empty)
	}
	if val.(string) != expectedVal {
		t.Errorf("`val` must be equal to `%v` not `%v`", expectedVal, val)
	}
	if t.Failed() {
		t.FailNow()
	}
}

func Test_UUIDMod_WithZeroVal(t *testing.T) {
	expectedVal := "00000000-0000-0000-0000-000000000000"
	uuid := dummyUUIDType{id: 0, str: expectedVal}
	v := reflect.ValueOf(uuid)

	val, empty, err := UUID(v)

	if err != nil {
		t.Errorf("`err` must be <nil>, not %v", err)
	}
	if !empty {
		t.Errorf("`empty` must be `true`, not `%v`", empty)
	}
	if val.(string) != expectedVal {
		t.Errorf("`val` must be equal to `%v` not `%v`", expectedVal, val)
	}
	if t.Failed() {
		t.FailNow()
	}
}

func Test_UUIDMod_WithNotImplementedIDMethod(t *testing.T) {
	v := reflect.ValueOf(0)
	expectedVal := ""
	expectedErr := errors.New("modifier:method ID not implemented")

	val, empty, err := UUID(v)

	if err == nil {
		t.Errorf("`err` must be `%v`", expectedErr)
	}
	if !empty {
		t.Errorf("`empty` must be `true`, not `%v`", empty)
	}
	if val != expectedVal {
		t.Errorf("`val` must be <nil>, not `%v`", val)
	}
	if t.Failed() {
		t.FailNow()
	}
}

type hasIDMethod struct{}

func (t hasIDMethod) ID() uint { return 1 }

func Test_UUIDMod_WithNotImplementedStringMethod(t *testing.T) {
	v := reflect.ValueOf(hasIDMethod{})
	expectedVal := ""
	expectedErr := errors.New("modifier:method String not implemented")

	val, empty, err := UUID(v)

	if err == nil {
		t.Errorf("`err` must be `%v`", expectedErr)
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf("expected `%v` error", expectedErr.Error())
	}
	if !empty {
		t.Errorf("`empty` must be `true`, not `%v`", empty)
	}
	if val != expectedVal {
		t.Errorf("`val` must be <nil>, not `%v`", val)
	}
	if t.Failed() {
		t.FailNow()
	}
}
