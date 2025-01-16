package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

// Uint64 that supports a null state.  Json serialization, SQL database
// storage and stringer interfaces are support
type Uint64 struct {
	// Val is the internal uint64 value
	Val uint64
	// valid is true if Uint64 is not NULL
	valid bool
}

// NewUint64 creates a new Uint64 from a standard library string
func NewUint64(u uint64) Uint64 {
	return Uint64{valid: true, Val: u}
}

// NewUint64s creates a new Uint64 from a standard library string, returning an
// error if the parameter value s cannot be converted correctly
func NewUint64s(s string) (Uint64, error) {
	u, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return Uint64{valid: false, Val: 0}, err
	}
	return Uint64{valid: true, Val: uint64(u)}, nil
}

// Set is a convenience method to set the value and set the valid flag to true
func (u *Uint64) Set(ui uint64) {
	u.Val = ui
	u.valid = true
}

// SetNull sets valid to false and Val to the zero value
func (u *Uint64) SetNull() {
	u.Val = 0
	u.valid = false
}

// Null implements the Nuller interface and returns the condition of is null
func (u Uint64) Null() bool {
	return !u.valid
}

// NotNull is true if valid
func (u Uint64) NotNull() bool {
	return u.valid
}

// Zero implements the Zeroer interface and returns the condition of is zero
func (u Uint64) Zero() bool {
	return u.Val == 0 && u.valid == true
}

// NonZero is true if valid and not the zero value
func (u Uint64) NonZero() bool {
	return u.Val != 0 && u.valid == true
}

// NullOrZero is true if not valid and is the zero value
func (u Uint64) NullOrZero() bool {
	return !u.valid || u.Val == 0
}

// Equals return true if the nullable is non-null and its value equals the
// parameter passed in
func (u Uint64) Equals(ui uint64) bool {
	return u.valid && (u.Val == ui)
}

// Merge sets the value of the receiver to the value passed if the value passed
// in is not null
func (u *Uint64) Merge(m Uint64) Uint64 {
	if m.Null() {
		return Uint64{}
	}
	u.Set(m.Val)
	return *u
}

// String implements the Stringer interface and returns the internal string
// so you can use a Uint64 in a fmt.Println statement for example
func (u Uint64) String() string {
	if u.valid == false {
		return "null"
	}
	return fmt.Sprintf("%d", u.Val)
}

// Verbose provides a string explaining the value of the Uint64
func (u Uint64) Verbose() string {
	if u.valid == false {
		return "Nil Uint64"
	}
	return fmt.Sprintf("Uint64 Value: %d", u.Val)
}

// #### This gives you proper json behavior ####

// UnmarshalJSON allows the Uint64 to be read in from json
func (u *Uint64) UnmarshalJSON(data []byte) error {
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
		u64, err = strconv.ParseUint(str, 10, 64)
		if err == nil {
			u.Val = u64
		} else {
			err = unmarshalTypeError(data, v, Uint64{}, err)
		}

	case map[string]interface{}:
		var n64 sql.NullInt64
		err = json.Unmarshal(data, &n64)
		if err == nil {
			u.Val, u.valid = uint64(n64.Int64), n64.Valid
		}
	case nil:
		u.valid = false
		return nil
	default:
		err = unmarshalTypeError(data, v, Uint64{}, nil)
	}
	u.valid = err == nil
	return err
}

// MarshalJSON writes out the Uint64 as json.
func (u Uint64) MarshalJSON() ([]byte, error) {
	if !u.valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%d", u.Val)), nil
}

// #### This gives you proper writes to a sql database ####

// Scan implements the Scanner interface for SQL retrieval.
func (u *Uint64) Scan(value interface{}) error {
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
	u.Val = uint64(nb.Int64)
	return nil
}

// Value implements the driver Valuer interface for SQL storage.
func (u Uint64) Value() (driver.Value, error) {
	if !u.valid {
		return nil, nil
	}
	return int64(u.Val), nil
}

func (u Uint64) EqualsI(in interface{}) (bool, error) {
	switch x := in.(type) {
	case uint64:
		return u.Equals(x), nil
	case Uint64:
		return u.Equals(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check equal interface: variable type invalid %+v,%T", x, x)
	}
}

func (u Uint64) GT(in uint64) bool {
	return u.valid && u.Val > in
}

func (u Uint64) LT(in uint64) bool {
	return u.valid && u.Val < in
}

func (u Uint64) GTE(in uint64) bool {
	return u.GT(in) || u.Equals(in)
}

func (u Uint64) LTE(in uint64) bool {
	return u.LT(in) || u.Equals(in)
}

func (u Uint64) GTI(in interface{}) (bool, error) {
	switch x := in.(type) {
	case uint64:
		return u.GT(x), nil
	case Uint64:
		return u.GT(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check greater than interface: variable type invalid %+v,%T", x, x)
	}
}

func (u Uint64) LTI(in interface{}) (bool, error) {
	switch x := in.(type) {
	case uint64:
		return u.LT(x), nil
	case Uint64:
		return u.LT(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check less than interface: variable type invalid %+v,%T", x, x)
	}
}

func (u Uint64) GTEI(in interface{}) (bool, error) {
	switch x := in.(type) {
	case uint64:
		return u.GTE(x), nil
	case Uint64:
		return u.GTE(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check greater than or equal interface: variable type invalid%+v,%T", x, x)
	}
}

func (u Uint64) LTEI(in interface{}) (bool, error) {
	switch x := in.(type) {
	case uint64:
		return u.LTE(x), nil
	case Uint64:
		return u.LTE(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check less than or equal interface: variable type invalid %+v,%T", x, x)
	}
}
