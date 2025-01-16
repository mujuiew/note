package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/TN-INCORPORATION/kit/v2/decimal"
)

// Dec5 that supports a null state.  Json serialization, SQL database
// storage and stringer interfaces are support
type Dec5 struct {
	// Val is the internal decimal.Dec5 value
	Val decimal.Dec5
	// valid is true id Dec5 is not NULL
	valid bool
}

// NewDec5 creates a new Dec5 from a decimal.Dec5
func NewDec5(d decimal.Dec5) Dec5 {
	return Dec5{valid: true, Val: d}
}

// NewDec5s creates a new Dec5 from a standard library string, returning an
// error if the parameter value s cannot be converted correctly
func NewDec5s(s string) (Dec5, error) {
	d := decimal.Dec5{}
	err := d.UnmarshalJSON([]byte(s))
	if err != nil {
		return Dec5{valid: false, Val: decimal.Dec5{}}, err
	}
	return Dec5{valid: true, Val: d}, nil
}

// Set is a convenience method to set the value and set the valid flag to true
func (d *Dec5) Set(de decimal.Dec5) {
	d.Val = de
	d.valid = true
}

// Null implements the Nuller interface and returns the condition of is null
func (d Dec5) Null() bool {
	return !d.valid
}

// NotNull is true if valid
func (d Dec5) NotNull() bool {
	return d.valid
}

// SetNull sets valid to false and Val to the zero value
func (d *Dec5) SetNull() {
	d.Val = decimal.Dec5{}
	d.valid = false
}

// Zero implements the Zeroer interface and returns the condition of is zero
func (d Dec5) Zero() bool {
	return d.Val.IsZero() && d.valid == true
}

// NonZero is true if valid and not the zero value
func (d Dec5) NonZero() bool {
	return !d.Val.IsZero() && d.valid == true
}

// NullOrZero is true if not valid and is the zero value
func (d Dec5) NullOrZero() bool {
	return !d.valid || d.Val.IsZero()
}

// Equals return true if the nullable is non-null and its value equals the
// parameter passed in
func (d Dec5) Equals(de decimal.Dec5) bool {
	return d.valid && (d.Val.Val == de.Val)
}

// EqualsI return true if the nullable is non-null and its value equals the
// interface parameter passed in
func (d Dec5) EqualsI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case decimal.Dec5:
		return d.Equals(x), nil
	case Dec5:
		return d.Equals(x.Val), nil
	default:
		err := fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

// GT return true if the nullable is non-null and its value greater the
// parameter passed in
func (d Dec5) GT(de decimal.Dec5) bool {
	return d.valid && d.Val.GT(de)
}

// GTI return true if the nullable is non-null and its value greater the
// interface parameter passed in
func (d Dec5) GTI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case decimal.Dec5:
		return d.GT(x), nil
	case Dec5:
		return d.GT(x.Val), nil
	default:
		err := fmt.Errorf("cannot check greater than interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

// GTE return true if the nullable is non-null and its value greater or eqaul the
// parameter passed in
func (d Dec5) GTE(de decimal.Dec5) bool {
	return d.GT(de) || d.Equals(de)
}

// GTEI return true if the nullable is non-null and its value greater than or equal the
// interface parameter passed in
func (d Dec5) GTEI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case decimal.Dec5:
		return d.GTE(x), nil
	case Dec5:
		return d.GTE(x.Val), nil
	default:
		err := fmt.Errorf("cannot check greater than or equal interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

// LT return true if the nullable is non-null and its value less than the
// parameter passed in
func (d Dec5) LT(de decimal.Dec5) bool {
	return d.valid && d.Val.LT(de)
}

// LTI return true if the nullable is non-null and its value less the
// interface parameter passed in
func (d Dec5) LTI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case decimal.Dec5:
		return d.LT(x), nil
	case Dec5:
		return d.LT(x.Val), nil
	default:
		err := fmt.Errorf("cannot check less than interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

// LTE return true if the nullable is non-null and its value less than or equal the
// parameter passed in
func (d Dec5) LTE(de decimal.Dec5) bool {
	return d.LT(de) || d.Equals(de)
}

// LTEI return true if the nullable is non-null and its value less than or equal the
// interface parameter passed in
func (d Dec5) LTEI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case decimal.Dec5:
		return d.LTE(x), nil
	case Dec5:
		return d.LTE(x.Val), nil
	default:
		err := fmt.Errorf("cannot check less than or equal interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

// Merge sets the value of the receiver to the value passed if the value passed
// in is not null
func (d *Dec5) Merge(m Dec5) Dec5 {
	if m.Null() {
		return Dec5{}
	}
	d.Set(m.Val)
	return *d
}

// String implements the Stringer interface and returns the internal string
// so you can use a Dec5 in a fmt.Println statement for example
func (d Dec5) String() string {
	if d.valid == false {
		return "null"
	}
	return d.Val.String()
}

// Verbose provides a string explaining the value of the Dec5
func (d Dec5) Verbose() string {
	if d.valid == false {
		return "Nil Dec5"
	}
	return fmt.Sprintf("Dec5 Value: %s", d.Val.String())
}

// #### This gives you proper json behavior ####

// UnmarshalJSON allows the Dec5 to be read in from json
func (d *Dec5) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		//validate fractional
		_, err = decimal.PartsFromString(string(data), decimal.Dec5Frac)
		if err != nil {
			return unmarshalTypeError(data, v, Dec5{}, err)
		}

		// Unmarshal again, directly to int64, to avoid intermediate float64
		d.Val = decimal.NewDec5f(float64(x))
	case string:
		str := string(x)
		if len(str) == 0 {
			d.valid = false
			return nil
		}
		d.Val, err = decimal.NewDec5s(str)
		if err != nil {
			return unmarshalTypeError(data, v, Dec5{}, err)
		}
	case nil:
		d.Val = decimal.Dec5{}
		d.valid = false
		return nil
	default:
		err = unmarshalTypeError(data, v, Dec5{}, nil)
	}
	d.valid = err == nil
	return err
}

// MarshalJSON writes out the Dec5 as json.
func (d Dec5) MarshalJSON() ([]byte, error) {
	if !d.valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%s", d.Val.String())), nil
}

// #### This gives you proper writes to a sql database ####

// Scan implements the Scanner interface for SQL retrieval.
func (d *Dec5) Scan(value interface{}) error {
	if value == nil {
		d.Val, d.valid = decimal.Dec5{}, false
		return nil
	}
	err := d.Val.Scan(value)
	if err != nil {
		d.Val, d.valid = decimal.Dec5{}, false
		return nil
	}
	d.valid = true
	return nil
}

// Value implements the driver Valuer interface for SQL storage.
func (d Dec5) Value() (driver.Value, error) {
	if !d.valid {
		return nil, nil
	}
	return d.Val.Value()
}
