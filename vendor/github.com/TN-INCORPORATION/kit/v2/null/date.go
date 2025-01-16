package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/TN-INCORPORATION/kit/v2/timezone"
	"time"

	"github.com/TN-INCORPORATION/kit/v2/date"
)

// Date that supports a null state.  Json serialization, SQL database
// storage and stringer interfaces are support
type Date struct {
	// Val is the internal date.Date value
	Val date.Date
	// valid is true it Time is not NULL
	valid bool
}

// NewDateNow creates a new null.Date that is initialized to the current time
// with a UTC location truncated to the beginning of the day
func NewDateNow() Date {
	return Date{valid: true, Val: date.NowDate()}
}

// NewDate creates a new Date
func NewDate(d date.Date) Date {
	return Date{valid: true, Val: d}
}

// NewDates creates a new Date, returning an error if the parameter
// value s cannot be converted correctly
func NewDates(s string) (Date, error) {
	d, err := date.NewDates(s)
	if err != nil {
		return Date{valid: false}, err
	}
	return Date{valid: true, Val: d}, nil
}

// Set is a convenience method to set the value and set the valid flag to true
func (d *Date) Set(da date.Date) {
	d.Val = da
	d.valid = true
}

// SetNull sets valid to false and Val to the zero value
func (d *Date) SetNull() {
	*d = Date{}
}

// Null implements the Nuller interface and returns the condition of is null
func (d Date) Null() bool {
	return !d.NotNull()
}

// NotNull is true if valid
func (d Date) NotNull() bool {
	return d.valid
}

// Zero implements the Zeroer interface and returns the condition of is zero
func (d Date) Zero() bool {
	return d.Val.IsZero() && d.valid == true
}

// NonZero is true if valid and not the zero value
func (d Date) NonZero() bool {
	return !d.Val.IsZero() && d.valid == true
}

// NullOrZero is true if not valid and is the zero value
func (d Date) NullOrZero() bool {
	return !d.valid || d.Val.IsZero()
}

// Equals return true if the nullable is non-null and its value equals the
// parameter passed in
func (d Date) Equals(da date.Date) bool {
	return d.valid && (d.Val.Equal(da))
}

func (d Date) EqualsI(v interface{}) (bool, error) {
	switch x := v.(type) {
	case date.Date:
		return d.Equals(x), nil
	case Date:
		return d.Equals(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check equals interface: variable type invalid %+v,%T", x, x)
	}
}

// After the Date d is after da.
func (d Date) After(da date.Date) bool {
	return d.valid && d.Val.After(da)
}

// Before reports whether the Date d is before da.
func (d Date) Before(da date.Date) bool {
	return d.valid && d.Val.Before(da)
}

func (d Date) GT(da date.Date) bool {
	return d.valid && d.Val.After(da)
}

func (d Date) LT(da date.Date) bool {
	return d.valid && d.Val.Before(da)
}

func (d Date) GTI(v interface{}) (bool, error) {
	switch x := v.(type) {
	case date.Date:
		return d.After(x), nil
	case Date:
		return d.After(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check greater than interface: variable type invalid %+v,%T", x, x)
	}
}

func (d Date) LTI(v interface{}) (bool, error) {
	switch x := v.(type) {
	case date.Date:
		return d.Before(x), nil
	case Date:
		return d.Before(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check less than interface: variable type invalid %+v,%T", x, x)
	}
}

func (d Date) GTEI(v interface{}) (bool, error) {
	switch x := v.(type) {
	case date.Date:
		return d.GTE(x), nil
	case Date:
		return d.GTE(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check greater than or equal interface: variable type invalid %+v,%T", x, x)
	}
}

func (d Date) GTE(da date.Date) bool {
	return d.Equals(da) || d.After(da)
}

func (d Date) LTEI(v interface{}) (bool, error) {
	switch x := v.(type) {
	case date.Date:
		return d.LTE(x), nil
	case Date:
		return d.LTE(x.Val), nil
	default:
		return false, fmt.Errorf("cannot check less than or equal interface: variable type invalid %+v,%T", x, x)
	}
}
func (d Date) LTE(da date.Date) bool {
	return d.Equals(da) || d.Before(da)
}

// Merge sets the value of the receiver to the value passed if the value passed
// in is not null
func (d *Date) Merge(m Date) Date {
	if m.Null() {
		return Date{}
	}
	d.Set(m.Val)
	return *d
}

// String implements the Stringer interface and returns the internal string
// so you can use a Time in a fmt.Println statement for example
func (d Date) String() string {
	if d.valid == false {
		return "null"
	}
	return d.Val.String()
}

// Verbose provides a string explaining the value of the Time
func (d Date) Verbose() string {
	if d.valid == false {
		return "Nil Date"
	} else if d.Val.IsZero() {
		return "Empty Date"
	}
	return fmt.Sprintf("Date Value: %s", d.Val.String())
}

// #### This gives you proper json behavior ####

// UnmarshalJSON allows the Date to be read in from json
func (d *Date) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case string:
		err = d.Val.UnmarshalJSON(data)
		if err != nil {
			err = unmarshalTypeError(data, v, Date{}, err)
		}

	case map[string]interface{}:
		da, daOK := x["Date"].(string)
		valid, validOK := x["valid"].(bool)
		if !daOK || !validOK {
			err = fmt.Errorf(`json: (%s) unmarshalling object into Go value of type null.Date requires key "Date" to be of type string and key "valid" to be of type bool; found %T and %T, respectively`, data, x["Time"], x["valid"])
			return unmarshalTypeError(data, v, Date{}, err)
		}
		err = d.Val.UnmarshalText([]byte(da))
		d.valid = valid
		return unmarshalTypeError(data, v, Date{}, err)
	case nil:
		d.valid = false
		return nil
	default:
		err = unmarshalTypeError(data, v, Date{}, nil)
	}
	d.valid = err == nil
	return err
}

// MarshalJSON writes out the Time as json.
func (d Date) MarshalJSON() ([]byte, error) {
	if !d.valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", d.Val.String())), nil
}

// #### This gives you proper writes to a sql database ####

// Scan implements the Scanner interface for SQL retrieval.
func (d *Date) Scan(value interface{}) error {
	var err error
	switch x := value.(type) {
	case date.Date:
		d.Val = x
	case time.Time:
		d.Val = date.NewDatet(x.In(timezone.GetTimeZone()))
	case nil:
		d.valid = false
		return nil
	default:
		err = fmt.Errorf("null: cannot scan type %T into null.Date: %v", value, value)
	}
	d.valid = err == nil
	return err
}

// Value implements the driver Valuer interface for SQL storage.
func (d Date) Value() (driver.Value, error) {
	if !d.valid {
		return nil, nil
	}
	return d.Val.Value()
}
