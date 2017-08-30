package grpcservice_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/quaponatech/golang-extensions/grpcservice"
	"github.com/quaponatech/golang-extensions/test"
)

/* GRPC SERVICE unit test suite */
func TestSuiteGRPCClient(t *testing.T) {
	t.Run("CloseFailsIfNotInitialized", func(t *testing.T) {
		tempClient := new(grpcservice.GRPCClient)
		err := tempClient.Close()
		test.AssertThat(t, err, "GRPC client: Is not initialized", "streq")
	})

	t.Run("ConnectingFailsOnEmptyTLS", func(t *testing.T) {
		// SetUp
		portCounter++
		port := mainPort + portCounter
		tempServer := grpcservice.NewGRPCServer(false, "", "", port)

		go func() {
			err := tempServer.Serve()
			test.AssertThat(t, err, nil)
		}()
		time.Sleep(10 * time.Microsecond)

		// Exercise + Verify
		tempClient := new(grpcservice.GRPCClient)
		err := tempClient.Connect(&grpcservice.ConnectionInfo{true, "", "",
			"localhost", fmt.Sprint(port), 10, 1, 1, "", ""})
		test.AssertThat(t, err,
			"tls: first record does not look like a TLS handshake",
			"contains", err.Error())

		// TearDown
		err = tempServer.Stop()
		test.AssertThat(t, err, nil)
	})

	t.Run("ConnectingFailsOnWrongTLS", func(t *testing.T) {
		// SetUp
		portCounter++
		port := mainPort + portCounter
		tempServer := grpcservice.NewGRPCServer(false, "", "", port)

		go func() {
			err := tempServer.Serve()
			test.AssertThat(t, err, nil)
		}()
		time.Sleep(10 * time.Microsecond)

		// Exercise + Verify
		tempClient := new(grpcservice.GRPCClient)
		err := tempClient.Connect(&grpcservice.ConnectionInfo{true,
			"/var/log/not-existing", ".file",
			"localhost", fmt.Sprint(port), 10, 0, 0, "", ""})
		test.AssertThat(t,
			err, "open /var/log/not-existing: no such file or directory", "streq")

		// TearDown
		err = tempServer.Stop()
		test.AssertThat(t, err, nil)
	})

	t.Run("ConnectingFailsOnMissingServer", func(t *testing.T) {
		// SetUp
		tempClient := new(grpcservice.GRPCClient)
		portCounter++

		// Exercise + Verify
		err := tempClient.Connect(&grpcservice.ConnectionInfo{false, "", "",
			"localhost", fmt.Sprint(mainPort + portCounter), 10, 0, 0, "", ""})
		test.AssertThat(t, err, "context deadline exceeded",
			"contains", err.Error())
	})

	t.Run("ConnectingSucceeds", func(t *testing.T) {
		// SetUp
		portCounter++
		port := mainPort + portCounter
		tempServer := grpcservice.NewGRPCServer(false, "", "", port)

		go func() {
			err := tempServer.Serve()
			test.AssertThat(t, err, nil)
		}()
		time.Sleep(10 * time.Microsecond)

		// Exercise + Verify
		tempClient := new(grpcservice.GRPCClient)
		err := tempClient.Connect(&grpcservice.ConnectionInfo{false, "", "",
			"localhost", fmt.Sprint(port), 1, 0, 0, "", ""})
		test.AssertThat(t, err, nil)

		// TearDown
		err = tempClient.Close()
		test.AssertThat(t, err, nil)

		err = tempServer.Stop()
		test.AssertThat(t, err, nil)
	})

	t.Run("GetConnectionReturnsCorrectClientConnectionToServer", func(t *testing.T) {
		// SetUp
		portCounter++
		port := mainPort + portCounter
		tempServer := grpcservice.NewGRPCServer(false, "", "", port)

		go func() {
			err := tempServer.Serve()
			test.AssertThat(t, err, nil)
		}()
		time.Sleep(10 * time.Microsecond)

		tempClient := new(grpcservice.GRPCClient)
		err := tempClient.Connect(&grpcservice.ConnectionInfo{false, "", "",
			"localhost", fmt.Sprint(port), 1, 0, 0, "", ""})
		test.AssertThat(t, err, nil)

		// Exercise + Verify
		conn := tempClient.GetConnection()
		test.AssertThat(t, grpcservice.GetInstanceFromClient(tempClient), conn)

		// TearDown
		err = tempClient.Close()
		test.AssertThat(t, err, nil)

		err = tempServer.Stop()
		test.AssertThat(t, err, nil)
	})
}
