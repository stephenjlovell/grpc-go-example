package main

import (
	"log"
	"net"

	calcpb "github.com/stephenjlovell/grpc-go-example/api/go/pkg/calcpb"
	calc "github.com/stephenjlovell/grpc-go-example/internal/calc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const (
	CertFile = "ssl/server.crt"
	KeyFile  = "ssl/server.pem"
)

func main() {
	listener := listenTo(calc.ListenAddress)
	grpcServer := grpc.NewServer(getCreds())

	calcpb.RegisterCalcServiceServer(grpcServer, &calc.Server{})

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}

func getCreds() grpc.ServerOption {
	creds, err := credentials.NewServerTLSFromFile(CertFile, KeyFile)
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
