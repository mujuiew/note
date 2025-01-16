package decimal

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"math"
)

// NewDec2 is a Go constructor to create a Dec2 from the
// whole and fractional parts.  The fractional part can be negative
// if the whole part is zero, otherwise the sign of whole will
// dictate the sign of the Dec2.  If the value of f is greater
// than 99, 99 will be returned for the fractional part
func NewDec2(w int64, f int64) Dec2 {

	if w != 0 {
		f = absolute(f)
	}
	val := w * Dec2Base

	// Clip the fractional part
	f = clip(f, 99)

	// create the complete value
	if w > 0 {
		val += f
	} else if w < 0 {
		val -= f
	} else {
		val = f
	}
	return Dec2{Val: val}
}

// NewDec2s constructor to create a Dec2  from string
func NewDec2s(s string) (Dec2, error) {
	ds, err := PartsFromString(s, Dec2Frac)
	if err != nil {
		return Dec2{}, err
	}
	// make the new dec
	if len(ds) == both {
		return NewDec2(ds[0], ds[1]), nil
	}
	return NewDec2(ds[0], 0), nil
}

// NewDec2f constructor to create an Dec2 from a float64.
func NewDec2f(f float64) Dec2 {
	return roundToDec2(f)
}

// NewDec2Dec5 is a  constructor to create a Dec2 from a Dec5
func NewDec2Dec5(d5 Dec5) Dec2 {
	rounder := int64(Dec2Dec5Rounder)
	if d5.LTZero() {
		rounder = int64(-Dec2Dec5Rounder)
	}
	return Dec2{Val: (d5.Val + rounder) / int64(Dec2Dec5Divisor)}
}

// NewDec2Dec8 is a  constructor to create a Dec2 from a Dec8
func NewDec2Dec8(d8 Dec8) Dec2 {
	rounder := int64(Dec2Dec8Rounder)
	if d8.LTZero() {
		rounder = int64(-Dec2Dec8Rounder)
	}
	return Dec2{Val: (d8.Val + rounder) / int64(Dec2Dec8Divisor)}
}

// NewDec2i is a  constructor to create a Dec2 from an int64
// 120 = 120.00
func NewDec2i(i int64) Dec2 {
	return Dec2{Val: i * int64(100)}
}

// NewDec2Raw is a  constructor to create a Dec2 from and int64 where
// 120 = 1.20
func NewDec2Raw(i int64) Dec2 {
	return Dec2{Val: i}
}

// roundToDec2 is a function to round off a float value to the right of the decimal point.
func roundToDec2(f float64) Dec2 {
	if f > 0 {
		return Dec2{Val: int64((f * float64(Dec2Base)) + roundfrac)}
	}
	return Dec2{Val: int64((f * float64(Dec2Base)) - roundfrac)}
}

// Multf multiplies itself with the passed in parameter f float64 and returns
// The returned Dec2 will be rounded to 2 digits of percision.
func (a Dec2) Multf(f float64) Dec2 {
	return a.Mult(NewDec2f(f))
}

// Mult multiplies itself with the passed in parameter d Dec2 and returns
// The returned Dec2 will be rounded to 2 digits of percision.
func (a Dec2) Mult(d Dec2) Dec2 {
	i := a.Val * d.Val
	p := i >= 0
	f := modulo(i, Dec2Base)
	i = i / Dec2Base
	if f >= 50 {
		if p {
			i = i + 1
		} else {
			i = i - 1
		}
	}
	return Dec2{Val: i}
}

// MultDiv multiplies itself by the numerator Dec5 and divides it by an integer
// value.  This is typically used in accrual.  The result is a Dec5 because
// the numerator is a Dec5 keeping the precision
func (a Dec2) MultDiv(num Dec5, den int64) Dec5 {
	// Convert the receiver value to a Dec5
	ab := NewDec5Dec2(a)
	mult := ab.Mult(num)
	// Convert the denominator to a Dec5
	den5 := NewDec5i(den)
	return mult.Div(den5)
}

// Divf divides itself by the passed in parameter f float64 and returns
// The returned Dec2 will be rounded to 2 digits of percision.
func (a Dec2) Divf(f float64) Dec2 {
	return a.Div(NewDec2f(f))
}

