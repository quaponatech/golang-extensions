package grpcservice

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
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
	KeyFile             string
	CaFile              string
}

// The GRPCClient is a struct defining all
// the contents needed to setup a connection to a grpc server.
type GRPCClient struct {
	connection *grpc.ClientConn
}

// GetConnection returns a pointer to the connection instance
func (grpcclient GRPCClient) GetConnection() *grpc.ClientConn {
	return grpcclient.connection
}

// Connect initializes a connection with a grpc server to use by a client.
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
	if info.TimeoutInMilliSecs != 0 {
		opts = append(opts, grpc.WithBlock())
		opts = append(opts, grpc.WithTimeout(
			time.Duration(info.TimeoutInMilliSecs)*time.Millisecond))
	}

	var retryAfterMilliSecs time.Duration
	if info.RetryAfterMilliSecs <= 0 {
		retryAfterMilliSecs = 1000
	} else {
		retryAfterMilliSecs = time.Duration(info.RetryAfterMilliSecs)
	}

	return dial(grpcclient, info, retryAfterMilliSecs, opts)
}

// ConnectMutual initializes a connection with a grpc server with the usage of mutual tls auth
func (grpcclient *GRPCClient) ConnectMutual(info *ConnectionInfo) error {
	log.Println("GRPC client: Initialize connection to grpc server")
	var opts []grpc.DialOption
	log.Println("GRPC client: Setup connection options")
	if info.UseTLS {
		// load peer cert/key, cacert
		peerCert, certKeyErr := tls.LoadX509KeyPair(info.CertFile, info.KeyFile)
		if certKeyErr != nil {
			log.Printf("GRPC client: load peer cert/key error: %v", certKeyErr)
			return certKeyErr
		}
		caCert, readCertErr := ioutil.ReadFile(info.CaFile)
		if readCertErr != nil {
			log.Printf("GRPC client: read ca cert file error: %v", readCertErr)
			return readCertErr
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		ta := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{peerCert},
			RootCAs:      caCertPool,
		})
		opts = append(opts, grpc.WithTransportCredentials(ta))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	if info.TimeoutInMilliSecs != 0 {
		opts = append(opts, grpc.WithBlock())
		opts = append(opts, grpc.WithTimeout(
			time.Duration(info.TimeoutInMilliSecs)*time.Millisecond))
	}

	var retryAfterMilliSecs time.Duration
	if info.RetryAfterMilliSecs <= 0 {
		retryAfterMilliSecs = 1000
	} else {
		retryAfterMilliSecs = time.Duration(info.RetryAfterMilliSecs)
	}

	return dial(grpcclient, info, retryAfterMilliSecs, opts)
}

func dial(
	grpcclient *GRPCClient,
	info *ConnectionInfo,
	retryAfterMilliSecs time.Duration,
	opts []grpc.DialOption,
) error {
	var err error
	var i int
	for i = 0; ; i++ {
		grpcclient.connection, err = grpc.Dial(info.IP+":"+info.Port, opts...)
		if err == nil {
			return nil
		}

		if i >= info.RetryTimes {
			break
		}
		time.Sleep(retryAfterMilliSecs * time.Millisecond)
		log.Println("GRPC client: Retrying after error: ", err)
	}
	return fmt.Errorf("GRPC client: After %d attempts, last error: %s", i, err)
}

//Close the grpc client connection
func (grpcclient GRPCClient) Close() error {
	if nil == grpcclient.connection {
		return errors.New("GRPC client: Is not initialized")
	}

	return grpcclient.connection.Close()
}
