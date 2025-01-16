package null

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TN-INCORPORATION/kit/v2/apptime"
	"github.com/TN-INCORPORATION/kit/v2/timezone"
	"time"
)

// HoursInDay is the number of hours in a day
const HoursInDay = 24

// MinutesInDay is the number of minutes in a day
const MinutesInDay = 1440

// SecondsInDay is the number of seconds in a day
const SecondsInDay = 86400

// MinutesInHour is the number of minutes in a hour
const MinutesInHour = 60

// SecondsInHour is the number of seconds in a hour
const SecondsInHour = 3600

// SecondsInMinute is the number of seconds in a minute
const SecondsInMinute = 60

// Time that supports a null state.  Json serialization, SQL database
// storage and stringer interfaces are support
// swagger:strfmt date-time
type Time struct {
	// Val is the internal time.Time value
	Val time.Time
	// valid is true it Time is not NULL
	valid bool
}

// NewTimeNow creates a new null.Time that is initialized to the current time
// with a UTC location
func NewTimeNow() Time {
	return Time{valid: true, Val: apptime.Now()}
}

// NewTime creates a new Time from a standard library time.Time
func NewTime(t time.Time) Time {
	return Time{valid: true, Val: t}
}

// NewTimes creates a new Time from a standard library string, returning an
// error if the parameter value s cannot be converted correctly
func NewTimes(s string) (Time, error) {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return Time{valid: false}, err
	}
	return Time{valid: true, Val: t}, nil
}

// Set is a convenience method to set the value and set the valid flag to true
func (t *Time) Set(ti time.Time) {
	t.Val = ti
	t.valid = true
}

// Null implements the Nuller interface and returns the condition of is null
func (t Time) Null() bool {
	return !t.valid
}

// NotNull is true if valid
func (t *Time) NotNull() bool {
	return t.valid
}

// SetNull sets valid to false and Val to the zero value
func (t *Time) SetNull() {
	t.Val = time.Time{}
	t.valid = false
}

// Zero implements the Zeroer interface and returns the condition of is zero
func (t Time) Zero() bool {
	return t.Val.IsZero() && t.valid == true
}

// NonZero is true if valid and not the zero value
func (t Time) NonZero() bool {
	return !t.Val.IsZero() && t.valid == true
}

// NullOrZero is true if not valid and is the zero value
func (t Time) NullOrZero() bool {
	return !t.valid || t.Val.IsZero()
}

// Equals return true if the nullable is non-null and its value equals the
// parameter passed in
func (t Time) Equals(ti time.Time) bool {
	return t.valid && (t.Val.Equal(ti))
}

// EqualsI return true if the nullable is non-null and its value equals the
// interface parameter passed in
func (t Time) EqualsI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case time.Time:
		return t.Equals(x), nil
	case Time:
		return t.Equals(x.Val), nil
	default:
		err := fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

// GT returns true if the DateTime t is after Time ti
func (t Time) GT(ti time.Time) bool {
	return t.valid && t.Val.After(ti)
}

