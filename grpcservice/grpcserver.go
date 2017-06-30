package grpcservice

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// The GRPCServer is a struct defining all the contents needed to setup a grpc server.
type GRPCServer struct {
	server    *grpc.Server
	listener  net.Listener
	isRunning bool
}

//NewGRPCServer initializes the server struct offering a service
func NewGRPCServer(
	useTLS bool, certFile string, keyFile string, port int) *GRPCServer {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if nil != err {
		log.Println(err)
		return nil
	}
	log.Print("GRPC server: Listening on port ", port)

	var opts []grpc.ServerOption

	if useTLS {
		log.Print("GRPC server: Prepare server options (with TLS)")
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			log.Println("GRPC server: Failed to generate credentials")
			listener.Close()
			return nil
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	} else {
		log.Print("GRPC server: Prepare server options (without TLS)")
	}

	log.Print("GRPC server: Creating new RPC server")
	server := grpc.NewServer(opts...)

	return &GRPCServer{server: server, listener: listener, isRunning: false}
}

// Serve registers the server as grpc server and starts it with the given listener
func (grpcserver *GRPCServer) Serve() error {

	if nil == grpcserver.server {
		return fmt.Errorf("GRPC server: Is not initialized")
	}
	if true == grpcserver.IsRunning() {
		return fmt.Errorf("GRPC server: Instance is already running")
	}

	grpcserver.isRunning = true
	grpcserver.server.Serve(grpcserver.listener)
	return nil
}

//Stop the grpc server
func (grpcserver *GRPCServer) Stop() error {
	if nil == grpcserver.server {
		return fmt.Errorf("GRPC server: Is not initialized")
	}

	if false == grpcserver.IsRunning() {
		return fmt.Errorf("GRPC server: Is not running")
	}

	grpcserver.server.Stop()
	grpcserver.listener.Close()
	grpcserver.server = nil
	grpcserver.listener = nil
	grpcserver.isRunning = false

	return nil
}

//IsRunning indicates if the server started listening properly
func (grpcserver GRPCServer) IsRunning() bool {
	return grpcserver.isRunning
}

//IsInitialized indicates if the server was initialized properly
func (grpcserver GRPCServer) IsInitialized() bool {
	return (nil != grpcserver.server)
}

//GetInstance returns a pointer to server instance
func (grpcserver GRPCServer) GetInstance() *grpc.Server {
	return grpcserver.server
}
