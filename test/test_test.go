package test

import (
	"testing"
)

func TestSuccessAssertThat(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		AssertThat(t, nil, nil)
		AssertThat(t, true, false, "not")
		AssertThat(t, "string1", "stri", "contains")
		AssertThat(t, "str", "str", "streq")
		AssertThat(t, "", "", "untyped")
		AssertThat(t, 2, 2, "debug")
		AssertThat(t, 2, 2, "default")
	})
}
