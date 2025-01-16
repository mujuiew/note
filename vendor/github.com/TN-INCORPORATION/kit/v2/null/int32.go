package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

// Int32 that supports a null state.  Json serialization, SQL database
// storage and stringer interfaces are support
type Int32 struct {
	// Val is the internal int32 value
	Val int32
	// valid is true if Int32 is not NULL
	valid bool
}

// NewInt32 creates a new Int32 from a standard library string
func NewInt32(i int32) Int32 {
	return Int32{valid: true, Val: i}
}

// NewInt32s creates a new Int32 from a standard library string, returning an
// error if the parameter value s cannot be converted correctly
func NewInt32s(s string) (Int32, error) {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return Int32{valid: false, Val: 0}, err
	}
	return Int32{valid: true, Val: int32(i)}, nil
}

// Set is a convenience method to set the value and set the valid flag to true
func (i *Int32) Set(in int32) {
	i.Val = in
	i.valid = true
}

// SetNull sets valid to false and Val to the zero value
func (i *Int32) SetNull() {
	i.Val = 0
	i.valid = false
}

// Null implements the Nuller interface and returns the condition of is null
func (i Int32) Null() bool {
	return !i.valid
}

// NotNull is true if valid
func (i Int32) NotNull() bool {
	return i.valid
}

// Zero implements the Zeroer interface and returns the condition of is zero
func (i Int32) Zero() bool {
	return i.Val == 0 && i.valid == true
}

// NonZero is true if valid and not the zero value
func (i Int32) NonZero() bool {
	return i.Val != 0 && i.valid == true
}

// NullOrZero is true if not valid and is the zero value
func (i Int32) NullOrZero() bool {
	return !i.valid || i.Val == 0
}

// Equals return true if the nullable is non-null and its value equals the
// parameter passed in
func (i Int32) Equals(in int32) bool {
	return i.valid && (i.Val == in)
}

// EqualsI return true if the nullable is non-null and its value equals the
// interface parameter passed in
func (i Int32) EqualsI(in interface{}) (bool, error) {
	switch x := in.(type) {
	case int32:
		return i.Equals(x), nil
	case Int32:
		return i.Equals(x.Val), nil
	default:
		err := fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

// GT return true if the nullable is non-null and its value greater the
// parameter passed in
func (i Int32) GT(in int32) bool {
	return i.valid && (i.Val > in)
}

// GTI return true if the nullable is non-null and its value greater the
// interface parameter passed in
func (i Int32) GTI(in interface{}) (bool, error) {
	switch x := in.(type) {
	case int32:
		return i.GT(x), nil
	case Int32:
		return i.GT(x.Val), nil
	default:
		err := fmt.Errorf("cannot check greater than interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

// GTE return true if the nullable is non-null and its value greater or eqaul the
// parameter passed in
func (i Int32) GTE(in int32) bool {
	return i.GT(in) || i.Equals(in)
}

// GTEI return true if the nullable is non-null and its value greater than or equal the
// interface parameter passed in
func (i Int32) GTEI(in interface{}) (bool, error) {
	switch x := in.(type) {
	case int32:
		return i.GTE(x), nil
	case Int32:
		return i.GTE(x.Val), nil
	default:
		err := fmt.Errorf("cannot check greater than or equal interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

// LT return true if the nullable is non-null and its value less than the
// parameter passed in
func (i Int32) LT(in int32) bool {
	return i.valid && (i.Val < in)
}

// LTI return true if the nullable is non-null and its value less the
// interface parameter passed in
func (i Int32) LTI(in interface{}) (bool, error) {
	switch x := in.(type) {
	case int32:
		return i.LT(x), nil
	case Int32:
		return i.LT(x.Val), nil
	default:
		err := fmt.Errorf("cannot check less than interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

// LTE return true if the nullable is non-null and its value less than or equal the
// parameter passed in
func (i Int32) LTE(in int32) bool {
	return i.LT(in) || i.Equals(in)
}

// LTEI return true if the nullable is non-null and its value less than or equal the
// interface parameter passed in
func (i Int32) LTEI(in interface{}) (bool, error) {
	switch x := in.(type) {
	case int32:
		return i.LTE(x), nil
	case Int32:
		return i.LTE(x.Val), nil
	default:
		err := fmt.Errorf("cannot check less than or equal interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

// Merge sets the value of the receiver to the value passed if the value passed
// in is not null
func (i *Int32) Merge(m Int32) Int32 {
	if m.Null() {
		return Int32{}
	}
	i.Set(m.Val)
	return *i
}

// String implements the Stringer interface and returns the internal string
// so you can use a Int32 in a fmt.Println statement for example
func (i Int32) String() string {
	if i.valid == false {
		return "null"
	}
	return fmt.Sprintf("%d", i.Val)
}

// Verbose provides a string explaining the value of the Int32
func (i Int32) Verbose() string {
	if i.valid == false {
		return "Nil Int32"
	}
	return fmt.Sprintf("Int32 Value: %d", i.Val)
}

// #### This gives you proper json behavior ####

// UnmarshalJSON allows the Int32 to be read in from json
func (i *Int32) UnmarshalJSON(data []byte) error {
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
		i64, err = strconv.ParseInt(str, 10, 32)
		if err == nil {
			i.Val = int32(i64)
		} else {
			err = unmarshalTypeError(data, v, Int32{}, err)
		}

	case map[string]interface{}:
		var n64 sql.NullInt64
		err = json.Unmarshal(data, &n64)
		if err == nil {
			i.Val, i.valid = int32(n64.Int64), n64.Valid
		}
	case nil:
		i.valid = false
		return nil
	default:
		err = unmarshalTypeError(data, v, Int32{}, nil)
	}
	i.valid = err == nil
	return err

}

// MarshalJSON writes out the Int32 as json.
func (i Int32) MarshalJSON() ([]byte, error) {
	if !i.valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%d", i.Val)), nil
}

// #### This gives you proper writes to a sql database ####

// Scan implements the Scanner interface for SQL retrieval.
func (i *Int32) Scan(value interface{}) error {
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
	i.Val = int32(nb.Int64)
	return nil
}

// Value implements the driver Valuer interface for SQL storage.
func (i Int32) Value() (driver.Value, error) {
	if !i.valid {
		return nil, nil
	}
	return int64(i.Val), nil
}
