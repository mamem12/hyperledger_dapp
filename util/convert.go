package util

import (
	"errors"
	"reflect"
	"strconv"
)

type source interface{}

func Convert(any interface{}) (source, error) {

	tt := reflect.TypeOf(any).Name()

	if tt == reflect.String.String() {

		t := any.(string)
		result, _ := strconv.Atoi(t)

		if result == 0 {
			return nil, errors.New("not positive arg from convert")
		}

		return result, nil

	} else if tt == reflect.Int.String() {

		t := any.(int)
		result := strconv.Itoa(t)

		return result, nil

	} else {
		return nil, errors.New("failed to chk type")
	}
}
