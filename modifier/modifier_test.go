package modifier

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_TimeMod(t *testing.T) {
	expected := "2019-04-22T15:49:32.556091+03:00"
	tm, _ := time.Parse(time.RFC3339Nano, expected)
	v := reflect.ValueOf(tm)

	val, empty, err := Time(v)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if empty {
		t.Error("`empty` must be `false`")
		t.FailNow()
	}
	if val.(string) != expected {
		t.Errorf("`time` must be equal with %v", val)
		t.FailNow()
	}
}

func Test_UUIDMod(t *testing.T) {
	expected := "6354e816-551d-11e9-92ee-acde48001122"
	v := reflect.ValueOf(uuid.MustParse(expected))

	val, empty, err := UUID(v)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if empty {
		t.Error("`empty` must be `false`")
		t.FailNow()
	}
	if val.(string) != expected {
		t.Errorf("`val` must be equal with %v", val)
		t.FailNow()
	}
}
