package main

import (
	"log"
	"net"

	calcpb "github.com/stephenjlovell/grpc-go-example/api/go/pkg/calcpb"
	calc "github.com/stephenjlovell/grpc-go-example/internal/calc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	CERT_FILE = "ssl/server.crt"
	KEY_FILE  = "ssl/server.pem"
)

func main() {
	listener := listenTo(calc.LISTEN_ADDRESS)
	grpcServer := grpc.NewServer(getCreds())

	calcpb.RegisterCalcServiceServer(grpcServer, &calc.Server{})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}

func getCreds() grpc.ServerOption {
	creds, err := credentials.NewServerTLSFromFile(CERT_FILE, KEY_FILE)
	if err != nil {
		log.Fatalln("failed to load certificate from file")
	}
	return grpc.Creds(creds)
}

func listenTo(address string) net.Listener {
	log.Println("calculating... beep boop...")
	// listen on the custom port for gRPC
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln("unable to connect to port")
	}
	return listener
}
