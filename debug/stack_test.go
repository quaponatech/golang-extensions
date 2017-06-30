package debug

import (
	"os"
	"testing"
)

func TestSuccessPrintStackTrace(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		out := PrintStackTrace(os.Stdout)
		if out == 0 {
			t.Fail()
		}
	})
}

func TestFailureNilPrintStackTrace(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		out := PrintStackTrace(nil)
		if out != 0 {
			t.Fail()
		}
	})
}

func TestFailureWriterInvalidPrintStackTrace(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		f, _ := os.Open("/this_file_should_not_exist")
		out := PrintStackTrace(f)
		if out != 0 {
			t.Fail()
		}
	})
}

func TestSuccessPrettyStackTraceString(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		o1 := PrettyStackTraceString(100)
		if o1 == "" {
			t.Fail()
		}

		o2 := PrettyStackTraceString(0)
		if o2 != "" {
			t.Fail()
		}
	})
}

func TestSuccessPrintCalledFunc(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		PrintCalledFunc(TestSuccessPrintCalledFunc)
	})
}
