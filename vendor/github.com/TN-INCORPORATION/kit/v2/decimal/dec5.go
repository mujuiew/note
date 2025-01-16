package decimal

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"math"
)

// NewDec5 is a constructor to create a Dec5 from the
// whole and fractional parts.  The fractional part can be negative
// if the whole part is zero, otherwise the sign of whole will
// dictate the sign of the Dec5.  If the value of f is greater
// than 99, 99 will be returned for the fractional part
func NewDec5(w int64, f int64) Dec5 {
	// Check if whole is equal to zero, if so sign is
	// determined by the fraction, otherwise whole does
	if w != 0 {
		f = absolute(f)
	}
	val := (w * Dec5Base)
	// Clip the fractional part
	f = clip(f, 99999)
	// create the complete value
	if w > 0 {
		val = val + f
	} else if w < 0 {
		val = val - f
	} else {
		val = f
	}
	return Dec5{Val: val}
}

// NewDec5s is a constructor to create a Dec5 from a string
func NewDec5s(s string) (Dec5, error) {
	ds, err := PartsFromString(s, Dec5Frac)
	if err != nil {
		return Dec5{}, err
	}
	// make the new dec
	if len(ds) == both {
		return NewDec5(ds[0], ds[1]), nil
	}
	return NewDec5(ds[0], 0), nil
}

// NewDec5f is a constuctor to create an Dec5 from a float64.
func NewDec5f(f float64) Dec5 {
	return roundToDec5(f)
}

// NewDec5Dec2 is a constructor to create a Dec5 from a Dec2
func NewDec5Dec2(d2 Dec2) Dec5 {
	return Dec5{Val: d2.Val * int64(Dec5Dec2Multiplier)}
}

// NewDec5Dec8 is a constructor to create a Dec5 from a Dec8
func NewDec5Dec8(d8 Dec8) Dec5 {
	rounder := int64(Dec5Dec8Rounder)
	if d8.LTZero() {
		rounder = int64(-Dec5Dec8Rounder)
	}
	return Dec5{Val: (d8.Val + rounder) / int64(Dec5Dec8Divisor)}
}

// NewDec5i is a constructor to create a Dec5 from an int64 where
// 120000 = 120000.00
func NewDec5i(i int64) Dec5 {
	return Dec5{Val: i * int64(100000)}
}

// NewDec5Raw is a constructor to create a Dec5 from and int64 where
// 120000 = 1.20000
func NewDec5Raw(i int64) Dec5 {
	return Dec5{Val: i}
}

// TruncDec2 returns a Dec2 containing the value of the original Dec5 truncated
// to a Dec (with no rounding) and a Dec5 with the last three digits of the
// original truncated Dec5
func (a Dec5) TruncDec2() (Dec2, Dec5) {
	neg := a.Val < 0
	d2 := Dec2{Val: a.Val / 1000}
	rem := modulo(a.Val, 1000)
	if neg {
		rem = -rem
	}
	d5 := Dec5{Val: rem}
	return d2, d5
}

// Round rounds the Dec5 to a Dec2 but uses direct rounding logic
// Info : https://en.wikipedia.org/wiki/Rounding
func (a Dec5) Round() Dec2 {
	rounder := int64(Dec2Dec5Rounder)
	if a.LTZero() {
		rounder = int64(-Dec2Dec5Rounder)
	}
	return Dec2{Val: (a.Val + rounder) / int64(Dec2Dec5Divisor)}
}

// RoundUp rounds the Dec5 to a Dec2 but uses direct rounding up logic
// Info : https://en.wikipedia.org/wiki/Rounding
func (a Dec5) RoundUp() Dec2 {
	f := a.FracInt()
	if modulo(f, 1000) == 0 || f/100 == 0 {
		return Dec2{Val: (a.Val) / int64(Dec2Dec5Divisor)}
	}
	rounder := int64(Dec2Dec5Rounder * 2)
	if a.LTZero() {
		rounder = -rounder
	}
	return Dec2{Val: (a.Val + rounder) / int64(Dec2Dec5Divisor)}
}

// RoundDown rounds the Dec5 to a Dec2 but uses direct rounding down logic
// Info : https://en.wikipedia.org/wiki/Rounding
func (a Dec5) RoundDown() Dec2 {
	return Dec2{Val: a.Val / int64(Dec2Dec5Divisor)}
}

// roundToDec5 is a helper function to round off a float value to the proper
// number of digits to the right of the decimal point.
func roundToDec5(f float64) Dec5 {
	if f > 0 {
		return Dec5{Val: int64((f * Dec5Base) + roundfrac)}
	}
	return Dec5{Val: int64((f * Dec5Base) - roundfrac)}
}

