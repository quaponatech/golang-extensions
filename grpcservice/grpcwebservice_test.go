package grpcservice_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/quaponatech/golang-extensions/grpcservice"
	"github.com/quaponatech/golang-extensions/server"
	"github.com/quaponatech/golang-extensions/test"
)

/*
 * GRPCWeb SERVICE unit test suite
 */

func TestSuiteGRPCWebService(t *testing.T) {
	var okName = "Test GRPCWeb Service"

	t.Run("SetupWrong", func(t *testing.T) {
		var (
			wrongName     = ""
			wrongChannel  chan bool
			wrongServer   *grpcservice.GRPCWebServer
			wrongLogger   *server.Logger
			failingLogger = server.NewLogger(okName, "/var/log/not-existing", ".file",
				make(chan server.Status), make(chan error), make(chan string),
				make(chan string), make(chan string), 0)
		)
		portCounter++
		var (
			okLogger  = server.New()
			okChannel = make(chan bool)
			okServer  = grpcservice.NewGRPCWebServer(false, "", "", mainPort+portCounter)
		)

		t.Run("WrongPrefix", func(t *testing.T) {
			tempService := new(grpcservice.GRPCWebService)
			err := tempService.Setup(wrongName, okServer, okLogger, okChannel)
			test.AssertThat(t, err, "Empty server name", "streq")
		})

		t.Run("WrongQuitChannel", func(t *testing.T) {
			tempService := new(grpcservice.GRPCWebService)
			err := tempService.Setup(okName, okServer, okLogger, wrongChannel)
			test.AssertThat(t, err, "Stop channel not initialized", "streq")
		})

		t.Run("WrongServer", func(t *testing.T) {
			tempService := new(grpcservice.GRPCWebService)
			err := tempService.Setup(okName, wrongServer, okLogger, okChannel)
			test.AssertThat(t, err, "GRPCWeb server not initialized", "streq")
		})

		t.Run("WrongLogger", func(t *testing.T) {
			tempService := new(grpcservice.GRPCWebService)
			err := tempService.Setup(okName, okServer, wrongLogger, okChannel)
			test.AssertThat(t, err, "Server logger not initialized", "streq")
		})

		t.Run("LoggerFailsOnStartWithInvalidFile", func(t *testing.T) {
			tempService := new(grpcservice.GRPCWebService)
			err := tempService.Setup(okName, okServer, failingLogger, okChannel)
			test.AssertThat(t, err, "Server logger failed on start: "+
				"Error: Creating log file path: "+
				"mkdir /var/log/not-existing: permission denied", "streq")
		})

		t.Run("StopFailsOnUninitializedService", func(t *testing.T) {
			tempService := new(grpcservice.GRPCWebService)
			err := tempService.Stop()
			test.AssertThat(t, err, "Service not initialized", "streq")
		})

		t.Run("ServeFailsOnUninitializedService", func(t *testing.T) {
			tempService := new(grpcservice.GRPCWebService)
			err := tempService.Serve()
			test.AssertThat(t, err, "Service not initialized", "streq")
		})

	})

	t.Run("SetupOK", func(t *testing.T) {

		t.Run("SetupSucceedsOnValidInputs", func(t *testing.T) {
			// Setup
			tempService := new(grpcservice.GRPCWebService)
			portCounter++
			var (
				okLogger = server.NewLogger(okName, "", "",
					make(chan server.Status), make(chan error), make(chan string),
					make(chan string), make(chan string), 0)
				okServer  = grpcservice.NewGRPCWebServer(false, "", "", mainPort+portCounter)
				okChannel = make(chan bool)
			)
			// Exercise + Verify
			err := tempService.Setup(okName, okServer, okLogger, okChannel)
			test.AssertThat(t, err, nil)

			// Teardown
			tempService.Stop()
		})

		t.Run("StopSucceedsOnNotRunningService", func(t *testing.T) {
			// Setup
			tempService := new(grpcservice.GRPCWebService)
			portCounter++
			err := tempService.Setup(okName,
				grpcservice.NewGRPCWebServer(false, "", "", mainPort+portCounter),
				server.NewLogger(okName, "", "",
					make(chan server.Status), make(chan error), make(chan string),
					make(chan string), make(chan string), 0),
				make(chan bool))
			test.AssertThat(t, err, nil)

			// Exercise + Verify
			tempService.Stop()
			test.AssertThat(t, err, nil)
		})

	})

	t.Run("ServingOK", func(t *testing.T) {

		t.Run("ServingSucceedsOnValidInputs", func(t *testing.T) {
			//t.Skipf("ServingSucceedsOnValidInputs: Needs to sleep a little.")
			// Setup
			tempService := new(grpcservice.GRPCWebService)
			portCounter++
			err := tempService.Setup(t.Name(),
				grpcservice.NewGRPCWebServer(false, "", "", mainPort+portCounter),
				server.NewLogger(t.Name(), "", "",
					make(chan server.Status), make(chan error), make(chan string),
					make(chan string), make(chan string), 0),
				make(chan bool))
			test.AssertThat(t, err, nil)

			// Exercise + Verify
			go func() {
				err := tempService.Serve()
				test.AssertThat(t, err, nil)
			}()
			time.Sleep(10 * time.Microsecond)

			// Teardown
			err = tempService.Stop()
			test.AssertThat(t, err, nil)
		})

		t.Run("ServingFailsOnGettingStartedTwice", func(t *testing.T) {
			//t.Skipf("ServingSucceedsOnValidInputs: Needs to sleep a little.")
			// Setup
			tempService := new(grpcservice.GRPCWebService)
			portCounter++
			err := tempService.Setup(t.Name(),
				grpcservice.NewGRPCWebServer(false, "", "", mainPort+portCounter),
				server.NewLogger(t.Name(), "", "",
					make(chan server.Status), make(chan error), make(chan string),
					make(chan string), make(chan string), 0),
				make(chan bool))
			test.AssertThat(t, err, nil)
			go func() {
				err := tempService.Serve()
				test.AssertThat(t, err, nil)
			}()
			time.Sleep(10 * time.Microsecond)

			// Exercise + Verify
			go func() {
				err := tempService.Serve()
				test.AssertThat(t, err, "Service already running", "streq")
			}()
			time.Sleep(10 * time.Microsecond)

			// Teardown
			err = tempService.Stop()
			test.AssertThat(t, err, nil)
		})

		t.Run("SetupFailsWhenAlreadyServing", func(t *testing.T) {
			// Setup
			tempService := new(grpcservice.GRPCWebService)
			portCounter++
			err := tempService.Setup(t.Name(),
				grpcservice.NewGRPCWebServer(false, "", "", mainPort+portCounter),
				server.NewLogger(t.Name(), "", "",
					make(chan server.Status), make(chan error), make(chan string),
					make(chan string), make(chan string), 0),
				make(chan bool))
			test.AssertThat(t, err, nil)

			go func() {
				err := tempService.Serve()
				test.AssertThat(t, err, nil)
			}()
			time.Sleep(100 * time.Microsecond)

			// Exercise + Verify
			err = tempService.Setup(okName, nil, nil, nil)
			test.AssertThat(t, err, "Service already running", "streq")

			// Teardown
			err = tempService.Stop()
			test.AssertThat(t, err, nil)
		})
	})

	t.Run("SetupFailsOnAlreadyStoppedGRPCWebServer", func(t *testing.T) {
		// Setup
		tempService := new(grpcservice.GRPCWebService)
		portCounter++
		tempServer := grpcservice.NewGRPCWebServer(false, "", "", mainPort+portCounter)

		// Exercise
		go func() {
			err := tempServer.Serve()
			test.AssertThat(t, err, nil)
		}()
		time.Sleep(10 * time.Microsecond)

		err := tempServer.Stop()
		test.AssertThat(t, err, nil)
		time.Sleep(10 * time.Microsecond)

		// Verify
		err = tempService.Setup(t.Name(),
			tempServer,
			server.NewLogger(t.Name(), "", "",
				make(chan server.Status), make(chan error), make(chan string),
				make(chan string), make(chan string), 0),
			make(chan bool))
		test.AssertThat(t, err, fmt.Errorf("GRPCWeb server not initialized"))

		go func() {
			err := tempService.Serve()
			test.AssertThat(t, err, "Service not initialized", "streq")
		}()
	})

	t.Run("StartServingFailsOnAlreadyStoppedGRPCWebServer", func(t *testing.T) {
		// Setup
		tempService := new(grpcservice.GRPCWebService)
		portCounter++
		tempServer := grpcservice.NewGRPCWebServer(false, "", "", mainPort+portCounter)

		// Exercise
		go func() {
			err := tempServer.Serve()
			test.AssertThat(t, err, nil)
		}()
		time.Sleep(10 * time.Microsecond)

		// Verify
		err := tempService.Setup(t.Name(),
			tempServer,
			server.NewLogger(t.Name(), "", "",
				make(chan server.Status), make(chan error), make(chan string),
				make(chan string), make(chan string), 0),
			make(chan bool))
		test.AssertThat(t, err, nil)

		err = tempServer.Stop()
		test.AssertThat(t, err, nil)
		time.Sleep(10 * time.Microsecond)

		go func() {
			err := tempService.Serve()
			test.AssertThat(t, err, "Service not initialized", "streq")
		}()
		time.Sleep(10 * time.Microsecond)

		// Teardown
		err = tempService.Stop()
		test.AssertThat(t, err, nil)
	})

	t.Run("AlreadyStartedServingFailsOnAlreadyStoppedGRPCWebServer", func(t *testing.T) {
		//t.Skipf("Does not work as intended - Passes but should fail")
		// Setup
		tempService := new(grpcservice.GRPCWebService)
		portCounter++
		tempServer := grpcservice.NewGRPCWebServer(false, "", "", mainPort+portCounter)

		err := tempService.Setup(t.Name(),
			tempServer,
			server.NewLogger(t.Name(), "", "",
				make(chan server.Status), make(chan error), make(chan string),
				make(chan string), make(chan string), 0),
			make(chan bool))
		test.AssertThat(t, err, nil)

		go func() {
			err := tempService.Serve()
			test.AssertThat(t, err, nil)
		}()
		time.Sleep(10 * time.Microsecond)

		// Exercise + Verify
		err = tempServer.Stop()
		test.AssertThat(t, err, nil)
		time.Sleep(10 * time.Microsecond)

		// Teardown
		err = tempService.Stop()
		test.AssertThat(t, err, nil)
	})
}
