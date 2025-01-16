package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

// Int64 that supports a null state.  Json serialization, SQL database
// storage and stringer interfaces are support
type Int64 struct {
	// Val is the internal int64 value
	Val int64
	// valid is true if Int64 is not NULL
	valid bool
}

// NewInt64 creates a new Int64 from a standard library string
func NewInt64(i int64) Int64 {
	return Int64{valid: true, Val: i}
}

// NewInt64s creates a new Int64 from a standard library string, returning an
// error if the parameter value s cannot be converted correctly
func NewInt64s(s string) (Int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return Int64{valid: false, Val: 0}, err
	}
	return Int64{valid: true, Val: int64(i)}, nil
}

// Set is a convenience method to set the value and set the valid flag to true
func (i *Int64) Set(in int64) {
	i.Val = in
	i.valid = true
}

// SetNull sets valid to false and Val to the zero value
func (i *Int64) SetNull() {
	i.Val = 0
	i.valid = false
}

// Null implements the Nuller interface and returns the condition of is null
func (i Int64) Null() bool {
	return !i.valid
}

// NotNull is true if valid
func (i Int64) NotNull() bool {
	return i.valid
}

// Zero implements the Zeroer interface and returns the condition of is zero
func (i Int64) Zero() bool {
	return i.Val == 0 && i.valid == true
}

// NonZero is true if valid and not the zero value
func (i Int64) NonZero() bool {
	return i.Val != 0 && i.valid == true
}

// NullOrZero is true if not valid and is the zero value
func (i Int64) NullOrZero() bool {
	return !i.valid || i.Val == 0
}

// Equals return true if the nullable is non-null and its value equals the
// parameter passed in
func (i Int64) Equals(in int64) bool {
	return i.valid && (i.Val == in)
}

// Merge sets the value of the receiver to the value passed if the value passed
// in is not null
func (i *Int64) Merge(m Int64) Int64 {
	if m.Null() {
		return Int64{}
	}
	i.Set(m.Val)
	return *i
}

// String implements the Stringer interface and returns the internal string
// so you can use a Int64 in a fmt.Println statement for example
func (i Int64) String() string {
	if i.valid == false {
		return "null"
	}
	return fmt.Sprintf("%d", i.Val)
}

// Verbose provides a string explaining the value of the Int64
func (i Int64) Verbose() string {
	if i.valid == false {
		return "Nil Int64"
	}
	return fmt.Sprintf("Int64 Value: %d", i.Val)
}

// #### This gives you proper json behavior ####

// UnmarshalJSON allows the Int64 to be read in from json
func (i *Int64) UnmarshalJSON(data []byte) error {
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
		i64, err = strconv.ParseInt(str, 10, 64)
		if err == nil {
			i.Val = i64
		} else {
			err = unmarshalTypeError(data, v, Int64{}, nil)
		}

	case map[string]interface{}:
		var n64 sql.NullInt64
		err = json.Unmarshal(data, &n64)
		if err == nil {
			i.Val, i.valid = n64.Int64, n64.Valid
		}
	case nil:
		i.valid = false
		return nil
	default:
		err = unmarshalTypeError(data, v, Int64{}, nil)
	}
	i.valid = err == nil
	return err

}

// MarshalJSON writes out the Int64 as json.
func (i Int64) MarshalJSON() ([]byte, error) {
	if !i.valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%d", i.Val)), nil
}

// #### This gives you proper writes to a sql database ####

// Scan implements the Scanner interface for SQL retrieval.
func (i *Int64) Scan(value interface{}) error {
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
	i.Val = int64(nb.Int64)
	return nil
}

// Value implements the driver Valuer interface for SQL storage.
func (i Int64) Value() (driver.Value, error) {
	if !i.valid {
		return nil, nil
	}
	return i.Val, nil
}

func (i Int64) EqualsI(in interface{}) (bool, error) {
	switch x := in.(type) {
	case int64:
		return i.Equals(x), nil
	case Int64:
		return i.Equals(x.Val), nil
	default:
		return false, fmt.Errorf("cannot compare value %+v equals to Go value of type null.Int64", i)
	}
}

func (i Int64) GT(in int64) bool {
	return i.valid && i.Val > in
}

func (i Int64) LT(in int64) bool {
	return i.valid && i.Val < in
}

func (i Int64) GTE(in int64) bool {
	return i.GT(in) || i.Equals(in)
}

func (i Int64) LTE(in int64) bool {
	return i.LT(in) || i.Equals(in)
}

func (i Int64) GTI(in interface{}) (bool, error) {
	switch x := in.(type) {
	case int64:
		return i.GT(x), nil
	case Int64:
		return i.GT(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check greater than interface: variable type invalid %+v,%T", in, in)
	}
}

func (i Int64) LTI(in interface{}) (bool, error) {
	switch x := in.(type) {
	case int64:
		return i.LT(x), nil
	case Int64:
		return i.LT(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check less than interface: variable type invalid %+v,%T", x, x)
	}
}

func (i Int64) GTEI(in interface{}) (bool, error) {
	switch x := in.(type) {
	case int64:
		return i.GTE(x), nil
	case Int64:
		return i.GTE(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check greater than or equal interface: variable type invalid %+v,%T", x, x)
	}
}

func (i Int64) LTEI(in interface{}) (bool, error) {
	switch x := in.(type) {
	case int64:
		return i.LTE(x), nil
	case Int64:
		return i.LTE(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check less than or equal interface: variable type invalid %+v,%T", x, x)
	}
}
