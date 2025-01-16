package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

// Int16 that supports a null state.  Json serialization, SQL database
// storage and stringer interfaces are support
// swagger:type int16
type Int16 struct {
	// Val is the internal int16 value
	Val int16
	// valid is true if Int16 is not NULL
	valid bool
}

// NewInt16 creates a new Int16 from a standard library string
func NewInt16(i int16) Int16 {
	return Int16{valid: true, Val: i}
}

// NewInt16s creates a new Int16 from a standard library string, returning an
// error if the parameter value s cannot be converted correctly
func NewInt16s(s string) (Int16, error) {
	i, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return Int16{valid: false, Val: 0}, err
	}
	return Int16{valid: true, Val: int16(i)}, nil
}

// Set is a convenience method to set the value and set the valid flag to true
func (i *Int16) Set(in int16) {
	i.Val = in
	i.valid = true
}

// SetNull sets valid to false and Val to the zero value
func (i *Int16) SetNull() {
	i.Val = 0
	i.valid = false
}

// Null implements the Nuller interface and returns the condition of is null
func (i Int16) Null() bool {
	return !i.valid
}

// NotNull is true if valid
func (i Int16) NotNull() bool {
	return i.valid
}

// Zero implements the Zeroer interface and returns the condition of is zero
func (i Int16) Zero() bool {
	return i.Val == 0 && i.valid == true
}

// NonZero is true if valid and not the zero value
func (i Int16) NonZero() bool {
	return i.Val != 0 && i.valid == true
}

// NullOrZero is true if not valid and is the zero value
func (i Int16) NullOrZero() bool {
	return !i.valid || i.Val == 0
}

// Equals return true if the nullable is non-null and its value equals the
// parameter passed in
func (i Int16) Equals(in int16) bool {
	return i.valid && (i.Val == in)
}

// EqualsI return true if the nullable is non-null and its value equals the
// interface parameter passed in
func (i Int16) EqualsI(in interface{}) (bool, error) {
	switch x := in.(type) {
	case int16:
		return i.Equals(x), nil
	case Int16:
		return i.Equals(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check equal interface: variable type invalid %+v,%T", x, x)
	}
}

// GT return true if the nullable is non-null and its value greater the
// parameter passed in
func (i Int16) GT(in int16) bool {
	return i.valid && (i.Val > in)
}

// GTI return true if the nullable is non-null and its value greater the
// interface parameter passed in
func (i Int16) GTI(in interface{}) (bool, error) {
	switch x := in.(type) {
	case int16:
		return i.GT(x), nil
	case Int16:
		return i.GT(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check greater than interface: variable type invalid %+v,%T", x, x)
	}
}

// GTE return true if the nullable is non-null and its value greater or equal the
// parameter passed in
func (i Int16) GTE(in int16) bool {
	return i.GT(in) || i.Equals(in)
}

// GTEI return true if the nullable is non-null and its value greater than or equal the
// interface parameter passed in
func (i Int16) GTEI(in interface{}) (bool, error) {
	switch x := in.(type) {
	case int16:
		return i.GTE(x), nil
	case Int16:
		return i.GTE(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check greater than or equal interface: variable type invalid %+v,%T", x, x)
	}
}

// LT return true if the nullable is non-null and its value less than the
// parameter passed in
func (i Int16) LT(in int16) bool {
	return i.valid && (i.Val < in)
}

// LTI return true if the nullable is non-null and its value less the
// interface parameter passed in
func (i Int16) LTI(in interface{}) (bool, error) {
	switch x := in.(type) {
	case int16:
		return i.LT(x), nil
	case Int16:
		return i.LT(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check less than interface: variable type invalid %+v,%T", x, x)
	}
}

// LTE return true if the nullable is non-null and its value less than or equal the
// parameter passed in
func (i Int16) LTE(in int16) bool {
	return i.LT(in) || i.Equals(in)
}

// LTEI return true if the nullable is non-null and its value less than or equal the
// interface parameter passed in
func (i Int16) LTEI(in interface{}) (bool, error) {
	switch x := in.(type) {
	case int16:
		return i.LTE(x), nil
	case Int16:
		return i.LTE(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check less than or equal interface: variable type invalid %+v,%T", x, x)
	}
}

// Merge sets the value of the receiver to the value passed if the value passed
// in is not null
func (i *Int16) Merge(m Int16) Int16 {
	if m.Null() {
		return Int16{}
	}
	i.Set(m.Val)
	return *i
}

// String implements the Stringer interface and returns the internal string
// so you can use a Int16 in a fmt.Println statement for example
func (i Int16) String() string {
	if i.valid == false {
		return "null"
	}
	return fmt.Sprintf("%d", i.Val)
}

// Verbose provides a string explaining the value of the Int16
func (i Int16) Verbose() string {
	if i.valid == false {
		return "Nil Int16"
	}
	return fmt.Sprintf("Int16 Value: %d", i.Val)
}

// #### This gives you proper json behavior ####

// UnmarshalJSON allows the Int16 to be read in from json
func (i *Int16) UnmarshalJSON(data []byte) error {
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
		i64, err = strconv.ParseInt(str, 10, 16)
		if err == nil {
			i.Val = int16(i64)
		} else {
			err = unmarshalTypeError(data, v, Int16{}, err)
		}

	case map[string]interface{}:
		var n64 sql.NullInt64
		err = json.Unmarshal(data, &n64)
		if err == nil {
			i.Val, i.valid = int16(n64.Int64), n64.Valid
		}
	case nil:
		i.valid = false
		return nil
	default:
		err = unmarshalTypeError(data, v, Int16{}, nil)
	}
	i.valid = err == nil
	return err

}

// MarshalJSON writes out the Int16 as json.
func (i Int16) MarshalJSON() ([]byte, error) {
	if !i.valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%d", i.Val)), nil
}

// #### This gives you proper writes to a sql database ####

// Scan implements the Scanner interface for SQL retrieval.
func (i *Int16) Scan(value interface{}) error {
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
	i.Val = int16(nb.Int64)
	return nil
}

// Value implements the driver Valuer interface for SQL storage.
func (i Int16) Value() (driver.Value, error) {
	if !i.valid {
		return nil, nil
	}
	return int64(i.Val), nil
}
