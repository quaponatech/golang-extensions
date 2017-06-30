package rtti

import (
	"testing"
)

func TestSuccessTestLogForUnderlyingNilObject(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		out, _ := TestLogForUnderlyingNilObject("test")
		if out != ValidUnderlyingValue {
			t.Fail()
		}

		out, _ = TestLogForUnderlyingNilObject("")
		if out != NilUnderlyingValue {
			t.Fail()
		}
	})
}

func TestSuccessTestLogForUnderlyingValue(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		out, _ := TestLogForUnderlyingValue("test")
		if out != ValidUnderlyingValue {
			t.Fail()
		}

		out, _ = TestLogForUnderlyingValue(nil)
		if out != NoUnderlyingValue {
			t.Fail()
		}
	})
}