// Multf multiplies itself with the passed in parameter f float64 and returns
// The returned Dec5 will be rounded to 5 digits of percision.
func (a Dec5) Multf(f float64) Dec5 {
	return a.Mult(NewDec5f(f))
}

// Mult multiplies itself with the passed in parameter d Dec5 and returns
// The returned Dec5 will be rounded to 5 digits of percision.
func (a Dec5) Mult(d Dec5) Dec5 {
	if a.Val == 0 || d.Val == 0 {
		return Dec5{}
	}
	i := a.Val * d.Val
	p := i >= 0
	f := modulo(i, Dec5Base)
	i = i / Dec5Base
	if f >= 50000 {
		if p {
			i = i + 1
		} else {
			i = i - 1
		}
	}
	return Dec5{Val: i}
}

// Divf divides itself by the passed in parameter f float64 and returns
// The returned Dec5 will be rounded to 5 digits of percision.
func (a Dec5) Divf(f float64) Dec5 {
	return a.Div(NewDec5f(f))
}

// Div divides itself by the passed in parameter d Dec5 and returns
// The returned Dec5 will be rounded to 5 digits of percision.
func (a Dec5) Div(d Dec5) Dec5 {
	// Use quoRem for division
	val := quoRem(a.Val, d.Val, Dec5Base)
	return Dec5{Val: val}
}

// DivSafef divides itself by the passed in parameter f float64 and returns
// The returned Dec5 will be rounded to 5 digits of percision.
func (a Dec5) DivSafef(f float64) (Dec5, error) {
	if f == 0 {
		return Dec5{}, ErrDivByZero
	}
	return a.Divf(f), nil
}

// DivSafe divides itself by the passed in parameter d Dec5 and returns
// The returned Dec5 will be rounded to 5 digits of percision.
func (a Dec5) DivSafe(d Dec5) (Dec5, error) {
	if d.Val == 0 {
		return Dec5{}, ErrDivByZero
	}
	return a.Div(d), nil
}

// Addf adds itself with the passed in parameter f float64 and returns
func (a Dec5) Addf(f float64) Dec5 {
	return a.Add(NewDec5f(f))
}

// Add adds itself with the passed in parameter d Dec5 and returns
func (a Dec5) Add(d Dec5) Dec5 {
	return Dec5{Val: a.Val + d.Val}
}

// Subf subtracts the passed in parameter f float64 from itself and returns
func (a Dec5) Subf(f float64) Dec5 {
	return a.Sub(NewDec5f(f))
}

// Sub subtracts the passed in parameter d Dec5 from itself and returns
func (a Dec5) Sub(d Dec5) Dec5 {
	return Dec5{Val: a.Val - d.Val}
}

// Whole returns a new Dec5 without the fractional part.
func (a Dec5) Whole() Dec5 {
	return Dec5{Val: (a.Val / Dec5Base) * Dec5Base}
}

// WholeInt returns the Whole part as an int64.
func (a Dec5) WholeInt() int64 {
	return int64(a.Val / Dec5Base)
}

// Frac returns a new Dec5 with just the fractional part and will always be positive.
func (a Dec5) Frac() Dec5 {
	return Dec5{Val: a.FracInt()}
}

// FracInt returns the fractional part as an int64 and will always be positive.
func (a Dec5) FracInt() int64 {
	i := int64(a.Val)
	if i < 0 {
		i = -i
	}
	return int64(i) % Dec5Base
}

// Modf takes in a float and excutes modulo
func (a Dec5) Modf(f float64) Dec5 {
	return a.Mod(NewDec5f(f))
}

// Mod returns a new Dec5 after excuting modulo.
func (a Dec5) Mod(d Dec5) Dec5 {
	return Dec5{Val: modulo(a.Val, d.Val)}
}

// ModSafef protects against n mod 0 is undefined, possibly resulting in a "Division by zero"
func (a Dec5) ModSafef(f float64) (Dec5, error) {
	if f == 0 {
		return Dec5{}, ErrDivByZero
	}
	return a.Modf(f), nil
}

// ModSafe protects against n mod 0 is undefined, possibly resulting in a "Division by zero"
func (a Dec5) ModSafe(d Dec5) (Dec5, error) {
	if d.Val == 0 {
		return Dec5{}, ErrDivByZero
	}
	return a.Mod(d), nil
}

// Float converts an Dec5 to a float64 with appropriate decimal point location.
// Casting an Dec5 directly may result in error, so please use the Float helper.
func (a Dec5) Float() float64 {
	return float64(a.Val) / float64(Dec5Base)
}

// StringInt returns a string using the internal representation of the Dec5.
func (a Dec5) StringInt() string {
	return fmt.Sprintf("%d", a.Val)
}

