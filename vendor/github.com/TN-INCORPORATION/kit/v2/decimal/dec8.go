package decimal

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"math"
)

// NewDec8 is a Go constructor to create a Dec8 from the
// whole and fractional parts.  The fractional part can be negative
// if the whole part is zero, otherwise the sign of whole will
// dictate the sign of the Dec8.  If the value of f is greater
// than 99, 99 will be returned for the fractional part
func NewDec8(w int64, f int64) Dec8 {
	// Check if whole is equal to zero, if so sign is
	// determined by the fraction, otherwise whole does
	if w != 0 {
		f = absolute(f)
	}
	val := w * Dec8Base
	// Clip the fractional part
	f = clip(f, 99999999)
	// create the complete value
	if w > 0 {
		val = val + f
	} else if w < 0 {
		val = val - f
	} else {
		val = f
	}
	return Dec8{Val: val}
}

// NewDec8s is a Go constructor to create a Dec8 from a string
// If the string is invalid, then zero will be returned
func NewDec8s(s string) (Dec8, error) {
	ds, err := PartsFromString(s, Dec8Frac)
	if err != nil {
		return Dec8{}, err
	}
	// make the new dec
	if len(ds) == both {
		return NewDec8(ds[0], ds[1]), nil
	}
	return NewDec8(ds[0], 0), nil
}

// NewDec8f is a Go constuctor to create an Dec8 from a float64.
func NewDec8f(f float64) Dec8 {
	return roundToDec8(f)
}

// NewDec8Dec2 is a Go constructor to create a Dec8 from a Dec2
func NewDec8Dec2(d2 Dec2) Dec8 {
	return Dec8{Val: d2.Val * int64(Dec8Dec2Multiplier)}
}

// NewDec8Dec5 is a Go constructor to create a Dec8 from a Dec5
func NewDec8Dec5(d5 Dec5) Dec8 {
	return Dec8{Val: d5.Val * int64(Dec8Dec5Multiplier)}
}

// NewDec8i is a Go constructor to create a Dec8 from an int64 where
// the int64 represents the whole part only 1 = 1.000000000
func NewDec8i(i int64) Dec8 {
	return Dec8{Val: i * int64(100000000)}
}

// NewDec8Raw is a Go constructor to create a Dec8 from and int64 where
// the int64 represents the whole and fractional part 1 = 0.00000001
func NewDec8Raw(i int64) Dec8 {
	return Dec8{Val: i}
}

// roundToDec8 is a helper function to round off a float value to the proper
// number of digits to the right of the decimal point.
func roundToDec8(f float64) Dec8 {
	if f > 0 {
		return Dec8{Val: int64((f * Dec8Base) + roundfrac)}
	}
	return Dec8{Val: int64((f * Dec8Base) - roundfrac)}
}

// Multf multiplies itself with the passed in parameter f float64 and returns
// The returned Dec8 will be rounded to 8 digits of percision.
func (a Dec8) Multf(f float64) Dec8 {
	return a.Mult(NewDec8f(f))
}

// Mult multiplies itself with the passed in parameter d Dec8 and returns
// The returned Dec8 will be rounded to 8 digits of percision.
func (a Dec8) Mult(d Dec8) Dec8 {
	if a.Val == 0 || d.Val == 0 {
		return Dec8{}
	}
	la := math.Log2(float64(absolute(a.Val)))
	ld := math.Log2(float64(absolute(d.Val)))
	// Will it fit in an int64
	if la+ld > 63 {
		return Dec8{}
	}
	i := a.Val * d.Val
	p := i >= 0
	f := modulo(i, Dec8Base)
	i = i / Dec8Base
	if f >= 50000000 {
		if p {
			i = i + 1
		} else {
			i = i - 1
		}
	}
	return Dec8{Val: i}
}

// Divf divides itself by the passed in parameter f float64 and returns
// The returned Dec8 will be rounded to 8 digits of percision.
func (a Dec8) Divf(f float64) Dec8 {
	return a.Div(NewDec8f(f))
}

