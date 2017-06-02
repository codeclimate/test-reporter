package formatters

import (
	"encoding/json"
	"strconv"
)

// NullInt adds an implementation for int
// that supports proper JSON encoding/decoding.
type NullInt struct {
	Int   int
	Valid bool // Valid is true if Int is not NULL
}

func (ns NullInt) Interface() interface{} {
	if !ns.Valid {
		return nil
	}
	return ns.Int
}

// NewNullInt returns a new, properly instantiated
// Int object.
func NewNullInt(i int) NullInt {
	return NullInt{Int: i, Valid: true}
}

// MarshalJSON marshals the underlying value to a
// proper JSON representation.
func (ns NullInt) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.Int)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON will unmarshal a JSON value into
// the propert representation of that value.
func (ns *NullInt) UnmarshalJSON(text []byte) error {
	if i, err := strconv.ParseInt(string(text), 10, strconv.IntSize); err == nil {
		ns.Valid = true
		ns.Int = int(i)
	}
	return nil
}

func (ns *NullInt) UnmarshalText(text []byte) error {
	return ns.UnmarshalJSON(text)
}
