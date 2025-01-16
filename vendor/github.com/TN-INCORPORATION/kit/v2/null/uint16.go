package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

// Uint16 that supports a null state.  Json serialization, SQL database
// storage and stringer interfaces are support
// swagger:type int16
type Uint16 struct {
	// Val is the internal uint16 value
	Val uint16
	// valid is true if Uint16 is not NULL
	valid bool
}

// NewUint16 creates a new Uint16 from a standard library string, returning an
// error if the parameter value s cannot be converted correctly
func NewUint16(u uint16) Uint16 {
	return Uint16{valid: true, Val: u}
}

// NewUint16s creates a new Uint16 from a standard library string
func NewUint16s(s string) (Uint16, error) {
	u, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return Uint16{valid: false, Val: 0}, err
	}
	return Uint16{valid: true, Val: uint16(u)}, nil
}

// Set is a convenience method to set the value and set the valid flag to true
func (u *Uint16) Set(ui uint16) {
	u.Val = ui
	u.valid = true
}

// SetNull sets valid to false and Val to the zero value
func (u *Uint16) SetNull() {
	u.Val = 0
	u.valid = false
}

// Null implements the Nuller interface and returns the condition of is null
func (u Uint16) Null() bool {
	return !u.valid
}

// NotNull is true if valid
func (u Uint16) NotNull() bool {
	return u.valid
}

// Zero implements the Zeroer interface and returns the condition of is zero
func (u Uint16) Zero() bool {
	return u.Val == 0 && u.valid == true
}

// NonZero is true if valid and not the zero value
func (u Uint16) NonZero() bool {
	return u.Val != 0 && u.valid == true
}

// NullOrZero is true if not valid and is the zero value
func (u Uint16) NullOrZero() bool {
	return !u.valid || u.Val == 0
}

// Equals return true if the nullable is non-null and its value equals the
// parameter passed in
func (u Uint16) Equals(ui uint16) bool {
	return u.valid && (u.Val == ui)
}

// Merge sets the value of the receiver to the value passed if the value passed
// in is not null
func (u *Uint16) Merge(m Uint16) Uint16 {
	if m.Null() {
		return Uint16{}
	}
	u.Set(m.Val)
	return *u
}

// String implements the Stringer interface and returns the internal string
// so you can use a Uint16 in a fmt.Println statement for example
func (u Uint16) String() string {
	if u.valid == false {
		return "null"
	}
	return fmt.Sprintf("%d", u.Val)
}

// Verbose provides a string explaining the value of the Uint16
func (u Uint16) Verbose() string {
	if u.valid == false {
		return "Nil Uint16"
	}
	return fmt.Sprintf("Uint16 Value: %d", u.Val)
}

// #### This gives you proper json behavior ####

// UnmarshalJSON allows the Uint16 to be read in from json
func (u *Uint16) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		// Unmarshal again, directly to int64, to avoid intermediate float64
		err = json.Unmarshal(data, &u.Val)
	case string:
		str := string(x)
		if len(str) == 0 {
			u.valid = false
			return nil
		}
		var u64 uint64
		u64, err = strconv.ParseUint(str, 10, 16)
		if err == nil {
			u.Val = uint16(u64)
		} else {
			err = unmarshalTypeError(data, v, Uint16{}, err)
		}

	case map[string]interface{}:
		var n64 sql.NullInt64
		err = json.Unmarshal(data, &n64)
		if err == nil {
			u.Val, u.valid = uint16(n64.Int64), n64.Valid
		}
	case nil:
		u.valid = false
		return nil
	default:
		err = unmarshalTypeError(data, v, Uint16{}, nil)
	}
	u.valid = err == nil
	return err
}

// MarshalJSON writes out the Uint16 as json.
func (u Uint16) MarshalJSON() ([]byte, error) {
	if !u.valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%d", u.Val)), nil
}

// #### This gives you proper writes to a sql database ####

// Scan implements the Scanner interface for SQL retrieval.
func (u *Uint16) Scan(value interface{}) error {
	if value == nil {
		u.Val, u.valid = 0, false
		return nil
	}
	var nb sql.NullInt64
	err := nb.Scan(value)
	if err != nil {
		u.Val, u.valid = 0, false
		return err
	}
	u.valid = nb.Valid
	u.Val = uint16(nb.Int64)
	return nil
}

// Value implements the driver Valuer interface for SQL storage.
func (u Uint16) Value() (driver.Value, error) {
	if !u.valid {
		return nil, nil
	}
	return int64(u.Val), nil
}
func (u Uint16) EqualsI(v interface{}) (bool, error) {
	switch x := v.(type) {
	case uint16:
		return u.Equals(x), nil
	case Uint16:
		return u.Equals(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
	}
}
func (u Uint16) GT(v uint16) bool {
	return u.valid && (u.Val > v)
}
func (u Uint16) GTI(v interface{}) (bool, error) {
	switch x := v.(type) {
	case uint16:
		return u.GT(x), nil
	case Uint16:
		return u.GT(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check greater than interface: variable type invalid %+v,%T", x, x)
	}
}

func (u Uint16) GTE(v uint16) bool {
	return u.GT(v) || u.Equals(v)
}
func (u Uint16) GTEI(v interface{}) (bool, error) {
	switch x := v.(type) {
	case uint16:
		return u.GTE(x), nil
	case Uint16:
		return u.GTE(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check greater than or equal interface: variable type invalid %+v,%T", x, x)
	}
}
func (u Uint16) LT(v uint16) bool {
	return u.valid && (u.Val < v)
}
func (u Uint16) LTI(v interface{}) (bool, error) {
	switch x := v.(type) {
	case uint16:
		return u.LT(x), nil
	case Uint16:
		return u.LT(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check less than interface: variable type invalid %+v,%T", x, x)
	}
}
func (u Uint16) LTE(v uint16) bool {
	return u.LT(v) || u.Equals(v)
}
func (u Uint16) LTEI(v interface{}) (bool, error) {
	switch x := v.(type) {
	case uint16:
		return u.LTE(x), nil
	case Uint16:
		return u.LTE(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check less than or equal interface: variable type invalid %+v,%T", x, x)
	}
}