// Div divides itself by the passed in parameter d Dec8 and returns
// The returned Dec8 will be rounded to 8 digits of percision.
func (a Dec8) Div(d Dec8) Dec8 {
	// Use quoRem for division
	val := quoRem(a.Val, d.Val, Dec8Base)
	return Dec8{Val: val}
}

// DivSafef divides itself by the passed in parameter f float64 and returns
// The returned Dec8 will be rounded to 8 digits of percision.
func (a Dec8) DivSafef(f float64) (Dec8, error) {
	if f == 0 {
		return Dec8{}, ErrDivByZero
	}
	return a.Divf(f), nil
}

// DivSafe divides itself by the passed in parameter d Dec8 and returns
// The returned Dec8 will be rounded to 8 digits of percision.
func (a Dec8) DivSafe(d Dec8) (Dec8, error) {
	if d.Val == 0 {
		return Dec8{}, ErrDivByZero
	}
	return a.Div(d), nil
}

// Addf adds itself with the passed in parameter f float64 and returns
func (a Dec8) Addf(f float64) Dec8 {
	return a.Add(NewDec8f(f))
}

// Add adds itself with the passed in parameter d Dec8 and returns
func (a Dec8) Add(d Dec8) Dec8 {
	return Dec8{Val: a.Val + d.Val}
}

// Subf subtracts the passed in parameter f float64 from itself and returns
func (a Dec8) Subf(f float64) Dec8 {
	return a.Sub(NewDec8f(f))
}

// Sub subtracts the passed in parameter d Dec8 from itself and returns
func (a Dec8) Sub(d Dec8) Dec8 {
	return Dec8{Val: a.Val - d.Val}
}

// Whole returns a new Dec8 without the fractional part.
func (a Dec8) Whole() Dec8 {
	return Dec8{Val: (a.Val / Dec8Base) * Dec8Base}
}

// WholeInt returns the Whole part as an int64.
func (a Dec8) WholeInt() int64 {
	return int64(a.Val / Dec8Base)
}

// Frac returns a new Dec8 with just the fractional part and will always be positive.
func (a Dec8) Frac() Dec8 {
	return Dec8{Val: a.FracInt()}
}

// FracInt returns the fractional part as an int64 and will always be positive.
func (a Dec8) FracInt() int64 {
	i := int64(a.Val)
	if i < 0 {
		i = -i
	}
	return int64(i) % Dec8Base
}

// Modf takes in a float and excutes modulo
func (a Dec8) Modf(f float64) Dec8 {
	return a.Mod(NewDec8f(f))
}

// Mod returns a new Dec8 after excuting modulo.
func (a Dec8) Mod(d Dec8) Dec8 {
	return Dec8{Val: modulo(a.Val, d.Val)}
}

// ModSafef protects against n mod 0 is undefined, possibly resulting in a "Division by zero"
func (a Dec8) ModSafef(f float64) (Dec8, error) {
	if f == 0 {
		return Dec8{}, ErrDivByZero
	}
	return a.Modf(f), nil
}

// ModSafe protects against n mod 0 is undefined, possibly resulting in a "Division by zero"
func (a Dec8) ModSafe(d Dec8) (Dec8, error) {
	if d.Val == 0 {
		return Dec8{}, ErrDivByZero
	}
	return a.Mod(d), nil
}

// Float converts an Dec8 to a float64 with appropriate decimal point location.
// Casting an Dec8 directly may result in error, so please use the Float helper.
func (a Dec8) Float() float64 {
	return float64(a.Val) / float64(Dec8Base)
}

// StringInt returns a string using the internal representation of the Dec8.
func (a Dec8) StringInt() string {
	return fmt.Sprintf("%d", a.Val)
}

// String implements the Stringer interface and returns a string using a
// 8.8 format that behaves like a floating point string representation.
func (a Dec8) String() string {
	var sign string
	wi := a.WholeInt()
	// If the whole part is negative, we have to compensate for the sign
	if wi == 0 && a.Val < 0 {
		sign = "-"
	}
	return fmt.Sprintf("%s%d.%08d", sign, wi, a.FracInt())
}

