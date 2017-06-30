package server

import (
	"testing"

	"github.com/quaponatech/golang-extensions/test"
)

func TestSuccessString(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		test.AssertThat(t, StateUndefined.String(), "StateUndefined")
		test.AssertThat(t, StateRunning.String(), "StateRunning")
		test.AssertThat(t, StateError.String(), "StateError")
	})
}
