package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

// Float32 that supports a null state.  Json serialization, SQL database
// storage and stringer interfaces are support
// swagger:type float32
type Float32 struct {
	// Val is the internal float32 value
	Val float32
	// valid is true if Float32 is not NULL
	valid bool
}

// NewFloat32 creates a new Float32 from a standard library float32
func NewFloat32(f float32) Float32 {
	return Float32{valid: true, Val: f}
}

// NewFloat32s creates a new Float32 from a standard library string, returning an
// error if the parameter value s cannot be converted correctly
func NewFloat32s(s string) (Float32, error) {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return Float32{valid: false, Val: 0.0}, err
	}
	return Float32{valid: true, Val: float32(f)}, nil
}

// Set is a convenience method to set the value and set the valid flag to true
func (f *Float32) Set(fl float32) {
	f.Val = fl
	f.valid = true
}

// SetNull sets valid to false and Val to the zero value
func (f *Float32) SetNull() {
	f.Val = 0.0
	f.valid = false
}

// Null implements the Nuller interface and returns the condition of is null
func (f Float32) Null() bool {
	return !f.valid
}

// NotNull is true if valid
func (f Float32) NotNull() bool {
	return f.valid
}

// Zero implements the Zeroer interface and returns the condition of is zero
func (f Float32) Zero() bool {
	return f.Val == 0.0 && f.valid == true
}

// NonZero is true if valid and not the zero value
func (f Float32) NonZero() bool {
	return f.Val != 0.0 && f.valid == true
}

// NullOrZero is true if not valid and is the zero value
func (f Float32) NullOrZero() bool {
	return !f.valid || f.Val == 0.0
}

// Equals return true if the nullable is non-null and its value equals the
// parameter passed in
func (f Float32) Equals(fl float32) bool {
	return f.valid && (f.Val == fl)
}

// Merge sets the value of the receiver to the value passed if the value passed
// in is not null
func (f *Float32) Merge(m Float32) Float32 {
	if m.Null() {
		return Float32{}
	}
	f.Set(m.Val)
	return *f
}

// String implements the Stringer interface and returns the internal string
// so you can use a Float32 in a fmt.Println statement for example
func (f Float32) String() string {
	if f.valid == false {
		return "null"
	}
	return strconv.FormatFloat(float64(f.Val), 'f', -1, 32)
}

// Verbose provides a string explaining the value of the Float32
func (f Float32) Verbose() string {
	if f.valid == false {
		return "Nil Float32"
	}
	return fmt.Sprintf("Float32 Value: %s", f.String())
}

// #### This gives you proper json behavior ####

// UnmarshalJSON allows the Float32 to be read in from json
func (f *Float32) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		f.Val = float32(x)
	case string:
		str := string(x)
		if len(str) == 0 {
			f.valid = false
			return nil
		}
		var f64 float64
		f64, err = strconv.ParseFloat(str, 32)
		if err == nil {
			f.Val = float32(f64)
		} else {
			err = unmarshalTypeError(data, v, Float32{}, err)
		}

	case map[string]interface{}:
		var n64 sql.NullFloat64
		err = json.Unmarshal(data, &n64)
		if err == nil {
			f.Val = float32(n64.Float64)
		}
	case nil:
		f.valid = false
		return nil
	default:
		err = unmarshalTypeError(data, v, Float32{}, nil)
	}
	f.valid = err == nil
	return err
}

// MarshalJSON writes out the Float32 as json.
func (f Float32) MarshalJSON() ([]byte, error) {
	if !f.valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatFloat(float64(f.Val), 'f', -1, 32)), nil
}

// #### This gives you proper writes to a sql database ####

// Scan implements the Scanner interface for SQL retrieval.
func (f *Float32) Scan(value interface{}) error {
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
	f.Val = float32(nb.Float64)
	return nil
}

// Value implements the driver Valuer interface for SQL storage.
func (f Float32) Value() (driver.Value, error) {
	if !f.valid {
		return nil, nil
	}
	return float64(f.Val), nil
}

func (f Float32) EqualsI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case float32:
		return f.Equals(x), nil
	case Float32:
		return f.Equals(x.Val), nil
	default:
		err := fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

func (f Float32) GT(fl float32) bool {
	return f.valid && f.Val > fl
}

func (f Float32) GTI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case float32:
		return f.GT(x), nil
	case Float32:
		return f.GT(x.Val), nil
	default:
		err := fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

func (f Float32) GTE(fl float32) bool {
	return f.GT(fl) || f.Equals(fl)
}

func (f Float32) GTEI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case float32:
		return f.GTE(x), nil
	case Float32:
		return f.GTE(x.Val), nil
	default:
		err := fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

func (f Float32) LT(fl float32) bool {
	return f.valid && f.Val < fl
}

func (f Float32) LTI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case float32:
		return f.LT(x), nil
	case Float32:
		return f.LT(x.Val), nil
	default:
		err := fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

func (f Float32) LTE(fl float32) bool {
	return f.LT(fl) || f.Equals(fl)
}

func (f Float32) LTEI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case float32:
		return f.LTE(x), nil
	case Float32:
		return f.LTE(x.Val), nil
	default:
		err := fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}