// Div divides itself by the passed in parameter d Dec2 and returns
// The returned Dec2 will be rounded to 2 digits of percision.
func (a Dec2) Div(d Dec2) Dec2 {
	i := a.Val * Dec2Base * Dec2Base / d.Val
	p := i >= 0
	f := modulo(i, Dec2Base)
	i = i / Dec2Base
	if f >= 50 {
		if p {
			i = i + 1
		} else {
			i = i - 1
		}
	}
	return Dec2{Val: i}
}

// DivSafef divides itself by the passed in parameter f float64 and returns
// a new Dec2 object and an error if divide by zero.
// The returned Dec2 will be rounded to 2 digits of percision.
func (a Dec2) DivSafef(f float64) (Dec2, error) {
	if f == 0 {
		return Dec2{}, ErrDivByZero
	}
	return a.Divf(f), nil
}

// DivSafe divides itself by the passed in parameter d Dec2 and returns
// a new Dec2 object and an error if divide by zero.
// The returned Dec2 will be rounded to 2 digits of percision.
func (a Dec2) DivSafe(d Dec2) (Dec2, error) {
	if d.Val == 0 {
		return Dec2{}, ErrDivByZero
	}
	return a.Div(d), nil
}

// Addf adds itself with the passed in parameter f float64 and returns
func (a Dec2) Addf(f float64) Dec2 {
	return a.Add(NewDec2f(f))
}

// Add adds itself with the passed in parameter d Dec2 and returns
func (a Dec2) Add(d Dec2) Dec2 {
	return Dec2{Val: a.Val + d.Val}
}

// Subf subtracts the passed in parameter f float64 from itself and returns
func (a Dec2) Subf(f float64) Dec2 {
	return a.Sub(NewDec2f(f))
}

// Sub subtracts the passed in parameter d Dec2 from itself and returns
func (a Dec2) Sub(d Dec2) Dec2 {
	return Dec2{Val: a.Val - d.Val}
}

// Whole returns a new Dec2 without the fractional part.
func (a Dec2) Whole() Dec2 {
	return Dec2{Val: (int64(a.Val) / int64(Dec2Base)) * int64(Dec2Base)}
}

// WholeInt returns the Whole part as an int64.
func (a Dec2) WholeInt() int64 {
	return int64(a.Val / Dec2Base)
}

// Frac returns a new Dec2 with just the fractional part and will always be positive.
func (a Dec2) Frac() Dec2 {
	return Dec2{Val: a.FracInt()}
}

// FracInt returns the fractional part as an int64 and will always be positive.
func (a Dec2) FracInt() int64 {
	i := int64(a.Val)
	if i < 0 {
		i = -i
	}
	return int64(i) % Dec2Base
}

// Modf takes in a float and excutes modulo
func (a Dec2) Modf(f float64) Dec2 {
	return a.Mod(NewDec2f(f))
}

// Mod returns a new Dec2 after excuting modulo.
func (a Dec2) Mod(d Dec2) Dec2 {
	return Dec2{Val: modulo(a.Val, d.Val)}
}

// ModSafef protects against n mod 0 is undefined, possibly resulting in a "Division by zero"
func (a Dec2) ModSafef(f float64) (Dec2, error) {
	if f == 0 {
		return Dec2{}, ErrDivByZero
	}
	return a.Modf(f), nil
}

// ModSafe protects against n mod 0 is undefined, possibly resulting in a "Division by zero"
func (a Dec2) ModSafe(d Dec2) (Dec2, error) {
	if d.Val == 0 {
		return Dec2{}, ErrDivByZero
	}
	return a.Mod(d), nil
}

// Float converts an Dec2 to a float64 with appropriate decimal point location.
// Casting an Dec2 directly may result in error, so please use the Float helper.
func (a Dec2) Float() float64 {
	return float64(a.Val) / float64(Dec2Base)
}

// StringInt returns a string using the internal representation of the Dec2.
func (a Dec2) StringInt() string {
	return fmt.Sprintf("%d", a.Val)
}

// String implements the Stringer interface and returns a string using a
// 14.2 format that behaves like a floating point string representation.
func (a Dec2) String() string {
	var sign string
	wi := a.WholeInt()
	// If the whole part is negative, we have to compensate for the sign
	if wi == 0 && a.Val < 0 {
		sign = "-"
	}
	return fmt.Sprintf("%s%d.%02d", sign, wi, a.FracInt())
}

