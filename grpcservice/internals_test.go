package grpcservice

import "google.golang.org/grpc"

//GetInstanceFromServer returns a pointer to server instance
func GetInstanceFromServer(grpcserver *GRPCServer) *grpc.Server {
	return grpcserver.server
}

//GetInstanceFromClient returns a pointer to server instance
func GetInstanceFromClient(grpcclient *GRPCClient) *grpc.ClientConn {
	return grpcclient.connection
}
