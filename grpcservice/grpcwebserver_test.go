package grpcservice_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/quaponatech/golang-extensions/grpcservice"
	"github.com/quaponatech/golang-extensions/test"
)

/*
 * GRPCWeb SERVICE unit test suite
 */

func TestSuiteGRPCWebServer(t *testing.T) {
	t.Run("CreateServerFailsOnWrongPort", func(t *testing.T) {
		tempServer := grpcservice.NewGRPCWebServer(false, "", "", 70000)
		err := tempServer.Serve()
		test.AssertThat(t, err, nil, "not")
	})

	t.Run("CreateServerFailsOnWrongTLS", func(t *testing.T) {
		portCounter++
		tempServer := grpcservice.NewGRPCWebServer(true, "", "", mainPort+portCounter)
		err := tempServer.Serve()
		test.AssertThat(t, err, nil, "not")
	})

	t.Run("ServeWithoutSetupFails", func(t *testing.T) {
		tempServer := &grpcservice.GRPCWebServer{}

		err := tempServer.Serve()
		test.AssertThat(t, err, fmt.Errorf("GRPCWeb server: Is not initialized"))
	})

	t.Run("ServingSucceeds", func(t *testing.T) {
		// SetUp
		portCounter++
		tempServer := grpcservice.NewGRPCWebServer(false, "", "", mainPort+portCounter)

		// Exercise + Verify
		go func() {
			err := tempServer.Serve()
			test.AssertThat(t, err, nil)
		}()
		time.Sleep(10 * time.Microsecond)

		// TearDown
		err := tempServer.Stop()
		test.AssertThat(t, err, nil)
	})

	t.Run("RunningTwiceIsNotPossbile", func(t *testing.T) {
		// SetUp
		portCounter++
		tempServer := grpcservice.NewGRPCWebServer(false, "", "", mainPort+portCounter)

		go func() {
			err := tempServer.Serve()
			test.AssertThat(t, err, nil)
		}()
		time.Sleep(10 * time.Microsecond)

		// Exercise + Verify
		go func() {
			err := tempServer.Serve()
			test.AssertThat(t, err, "GRPCWeb server: Instance is already running", "streq")
		}()
		time.Sleep(10 * time.Microsecond)

		// TearDown
		err := tempServer.Stop()
		test.AssertThat(t, err, nil)
	})

	t.Run("StopFailsWhenServerIsNotInitialized", func(t *testing.T) {
		// Setup
		portCounter++
		tempServer := &grpcservice.GRPCWebServer{}

		// Exercise + Verify
		err := tempServer.Stop()
		test.AssertThat(t, err, "GRPCWeb server: Is not initialized", "streq")
	})

	t.Run("StopFailsWhenServerIsNotRunning", func(t *testing.T) {
		// Setup
		portCounter++
		tempServer := grpcservice.NewGRPCWebServer(false, "", "", mainPort+portCounter)

		// Exercise + Verify
		err := tempServer.Stop()
		test.AssertThat(t, err, "GRPCWeb server: Is not running", "streq")
	})
}
