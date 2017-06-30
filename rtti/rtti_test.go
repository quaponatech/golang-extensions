package rtti

import (
	"testing"
)

func TestSuccessGetFunction(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		fun, ptr := GetFunction(0)
		if fun == nil || ptr == 0 {
			t.Fail()
		}
	})
}

func TestFailureGetFunction(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		fun, ptr := GetFunction(-1)
		if fun != nil && ptr != 0 {
			t.Fail()
		}
	})
}

func TestSuccessGetCurrentFunctionName(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		name := GetCurrentFunctionName(false)
		if len(name) == 0 {
			t.Fail()
		}
	})
}

func TestSuccessGetSpecificFunctionName(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		name := GetSpecificFunctionName(1, false)
		if name != "testing.tRunner" {
			t.Fail()
		}
		GetSpecificFunctionName(1, true)
	})
}

func TestFailureGetSpecificFunctionName(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		name := GetSpecificFunctionName(100, false)
		if name != "" {
			t.Fail()
		}
	})
}

func TestSuccessGetFunctionName(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		name := GetFunctionName(t.Skipped)
		if name != "testing.(*common).Skipped-fm" {
			t.Fail()
		}
	})
}

func TestFailureGetFunctionName(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		name := GetFunctionName(nil)
		if name != "" {
			t.Fail()
		}
	})
}
