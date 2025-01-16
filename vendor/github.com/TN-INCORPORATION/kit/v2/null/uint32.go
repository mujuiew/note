package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

// Uint32 that supports a null state.  Json serialization, SQL database
// storage and stringer interfaces are support
// swagger:type int32
type Uint32 struct {
	// Val is the internal uint32 value
	Val uint32
	// valid is true if Uint32 is not NULL
	valid bool
}

// NewUint32 creates a new Uint32 from a standard library string
func NewUint32(u uint32) Uint32 {
	return Uint32{valid: true, Val: u}
}

// NewUint32s creates a new Uint32 from a standard library string, returning an
// error if the parameter value s cannot be converted correctly
func NewUint32s(s string) (Uint32, error) {
	u, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return Uint32{valid: false, Val: 0}, err
	}
	return Uint32{valid: true, Val: uint32(u)}, nil
}

// Set is a convenience method to set the value and set the valid flag to true
func (u *Uint32) Set(ui uint32) {
	u.Val = ui
	u.valid = true
}

// SetNull sets valid to false and Val to the zero value
func (u *Uint32) SetNull() {
	u.Val = 0
	u.valid = false
}

// Null implements the Nuller interface and returns the condition of is null
func (u Uint32) Null() bool {
	return !u.valid
}

// NotNull is true if valid
func (u Uint32) NotNull() bool {
	return u.valid
}

// Zero implements the Zeroer interface and returns the condition of is zero
func (u Uint32) Zero() bool {
	return u.Val == 0 && u.valid == true
}

// NonZero is true if valid and not the zero value
func (u Uint32) NonZero() bool {
	return u.Val != 0 && u.valid == true
}

// NullOrZero is true if not valid and is the zero value
func (u Uint32) NullOrZero() bool {
	return !u.valid || u.Val == 0
}

// Equals return true if the nullable is non-null and its value equals the
// parameter passed in
func (u Uint32) Equals(ui uint32) bool {
	return u.valid && (u.Val == ui)
}

// EqualsI return true if the nullable is non-null and its value equals the
// interface parameter passed in
func (u Uint32) EqualsI(ui interface{}) (bool, error) {
	switch x := ui.(type) {
	case uint32:
		return u.Equals(x), nil
	case Uint32:
		return u.Equals(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check equal interface: variable type invalid %+v", x)
	}
}

// GT return true if the nullable is non-null and its value greater the
// parameter passed in
func (u Uint32) GT(ui uint32) bool {
	return u.valid && (u.Val > ui)
}

// GTI return true if the nullable is non-null and its value greater the
// interface parameter passed in
func (u Uint32) GTI(ui interface{}) (bool, error) {
	switch x := ui.(type) {
	case uint32:
		return u.GT(x), nil
	case Uint32:
		return u.GT(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check greater than interface: variable type invalid %+v", x)
	}
}

// GTE return true if the nullable is non-null and its value greater or equal the
// parameter passed in
func (u Uint32) GTE(ui uint32) bool {
	return u.GT(ui) || u.Equals(ui)
}

// GTEI return true if the nullable is non-null and its value greater than or equal the
// interface parameter passed in
func (u Uint32) GTEI(ui interface{}) (bool, error) {
	switch x := ui.(type) {
	case uint32:
		return u.GTE(x), nil
	case Uint32:
		return u.GTE(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check greater than or equal interface: variable type invalid %+v,%T", x, x)
	}
}

// LT return true if the nullable is non-null and its value less than the
// parameter passed in
func (u Uint32) LT(ui uint32) bool {
	return u.valid && (u.Val < ui)
}

// LTI return true if the nullable is non-null and its value less the
// interface parameter passed in
func (u Uint32) LTI(ui interface{}) (bool, error) {
	switch x := ui.(type) {
	case uint32:
		return u.LT(x), nil
	case Uint32:
		return u.LT(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check less than interface: variable type invalid %+v,%T", x, x)
	}
}

// LTE return true if the nullable is non-null and its value less than or equal the
// parameter passed in
func (u Uint32) LTE(ui uint32) bool {
	return u.LT(ui) || u.Equals(ui)
}

// LTEI return true if the nullable is non-null and its value less than or equal the
// interface parameter passed in
func (u Uint32) LTEI(ui interface{}) (bool, error) {
	switch x := ui.(type) {
	case uint32:
		return u.LTE(x), nil
	case Uint32:
		return u.LTE(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check less than or equal interface: variable type invalid %+v,%T", x, x)
	}
}

// Merge sets the value of the receiver to the value passed if the value passed
// in is not null
func (u *Uint32) Merge(m Uint32) Uint32 {
	if m.Null() {
		return Uint32{}
	}
	u.Set(m.Val)
	return *u
}

// String implements the Stringer interface and returns the internal string
// so you can use a Uint32 in a fmt.Println statement for example
func (u Uint32) String() string {
	if u.valid == false {
		return "null"
	}
	return fmt.Sprintf("%d", u.Val)
}

// Verbose provides a string explaining the value of the Uint32
func (u Uint32) Verbose() string {
	if u.valid == false {
		return "Nil Uint32"
	}
	return fmt.Sprintf("Uint32 Value: %d", u.Val)
}

// #### This gives you proper json behavior ####

// UnmarshalJSON allows the Uint32 to be read in from json
func (u *Uint32) UnmarshalJSON(data []byte) error {
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
		u64, err = strconv.ParseUint(str, 10, 32)
		if err == nil {
			u.Val = uint32(u64)
		} else {
			err = unmarshalTypeError(data, v, Uint32{}, err)
		}

	case map[string]interface{}:
		var n64 sql.NullInt64
		err = json.Unmarshal(data, &n64)
		if err == nil {
			u.Val, u.valid = uint32(n64.Int64), n64.Valid
		}
	case nil:
		u.valid = false
		return nil
	default:
		err = unmarshalTypeError(data, v, Uint32{}, nil)
	}
	u.valid = err == nil
	return err
}

// MarshalJSON writes out the Uint32 as json.
func (u Uint32) MarshalJSON() ([]byte, error) {
	if !u.valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%d", u.Val)), nil
}

// #### This gives you proper writes to a sql database ####

// Scan implements the Scanner interface for SQL retrieval.
func (u *Uint32) Scan(value interface{}) error {
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
	u.Val = uint32(nb.Int64)
	return nil
}

// Value implements the driver Valuer interface for SQL storage.
func (u Uint32) Value() (driver.Value, error) {
	if !u.valid {
		return nil, nil
	}
	return int64(u.Val), nil
}
