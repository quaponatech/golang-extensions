package server

import (
	"testing"

	"github.com/quaponatech/golang-extensions/test"
)

func TestSuccessNewLogger(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		test.AssertThat(t, New(), nil, "not")

		logger := NewLogger(
			"serverName",
			"",
			"",
			make(chan Status),
			make(chan error),
			make(chan string),
			make(chan string),
			make(chan string),
			0)
		test.AssertThat(t, logger.StartLogger(), nil)
		logger.StopLogger()
		logger.WaitGroup.Wait()

		logger = NewLogger(
			"serverName",
			"test",
			"test",
			nil,
			nil,
			nil,
			nil,
			nil,
			0)
		test.AssertThat(t, logger, nil, "not")
		test.AssertThat(t, logger.StartLogger(), nil)
		logger.StopLogger()
		logger.WaitGroup.Wait()
	})
}

func TestFailureNewLogger(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		logger := NewLogger(
			"",
			"",
			"",
			nil,
			nil,
			nil,
			nil,
			nil,
			0)
		test.AssertThat(t, logger, nil, "not")

		logger = NewLogger(
			"name",
			"/this_dir_is_not_writeable",
			"/this_dir_is_not_writeable",
			nil,
			nil,
			nil,
			nil,
			nil,
			0)
		test.AssertThat(t, logger, nil, "not")
	})
}
