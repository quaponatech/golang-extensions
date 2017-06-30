package grpcservice_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/quaponatech/golang-extensions/grpcservice"
	"github.com/quaponatech/golang-extensions/test"
)

// GRPC SERVICE unit test suite

func TestSuiteGRPCServer(t *testing.T) {
	t.Run("CreateServerFailsOnWrongPort", func(t *testing.T) {
		var nilServer *grpcservice.GRPCServer
		tempServer := grpcservice.NewGRPCServer(false, "", "", 70000)
		test.AssertThat(t, tempServer, nilServer)
	})

	t.Run("CreateServerFailsOnWrongTLS", func(t *testing.T) {
		portCounter++
		var nilServer *grpcservice.GRPCServer
		tempServer := grpcservice.NewGRPCServer(true, "", "", mainPort+portCounter)
		test.AssertThat(t, tempServer, nilServer)
	})

	t.Run("ServeWithoutSetupFails", func(t *testing.T) {
		tempServer := &grpcservice.GRPCServer{}

		err := tempServer.Serve()
		test.AssertThat(t, err, fmt.Errorf("GRPC server: Is not initialized"))
	})

	t.Run("ServingSucceeds", func(t *testing.T) {
		// SetUp
		portCounter++
		tempServer := grpcservice.NewGRPCServer(false, "", "", mainPort+portCounter)

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
		tempServer := grpcservice.NewGRPCServer(false, "", "", mainPort+portCounter)

		go func() {
			err := tempServer.Serve()
			test.AssertThat(t, err, nil)
		}()
		time.Sleep(10 * time.Microsecond)

		// Exercise + Verify
		go func() {
			err := tempServer.Serve()
			test.AssertThat(t, err, "GRPC server: Instance is already running", "streq")
		}()
		time.Sleep(10 * time.Microsecond)

		// TearDown
		err := tempServer.Stop()
		test.AssertThat(t, err, nil)
	})

	t.Run("StopFailsWhenServerIsNotInitialized", func(t *testing.T) {
		// Setup
		portCounter++
		tempServer := &grpcservice.GRPCServer{}

		// Exercise + Verify
		err := tempServer.Stop()
		test.AssertThat(t, err, "GRPC server: Is not initialized", "streq")
	})

	t.Run("StopFailsWhenServerIsNotRunning", func(t *testing.T) {
		// Setup
		portCounter++
		tempServer := grpcservice.NewGRPCServer(false, "", "", mainPort+portCounter)

		// Exercise + Verify
		err := tempServer.Stop()
		test.AssertThat(t, err, "GRPC server: Is not running", "streq")
	})

	t.Run("GetInstanceReturnsCorrectServerInstance", func(t *testing.T) {
		// SetUp
		portCounter++
		tempServer := grpcservice.NewGRPCServer(false, "", "", mainPort+portCounter)

		// Exercise + Verify
		instance := tempServer.GetInstance()

		// TearDown
		test.AssertThat(t, grpcservice.GetInstanceFromServer(tempServer), instance)
	})
}
