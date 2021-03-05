package main

import (
	"context"
	"fmt"
	"log"
	"net"

	greetpb "github.com/stephenjlovell/grpc-go-example/api/go/pkg/greetpb"

	"google.golang.org/grpc"
)

const (
	LISTEN_ADDRESS = "0.0.0.0:50051"
)

// GreetServer is a placeholder for where our server logic would reside.
type GreetServer struct {
	greetpb.UnimplementedGreetServiceServer
}

// Greet generates a response to the rpc call
func (s *GreetServer) Greet(ctx context.Context, pb *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	firstName := pb.GetGreeting().GetFirstName()
	lastName := pb.GetGreeting().GetLastName()
	response := "Hello " + firstName + " " + lastName + "!"
	return &greetpb.GreetResponse{
		Response: response,
	}, nil
}

func main() {
	fmt.Println("greetings!")
	// listen on the default port for gRPC
	listener, err := net.Listen("tcp", LISTEN_ADDRESS)
	if err != nil {
		log.Fatalln("unable to connect to port")
	}

	grpcServer := grpc.NewServer()

	greetpb.RegisterGreetServiceServer(grpcServer, &GreetServer{})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
