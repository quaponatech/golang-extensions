package grpcservice

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
)

// The GRPCWebServer is a struct defining all the contents needed to setup a grpc server.
type GRPCWebServer struct {
	innerServer *grpc.Server
	server      *grpcweb.WrappedGrpcServer
	isRunning   bool
	useTLS      bool
	certFile    string
	keyFile     string
	port        int
}

// NewGRPCWebServer initializes a server struct offering a GRPCWeb Web service
func NewGRPCWebServer(useTLS bool, certFile string, keyFile string, port int) *GRPCWebServer {
	server := grpc.NewServer()
	webServer := grpcweb.WrapServer(server)
	log.Print("GRPC Web server: Listening on port ", port)
	if useTLS {
		log.Print("GRPC Web server: Preparing server (with TLS)")
	} else {
		log.Print("GRPC Web server: Preparing server (without TLS)")
	}
	return &GRPCWebServer{
		server:      webServer,
		innerServer: server,
		useTLS:      useTLS,
		certFile:    certFile,
		keyFile:     keyFile,
		port:        port,
	}
}

// Serve registers the server as grpc server
func (grpcserver *GRPCWebServer) Serve() error {

	if grpcserver.server == nil {
		return fmt.Errorf("GRPCWeb server: Is not initialized")
	}
	if grpcserver.IsRunning() {
		return fmt.Errorf("GRPCWeb server: Instance is already running")
	}

	grpcserver.isRunning = true
	grpcHandler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		// Allow origins
		if origin := req.Header.Get("Origin"); origin != "" {
			resp.Header().Set("Access-Control-Allow-Origin", origin)
			resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			resp.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, x-grpc-web")
		}
		// Stop here if its Preflighted OPTIONS request
		if req.Method == "OPTIONS" {
			return
		}

		if grpcserver.server.IsGrpcWebRequest(req) {
			// Answer with GRPC
			grpcserver.server.ServeHTTP(resp, req)
		} else {
			// Fall back to other servers
			http.DefaultServeMux.ServeHTTP(resp, req)
		}
	})
	var err error
	handler := handlers.LoggingHandler(os.Stdout, grpcHandler)

	if !grpcserver.useTLS {
		err = http.ListenAndServe(fmt.Sprintf(":%v", grpcserver.port), handler)
	} else {
		err = http.ListenAndServeTLS(fmt.Sprintf(":%v", grpcserver.port),
			grpcserver.certFile, grpcserver.keyFile, handler)
	}
	return err
}

// Stop the grpc server
func (grpcserver *GRPCWebServer) Stop() error {
	if grpcserver.server == nil {
		return fmt.Errorf("GRPCWeb server: Is not initialized")
	}

	if !grpcserver.IsRunning() {
		return fmt.Errorf("GRPCWeb server: Is not running")
	}

	grpcserver.innerServer.Stop()
	grpcserver.server = nil
	grpcserver.isRunning = false

	return nil
}

// IsRunning indicates if the server started listening properly
func (grpcserver GRPCWebServer) IsRunning() bool {
	return grpcserver.isRunning
}

// IsInitialized indicates if the server was initialized properly
func (grpcserver GRPCWebServer) IsInitialized() bool {
	return (grpcserver.server != nil)
}

// GetInstance returns a pointer to server instance
func (grpcserver GRPCWebServer) GetInstance() *grpcweb.WrappedGrpcServer {
	return grpcserver.server
}

// GetInnerInstance returns a pointer to server instance
func (grpcserver GRPCWebServer) GetInnerInstance() *grpc.Server {
	return grpcserver.innerServer
}
