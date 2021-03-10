package shared

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	caCertFile    = "ssl/ca.crt"
	listenAddress = "localhost:50051"
	certFile      = "ssl/server.crt"
	keyFile       = "ssl/server.pem"
)

func Connect() *grpc.ClientConn {
	creds, sslErr := credentials.NewClientTLSFromFile(caCertFile, "")
	if sslErr != nil {
		log.Fatalf("Failed to load CA trust certificate: %v", sslErr)
	}
	cc, err := grpc.Dial(listenAddress, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v\n", err)
	}
	return cc
}

func GetCreds() grpc.ServerOption {
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalln("failed to load certificate from file")
	}
	return grpc.Creds(creds)
}

func ListenTo() net.Listener {
	// listen on the custom port for gRPC
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalln("unable to connect to port")
	}
	return listener
}
