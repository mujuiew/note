package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

// String that supports a null state.  Json serialization, SQL database
// storage and stringer interfaces are support
// swagger:type string
type String struct {
	// Val is the internal string value
	Val string
	// valid is true if String is not NULL
	valid bool
}

// NewString creates a new String from a standard library string
func NewString(s string) String {
	return String{valid: true, Val: s}
}

// NewStrings creates a new String from a standard library string.  This method
// is here to support code generation, otherwise it is unnecessary.
func NewStrings(s string) (String, error) {
	return String{valid: true, Val: s}, nil
}

// Set is a convenience method to set the value and set the valid flag to true
func (s *String) Set(st string) {
	s.Val = st
	s.valid = true
}

// SetNull sets valid to false and Val to the zero value
func (s *String) SetNull() {
	s.Val = ""
	s.valid = false
}

// Null implements the Nuller interface and returns the condition of is null
func (s String) Null() bool {
	return !s.valid
}

// NotNull is true if valid
func (s String) NotNull() bool {
	return s.valid
}

// Zero implements the Zeroer interface and returns the condition of is zero
func (s String) Zero() bool {
	return s.Val == "" && s.valid == true
}

// NonZero is true if valid and not the zero value
func (s String) NonZero() bool {
	return s.Val != "" && s.valid == true
}

// NullOrZero is true if not valid and is the zero value
func (s String) NullOrZero() bool {
	return !s.valid || s.Val == ""
}

// Equals return true if the nullable is non-null and its value equals the
// parameter passed in
func (s String) Equals(st string) bool {
	return s.valid && (strings.Compare(s.Val, st) == 0)
}

// Merge sets the value of the receiver to the value passed if the value passed
// in is not null
func (s *String) Merge(m String) String {
	if m.Null() {
		return String{}
	}
	s.Set(m.Val)
	return *s
}

// String implements the Stringer interface and returns the internal string
// so you can use a String in a fmt.Println statement for example
func (s String) String() string {
	if s.valid == false {
		return "null"
	}
	return s.Val
}

// Verbose provides a string explaining the value of the String
func (s String) Verbose() string {
	if s.valid == false {
		return "Nil String"
	} else if s.Val == "" {
		return "Empty String"
	}
	return fmt.Sprintf("String Value: %s", s.Val)
}

// #### This gives you proper json behavior ####

// UnmarshalJSON allows the String to be read in from json
func (s *String) UnmarshalJSON(b []byte) error {
	if b == nil || len(b) == 0 {
		s.Val, s.valid = "", false
		return nil
	}
	var err error
	var v interface{}
	if err = json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		s.Val = x
		s.valid = true
	case map[string]interface{}:
		err = unmarshalTypeError(b, v, String{}, nil)
	case nil:
		s.valid = false
		return nil
	default:
		err = unmarshalTypeError(b, v, String{}, nil)
	}
	s.valid = err == nil
	return err
}

// MarshalJSON writes out the String as json.
func (s String) MarshalJSON() ([]byte, error) {
	if s.valid == false {
		return []byte("null"), nil
	}
	// str := "\"" + s.Val + "\""
	// return []byte(str), nil
	return json.Marshal(s.Val)
}

// #### This gives you proper writes to a sql database ####

// Scan implements the Scanner interface for SQL retrieval.
func (s *String) Scan(value interface{}) error {
	if value == nil {
		s.Val, s.valid = "", false
		return nil
	}
	var ns sql.NullString
	err := ns.Scan(value)
	if err != nil {
		s.Val, s.valid = "", false
		return err
	}
	s.valid = ns.Valid
	s.Val = ns.String
	return nil
}

// Value implements the driver Valuer interface for SQL storage.
func (s String) Value() (driver.Value, error) {
	if !s.valid {
		return nil, nil
	}
	return s.Val, nil
}

func (s String) EqualsI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case string:
		return s.Equals(x), nil
	case String:
		return s.Equals(x.Val), nil
	default:
		err := fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

func (s String) GT(str string) bool {
	return s.valid && s.Val > str
}

func (s String) GTI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case string:
		return s.GT(x), nil
	case String:
		return s.GT(x.Val), nil
	default:
		err := fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

func (s String) GTE(str string) bool {
	return s.GT(str) || s.Equals(str)
}

func (s String) GTEI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case string:
		return s.GTE(x), nil
	case String:
		return s.GTE(x.Val), nil
	default:
		err := fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

func (s String) LT(str string) bool {
	return s.valid && s.Val < str
}

func (s String) LTI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case string:
		return s.LT(x), nil
	case String:
		return s.LT(x.Val), nil
	default:
		err := fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

func (s String) LTE(str string) bool {
	return s.LT(str) || s.Equals(str)
}

func (s String) LTEI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case string:
		return s.LTE(x), nil
	case String:
		return s.LTE(x.Val), nil
	default:
		err := fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

func (s String) Like(str string) bool {
	isPrefix := false
	isSuffix := false
	if strings.HasPrefix(str, "%") {
		str = str[1:]
		isSuffix = true
	}

	if strings.HasSuffix(str, "%") {
		str = str[:len(str)-1]
		isPrefix = true
	}

	if isPrefix {
		return s.valid && strings.HasPrefix(s.Val, str)
	} else if isSuffix {
		return s.valid && strings.HasSuffix(s.Val, str)
	} else if isPrefix && isSuffix {
		return s.valid && strings.Contains(s.Val, str)
	} else {
		return s.Equals(str)
	}
}

func (s String) LikeI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case string:
		return s.Like(x), nil
	case String:
		return s.Like(x.Val), nil
	default:
		err := fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}
