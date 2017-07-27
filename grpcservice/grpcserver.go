package grpcservice

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
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

// NewGRPCServer initializes the server struct offering a service
func NewGRPCServer(useTLS bool, certFile string, keyFile string, port int) *GRPCServer {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
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

// NewMutualGRPCServer initializes the server struct including mutual tls auth for offering a service
func NewMutualGRPCServer(useTLS bool, certFile string, keyFile string, caFile string, port int) *GRPCServer {
	var opts []grpc.ServerOption
	if useTLS {
		// load peer cert/key, ca cert
		peerCert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Printf("GRPC server: load peer cert/key error: %v", err)
			return nil
		}
		caCert, err := ioutil.ReadFile(caFile)
		if err != nil {
			log.Printf("GRPC server: read ca cert file error: %v", err)
			return nil
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		ta := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{peerCert},
			ClientCAs:    caCertPool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		})
		opts = []grpc.ServerOption{grpc.Creds(ta)}
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("GRPC server: listen on port %d error:%v", port, err)
		return nil
	}

	server := grpc.NewServer(opts...)
	return &GRPCServer{server: server, listener: listener, isRunning: false}
}

// Serve registers the server as grpc server and starts it with the given listener
func (grpcserver *GRPCServer) Serve() error {

	if grpcserver.server == nil {
		return fmt.Errorf("GRPC server: Is not initialized")
	}
	if grpcserver.IsRunning() == true {
		return fmt.Errorf("GRPC server: Instance is already running")
	}

	grpcserver.isRunning = true
	grpcserver.server.Serve(grpcserver.listener)
	return nil
}

//Stop the grpc server
func (grpcserver *GRPCServer) Stop() error {
	if grpcserver.server == nil {
		return fmt.Errorf("GRPC server: Is not initialized")
	}

	if grpcserver.IsRunning() == false {
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