// GTI return true if the DateTime t is after interface i
// interface parameter passed in
func (t Time) GTI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case time.Time:
		return t.GT(x), nil
	case Time:
		return t.GT(x.Val), nil
	default:
		err := fmt.Errorf("cannot check greater than interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

// GTE returns true if the DateTime t is after or equal to time ti
func (t Time) GTE(ti time.Time) bool {
	return t.GT(ti) || t.Equals(ti)
}

// GTEI return true if the DateTime t is after or equal to interface time i
// interface parameter passed in
func (t Time) GTEI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case time.Time:
		return t.GTE(x), nil
	case Time:
		return t.GTE(x.Val), nil
	default:
		err := fmt.Errorf("cannot check greater than or equal interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

// LT returns true if the DateTime t is before Time t
func (t Time) LT(ti time.Time) bool {
	return t.valid && t.Val.Before(ti)
}

// LTI return true if the DateTime t is before interface Time i
// interface parameter passed in
func (t Time) LTI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case time.Time:
		return t.LT(x), nil
	case Time:
		return t.LT(x.Val), nil
	default:
		err := fmt.Errorf("cannot check less than interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

// LTE returns true if the DateTime t is before or equal to Time ti
func (t Time) LTE(ti time.Time) bool {
	return t.LT(ti) || t.Equals(ti)
}

// LTEI return true if the DateTime t is before or equal to interface Time i
// interface parameter passed in
func (t Time) LTEI(i interface{}) (bool, error) {
	switch x := i.(type) {
	case time.Time:
		return t.LTE(x), nil
	case Time:
		return t.LTE(x.Val), nil
	default:
		err := fmt.Errorf("cannot check less than or equal interface: variable type invalid %+v,%T", x, x)
		return false, err
	}
}

// NotEqual returns true if the DateTime t is not equal to Time ti
func (t Time) NotEqual(ti time.Time) bool {
	return t.valid && !t.Val.Equal(ti)
}

// Merge sets the value of the receiver to the value passed if the value passed
// in is not null
func (t *Time) Merge(m Time) Time {
	if m.Null() {
		return Time{}
	}
	t.Set(m.Val)
	return *t
}

// String implements the Stringer interface and returns the internal string
// so you can use a Time in a fmt.Println statement for example
func (t Time) String() string {
	if t.valid == false {
		return "null"
	}
	return t.Val.Format(time.RFC3339Nano)
}

// Verbose provides a string explaining the value of the Time
func (t Time) Verbose() string {
	if t.valid == false {
		return "Nil Time"
	} else if t.Val.IsZero() {
		return "Empty Time"
	}
	return fmt.Sprintf("Time Value: %s", t.Val.Format(time.RFC3339Nano))
}

// Now sets the Time to the current time with a UTC location
func (t *Time) Now() Time {
	tn := NewTimeNow()
	t.valid = true
	t.Val = tn.Val
	return tn
}

// Date returns the Time struct truncated to the beginning of the day therefore
// the time component is zeroed out
func (t Time) Date() Time {
	if t.Null() {
		return Time{valid: false, Val: time.Time{}}
	}
	yy, mm, dd := t.Val.Date()
	return Time{valid: true, Val: time.Date(yy, mm, dd, 0, 0, 0, 0, t.Val.Location())}
}

// #### This gives you proper json behavior ####

// UnmarshalJSON allows the Time to be read in from json
func (t *Time) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		err = t.Val.UnmarshalJSON(data)
		switch v := err.(type) {
		case *time.ParseError:
			info := "timestamp needs to conform to IETF RFC3339, time zones optional"
			desc := fmt.Sprintf("json: Value(%s) - %s Format(%s) - %s", v.Value, info, v.Layout, v.Error())
			err = errors.New(desc)
			err = unmarshalTypeError(data, v, Time{}, err)
		}
	case map[string]interface{}:
		ti, tiOK := x["Time"].(string)
		valid, validOK := x["valid"].(bool)
		if !tiOK || !validOK {
			err = fmt.Errorf(`json: (%s) unmarshalling object into Go value of type null.Time requires key "Time" to be of type string and key "valid" to be of type bool; found %T and %T, respectively`, data, x["Time"], x["valid"])
			return unmarshalTypeError(data, v, Time{}, err)
		}
		err = t.Val.UnmarshalText([]byte(ti))
		t.valid = valid
		return err
	case nil:
		t.valid = false
		return nil
	default:
		err = unmarshalTypeError(data, v, Time{}, nil)
	}
	t.valid = err == nil
	return err
}

// MarshalJSON writes out the Time as json.
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", t.Val.Format(time.RFC3339Nano))), nil
}

// #### This gives you proper writes to a sql database ####

// Scan implements the Scanner interface for SQL retrieval.
func (t *Time) Scan(value interface{}) error {
	var err error
	switch x := value.(type) {
	case time.Time:
		t.Val = x.In(timezone.GetTimeZone())
	case nil:
		t.valid = false
		return nil
	default:
		err = fmt.Errorf("null: cannot scan type %T into null.Time: %v", value, value)
	}
	t.valid = err == nil
	return err
}

// Value implements the driver Valuer interface for SQL storage.
func (t Time) Value() (driver.Value, error) {
	if !t.valid {
		return nil, nil
	}
	return t.Val.In(timezone.GetTimeZone()), nil
}
