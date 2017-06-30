package test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	debugging "github.com/quaponatech/golang-extensions/debug"
	"github.com/quaponatech/golang-extensions/rtti"
)

func log(tb testing.TB, observed, expected interface{},
	msgObservedType, msgExpectedType string,
	debug bool, options ...string) {
	_, file, line, _ := runtime.Caller(3)
	msgPrefix := fmt.Sprintf("%s:%d ", filepath.Base(file), line)

	tb.Logf(msgPrefix+"Observed: %T%s", observed, msgObservedType)
	tb.Logf(msgPrefix+"Observed: %v", observed)
	tb.Logf(msgPrefix+"Expected: %v", expected)
	tb.Logf(msgPrefix+"Expected: %T%s", expected, msgExpectedType)
	tb.Logf(msgPrefix+"Options: %#v", options)

	if debug {
		tb.Log(debugging.PrettyStackTraceString(2048))
	}
}

func verify(tb testing.TB, observed, expected interface{},
	options ...string) bool {

	var invert, stringsContain, stringsEqual, untyped, debug bool
	var msgObservedType, msgExpectedType string

	for _, option := range options {
		switch {
		case "not" == option || "!" == option:
			invert = true
			// Result  Invert  | Return
			// true    false   | true
			// false   false   | false
			// true    true    | false
			// false   true    | true
			// Operation to use: Return := (Result != Invert)
		case "contains" == option || "c" == option:
			stringsContain = true
		case "streq" == option || "==" == option:
			stringsEqual = true
		case "untyped" == option:
			untyped = true
		case "debug" == option:
			debug = true
		default: // COMMENT
			_, file, line, _ := runtime.Caller(2)
			tb.Logf("%s:%d COMMENT: %s", filepath.Base(file), line, option)
		}
	}

	if untyped {
		if invert != (fmt.Sprint(observed) == fmt.Sprint(expected)) {
			return true
		}
		log(tb, observed, expected, msgObservedType, msgExpectedType, debug, options...)
		return false
	}

	observedValue, observedMessage := rtti.TestLogForUnderlyingValue(observed)
	expectedValue, expectedMessage := rtti.TestLogForUnderlyingValue(expected)
	if observedValue < rtti.ValidUnderlyingValue ||
		expectedValue < rtti.ValidUnderlyingValue {

		if invert != (observedValue == expectedValue) {
			return true
		}
		msgObservedType = msgObservedType + "  [" + observedMessage + "]"
		msgExpectedType = msgExpectedType + "  [" + expectedMessage + "]"
	} else {
		if stringsContain {
			if invert != strings.Contains(fmt.Sprint(observed), fmt.Sprint(expected)) {
				return true
			}
		} else if stringsEqual {
			if invert != (fmt.Sprint(observed) == fmt.Sprint(expected)) {
				return true
			}
		} else if invert != reflect.DeepEqual(expected, observed) {
			return true
		}
	}

	log(tb, observed, expected, msgObservedType, msgExpectedType, debug, options...)
	return false
}

// ExpectThat compares an actual type and its content
// with the given resulting one (expected).
// Detects empty and nilled underlying values/objects
// and gives advices about them per given actual/expected type.
// Detects nilled types with no underlying value/object
// and gives advices about them per given actual/expected type.
// Compares types with valid underlying value/object
// if they are comparable and decides about equality.
// If non comparable types, it complains by giving back values and types.
// Prints type dependent detailed result and error messages.
// Not test ending.
func ExpectThat(tb testing.TB, actual, expected interface{}, options ...string) bool {
	if verify(tb, actual, expected, options...) {
		return true
	}
	tb.Fail()
	return false
}

// AssertThat compares an expected type with the actual resulting one.
// Delegates to ExpectThat and reuses its logic.
// Test ending.
func AssertThat(tb testing.TB, actual, expected interface{}, options ...string) bool {
	if verify(tb, actual, expected, options...) {
		return true
	}
	tb.FailNow()
	return false
}
