package hashing

import (
	"testing"

	"github.com/quaponatech/golang-extensions/test"
)

func TestSuccessBuildPasswordHash(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		out, err := BuildPasswordHash("username", "password", "salt1", "salt2")
		test.AssertThat(t, err, nil)
		test.AssertThat(t, out, "b7627008958dd3e311329bfe7be9245e1e1b656e669e8647de62d15f53754882")
	})
}

func TestFailureBuildPasswordHash(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		out, err := BuildPasswordHash("", "", "", "")
		test.AssertThat(t, out, "")
		test.AssertThat(t, err, nil, "not")
	})
}

func TestSuccessBuildHashWithOneSaltAndHmac(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		out, err := BuildHashWithOneSaltAndHmac("hmac", "salt")
		test.AssertThat(t, out, "0Op+9JpUO6+2GdpWkB2zLHVOt/DFOUKta3IGKjiSBm0=")
		test.AssertThat(t, err, nil)
	})
}

func TestFailureBuildHashWithOneSaltAndHmac(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		out, err := BuildHashWithOneSaltAndHmac("", "salt")
		test.AssertThat(t, out, "")
		test.AssertThat(t, err, nil, "not")
	})
}
