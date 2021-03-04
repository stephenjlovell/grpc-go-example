package main

import (
	"fmt"
	"log"
	"net"

	greetpb "github.com/stephenjlovell/grpc-go-example/api/go/pkg"

	"google.golang.org/grpc"
)

const (
	LISTEN_ADDRESS = "0.0.0.0:50051"
)

type Server struct {
}

// func (s *Server) Greet() {

// }

func main() {
	fmt.Println("greetings!")
	// listen on the default port for gRPC
	listener, err := net.Listen("tcp", LISTEN_ADDRESS)
	if err != nil {
		log.Fatalln("unable to connect to port")
	}

	grpc.NewServer()
	server := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(server, &Server{})

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
