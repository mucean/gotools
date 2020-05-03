package store

import (
	"fmt"
	"reflect"
)

// A InvalidBindError represent a valid type passed to Bind.
type InvalidBindError struct {
	Val reflect.Value
}

func (e *InvalidBindError) Error() string {
	tpl := "Store: Bind(%s)"
	if e.Val.Kind() == reflect.Invalid {
		return fmt.Sprintf(tpl, "nil")
	}

	if e.Val.Kind() != reflect.Ptr {
		return fmt.Sprintf(tpl, fmt.Sprintf("non-pointer %s", e.Val.Type().String()))
	}
	return fmt.Sprintf(tpl, fmt.Sprintf("%s pointer is nil", e.Val.Type().String()))
}

// A KeyNotExistError represent that the key in the context can not be found.
type KeyNotExistError struct {
	Key interface{}
}

func (e *KeyNotExistError) Error() string {
	return fmt.Sprintf("key %#v is not found", e.Key)
}

// A BindTypeNotMatchError represent that the type of actual value does not correspond
// with the type of expect value
type BindTypeNotMatchError struct {
	Expect reflect.Type
	Actual reflect.Type
}

func (e *BindTypeNotMatchError) Error() string {
	return fmt.Sprintf("expect type: %s, actual type: %s", e.Expect.String(), e.Actual.String())
}

// Store is a map wrapper that can contain any type value indexed by any type key
type Store struct {
	s map[interface{}]interface{}
}

// New return Store struct
func New() Store {
	return Store{
		s: make(map[interface{}]interface{}),
	}
}

// Add return Store struct pointer
// insert a value into Store container
func (s *Store) Add(k, v interface{}) *Store {
	s.s[k] = v
	return s
}

// Append return Store struct pointer
// this method will grew by append a bunch of values, the value has existed
// will be replaced
func (s *Store) Append(vs map[interface{}]interface{}) *Store {
	if len(vs) == 0 {
		return s
	}
	for k, v := range vs {
		s.Add(k, v)
	}
	return s
}

// Append return Store struct pointer
// this method will replace the exist store map
func (s *Store) With(vs map[interface{}]interface{}) *Store {
	if len(vs) == 0 {
		return s
	}
	s.s = vs
	return s
}

// Get return the value correspond to the key
func (s *Store) Get(k interface{}) (v interface{}, ok bool) {
	v, ok = s.s[k]
	return
}

// Underlying return the value underlying the Store
func (s *Store) Underlying() map[interface{}]interface{} {
	return s.s
}

// MustGet return the value correspond to the key, if the key don't exist
// the method will panic
func (s *Store) MustGet(k interface{}) interface{} {
	v, ok := s.Get(k)
	if !ok {
		panic(KeyNotExistError{Key: k})
	}
	return v
}

// Bind return a error when the flowing thing happened
// 1. key is not existed
// 2. passed bind value is not a pointer to it
// 3. actual value type is not corresponded to the expect value type
func (s *Store) Bind(k interface{}, v interface{}) error {
	val, ok := s.Get(k)
	if !ok {
		return &KeyNotExistError{Key: k}
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidBindError{Val: rv}
	}
	rVal := reflect.ValueOf(val)
	if rv.Elem().Kind() != rVal.Kind() {
		return &BindTypeNotMatchError{Expect: rv.Elem().Type(), Actual: rVal.Type()}
	}
	rv.Elem().Set(rVal)
	return nil
}
