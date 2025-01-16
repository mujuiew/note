package null

import (
	"bytes"
	"math"
)

// StripNullJSON scans through the passed in byte slice containing json
// and strips out null values and returns a byte slice.  The byte slice returned
// is the same byte slice if the passed in one is null or empty or if there were
// no instances of null.  Otherwise, it is a new slice.
func StripNullJSONDeprecated(b []byte) []byte {
	if nil == b || len(b) <= minJSONSizeWithNull {
		return b
	}
	// Hold all of the removals
	rems := make([]removeIndex, 0, 20)

	// number of json string deep inside json object
	depthOfJSONStr := 0

	// find all of the nulls
	for i := 0; i < len(b)-len(null); i++ {
		window := b[i : i+len(null)]

		if i-1 >= 0 && b[i] == leftCurly && b[i-1] == quote { //    "{
			depthOfJSONStr++
		}
		if (i+1 < len(b) && b[i] == rightCurly && b[i+1] == backslash) || //    }\
			(i+1 < len(b) && b[i] == rightCurly && b[i+1] == quote) { //    }"
			depthOfJSONStr--
		}

		// We found a null
		if bytes.Compare(window, null) == 0 {
			// start next window past the null
			quoteCount := 0
			// Now search backwards for the quoted identifier
			for j := i; j >= 0; j-- {
				// eat two quotes
				if b[j] == quote {
					quoteCount = quoteCount + 1
					if quoteCount == 2 {
						start := j
						end := i + len(null) - 1

						// If there was a backslash before the quote escaping it, because we are inside a json object embedded in a
						// string property, eat the backslash as leaving it means we'd have malformed json
						if j-1 >= 0 && b[j-1] == backslash {
							// calculate number of backslash in json string in case of many json string deep inside json
							backslashCount := int(math.Pow(2, float64(depthOfJSONStr))) - 1 // 2^n - 1
							start = j - backslashCount

							// check if a previous value is a comma, is it
							if j-(backslashCount+1) >= 0 && b[j-(backslashCount+1)] == comma {
								start = start - 1
							}
						}

						// If there is a comma before us, eat it
						if j-1 >= 0 && b[j-1] == comma {
							start = j - 1
						}
						// If there is a { before us, eat the comma after us if it is there
						if j-1 >= 0 && b[j-1] == leftCurly {
							if end+1 < len(b) && b[end+1] == comma {
								end = end + 1
							}
						}
						ri := removeIndex{Start: start, End: end}
						rems = append(rems, ri)
						break
					}
				}
			}
			i = i + len(null) - 1 // move i to the last [l] in [:][n][u][l][l] , and  plus i in for loop to continue next windows
		}
	}
	// We did not find any
	if len(rems) == 0 {
		return b
	}
	// Build up the stripped down json with nulls removed by looping through
	// the remove indexes and take all the data in between
	res := make([]byte, 0, len(b))
	index := 0
	for _, rem := range rems {
		// Check to see if the removals are contiguous and if so, then just
		// ignore adding anything in between them
		if index < rem.Start {
			i := rem.Start - 1
			for ; i >= 0 && i >= index; i-- {
				if b[i] == backslash {
					continue
				}
				if b[i] == comma {
					continue
				}
				if b[i] == leftCurly {
					break
				}
				break
			}
			rem.Start = i + 1

			// remove {, combinations
			if len(res) > 0 && res[len(res)-1] == leftCurly && b[index] == comma {
				index = index + 1
			}
			res = append(res, b[index:rem.Start]...)
		}
		index = rem.End + 1
	}
	res = append(res, b[index:len(b)]...)
	// final hack to remove {, if there are any.  It will loop until all are gone
	upper := len(res) - 3
	i := 0
	found := false
	for i < upper {
		res, found = stripLBracketComma(res, i, upper)
		if found {
			i = i + 1
		} else {
			break
		}
	}
	return res
}