// UnmarshalJSON reads from a byte buffer b to extract a float and
// convert the value to a Dec2
func (a *Dec2) UnmarshalJSON(b []byte) error {
	*a = Dec2{}
	ds, err := PartsFromString(string(b), Dec2Frac)
	if err != nil {
		return errors.New(fmt.Sprintf("json: value (%s) Dec2 invalid format -  %s", string(b), err.Error()))
	}
	// make the new dec
	if len(ds) == 1 {
		*a = NewDec2(ds[0], 0)
	} else if len(ds) == 2 {
		*a = NewDec2(ds[0], ds[1])
	}
	return nil
}

// MarshalJSON writes a float value representation of the Dec2
// to the marshaller
func (a Dec2) MarshalJSON() ([]byte, error) {
	// Note, could also use
	// f := a.Float()
	// but the string seems safer, more testing
	return []byte(a.String()), nil
}

// IsZero determines if the value of the Dec2 is zero
func (a Dec2) IsZero() bool {
	return a.Val == 0
}

// GTZero determines if the value of the Dec2 is greater than zero
func (a Dec2) GTZero() bool {
	return a.Val > 0
}

// LTZero determines if the value of the Dec2 is less than zero
func (a Dec2) LTZero() bool {
	return a.Val < 0
}

// LTf deterimines if the Dec2 is less than the parameter f float64
func (a Dec2) LTf(f float64) bool {
	return a.Val < NewDec2f(f).Val
}

// LT deterimines if the Dec2 is less than the parameter d Dec2
func (a Dec2) LT(d Dec2) bool {
	return a.Val < d.Val
}

// GTf deterimines if the Dec2 is greater than the parameter f float64
func (a Dec2) GTf(f float64) bool {
	return a.Val > NewDec2f(f).Val
}

// GT deterimines if the Dec2 is greater than the parameter d Dec2
func (a Dec2) GT(d Dec2) bool {
	return a.Val > d.Val
}

// Abs returns the absolute value of the Dec2 as a new Dec2
func (a Dec2) Abs() Dec2 {
	if a.Val >= 0 {
		return Dec2{Val: a.Val}
	}
	return Dec2{Val: -a.Val}
}

// Neg returns the negative value of the Dec2 as a new Dec2
// or -Dec2
func (a Dec2) Neg() Dec2 {
	return Dec2{Val: -a.Val}
}

// Pow returns the Dec2 to the exp power
func (a Dec2) PowSafe(exp float64) (Dec2, error) {
	p := math.Pow(a.Float(), exp)
	if math.IsNaN(p) {
		return Dec2{}, errors.New("NaN")
	}
	if math.IsInf(p, 0) {
		return Dec2{}, errors.New("+/-Inf")
	}
	return NewDec2f(p), nil
}

// Pow returns the Dec2 to the exp power (exp specified as a Dec2)
func (a Dec2) Pow2Safe(exp Dec2) (Dec2, error) {
	return a.PowSafe(exp.Float())
}

// Scan implements the sql.Scanner interface for database deserialization.
func (a *Dec2) Scan(value interface{}) error {
	// first try to see if the data is stored in database as a Numeric datatype
	switch v := value.(type) {

	case float32:
		*a = NewDec2f(float64(v))
		return nil

	case float64:
		*a = NewDec2f(v)
		return nil

	case int64:
		*a = Dec2{Val: v}
		return nil

	default:
		// default is trying to interpret value stored as string
		str, err := scanAsString(v)
		if err != nil {
			return err
		}
		temp, err := NewDec2s(str)
		if err != nil {
			return err
		}
		*a = temp
		return err
	}
}

// Value implements the driver.Valuer interface for database serialization.
func (a Dec2) Value() (driver.Value, error) {
	return a.String(), nil
}

// ThaiBathWord - Thai word
func (a Dec2) ThaiBathWord() string {
	s, err := fullNumToWords(a.String())
	if err != nil {
		return "Word Error"
	}
	return s
}

// StringHuman - String for human reading
func (a Dec2) StringHuman() string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%.2f", a.Float())
}

// Dec2Slice is a slice of Dec2s for the purpose of sorting. It is modeled
// after the standard library's sort.Float64Slice
type Dec2Slice []Dec2

func (a Dec2Slice) Len() int           { return len(a) }
func (a Dec2Slice) Less(i, j int) bool { return a[i].Val < a[j].Val }
func (a Dec2Slice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
