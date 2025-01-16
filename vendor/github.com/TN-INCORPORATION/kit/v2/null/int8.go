package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

// Int8 that supports a null state.  Json serialization, SQL database
// storage and stringer interfaces are support
// swagger:type int8
type Int8 struct {
	// Val is the internal int8 value
	Val int8
	// valid is true if Int8 is not NULL
	valid bool
}

// NewInt8 creates a new Int8 from a standard library int8
func NewInt8(i int8) Int8 {
	return Int8{valid: true, Val: i}
}

// NewInt8s creates a new Int8 from a standard library string, returning an
// error if the parameter value s cannot be converted correctly
func NewInt8s(s string) (Int8, error) {
	i, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return Int8{valid: false, Val: 0}, err
	}
	return Int8{valid: true, Val: int8(i)}, nil
}

// Set is a convenience method to set the value and set the valid flag to true
func (i *Int8) Set(in int8) {
	i.Val = in
	i.valid = true
}

// SetNull sets valid to false and Val to the zero value
func (i *Int8) SetNull() {
	i.Val = 0
	i.valid = false
}

// Null implements the Nuller interface and returns the condition of is null
func (i Int8) Null() bool {
	return !i.valid
}

// NotNull is true if valid
func (i Int8) NotNull() bool {
	return i.valid
}

// Zero implements the Zeroer interface and returns the condition of is zero
func (i Int8) Zero() bool {
	return i.Val == 0 && i.valid == true
}

// NonZero is true if valid and not the zero value
func (i Int8) NonZero() bool {
	return i.Val != 0 && i.valid == true
}

// NullOrZero is true if not valid and is the zero value
func (i Int8) NullOrZero() bool {
	return !i.valid || i.Val == 0
}

// Equals return true if the nullable is non-null and its value equals the
// parameter passed in
func (i Int8) Equals(in int8) bool {
	return i.valid && (i.Val == in)
}

// Merge sets the value of the receiver to the value passed if the value passed
// in is not null
func (i *Int8) Merge(m Int8) Int8 {
	if m.Null() {
		return Int8{}
	}
	i.Set(m.Val)
	return *i
}

// String implements the Stringer interface and returns the internal string
// so you can use a Int8 in a fmt.Println statement for example
func (i Int8) String() string {
	if i.valid == false {
		return "null"
	}
	return fmt.Sprintf("%d", i.Val)
}

// Verbose provides a string explaining the value of the Int8
func (i Int8) Verbose() string {
	if i.valid == false {
		return "Nil Int8"
	}
	return fmt.Sprintf("Int8 Value: %d", i.Val)
}

// #### This gives you proper json behavior ####

// UnmarshalJSON allows the Int8 to be read in from json
func (i *Int8) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		// Unmarshal again, directly to int64, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Val)
	case string:
		str := string(x)
		if len(str) == 0 {
			i.valid = false
			return nil
		}
		var i64 int64
		i64, err = strconv.ParseInt(str, 10, 8)
		if err == nil {
			i.Val = int8(i64)
		} else {
			err = unmarshalTypeError(data, v, Int8{}, err)
		}

	case map[string]interface{}:
		var n64 sql.NullInt64
		err = json.Unmarshal(data, &n64)
		if err == nil {
			i.Val, i.valid = int8(n64.Int64), n64.Valid
		}
	case nil:
		i.valid = false
		return nil
	default:
		err = unmarshalTypeError(data, v, Int8{}, nil)
	}
	i.valid = err == nil
	return err
}

// MarshalJSON writes out the Int8 as json.
func (i Int8) MarshalJSON() ([]byte, error) {
	if !i.valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%d", i.Val)), nil
}

// #### This gives you proper writes to a sql database ####

// Scan implements the Scanner interface for SQL retrieval.
func (i *Int8) Scan(value interface{}) error {
	if value == nil {
		i.Val, i.valid = 0, false
		return nil
	}
	var nb sql.NullInt64
	err := nb.Scan(value)
	if err != nil {
		i.Val, i.valid = 0, false
		return err
	}
	i.valid = nb.Valid
	i.Val = int8(nb.Int64)
	return nil
}

// Value implements the driver Valuer interface for SQL storage.
func (i Int8) Value() (driver.Value, error) {
	if !i.valid {
		return nil, nil
	}
	return int64(i.Val), nil
}
func (i Int8) EqualsI(v interface{}) (bool, error) {
	switch x := v.(type) {
	case int8:
		return i.Equals(x), nil
	case Int8:
		return i.Equals(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
	}
}
func (i Int8) GT(v int8) bool {
	return i.valid && (i.Val > v)
}
func (i Int8) GTI(v interface{}) (bool, error) {
	switch x := v.(type) {
	case int8:
		return i.GT(x), nil
	case Int8:
		return i.GT(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check greater than interface: variable type invalid %+v,%T", x, x)
	}
}

func (i Int8) GTE(v int8) bool {
	return i.GT(v) || i.Equals(v)
}

func (i Int8) GTEI(v interface{}) (bool, error) {
	switch x := v.(type) {
	case int8:
		return i.GTE(x), nil
	case Int8:
		return i.GTE(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check greater than or equal interface: variable type invalid %+v,%T", x, x)
	}
}
func (i Int8) LT(v int8) bool {
	return i.valid && (i.Val < v)
}
func (i Int8) LTI(v interface{}) (bool, error) {
	switch x := v.(type) {
	case int8:
		return i.LT(x), nil
	case Int8:
		return i.LT(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check less than interface: variable type invalid %+v,%T", x, x)
	}
}
func (i Int8) LTE(v int8) bool {
	return i.LT(v) || i.Equals(v)
}
func (i Int8) LTEI(v interface{}) (bool, error) {
	switch x := v.(type) {
	case int8:
		return i.LTE(x), nil
	case Int8:
		return i.LTE(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check less than or equal interface: variable type invalid %+v,%T", x, x)
	}
}
