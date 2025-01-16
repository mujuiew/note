package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/TN-INCORPORATION/kit/v2/decimal"
)

// Dec8 that supports a null state.  Json serialization, SQL database
// storage and stringer interfaces are support
// swagger:type float64
type Dec8 struct {
	// Val is the internal decimal.Dec8 value
	Val decimal.Dec8
	// valid is true id Dec8 is not NULL
	valid bool
}

// NewDec8 creates a new Dec8 from a decimal.Dec8
func NewDec8(d decimal.Dec8) Dec8 {
	return Dec8{valid: true, Val: d}
}

// NewDec8s creates a new Dec8 from a standard library string, returning an
// error if the parameter value s cannot be converted correctly
func NewDec8s(s string) (Dec8, error) {
	d := decimal.Dec8{}
	err := d.UnmarshalJSON([]byte(s))
	if err != nil {
		return Dec8{valid: false, Val: decimal.Dec8{}}, err
	}
	return Dec8{valid: true, Val: d}, nil
}

// Set is a convenience method to set the value and set the valid flag to true
func (d *Dec8) Set(de decimal.Dec8) {
	d.Val = de
	d.valid = true
}

// Null implements the Nuller interface and returns the condition of is null
func (d Dec8) Null() bool {
	return !d.valid
}

// NotNull is true if valid
func (d Dec8) NotNull() bool {
	return d.valid
}

// SetNull sets valid to false and Val to the zero value
func (d *Dec8) SetNull() {
	d.Val = decimal.Dec8{}
	d.valid = false
}

// Zero implements the Zeroer interface and returns the condition of is zero
func (d Dec8) Zero() bool {
	return d.Val.IsZero() && d.valid == true
}

// NonZero is true if valid and not the zero value
func (d Dec8) NonZero() bool {
	return !d.Val.IsZero() && d.valid == true
}

// NullOrZero is true if not valid and is the zero value
func (d Dec8) NullOrZero() bool {
	return !d.valid || d.Val.IsZero()
}

// Equals return true if the nullable is non-null and its value equals the
// parameter passed in
func (d Dec8) Equals(de decimal.Dec8) bool {
	return d.valid && (d.Val.Val == de.Val)
}

// EqualsI return true if the nullable is non-null and its value equals the
// interface parameter passed in
func (d Dec8) EqualsI(dei interface{}) (bool, error) {
	switch x := dei.(type) {
	case decimal.Dec8:
		return d.Equals(x), nil
	case Dec8:
		return d.Equals(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
	}
}

// GT return true if the nullable is non-null and its value greater the
// parameter passed in
func (d Dec8) GT(de decimal.Dec8) bool {
	return d.valid && d.Val.GT(de)
}

// GTI return true if the nullable is non-null and its value greater the
// interface parameter passed in
func (d Dec8) GTI(dei interface{}) (bool, error) {
	switch x := dei.(type) {
	case decimal.Dec8:
		return d.GT(x), nil
	case Dec8:
		return d.GT(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check greater than interface: variable type invalid %+v,%T", x, x)
	}
}

// GTE return true if the nullable is non-null and its value greater or equal the
// parameter passed in
func (d Dec8) GTE(de decimal.Dec8) bool {
	return d.GT(de) || d.Equals(de)
}

// GTEI return true if the nullable is non-null and its value greater than or equal the
// interface parameter passed in
func (d Dec8) GTEI(dei interface{}) (bool, error) {
	switch x := dei.(type) {
	case decimal.Dec8:
		return d.GTE(x), nil
	case Dec8:
		return d.GTE(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check greater than or equal interface: variable type invalid %+v,%T", x, x)
	}
}

// LT return true if the nullable is non-null and its value less than the
// parameter passed in
func (d Dec8) LT(de decimal.Dec8) bool {
	return d.valid && d.Val.LT(de)
}

// LTI return true if the nullable is non-null and its value less the
// interface parameter passed in
func (d Dec8) LTI(dei interface{}) (bool, error) {
	switch x := dei.(type) {
	case decimal.Dec8:
		return d.LT(x), nil
	case Dec8:
		return d.LT(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check less than interface: variable type invalid %+v,%T", x, x)
	}
}

// LTE return true if the nullable is non-null and its value less than or equal the
// parameter passed in
func (d Dec8) LTE(de decimal.Dec8) bool {
	return d.LT(de) || d.Equals(de)
}

// LTEI return true if the nullable is non-null and its value less than or equal the
// interface parameter passed in
func (d Dec8) LTEI(dei interface{}) (bool, error) {
	switch x := dei.(type) {
	case decimal.Dec8:
		return d.LTE(x), nil
	case Dec8:
		return d.LTE(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check less than or equal interface: variable type invalid %+v,%T", x, x)
	}
}

// Merge sets the value of the receiver to the value passed if the value passed
// in is not null
func (d *Dec8) Merge(m Dec8) Dec8 {
	if m.Null() {
		return Dec8{}
	}
	d.Set(m.Val)
	return *d
}

// String implements the Stringer interface and returns the internal string
// so you can use a Dec8 in a fmt.Println statement for example
func (d Dec8) String() string {
	if d.valid == false {
		return "null"
	}
	return d.Val.String()
}

// Verbose provides a string explaining the value of the Dec8
func (d Dec8) Verbose() string {
	if d.valid == false {
		return "Nil Dec8"
	}
	return fmt.Sprintf("Dec8 Value: %s", d.Val.String())
}

// #### This gives you proper json behavior ####

// UnmarshalJSON allows the Dec8 to be read in from json
func (d *Dec8) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		//validate fractional
		_, err = decimal.PartsFromString(string(data), decimal.Dec8Frac)
		if err != nil {
			return unmarshalTypeError(data, v, Dec8{}, err)
		}

		// Unmarshal again, directly to int64, to avoid intermediate float64
		d.Val = decimal.NewDec8f(float64(x))
	case string:
		str := string(x)
		if len(str) == 0 {
			d.valid = false
			return nil
		}
		d.Val, err = decimal.NewDec8s(str)
		if err != nil {
			return unmarshalTypeError(data, v, Dec8{}, err)
		}
	case nil:
		d.Val = decimal.Dec8{}
		d.valid = false
		return nil
	default:
		err = unmarshalTypeError(data, v, Dec8{}, nil)
	}
	d.valid = err == nil
	return err
}

// MarshalJSON writes out the Dec8 as json.
func (d Dec8) MarshalJSON() ([]byte, error) {
	if !d.valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%s", d.Val.String())), nil
}

// #### This gives you proper writes to a sql database ####

// Scan implements the Scanner interface for SQL retrieval.
func (d *Dec8) Scan(value interface{}) error {
	if value == nil {
		d.Val, d.valid = decimal.Dec8{}, false
		return nil
	}
	err := d.Val.Scan(value)
	if err != nil {
		d.Val, d.valid = decimal.Dec8{}, false
		return nil
	}
	d.valid = true
	return nil
}

// Value implements the driver Valuer interface for SQL storage.
func (d Dec8) Value() (driver.Value, error) {
	if !d.valid {
		return nil, nil
	}
	return d.Val.Value()
}
