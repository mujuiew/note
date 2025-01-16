// Package null provides a way to handle null without using pointers for
// all of the golang types plus Dec2, Dec5 and Dec8
package null

import (
	"bytes"
	"encoding/json"
	"reflect"
)

// Nuller supports identifying values as nullable
type Nuller interface {
	Null() bool
}

// Zeroer supports identifying values as the go zero value
type Zeroer interface {
	Zero() bool
}

// Constants that are used by StripNullJSON
const minJSONSizeWithNull = 7

var (
	// null is a string of bytes containing null in UTF8
	null = []byte(":null")
	// quote contains a quote in a byte value in UTF8
	quote = []byte("\"")[0]
	// comma contains a comma in a byte value in UTF8
	comma = []byte(",")[0]
	// leftCurly contains a leftCurly in a byte value in UTF8
	leftCurly = []byte("{")[0]
	// rightCurly contains a rightCurly in a byte value in UTF8
	rightCurly = []byte("}")[0]
	// minimum size of json with a null = "":null

	// backslash contains a backslash in a byte value in UTF8
	backslash = []byte(`\`)[0]
)

// removeIndex is a structure that holds the byte slice start and end index of
// bytes that should be removed from a target byte slice.  It is used by
// StripNullJSON
type removeIndex struct {
	Start int
	End   int
}

// StripNullJSON scans through the passed in byte slice containing json
// and strips out null values and returns a byte slice.  The byte slice returned
// is the same byte slice if the passed in one is null or empty or if there were
// no instances of null.  Otherwise, it is a new slice.
func StripNullJSON(b []byte) []byte {
	stripJson, _ := StripNullJSON2(b)

	return stripJson
}

type ignoreIdx struct {
	start int
	end   int
}

func StripNullJSON2(b []byte) ([]byte, error) {

	ignoreIdxs := getIgnoreIdxs(b)
	if len(ignoreIdxs) == 0 {
		return b, nil
	}
	var buff bytes.Buffer
	buff.Grow(len(b))
	idx_start := 0
	for _, ignoreIdx := range ignoreIdxs {
		if isEndBracket(b[idx_start]) && buff.Len() > 0 {
			if (buff.Bytes())[buff.Len()-1] == ',' {
				buff.Truncate(buff.Len() - 1)
			}
		}
		if idx_start >= ignoreIdx.start {
			idx_start = ignoreIdx.end + 1
			continue
		}
		_, err := buff.Write(b[idx_start:ignoreIdx.start])
		if err != nil {
			return b, err
		}
		idx_start = ignoreIdx.end + 1
	}
	if idx_start < len(b) {
		if isEndBracket(b[idx_start]) && buff.Len() > 0 {
			if (buff.Bytes())[buff.Len()-1] == ',' {
				buff.Truncate(buff.Len() - 1)
			}
		}
		_, err := buff.Write(b[idx_start:])
		if err != nil {
			return b, err
		}
	}
	return buff.Bytes(), nil
}

func getIgnoreIdxs(b []byte) []ignoreIdx {
	ignoreIdxs := make([]ignoreIdx, 0, 20)
	startTrimIdx := -1
	for i := 0; i < len(b)-len(null); i++ {
		if isStartJsonFieldValue(b[i]) {
			startTrimIdx = i
		}
		if bytes.Compare(null, b[i:i+len(null)]) == 0 {
			startIdx, endIdx := startTrimIdx, -1
			for j := i + len(null); j < len(b); j++ {
				if isEndJsonFieldValue(b[j]) {
					endIdx = j
					break
				}
			}
			if startIdx >= 0 && endIdx >= 0 {
				if b[startIdx] == ',' && b[endIdx] == ',' {
					startIdx++ // ,"a":null, will ignore "a":null, result should be ,
				} else if isStartBracket(b[startIdx]) && b[endIdx] == ',' {
					startIdx++ // {"a":null, will ignore "a":null, result should be {
				} else if b[startIdx] == ',' && isEndBracket(b[endIdx]) {
					endIdx-- // ,"a":null} will ignore ,"a":null result should be }
				} else if isStartBracket(b[startIdx]) && isEndBracket(b[endIdx]) {
					startIdx++ // {"a":null} will ignore "a":null result should be {}
					endIdx--
				}
				ignoreIdxs = append(ignoreIdxs, ignoreIdx{startIdx, endIdx})
			}
		}
	}
	return ignoreIdxs
}

func isStartJsonFieldValue(b byte) bool {
	return isStartBracket(b) || b == ','
}

func isEndJsonFieldValue(b byte) bool {
	return isEndBracket(b) || b == ','
}

func isEndBracket(b byte) bool {
	return b == '}' || b == ']'
}

func isStartBracket(b byte) bool {
	return b == '{' || b == '['
}

func StripNullJSONNew(b []byte) []byte {
	return []byte(trimObjectString(string(b)))
}
func trimObjectString(object string) string {
	var mappingString map[string]interface{}
	var arrayString []interface{}
	value := reflect.ValueOf(object)

	json.Unmarshal([]byte(value.String()), &mappingString)
	if len(mappingString) != 0 {
		st, err := json.Marshal(trimMap(mappingString))
		if err != nil {
			return ""
		}
		return string(st)
	}
	json.Unmarshal([]byte(value.String()), &arrayString)
	if len(arrayString) != 0 {
		st, err := json.Marshal(trimSlice(arrayString))
		if err != nil {
			return ""
		}
		return string(st)
	}
	if value.String() != "" {
		return object
	}
	return ""
}
func trimMap(object map[string]interface{}) map[string]interface{} {
	for key := range object {
		typ := reflect.TypeOf(object[key])
		if object[key] == nil {
			delete(object, key)
			continue
		}
		switch typ.Kind().String() {
		case "string":
			object[key] = trimObjectString(object[key].(string))
		case "map":
			object[key] = trimMap(object[key].(map[string]interface{}))
		case "slice":
			object[key] = trimSlice(object[key].([]interface{}))
		}
		if object[key] == nil {
			delete(object, key)
		}
	}
	return object
}
func removeArr(s []interface{}, key int) []interface{} {
	s[key] = s[len(s)-1]
	return s[:len(s)-1]
}
func trimSlice(object []interface{}) []interface{} {
	for key := range object {
		typ := reflect.TypeOf(object[key])
		switch typ.Kind().String() {
		case "string":
			object[key] = trimObjectString(object[key].(string))
		case "map":
			object[key] = trimMap(object[key].(map[string]interface{}))
		case "slice":
			object[key] = trimSlice(object[key].([]interface{}))
		}
		if object[key] == nil {
			removeArr(object, key)
		}
	}
	return object
}

func stripLBracketComma(b []byte, index int, upper int) ([]byte, bool) {
	var res []byte
	found := false
	for i := index; i < upper; i++ {
		if b[i] == leftCurly && b[i+1] == comma {
			res = append(b[:i+1], b[i+2:]...)
			found = true
			break
		}
	}
	if !found {
		res = b
	}
	return res, found
}
