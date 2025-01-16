package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/TN-INCORPORATION/kit/v2/decimal"
)

// Dec2 that supports a null state.  Json serialization, SQL database
// storage and stringer interfaces are support
type Dec2 struct {
	// Val is the internal decimal.Dec2 value
	Val decimal.Dec2
	// valid is true id Dec2 is not NULL
	valid bool
}

// NewDec2 creates a new Dec2 from a decimal.Dec2
func NewDec2(d decimal.Dec2) Dec2 {
	return Dec2{valid: true, Val: d}
}

// NewDec2s creates a new Dec2 from a standard library string, returning an
// error if the parameter value s cannot be converted correctly
func NewDec2s(s string) (Dec2, error) {
	d := decimal.Dec2{}
	err := d.UnmarshalJSON([]byte(s))
	if err != nil {
		return Dec2{valid: false, Val: decimal.Dec2{}}, err
	}
	return Dec2{valid: true, Val: d}, nil
}

// Set is a convenience method to set the value and set the valid flag to true
func (d *Dec2) Set(de decimal.Dec2) {
	d.Val = de
	d.valid = true
}

// Null implements the Nuller interface and returns the condition of is null
func (d Dec2) Null() bool {
	return !d.valid
}

// NotNull is true if valid
func (d Dec2) NotNull() bool {
	return d.valid
}

// SetNull sets valid to false and Val to the zero value
func (d *Dec2) SetNull() {
	d.Val = decimal.Dec2{}
	d.valid = false
}

// Zero implements the Zeroer interface and returns the condition of is zero
func (d Dec2) Zero() bool {
	return d.Val.IsZero() && d.valid == true
}

// NonZero is true if valid and not the zero value
func (d Dec2) NonZero() bool {
	return !d.Val.IsZero() && d.valid == true
}

// NullOrZero is true if not valid and is the zero value
func (d Dec2) NullOrZero() bool {
	return !d.valid || d.Val.IsZero()
}

// Equals return true if the nullable is non-null and its value equals the
// parameter passed in
func (d Dec2) Equals(de decimal.Dec2) bool {
	return d.valid && (d.Val.Val == de.Val)
}

// Merge sets the value of the receiver to the value passed if the value passed
// in is not null
func (d *Dec2) Merge(m Dec2) Dec2 {
	if m.Null() {
		return Dec2{}
	}
	d.Set(m.Val)
	return *d
}

// String implements the Stringer interface and returns the internal string
// so you can use a Dec2 in a fmt.Println statement for example
func (d Dec2) String() string {
	if d.valid == false {
		return "null"
	}
	return d.Val.String()
}

// Verbose provides a string explaining the value of the Dec2
func (d Dec2) Verbose() string {
	if d.valid == false {
		return "Nil Dec2"
	}
	return fmt.Sprintf("Dec2 Value: %s", d.Val.String())
}

// #### This gives you proper json behavior ####

// UnmarshalJSON allows the Dec2 to be read in from json
func (d *Dec2) UnmarshalJSON(data []byte) error {

	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		//validate fractional
		_, err = decimal.PartsFromString(string(data), decimal.Dec2Frac)
		if err != nil {
			return unmarshalTypeError(data, v, Dec2{}, err)
		}

		// Unmarshal again, directly to int64, to avoid intermediate float64
		d.Val = decimal.NewDec2f(float64(x))
	case string:
		str := string(x)
		if len(str) == 0 {
			d.valid = false
			return nil
		}
		d.Val, err = decimal.NewDec2s(str)
		if err != nil {
			return unmarshalTypeError(data, v, Dec2{}, err)
		}
	case nil:
		d.Val = decimal.Dec2{}
		d.valid = false
		return nil
	default:
		err = unmarshalTypeError(data, v, Dec2{}, nil)
	}
	d.valid = err == nil
	return err
}

// MarshalJSON writes out the Dec2 as json.
func (d Dec2) MarshalJSON() ([]byte, error) {
	if !d.valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%s", d.Val.String())), nil
}

// #### This gives you proper writes to a sql database ####

// Scan implements the Scanner interface for SQL retrieval.
func (d *Dec2) Scan(value interface{}) error {
	if value == nil {
		d.Val, d.valid = decimal.Dec2{}, false
		return nil
	}
	err := d.Val.Scan(value)
	if err != nil {
		d.Val, d.valid = decimal.Dec2{}, false
		return nil
	}
	d.valid = true
	return nil
}

// Value implements the driver Valuer interface for SQL storage.
func (d Dec2) Value() (driver.Value, error) {
	if !d.valid {
		return nil, nil
	}
	return d.Val.Value()
}

func (d Dec2) EqualsI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case decimal.Dec2:
		return d.Equals(x), nil
	case Dec2:
		return d.Equals(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
	}
}

func (d Dec2) GT(de decimal.Dec2) bool {
	return d.valid && d.Val.GT(de)
}

func (d Dec2) LT(de decimal.Dec2) bool {
	return d.valid && d.Val.LT(de)
}

func (d Dec2) GTE(de decimal.Dec2) bool {
	return d.GT(de) || d.Equals(de)
}

func (d Dec2) LTE(de decimal.Dec2) bool {
	return d.LT(de) || d.Equals(de)
}

func (d Dec2) GTI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case decimal.Dec2:
		return d.GT(x), nil
	case Dec2:
		return d.GT(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check greater than interface: variable type invalid %+v,%T", x, x)
	}
}

func (d Dec2) LTI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case decimal.Dec2:
		return d.LT(x), nil
	case Dec2:
		return d.LT(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check less than interface: variable type invalid %+v,%T", x, x)
	}
}

func (d Dec2) GTEI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case decimal.Dec2:
		return d.GTE(x), nil
	case Dec2:
		return d.GTE(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check greater than or equal interface: variable type invalid %+v,%T", x, x)
	}
}

func (d Dec2) LTEI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case decimal.Dec2:
		return d.LTE(x), nil
	case Dec2:
		return d.LTE(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check less than or equal interface: variable type invalid %+v,%T", x, x)
	}
}
