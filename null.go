package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
)

type (
	// Var[T] defines a variable with generic field of T type
	Var[T any] struct {
		set   bool // tells if the value was set
		valid bool // tells if the value is NULL
		value T    // the value itself
	}

	// an internal interface that helps recognizing any nullable variables
	nullVar interface {
		isSet() bool
		getVal() any
	}

	// an exported interface that helps recognizing which structs can be filtered
	// recursively by FilterStruct
	Filterable interface {
		__()
	}
)

var (
	nullBytes = []byte("null")
)

// Set sets the value
func (v *Var[T]) Set(value T) {
	v.set = true
	v.valid = true
	v.value = value
}

// Unset unsets the value
func (v *Var[T]) Unset() {
	var def T
	v.set = false
	v.valid = false
	v.value = def
}

// SetNil sets the value to NULL
func (v *Var[T]) SetNil() {
	var def T
	v.set = true
	v.valid = false
	v.value = def
}

// Val returns the value
func (v Var[T]) Val() T {
	return v.value
}

// IsSet returns if the value was set
func (v Var[T]) IsSet() bool {
	return v.set
}

// Valid returns if the value is NULL
func (v Var[T]) Valid() bool {
	return v.valid
}

// MarshalJSON implements the json.Marshaler interface
func (v Var[T]) MarshalJSON() ([]byte, error) {
	if !v.valid || !v.set {
		return nullBytes, nil
	}

	return json.Marshal(v.value)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (v *Var[T]) UnmarshalJSON(data []byte) error {
	var def T
	v.set = true
	v.value = def

	if bytes.Equal(nullBytes, data) {
		v.valid = false
		return nil
	}

	v.valid = true
	return json.Unmarshal(data, &v.value)
}

// Value implements the sql package's driver.Valuer interface
func (v Var[T]) Value() (driver.Value, error) {
	if !v.set || !v.valid {
		return nil, nil
	}

	switch val := any(v.value).(type) {
	case driver.Valuer:
		return val.Value()
	default:
		return v.value, nil
	}
}

// Scan implements the sql.Scanner interface
func (v *Var[T]) Scan(src any) error {
	var def T
	v.set = true

	if src == nil {
		v.valid, v.value = false, def
		return nil
	}

	v.valid = true
	return convertAssign(&v.value, src)
}

// isSet implements the nullVar interface for internal usage
func (v Var[T]) isSet() bool {
	return v.set
}

// getVal implements the nullVar interface for internal usage
func (v Var[T]) getVal() any {
	if !v.set || !v.valid {
		return nil
	}

	return v.value
}