// String implements the Stringer interface and returns a string using a
// 11.5 format that behaves like a floating point string representation.
func (a Dec5) String() string {
	var sign string
	wi := a.WholeInt()
	// If the whole part is negative, we have to compensate for the sign
	if wi == 0 && a.Val < 0 {
		sign = "-"
	}
	return fmt.Sprintf("%s%d.%05d", sign, wi, a.FracInt())
}

// UnmarshalJSON reads from a byte buffer b to extract a float and
// convert the value to a Dec5
func (a *Dec5) UnmarshalJSON(b []byte) error {
	*a = Dec5{}
	ds, err := PartsFromString(string(b), Dec5Frac)
	if err != nil {
		return errors.New(fmt.Sprintf("json: value (%s) Dec5 invalid format -  %s", string(b), err.Error()))

	}
	// make the new dec
	if len(ds) == 1 {
		*a = NewDec5(ds[0], 0)
	} else if len(ds) == 2 {
		*a = NewDec5(ds[0], ds[1])
	}
	return nil
}

// MarshalJSON writes a float value representation of the Dec5
// to the marshaller
func (a Dec5) MarshalJSON() ([]byte, error) {
	// Note, could also use
	// f := a.Float()
	// but the string seems safer, more testing
	return []byte(a.String()), nil
}

// IsZero determines if the value of the Dec5 is zero
func (a Dec5) IsZero() bool {
	return a.Val == 0
}

// GTZero determines if the value of the Dec5 is greater than zero
func (a Dec5) GTZero() bool {
	return a.Val > 0
}

// LTZero determines if the value of the Dec5 is less than zero
func (a Dec5) LTZero() bool {
	return a.Val < 0
}

// LTf deterimines if the Dec5 is less than the parameter f float64
func (a Dec5) LTf(f float64) bool {
	return a.Val < NewDec5f(f).Val
}

// LT deterimines if the Dec5 is less than the parameter d Dec5
func (a Dec5) LT(d Dec5) bool {
	return a.Val < d.Val
}

// GTf deterimines if the Dec5 is greater than the parameter f float64
func (a Dec5) GTf(f float64) bool {
	return a.Val > NewDec5f(f).Val
}

// GT deterimines if the Dec5 is greater than the parameter d Dec5
func (a Dec5) GT(d Dec5) bool {
	return a.Val > d.Val
}

// Abs returns the absolute value of the Dec5 as a new Dec5
func (a Dec5) Abs() Dec5 {
	if a.Val >= 0 {
		return Dec5{Val: a.Val}
	}
	return Dec5{Val: -a.Val}
}

// Neg returns the negative value of the Dec5 as a new Dec5
// or -Dec5
func (a Dec5) Neg() Dec5 {
	return Dec5{Val: -a.Val}
}

// Pow returns the Dec5 to the exp power
func (a Dec5) PowSafe(exp float64) (Dec5, error) {
	p := math.Pow(a.Float(), exp)
	if math.IsNaN(p) {
		return Dec5{}, errors.New("NaN")
	}
	if math.IsInf(p, 0) {
		return Dec5{}, errors.New("+/-Inf")
	}
	return NewDec5f(p), nil
}

// Pow returns the Dec5 to the exp power (exp specified as a Dec5)
func (a Dec5) Pow5Safe(exp Dec5) (Dec5, error) {
	return a.PowSafe(exp.Float())
}

// Scan implements the sql.Scanner interface for database deserialization.
func (a *Dec5) Scan(value interface{}) error {
	// first try to see if the data is stored in database as a Numeric datatype
	switch v := value.(type) {

	case float32:
		*a = NewDec5f(float64(v))
		return nil

	case float64:
		*a = NewDec5f(v)
		return nil

	case int64:
		*a = Dec5{Val: v}
		return nil

	default:
		// default is trying to interpret value stored as string
		str, err := scanAsString(v)
		if err != nil {
			return err
		}
		temp, err := NewDec5s(str)
		if err != nil {
			return err
		}
		*a = temp
		return err
	}
}

// Value implements the driver.Valuer interface for database serialization.
func (a Dec5) Value() (driver.Value, error) {
	return a.String(), nil
}

// StringHuman - String for human reading
func (a Dec5) StringHuman() string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%.5f", a.Float())
}

// Dec5Slice is a slice of Dec5s for the purpose of sorting. It is modeled
// after the standard library's sort.Float64Slice
type Dec5Slice []Dec5

func (a Dec5Slice) Len() int           { return len(a) }
func (a Dec5Slice) Less(i, j int) bool { return a[i].Val < a[j].Val }
func (a Dec5Slice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
