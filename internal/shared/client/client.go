package client

import (
	"log"

	"github.com/stephenjlovell/grpc-go-example/internal/calc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	caCertFile = "ssl/ca.crt"
)

func Connect() *grpc.ClientConn {
	creds, sslErr := credentials.NewClientTLSFromFile(caCertFile, "")
	if sslErr != nil {
		log.Fatalf("Failed to load CA trust certificate: %v", sslErr)
	}
	cc, err := grpc.Dial(calc.ListenAddress, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v\n", err)
	}
	return cc
}
