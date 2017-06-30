package grpcservice

import (
	"fmt"
	"log"

	"github.com/quaponatech/golang-extensions/server"
)

// GRPCWebService defines anything necessary to setup, run and stop a general grpc server
type GRPCWebService struct {
	Prefix        string
	isInitialized bool
	isRunning     bool
	*GRPCWebServer
	*server.Logger

	// StopChannel is the channel the main waits for to quit the server execution
	StopChannel chan bool
}

//Setup the service
func (g *GRPCWebService) Setup(serverName string, grpcServer *GRPCWebServer,
	serverLogger *server.Logger, stopChan chan bool) error {

	if g.isRunning {
		return fmt.Errorf("Service already running")
	}

	if serverName == "" {
		return fmt.Errorf("Empty server name")
	}

	logPrefix := serverName + " - " + "[SERVICE] "
	log.Println(logPrefix + "Setting up Service")
	g.Prefix = logPrefix

	if nil == stopChan {
		return fmt.Errorf("Stop channel not initialized")
	}
	g.StopChannel = stopChan
	if nil == grpcServer || !grpcServer.IsInitialized() {
		return fmt.Errorf("GRPCWeb server not initialized")
	}
	g.GRPCWebServer = grpcServer
	if nil == serverLogger {
		return fmt.Errorf("Server logger not initialized")
	}
	g.Logger = serverLogger
	if err := g.StartLogger(); nil != err {
		return fmt.Errorf("Server logger failed on start: %v", err)
	}

	g.LogChan <- "Initialized Server"
	g.StatusChan <- server.StateInitialized
	g.isInitialized = true

	return nil
}

//Serve the service
func (g *GRPCWebService) Serve() error {
	if nil == g.GRPCWebServer ||
		!g.GRPCWebServer.IsInitialized() || !g.isInitialized {
		return fmt.Errorf("Service not initialized")
	}
	if g.isRunning {
		return fmt.Errorf("Service already running")
	}
	go func() {
		if err := g.GRPCWebServer.Serve(); nil != err {
			g.ErrorChan <- err
			g.StatusChan <- server.StateError
			g.Stop()
			return
		}
	}()
	g.StatusChan <- server.StateRunning
	g.isRunning = true

	for {
		stopped, ok := <-g.StopChannel
		if !ok {
			log.Println(g.Prefix + "Server shutdown unexpectedly")
			break
		}
		if stopped {
			break
		}
	}
	g.Logger.WaitGroup.Wait()
	return nil
}

//Stop the service
func (g *GRPCWebService) Stop() error {
	if !g.isInitialized {
		return fmt.Errorf("Service not initialized")
	}

	g.WarningChan <- "Shutting down"

	g.StatusChan <- server.StateStopping
	g.LogChan <- "Stopping GRPCWeb Server"
	if nil == g.GRPCWebServer || g.isRunning {
		g.GRPCWebServer.Stop()
		g.GRPCWebServer = nil
		g.isRunning = false

		g.StopChannel <- true
		close(g.StopChannel)
	}
	g.StatusChan <- server.StateStopped

	g.LogChan <- "Shutdown Log Environment"
	g.isInitialized = false
	g.StopLogger()

	log.Println(g.Prefix + "Shutted down")
	return nil
}
