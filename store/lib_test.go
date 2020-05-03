package store

import (
	"fmt"
	"reflect"
	"testing"
)

var k = "test"
var v = "hello"
var ori = map[interface{}]interface{}{k: v}

func TestNew(t *testing.T) {
	s := New()
	if s.s == nil {
		t.Errorf("s attribute must be initialization")
		t.FailNow()
	}
}

func TestStore_Add(t *testing.T) {
	s := New()
	key := "test"
	val := "hello"
	s.Add(key, val)
	if !reflect.DeepEqual(val, s.s[key]) {
		t.Errorf("left value: %#v, right value: %#v", val, s.s[key])
		t.FailNow()
	}

	if len(s.s) != 1 {
		t.Errorf("underlying map length is not 1")
		t.FailNow()
	}
}

func TestStore_Get(t *testing.T) {
	s := New()
	key := "test"
	val := "hello"

	v, ok := s.Get(key)

	if ok || v != nil {
		t.Errorf("empty underlying map of store must can not get a value")
		t.FailNow()
	}

	v, ok = s.Add(key, val).Get(key)

	if !ok || !reflect.DeepEqual(val, v) {
		t.Errorf("Get method must can retrieve a value by the key that had stored a value before")
		t.FailNow()
	}

	k1 := "1"
	k2 := 1

	v, ok = s.Add(k1, val).Get(k2)

	if ok || v != nil {
		t.Errorf("diff type value must can not retrieve a value from the Store")
		t.FailNow()
	}
}

func TestStore_Append(t *testing.T) {
	old_key := "test"
	old_key2 := "next"
	old_val := "hello"
	s := New()
	s.Add(old_key, old_val).Add(old_key2, old_val)

	oldS := map[interface{}]interface{}{old_key: old_val, old_key2: old_val}
	appendAndExpects := [][2]map[interface{}]interface{}{
		{nil, oldS},
		{map[interface{}]interface{}{}, oldS},
		{
			map[interface{}]interface{}{
				old_key: old_key,
				old_val: old_key,
			},
			map[interface{}]interface{}{
				old_key:  old_key,
				old_val:  old_key,
				old_key2: old_val,
			},
		},
	}

	for k, v := range appendAndExpects {
		s.Append(v[0])
		if !reflect.DeepEqual(s.s, v[1]) {
			t.Errorf("round %d, left value: %#v, right value: %#v", k, s.s, v[1])
			t.FailNow()
		}
	}
}

func TestStore_With(t *testing.T) {
	s := New()
	s.Add(k, v)
	ori := map[interface{}]interface{}{k: v}
	appendAndExpects := [][2]map[interface{}]interface{}{
		{nil, ori},
		{map[interface{}]interface{}{}, ori},
		{
			map[interface{}]interface{}{
				v: k,
			},
			map[interface{}]interface{}{
				v: k,
			},
		},
	}

	for k, v := range appendAndExpects {
		s.With(v[0])
		if !reflect.DeepEqual(s.s, v[1]) {
			t.Errorf("round %d, left value: %#v, right value: %#v", k, s.s, v[1])
			t.FailNow()
		}
	}
}

func TestStore_Underlying(t *testing.T) {
	s := New()
	s.Add(k, v)
	if !reflect.DeepEqual(s.Underlying(), ori) {
		t.Errorf("left is: %#v, right is: %#v", s.Underlying(), ori)
	}
}

func TestStore_MustGet_Get(t *testing.T) {
	s := New()
	s.Add(k, v)
	if !reflect.DeepEqual(s.MustGet(k), v) {
		t.Errorf("left is: %#v, right is: %#v", s.MustGet(k), v)
	}
}

func TestStore_MustGet_Panic(t *testing.T) {
	s := New()
	assertPanic(t, func() {
		s.MustGet(k)
	})
}

func TestStore_Bind(t *testing.T) {
	s := New()
	k1 := "test"
	v1 := "hello"
	s.Add(k1, v1)

	_, ok := s.Bind("test1", nil).(*KeyNotExistError)
	if !ok {
		t.Errorf("it is not KeyNotExistError error")
		t.FailNow()
	}

	var nil_pointer *string
	invalidTypes := []interface{}{nil, k1, nil_pointer}
	for k, v := range invalidTypes {
		_, ok := s.Bind(k1, v).(*InvalidBindError)
		if !ok {
			t.Errorf("round %d, it is not InvalidBindError error", k)
			t.FailNow()
		}
	}

	var invalidBindValue int
	_, ok = s.Bind(k1, &invalidBindValue).(*BindTypeNotMatchError)
	if !ok {
		t.Errorf("it is not BindTypeNotMatchError error")
		t.FailNow()
	}

	var bindValue string
	e := s.Bind(k1, &bindValue)
	if e != nil {
		t.Errorf("the error must be nil, now is: %#v", e)
		t.FailNow()
	} else if !reflect.DeepEqual(v1, bindValue) {
		t.Errorf("left is: %s, right is: %s", v1, bindValue)
	}
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if e := recover(); e == nil {
			t.Errorf("this code did not panic")
		}
	}()
	f()
}
