package rtti

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
)

// GetFunction returns the name and uintptr for one of its calling function
// which is specified by a number:
// 0 for the directly calling function,
// increase for deeper calling functions
func GetFunction(callers int) (*runtime.Func, uintptr) {
	if callers < 0 {
		return nil, 0
	}
	pc := make([]uintptr, 10)
	runtime.Callers(2+callers, pc)
	return runtime.FuncForPC(pc[0]), pc[0]
}

// GetCurrentFunctionName returns a string of the calling
// function containing its file, its line and its name (If withLocation is true)
// or only its name (If withLocation is false):
// Callers:
// 0 for the directly calling function,
// increase for deeper calling functions.
func GetCurrentFunctionName(withLocation bool) string {
	return GetSpecificFunctionName(1, withLocation)
}

// GetSpecificFunctionName returns a string for the by "callers" number specified
// function containing its file, its line and its name (If withLocation is true)
// or only its name (If withLocation is false):
// Callers:
// 0 for the directly calling function,
// increase for deeper calling functions.
func GetSpecificFunctionName(callers int, withLocation bool) string {
	function, funcptr := GetFunction(callers + 1)
	if nil == function {
		return ""
	}

	if withLocation {
		file, line := function.FileLine(funcptr)
		return fmt.Sprintf("%s:%d %s", filepath.Base(file), line, function.Name())
	}
	return function.Name()
}

// GetFunctionName returns path and name of given function
func GetFunctionName(i interface{}) string {

	value := reflect.ValueOf(i)
	if !value.IsValid() {
		return ""
	}
	function := runtime.FuncForPC(value.Pointer())
	if nil == function {
		return ""
	}
	return function.Name()
}
