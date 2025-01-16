package null

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func unmarshalTypeError(rawInput []byte, inputType any, expectType any, externalErr error) error {
	// In case has external error message
	if externalErr != nil {
		return &json.UnmarshalTypeError{
			Value: fmt.Sprintf("[convert-fn-error: %s] (%s)(%s)", externalErr.Error(), rawInput, reflect.TypeOf(inputType).String()),
			Type:  reflect.TypeOf(expectType),
		}
	}

	return &json.UnmarshalTypeError{
		Value: fmt.Sprintf("(%s)(%s)", rawInput, reflect.TypeOf(inputType).String()),
		Type:  reflect.TypeOf(expectType),
	}
}
