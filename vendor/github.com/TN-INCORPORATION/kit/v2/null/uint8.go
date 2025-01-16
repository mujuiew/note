package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

// Uint8 that supports a null state.  Json serialization, SQL database
// storage and stringer interfaces are support
// swagger:type int8
type Uint8 struct {
	// Val is the internal uint8 value
	Val uint8
	// valid is true if Uint8 is not NULL
	valid bool
}

// NewUint8 creates a new Uint8 from a standard library string
func NewUint8(u uint8) Uint8 {
	return Uint8{valid: true, Val: u}
}

// NewUint8s creates a new Uint8 from a standard library string, returning an
// error if the parameter value s cannot be converted correctly
func NewUint8s(s string) (Uint8, error) {
	u, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return Uint8{valid: false, Val: 0}, err
	}
	return Uint8{valid: true, Val: uint8(u)}, nil
}

// Set is a convenience method to set the value and set the valid flag to true
func (u *Uint8) Set(ui uint8) {
	u.Val = ui
	u.valid = true
}

// SetNull sets valid to false and Val to the zero value
func (u *Uint8) SetNull() {
	u.Val = 0
	u.valid = false
}

// Null implements the Nuller interface and returns the condition of is null
func (u Uint8) Null() bool {
	return !u.valid
}

// NotNull is true if valid
func (u Uint8) NotNull() bool {
	return u.valid
}

// Zero implements the Zeroer interface and returns the condition of is zero
func (u Uint8) Zero() bool {
	return u.Val == 0 && u.valid == true
}

// NonZero is true if valid and not the zero value
func (u Uint8) NonZero() bool {
	return u.Val != 0 && u.valid == true
}

// NullOrZero is true if not valid and is the zero value
func (u Uint8) NullOrZero() bool {
	return !u.valid || u.Val == 0
}

// Equals return true if the nullable is non-null and its value equals the
// parameter passed in
func (u Uint8) Equals(ui uint8) bool {
	return u.valid && (u.Val == ui)
}

func (u Uint8) EqualsI(uI interface{}) (bool, error) {
	switch x := uI.(type) {
	case uint8:
		return u.Equals(x), nil
	case Uint8:
		return u.Equals(x.Val), nil
	default:
		err := fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

func (u Uint8) GT(ui uint8) bool {
	return u.valid && (u.Val > ui)
}

func (u Uint8) GTI(uI interface{}) (bool, error) {
	switch x := uI.(type) {
	case uint8:
		return u.GT(x), nil
	case Uint8:
		return u.GT(x.Val), nil
	default:
		err := fmt.Errorf("cannot check greater than interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

func (u Uint8) GTE(ui uint8) bool {
	return u.GT(ui) || u.Equals(ui)
}

func (u Uint8) GTEI(uI interface{}) (bool, error) {
	switch x := uI.(type) {
	case uint8:
		return u.GTE(x), nil
	case Uint8:
		return u.GTE(x.Val), nil
	default:
		err := fmt.Errorf("cannot check greater than or equal interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

func (u Uint8) LT(ui uint8) bool {
	return u.valid && (u.Val < ui)
}

func (u Uint8) LTI(uI interface{}) (bool, error) {
	switch x := uI.(type) {
	case uint8:
		return u.LT(x), nil
	case Uint8:
		return u.LT(x.Val), nil
	default:
		err := fmt.Errorf("cannot check less than interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

func (u Uint8) LTE(ui uint8) bool {
	return u.LT(ui) || u.Equals(ui)
}

func (u Uint8) LTEI(uI interface{}) (bool, error) {
	switch x := uI.(type) {
	case uint8:
		return u.LTE(x), nil
	case Uint8:
		return u.LTE(x.Val), nil
	default:
		err := fmt.Errorf("cannot check less than or equal interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

// Merge sets the value of the receiver to the value passed if the value passed
// in is not null
func (u *Uint8) Merge(m Uint8) Uint8 {
	if m.Null() {
		return Uint8{}
	}
	u.Set(m.Val)
	return *u
}

// String implements the Stringer interface and returns the internal string
// so you can use a Uint8 in a fmt.Println statement for example
func (u Uint8) String() string {
	if u.valid == false {
		return "null"
	}
	return fmt.Sprintf("%d", u.Val)
}

// Verbose provides a string explaining the value of the Uint8
func (u Uint8) Verbose() string {
	if u.valid == false {
		return "Nil Uint8"
	}
	return fmt.Sprintf("Uint8 Value: %d", u.Val)
}

// #### This gives you proper json behavior ####

// UnmarshalJSON allows the Uint8 to be read in from json
func (u *Uint8) UnmarshalJSON(data []byte) error {
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
		u64, err = strconv.ParseUint(str, 10, 8)
		if err == nil {
			u.Val = uint8(u64)
		} else {
			err = unmarshalTypeError(data, v, Uint8{}, err)
		}

	case map[string]interface{}:
		var n64 sql.NullInt64
		err = json.Unmarshal(data, &n64)
		if err == nil {
			u.Val, u.valid = uint8(n64.Int64), n64.Valid
		}
	case nil:
		u.valid = false
		return nil
	default:
		err = unmarshalTypeError(data, v, Uint8{}, nil)
	}
	u.valid = err == nil
	return err
}

// MarshalJSON writes out the Uint8 as json.
func (u Uint8) MarshalJSON() ([]byte, error) {
	if !u.valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%d", u.Val)), nil
}

// #### This gives you proper writes to a sql database ####

// Scan implements the Scanner interface for SQL retrieval.
func (u *Uint8) Scan(value interface{}) error {
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
	u.Val = uint8(nb.Int64)
	return nil
}

// Value implements the driver Valuer interface for SQL storage.
func (u Uint8) Value() (driver.Value, error) {
	if !u.valid {
		return nil, nil
	}
	return int64(u.Val), nil
}
