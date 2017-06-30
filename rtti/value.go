/*
	Reflecting consistency of given interface objects
	- by testing them for their underlying value
	- returning numbered error code and string in any case
*/

package rtti

import (
	"fmt"
	"reflect"
)

// Golang enum
const (
	NoUnderlyingValue    = 0
	NilUnderlyingValue   = 1
	ValidUnderlyingValue = 2
)

// TestLogForUnderlyingNilObject tests for ...
// ... empty and nilled underlying values/objects of a given interface and
// ... returns message and comparable value about given interface type.
func TestLogForUnderlyingNilObject(i interface{}) (int, string) {
	if i == reflect.Zero(reflect.TypeOf(i)).Interface() {
		return NilUnderlyingValue,
			fmt.Sprint("Has underlying nil or empty value")
	}
	return ValidUnderlyingValue,
		fmt.Sprint("Has underlying value")
}

// TestLogForUnderlyingValue tests for ...
// ... nilled types with no underlying value/object and
// ... delegates test to testLogForUnderlyingNilObject,
// if resolvable, non nil, non empty type.
// ... returns message and comparable value about given interface type.
func TestLogForUnderlyingValue(i interface{}) (int, string) {
	if i != nil {
		return TestLogForUnderlyingNilObject(i)
	}
	return NoUnderlyingValue,
		fmt.Sprintf("Has no underlying value")
}
