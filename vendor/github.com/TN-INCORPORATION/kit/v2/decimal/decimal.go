// Package decimal values to aviod floating issue
// with 2, 5 and 8 digits of precision.
package decimal

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"unicode"
)

// Dec2 represents a 14.2 digit fixed point of a currency
// amount able to represent exact values with a range of
// -92,233,720,368,547.75 to 92,233,720,368,547.75.  The 14.2 can be
// thought of as the Whole being in the first 30 digits and the Fractional (Frac)
// in the last 2 digits ranging from .00 to 0.99.
// swagger:type float64
type Dec2 struct {
	Val int64
}

// Dec5 represents a 11.5 digit fixed point of a currency
// amount able to represent exact values with a range of
// -92,233,720,368,547.75807 to 92,233,720,368,547.75806.  The 11.5 can be
// thought of as the Whole being in the first 30 digits and the Fractional (Frac)
// in the last 5 digits ranging from .00000 to 0.99999.
// swagger:type float64
type Dec5 struct {
	Val int64
}

// Dec8 represents a 8.8 digit fixed point of a currency
// amount able to represent exact values with a range of
// -92,233,720,368.54775807 to 92,233,720,368.54775806.  The 8.8 can be
// thought of as the Whole being in the first 30 digits and the Fractional (Frac)
// in the last 8 digits ranging from .00000000 to 0.99999999.
// swagger:type float64
type Dec8 struct {
	Val int64
}

// decimal constants to determine fixed point precision and rounding locations.
const (
	// Dec2Frac specifies that there are 2 digits of fractional precision
	Dec2Frac = 2
	// Dec2Base specifies the decimal point location to be at two decimal points.
	Dec2Base = 100
	// Dec5Frac specifies that there are 5 digits of fractional precision
	Dec5Frac = 5
	// Dec5Base specifies the decimal point location to be at five decimal points.
	Dec5Base = 100000
	// Dec8Frac specifies that there are 8 digits of fractional precision
	Dec8Frac = 8
	// Dec8Base specifies the decimal point location to be at eight decimal points.
	Dec8Base = 100000000
	// roundfrac determines around which part of the fraction to round.  It is set
	// half way per the normal rounding rules.
	roundfrac = 0.5
	// negative sign character
	negSign = "-"
	// Period character
	period = "."
	// whole part, fractional part, both
	whole, fractional, both = 0, 1, 2
	// Dec2Dec5Rounder is an offset added to a value to round correctly from Dec5 to Dec2
	Dec2Dec5Rounder = 500
	// Dec2Dec5Divisor is a dividend used in rounding Dec5 to Dec2
	Dec2Dec5Divisor = 1000
	// Dec2Dec8Rounder is an offset added to a value to round correctly from Dec8 to Dec2
	Dec2Dec8Rounder = 500000
	// Dec2Dec8Divisor is a dividend used in rounding Dec8 to Dec2
	Dec2Dec8Divisor = 1000000
	// Dec5Dec2Multiplier is a multipler to change a Dec2 to a Dec5
	Dec5Dec2Multiplier = 1000
	// Dec5Dec8Rounder is an offset added to a value to round correctly from Dec8 to Dec5
	Dec5Dec8Rounder = 500
	// Dec5Dec8Divisor is a dividend used in rounding Dec8 to Dec5
	Dec5Dec8Divisor = 1000
	// Dec8Dec2Multiplier is a multipler to change a Dec2 to a Dec8
	Dec8Dec2Multiplier = 1000000
	// Dec8Dec5Multiplier is a multipler to change a Dec5 to a Dec8
	Dec8Dec5Multiplier = 1000
	// Dec2iMultiplier is a multiplier to change an integer to a Dec2
	Dec2iMultiplier = 100
	// Dec5iMultiplier is a multiplier to change an integer to a Dec5
	Dec5iMultiplier = 100000
	// Dec8iMultiplier is a multiplier to change an integer to a Dec8
	Dec8iMultiplier = 100000000
)

var (
	// Dec2Zero is the zero value of a Dec2
	Dec2Zero = Dec2{}
	// Dec5Zero is the zero value of a Dec2
	Dec5Zero = Dec5{}
	// Dec8Zero is the zero value of a Dec2
	Dec8Zero = Dec8{}
	// MaxDec2 represents the larget positive Dec2.
	MaxDec2 = Dec2{Val: math.MaxInt64 / 1000}
	// MinDec2 represents the smallest negative Dec2.
	MinDec2 = Dec2{Val: math.MinInt64 / 1000}
	// MaxDec5 represents the larget positive Dec5.
	MaxDec5 = Dec5{Val: math.MaxInt64 / 1000}
	// MinDec5 represents the smallest negative Dec5.
	MinDec5 = Dec5{Val: math.MinInt64 / 1000}
	// MaxDec8 represents the larget positive Dec8.
	MaxDec8 = Dec8{Val: math.MaxInt64 / 1000}
	// MinDec8 represents the smallest negative Dec8.
	MinDec8 = Dec8{Val: math.MinInt64 / 1000}
)

