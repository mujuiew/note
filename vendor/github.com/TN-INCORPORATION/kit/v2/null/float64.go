package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

// Float64 that supports a null state.  Json serialization, SQL database
// storage and stringer interfaces are support
type Float64 struct {
	// Val is the internal float64 value
	Val float64
	// valid is true if Float64 is not NULL
	valid bool
}

// NewFloat64 creates a new Float64 from a standard library float64
func NewFloat64(f float64) Float64 {
	return Float64{valid: true, Val: f}
}

// NewFloat64s creates a new Float64 from a standard library string, returning an
// error if the parameter value s cannot be converted correctly
func NewFloat64s(s string) (Float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return Float64{valid: false, Val: 0.0}, err
	}
	return Float64{valid: true, Val: f}, nil
}

// Set is a convenience method to set the value and set the valid flag to true
func (f *Float64) Set(fl float64) {
	f.Val = fl
	f.valid = true
}

// SetNull sets valid to false and Val to the zero value
func (f *Float64) SetNull() {
	f.Val = 0.0
	f.valid = false
}

// Null implements the Nuller interface and returns the condition of is null
func (f Float64) Null() bool {
	return !f.valid
}

// NotNull is true if valid
func (f Float64) NotNull() bool {
	return f.valid
}

// Zero implements the Zeroer interface and returns the condition of is zero
func (f Float64) Zero() bool {
	return f.Val == 0.0 && f.valid == true
}

// NonZero is true if valid and not the zero value
func (f Float64) NonZero() bool {
	return f.Val != 0.0 && f.valid == true
}

// NullOrZero is true if not valid and is the zero value
func (f Float64) NullOrZero() bool {
	return !f.valid || f.Val == 0.0
}

// Equals return true if the nullable is non-null and its value equals the
// parameter passed in
func (f Float64) Equals(fl float64) bool {
	return f.valid && (f.Val == fl)
}

func (f Float64) EqualsI(flI interface{}) (bool, error) {
	switch x := flI.(type) {
	case float64:
		return f.Equals(x), nil
	case Float64:
		return f.Equals(x.Val), nil
	default:
		err := fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

func (f Float64) GT(fl float64) bool {
	return f.valid && (f.Val > fl)
}

func (f Float64) GTI(flI interface{}) (bool, error) {
	switch x := flI.(type) {
	case float64:
		return f.GT(x), nil
	case Float64:
		return f.GT(x.Val), nil
	default:
		err := fmt.Errorf("cannot check greater than interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

func (f Float64) GTE(fl float64) bool {
	return f.GT(fl) || f.Equals(fl)
}

func (f Float64) GTEI(flI interface{}) (bool, error) {
	switch x := flI.(type) {
	case float64:
		return f.GTE(x), nil
	case Float64:
		return f.GTE(x.Val), nil
	default:
		err := fmt.Errorf("cannot check greater than or equal interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

func (f Float64) LT(fl float64) bool {
	return f.valid && (f.Val < fl)
}

func (f Float64) LTI(flI interface{}) (bool, error) {
	switch x := flI.(type) {
	case float64:
		return f.LT(x), nil
	case Float64:
		return f.LT(x.Val), nil
	default:
		err := fmt.Errorf("cannot check less than interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

func (f Float64) LTE(fl float64) bool {

	return f.LT(fl) || f.Equals(fl)
}

func (f Float64) LTEI(flI interface{}) (bool, error) {
	switch x := flI.(type) {
	case float64:
		return f.LTE(x), nil
	case Float64:
		return f.LTE(x.Val), nil
	default:
		err := fmt.Errorf("cannot check less than or equal interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

// Merge sets the value of the receiver to the value passed if the value passed
// in is not null
func (f *Float64) Merge(m Float64) Float64 {
	if m.Null() {
		return Float64{}
	}
	f.Set(m.Val)
	return *f
}

// String implements the Stringer interface and returns the internal string
// so you can use a Float64 in a fmt.Println statement for example
func (f Float64) String() string {
	if f.valid == false {
		return "null"
	}
	return strconv.FormatFloat(f.Val, 'f', -1, 64)
}

// Verbose provides a string explaining the value of the Float64
func (f Float64) Verbose() string {
	if f.valid == false {
		return "Nil Float64"
	}
	return fmt.Sprintf("Float64 Value: %s", f.String())
}

// #### This gives you proper json behavior ####

// UnmarshalJSON allows the Float64 to be read in from json
func (f *Float64) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		f.Val = float64(x)
	case string:
		str := string(x)
		if len(str) == 0 {
			f.valid = false
			return nil
		}
		var f64 float64
		f64, err = strconv.ParseFloat(str, 64)
		if err == nil {
			f.Val = f64
		} else {
			err = unmarshalTypeError(data, v, Float64{}, err)
		}

	case map[string]interface{}:
		var n64 sql.NullFloat64
		err = json.Unmarshal(data, &n64)
		if err == nil {
			f.Val = n64.Float64
		}
	case nil:
		f.valid = false
		return nil
	default:
		err = unmarshalTypeError(data, v, Float64{}, nil)
	}
	f.valid = err == nil
	return err
}

// MarshalJSON writes out the Float64 as json.
func (f Float64) MarshalJSON() ([]byte, error) {
	if !f.valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatFloat(f.Val, 'f', -1, 64)), nil
}

// #### This gives you proper writes to a sql database ####

// Scan implements the Scanner interface for SQL retrieval.
func (f *Float64) Scan(value interface{}) error {
	if value == nil {
		f.Val, f.valid = 0.0, false
		return nil
	}
	var nb sql.NullFloat64
	err := nb.Scan(value)
	if err != nil {
		f.Val, f.valid = 0.0, false
		return err
	}
	f.valid = nb.Valid
	f.Val = nb.Float64
	return nil
}

// Value implements the driver Valuer interface for SQL storage.
func (f Float64) Value() (driver.Value, error) {
	if !f.valid {
		return nil, nil
	}
	return f.Val, nil
}
