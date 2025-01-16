// Package date provides a date custom type.
package date

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/TN-INCORPORATION/kit/v2/apptime"
	"github.com/TN-INCORPORATION/kit/v2/timezone"
	"strings"
	"time"
)

// DateFormat is RFC3339 truncated for just date
const DateFormat = "2006-01-02"

// Date represents the year, month and day
type Date struct {
	time.Time
}

// NewDate returns the Date corresponding to yyyy-mm-dd
func NewDate(year int, month time.Month, day int) Date {
	t := time.Date(year, month, day, 0, 0, 0, 0, timezone.GetTimeZone())
	return Date{Time: t}
}

// NewDatet is a go constructor to create a new date based on a time object
func NewDatet(t time.Time) Date {
	return Trunc(t)
}

// NewDates is a go constructor to create a new date based on a date string
// of the form "2006-01-02"
func NewDates(ds string) (Date, error) {
	t, err := time.ParseInLocation(DateFormat, ds, timezone.GetTimeZone())
	if err != nil {
		return Date{}, err
	}
	return Date{Time: t}, nil
}

// NowDate is a go constructor to create a new date based on time.Now()
func NowDate() Date {
	t2 := TruncTime(apptime.Now())
	return Date{Time: t2}
}

// Trunc truncates the time.Time into a Date
func Trunc(t time.Time) Date {
	t2 := TruncTime(t)
	return Date{Time: t2}
}

// TruncTime truncates the time.Time into zero time
func TruncTime(t time.Time) time.Time {
	yy, mm, dd := t.Date()
	return time.Date(yy, mm, dd, 0, 0, 0, 0, t.Location())
}

// ToTime returns a time.Time with zero time  based on the Date
func (d Date) ToTime() time.Time {
	return d.Time
}

// After the Date d is after u.
func (d Date) After(u Date) bool {
	return d.Time.After(u.Time)
}

// Before reports whether the Date d is before u.
func (d Date) Before(u Date) bool {
	return d.Time.Before(u.Time)
}

// Equal d and u represent the same time instant.
func (d Date) Equal(u Date) bool {
	dY, dM, dD := d.Date()
	tY, tM, tD := u.Date()
	isSameDate := dY == tY && dM == tM && dD == tD
	return isSameDate
}

// String returns the date formatted
func (d Date) String() string {
	return d.Time.Format(DateFormat)
}

// IsZero reports whether t represents the zero time instant,
// January 1, year 1, 00:00:00 UTC.
func (d Date) IsZero() bool {
	return d.Time.IsZero()
}

// Date returns the year, month, and day in which d occurs.
func (d Date) Date() (year int, month time.Month, day int) {
	year = d.Year()
	month = d.Month()
	day = d.Day()
	return
}

// Year returns the year in which d occurs.
func (d Date) Year() int {
	return d.Time.Year()
}

// Month returns the month of the year specified by d.
func (d Date) Month() time.Month {
	return d.Time.Month()
}

// Day returns the day of the month specified by d.
func (d Date) Day() int {
	return d.Time.Day()
}

// Weekday returns the day of the week specified by d.
func (d Date) Weekday() time.Weekday {
	return d.Time.Weekday()
}

// ISOWeek returns the ISO 8601 year and week number in which t occurs.
// Week ranges from 1 to 53. Jan 01 to Jan 03 of year n might belong to
// week 52 or 53 of year n-1, and Dec 29 to Dec 31 might belong to week 1
// of year n+1.
func (d Date) ISOWeek() (year, week int) {
	return d.Time.ISOWeek()
}

// YearDay returns the day of the year specified by d, in the range [1,365] for non-leap years,
// and [1,366] in leap years.
func (d Date) YearDay() int {
	return d.Time.YearDay()
}

// Sub returns the duration d-u. If the result exceeds the maximum (or minimum)
// value that can be stored in a Duration, the maximum (or minimum) duration
// will be returned.
// To compute u-d for a duration d, use u.Add(-d).
func (d Date) Sub(u Date) time.Duration {
	return d.Time.Sub(u.Time)
}

func date(t time.Time) Date {
	return Date{Time: t}
}

// IsLeapYear specifies is the year of this date is a leap year
func (d Date) IsLeapYear() bool {
	return IsLeapYear(d.Time.Year())
}

// IsLeapYear specifies is the year of this date is a leap year
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// AddDate returns the time corresponding to adding the
func (d Date) AddDate(years int, months int, days int) Date {
	return date(d.Time.AddDate(years, months, days))
}

// MarshalJSON implements the json.Marshaler interface.
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte("\"" + d.String() + "\""), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Date) UnmarshalJSON(data []byte) error {
	if data == nil {
		return nil
	}
	s := string(data)
	s = strings.TrimPrefix(s, "\"")
	s = strings.TrimSuffix(s, "\"")
	t, err := time.ParseInLocation(DateFormat, s, timezone.GetTimeZone())
	if err != nil {
		return errors.New(fmt.Sprintf("json: value (%s) Date format needs to conform to RFC3339 truncated for just date (%s)", s, DateFormat))
	}
	*d = NewDatet(t)
	return nil
}

// Scan implements the Scanner interface for SQL retrieval.
func (d *Date) Scan(value interface{}) error {
	var err error
	switch x := value.(type) {
	case time.Time:
		d.Time = TruncTime(x.In(timezone.GetTimeZone()))
	case string:
		d.Time, err = time.ParseInLocation(DateFormat, x, timezone.GetTimeZone())
	default:
		err = fmt.Errorf("null: cannot scan type %T into date.Date: %v", value, value)
	}
	return err
}

// Value implements the driver Valuer interface for SQL storage.
func (d Date) Value() (driver.Value, error) {
	return d.Time, nil
}