var (
	// ErrDivByZero is returned if you try to divide by zero.  Panic is avoided.
	ErrDivByZero = errors.New("Divide by zero is not allowed")
)

// PartsFromString takes a string and the number of fractional digits
// and returns the whole and fractional part in the int64 slice
// whole is index 0 and fractional is in index 1 if it exists
func PartsFromString(s string, f int) ([]int64, error) {

	var ds []int64
	s = strings.TrimSpace(s)
	// Empty string is considered an error
	if len(s) == 0 {
		return ds, errors.New("Missing currency value")
	}
	err := checkForInvalidChars(s)
	if err != nil {
		return ds, err
	}
	neg := false
	if strings.HasPrefix(s, negSign) {
		neg = true
		s = strings.TrimPrefix(s, negSign)
		if len(s) == 0 {
			return ds, errors.New("Invalid currency value")
		}
	}
	// If period is last, trim it, even if more than one
	s = strings.TrimRight(s, period)
	// Determine how many . periods we have if any
	vs := strings.Split(s, ".")
	if len(vs) > 2 || len(vs) == 0 {
		return ds, errors.New("Invalid currency value")
	}
	// Verify each segment is a number
	for i, v := range vs {
		if len(v) == 0 {
			ds = append(ds, 0)
			continue
		}
		d, err := strconv.Atoi(v)
		d64 := int64(d)
		if err != nil {
			return ds, errors.New("Invalid currency value")
		}
		// If this is the second value, then it represents the fraction
		if i == fractional {
			// Fraction must be between 1 and f (number of fractional digits)
			l := len(v)
			if l > f {
				return ds, errors.New("Invalid Currency Precision")
			}
			// This really means we have something like .6 which
			// is not the full number of fractional digits. The
			// int should be multiplied by some factor of 10 to
			// get the appropriate value
			if l < f {
				d64 = int64(d) * pow10(f-l)
			}
		}
		ds = append(ds, d64)
	}
	if neg == true {
		if ds[whole] != 0 {
			ds[whole] = -ds[whole]
		} else if len(ds) == both {
			ds[fractional] = -ds[fractional]
		}
	}
	return ds, nil
}

// checkForInvalidChars checks to make sure that the string is composed of valid
// characters that can make up a string that could be converted to a number
func checkForInvalidChars(s string) error {
	// First check for any invalid characters
	for _, r := range s {
		// If it is not a digit or a decimal point or a negative or positive sign
		if !(unicode.IsDigit(r) || r == '.' || r == '-' || r == '+') {
			return errors.New("Not a valid number")
		}
	}
	// Check for multiple decimals
	if strings.Count(s, ".") > 1 {
		return errors.New("Not a valid number")
	}
	return nil
}

func pow10(pow int) int64 {
	base10 := int64(1)
	if pow <= 0 {
		return base10
	}
	for i := 0; i < pow; i++ {
		base10 = base10 * 10
	}
	return base10
}

// modulo returns a modulo that takes into consideration
// the sign of the value
func modulo(i int64, b int64) int64 {
	if i < 0 {
		i = -i
	}
	return i % b
}

// absolute is an int based absolute value
func absolute(i int64) int64 {
	if i >= 0 {
		return i
	}
	return -i
}

func quoRem(a int64, d int64, b int64) int64 {
	i1, i2, base := big.NewInt(a), big.NewInt(d), big.NewInt(b)
	i1.Mul(i1, base)
	var rem big.Int
	i1.QuoRem(i1, i2, &rem)
	// Double the modulus and see if it is greater than the divisor
	// If so, we want to round up
	r := absolute(rem.Int64()) << 1
	up := r >= absolute(d)
	// Get result and negative status
	val := i1.Int64()
	neg := val < 0
	// If we need to round, round up for positive, down for negative
	if up {
		if neg {
			val = val - 1
		} else {
			val = val + 1
		}
	}
	return val
}

// Clip the fractional part in i to c
func clip(i int64, c int64) int64 {
	r := i
	// Clip the fractional part
	if i > 0 && i > c {
		r = c
	} else if i < 0 && i < -c {
		r = -c
	}
	return r
}

// Database scan as a string value.  It will unquote if quoted.
func scanAsString(value interface{}) (string, error) {
	var bytes []byte

	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return "", fmt.Errorf("Could not convert value '%+v' to byte array of type '%T'",
			value, value)
	}

	// If the amount is quoted, strip the quotes
	if len(bytes) > 2 && bytes[0] == '"' && bytes[len(bytes)-1] == '"' {
		bytes = bytes[1 : len(bytes)-1]
	}
	return string(bytes), nil
}