// UnmarshalJSON reads from a byte buffer b to extract a float and
// convert the value to a Dec8
func (a *Dec8) UnmarshalJSON(b []byte) error {
	*a = Dec8{}
	ds, err := PartsFromString(string(b), Dec8Frac)
	if err != nil {
		return errors.New(fmt.Sprintf("json: value (%s) Dec8 invalid format -  %s", string(b), err.Error()))

	}
	// make the new dec
	if len(ds) == 1 {
		*a = NewDec8(ds[0], 0)
	} else if len(ds) == 2 {
		*a = NewDec8(ds[0], ds[1])
	}
	return nil
}

// MarshalJSON writes a float value representation of the Dec8
// to the marshaller
func (a Dec8) MarshalJSON() ([]byte, error) {
	// Note, could also use
	// f := a.Float()
	// but the string seems safer, more testing
	return []byte(a.String()), nil
}

// IsZero determines if the value of the Dec8 is zero
func (a Dec8) IsZero() bool {
	return a.Val == 0
}

// GTZero determines if the value of the Dec8 is greater than zero
func (a Dec8) GTZero() bool {
	return a.Val > 0
}

// LTZero determines if the value of the Dec8 is less than zero
func (a Dec8) LTZero() bool {
	return a.Val < 0
}

// LTf deterimines if the Dec8 is less than the parameter f float64
func (a Dec8) LTf(f float64) bool {
	return a.Val < NewDec8f(f).Val
}

// LT deterimines if the Dec8 is less than the parameter d Dec8
func (a Dec8) LT(d Dec8) bool {
	return a.Val < d.Val
}

// GTf deterimines if the Dec8 is greater than the parameter f float64
func (a Dec8) GTf(f float64) bool {
	return a.Val > NewDec8f(f).Val
}

// GT deterimines if the Dec8 is greater than the parameter d Dec8
func (a Dec8) GT(d Dec8) bool {
	return a.Val > d.Val
}

// Abs returns the absolute value of the Dec8 as a new Dec8
func (a Dec8) Abs() Dec8 {
	if a.Val >= 0 {
		return Dec8{Val: a.Val}
	}
	return Dec8{Val: -a.Val}
}

// Neg returns the negative value of the Dec8 as a new Dec8
// or -Dec8
func (a Dec8) Neg() Dec8 {
	return Dec8{Val: -a.Val}
}

// Pow returns the Dec8 to the exp power
func (a Dec8) PowSafe(exp float64) (Dec8, error) {
	p := math.Pow(a.Float(), exp)
	if math.IsNaN(p) {
		return Dec8{}, errors.New("NaN")
	}
	if math.IsInf(p, 0) {
		return Dec8{}, errors.New("+/-Inf")
	}
	return NewDec8f(p), nil
}

// Pow returns the Dec8 to the exp power (exp specified as a Dec8)
func (a Dec8) Pow8Safe(exp Dec8) (Dec8, error) {
	return a.PowSafe(exp.Float())
}

// Scan implements the sql.Scanner interface for database deserialization.
func (a *Dec8) Scan(value interface{}) error {
	// first try to see if the data is stored in database as a Numeric datatype
	switch v := value.(type) {

	case float32:
		*a = NewDec8f(float64(v))
		return nil

	case float64:
		*a = NewDec8f(v)
		return nil

	case int64:
		*a = Dec8{Val: v}
		return nil

	default:
		// default is trying to interpret value stored as string
		str, err := scanAsString(v)
		if err != nil {
			return err
		}
		temp, err := NewDec8s(str)
		if err != nil {
			return err
		}
		*a = temp
		return err
	}
}

// Value implements the driver.Valuer interface for database serialization.
func (a Dec8) Value() (driver.Value, error) {
	return a.String(), nil
}

// StringHuman - String for human reading
func (a Dec8) StringHuman() string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%.8f", a.Float())
}

// Dec8Slice is a slice of Dec8s for the purpose of sorting. It is modeled
// after the standard library's sort.Float64Slice
type Dec8Slice []Dec8

func (a Dec8Slice) Len() int           { return len(a) }
func (a Dec8Slice) Less(i, j int) bool { return a[i].Val < a[j].Val }
func (a Dec8Slice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
