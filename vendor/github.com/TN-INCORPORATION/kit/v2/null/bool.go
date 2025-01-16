package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

// Bool that supports a null state.  Json serialization, SQL database
// storage and stringer interfaces are support
type Bool struct {
	// Val is the internal string value
	Val bool
	// valid is true if Bool is not NULL
	valid bool
}

// Set is a convenience method to set the value and set the valid flag to true
func (b *Bool) Set(bo bool) {
	b.Val = bo
	b.valid = true
}

// SetNull sets valid to false and Val to the zero value
func (b *Bool) SetNull() {
	b.Val = false
	b.valid = false
}

// NewBool creates a new Bool from a standard library bool
func NewBool(b bool) Bool {
	return Bool{valid: true, Val: b}
}

// NewBools creates a new Bool from a standard library string, returning an
// error if the parameter value s cannot be converted correctly
func NewBools(s string) (Bool, error) {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return Bool{valid: false, Val: false}, err
	}
	return Bool{valid: true, Val: b}, nil
}

// Null implements the Nuller interface and returns the condition of is null
func (b Bool) Null() bool {
	return !b.valid
}

// NotNull is true if valid
func (b Bool) NotNull() bool {
	return b.valid
}

// Zero implements the Zeroer interface and returns the condition of is zero
func (b Bool) Zero() bool {
	return b.Val == false && b.valid == true
}

// NonZero is true if valid and not the zero value
func (b Bool) NonZero() bool {
	return b.Val != false && b.valid == true
}

// NullOrZero is true if not valid and is the zero value
func (b Bool) NullOrZero() bool {
	return !b.valid || b.Val == false
}

// Equals return true if the nullable is non-null and its value equals the
// parameter passed in
func (b Bool) Equals(bo bool) bool {
	return b.valid && (b.Val == bo)
}

// Merge sets the value of the receiver to the value passed if the value passed
// in is not null
func (b *Bool) Merge(m Bool) Bool {
	if m.Null() {
		return Bool{}
	}
	b.Set(m.Val)
	return *b
}

// String implements the stringer interface and returns the internal string
// so you can use a Bool in a fmt.Println statement for example
func (b Bool) String() string {
	if b.valid == false {
		return "null"
	}
	if b.Val {
		return "true"
	}
	return "false"
}

// Verbose provides a string explaining the value of the Bool
func (b Bool) Verbose() string {
	if b.valid == false {
		return "Nil Bool"
	}
	return fmt.Sprintf("Bool Value: %t", b.Val)
}

// #### This gives you proper json behavior ####

// UnmarshalJSON allows the Bool to be read in from json
func (b *Bool) UnmarshalJSON(bytes []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(bytes, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case bool:
		b.Val = x
	case map[string]interface{}:
		err = json.Unmarshal(bytes, &b.Val)
	case nil:
		b.valid = false
		return nil
	default:
		err = unmarshalTypeError(bytes, v, Bool{}, nil)
	}
	b.valid = err == nil
	return err
}

// MarshalJSON writes out the Bool as json.
func (b Bool) MarshalJSON() ([]byte, error) {
	if !b.valid {
		return []byte("null"), nil
	}
	if !b.Val {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

// #### This gives you proper writes to a sql database ####

// Scan implements the scanner interface for SQL retrieval.
func (b *Bool) Scan(value interface{}) error {
	if value == nil {
		b.Val, b.valid = false, false
		return nil
	}
	var nb sql.NullBool
	err := nb.Scan(value)
	if err != nil {
		b.Val, b.valid = false, false
		return err
	}
	b.valid = nb.Valid
	b.Val = nb.Bool
	return nil
}

// Value implements the driver Valuer interface for SQL storage.
func (b Bool) Value() (driver.Value, error) {
	if !b.valid {
		return nil, nil
	}
	return b.Val, nil
}

// EqualsI - return true if the nullable is non-null and its value equals the
// parameter passed in
func (b Bool) EqualsI(v interface{}) (bool, error) {
	switch x := v.(type) {
	case bool:
		return b.Equals(x), nil
	case Bool:
		return b.Equals(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
	}
}

func (b Bool) GT(v bool) bool {
	var result bool
	if v == b.Val { // equal
		result = false

	} else {
		if v { // b=false    v = true    v GT b
			result = false
		} else { // b=true    v = false  v LT b
			result = true
		}

	}

	return b.valid && result
}
func (b Bool) GTI(v interface{}) (bool, error) {
	switch x := v.(type) {
	case bool:
		return b.GT(x), nil
	case Bool:
		return b.GT(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check greater than interface: variable type invalid %+v,%T", x, x)
	}
}

func (b Bool) GTE(v bool) bool {
	return b.GT(v) || b.Equals(v)
}

func (b Bool) GTEI(v interface{}) (bool, error) {
	switch x := v.(type) {
	case bool:
		return b.GTE(x), nil
	case Bool:
		return b.GTE(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check greater than or equal interface: variable type invalid %+v,%T", x, x)
	}
}
func (b Bool) LT(v bool) bool {
	var result bool
	if v == b.Val { // equal
		result = false

	} else {
		if v { // b=false    v = true    v GT b
			result = true
		} else { // b=true    v = false  v LT b
			result = false
		}

	}

	return b.valid && result
}
func (b Bool) LTI(v interface{}) (bool, error) {
	switch x := v.(type) {
	case bool:
		return b.LT(x), nil
	case Bool:
		return b.LT(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check less than interface: variable type invalid %+v,%T", x, x)
	}
}
func (b Bool) LTE(v bool) bool {
	return b.LT(v) || b.Equals(v)
}
func (b Bool) LTEI(v interface{}) (bool, error) {
	switch x := v.(type) {
	case bool:
		return b.LTE(x), nil
	case Bool:
		return b.LTE(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check less than or equal interface: variable type invalid %+v,%T", x, x)
	}
}
