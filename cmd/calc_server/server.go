package main

import (
	"fmt"
	"log"
	"net"

	calcpb "github.com/stephenjlovell/grpc-go-example/api/go/pkg/calcpb"
	calc "github.com/stephenjlovell/grpc-go-example/internal/calc_server"

	"google.golang.org/grpc"
)

const (
	LISTEN_ADDRESS = "0.0.0.0:50052"
)

func main() {
	listener := listenTo(LISTEN_ADDRESS)
	grpcServer := grpc.NewServer()

	calcpb.RegisterCalcServiceServer(grpcServer, &calc.CalcServer{})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}

func listenTo(address string) net.Listener {
	fmt.Println("calculating... beep boop...")
	// listen on the custom port for gRPC
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln("unable to connect to port")
	}
	return listener
}
