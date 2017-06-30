package grpcservice

import (
	"errors"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// ConnectionInfo describes the information necessary to connect to a grpc service
type ConnectionInfo struct {
	UseTLS              bool
	CertFile            string
	ServerHostName      string
	IP                  string
	Port                string
	TimeoutInMilliSecs  int
	RetryTimes          int
	RetryAfterMilliSecs int
}

// The GRPCClient is a struct defining all
//the contents needed to setup a connection to a grpc server.
type GRPCClient struct {
	connection *grpc.ClientConn
}

//GetConnection returns a pointer to the connection instance
func (grpcclient GRPCClient) GetConnection() *grpc.ClientConn {
	return grpcclient.connection
}

//Connect initializes a connection with a grpc server to use by a client.
func (grpcclient *GRPCClient) Connect(info *ConnectionInfo) error {

	log.Println("GRPC client: Initialize connection to grpc server")

	var opts []grpc.DialOption
	var err error

	log.Println("GRPC client: Setup connection options")
	if info.UseTLS {
		log.Println("GRPC client: Setup TLS connection")

		var creds credentials.TransportCredentials
		if info.CertFile != "" {
			creds, err = credentials.
				NewClientTLSFromFile(info.CertFile, "")
			if err != nil {
				log.Printf("GRPC client: Failed to create TLS credentials %v", err)
				return err
			}
		} else {
			creds = credentials.NewClientTLSFromCert(nil, "")
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	if 0 != info.TimeoutInMilliSecs {
		opts = append(opts, grpc.WithBlock())
		opts = append(opts, grpc.WithTimeout(
			time.Duration(info.TimeoutInMilliSecs)*time.Millisecond))
	}

	var retryAfterMilliSecs time.Duration
	if 0 >= info.RetryAfterMilliSecs {
		retryAfterMilliSecs = 1000
	} else {
		retryAfterMilliSecs = time.Duration(info.RetryAfterMilliSecs)
	}
	log.Println("GRPC client: Connect to server")
	var i int
	for i = 0; ; i++ {
		grpcclient.connection, err = grpc.Dial(info.IP+":"+info.Port, opts...)
		if err == nil {
			return nil
		}

		if i >= info.RetryTimes {
			break
		}
		time.Sleep(time.Duration(retryAfterMilliSecs) * time.Millisecond)
		log.Println("GRPC client: Retrying after error: ", err)
	}
	return fmt.Errorf(
		"GRPC client: After %d attempts, last error: %s", i, err)
}

//Close the grpc client connection
func (grpcclient GRPCClient) Close() error {
	if nil == grpcclient.connection {
		return errors.New("GRPC client: Is not initialized")
	}

	return grpcclient.connection.Close()
}
