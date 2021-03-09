package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	greetpb "github.com/stephenjlovell/grpc-go-example/api/go/pkg/greetpb"

	"google.golang.org/grpc"
)

const (
	LISTEN_ADDRESS = "0.0.0.0:50051"
)

// GreetServer is a placeholder for where our server logic would reside.
type GreetServer struct {
	// this is awkward but necessary to provide guarantees about our interface to greetpb.RegisterGreetServiceServer
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

// GreetManyTimes implements a streaming service
func (s *GreetServer) GreetManyTimes(pb *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	firstName := pb.GetGreeting().GetFirstName()
	lastName := pb.GetGreeting().GetLastName()

	for i := 1; i <= 10; i++ {
		str := strconv.Itoa(i) + ": Hello " + firstName + " " + lastName + "!"
		resp := &greetpb.GreetManyTimesResponse{
			Response: str,
		}
		stream.Send(resp)
		time.Sleep(1 * time.Second)
	}
	return nil
}

func main() {
	listener := listenTo(LISTEN_ADDRESS)
	grpcServer := grpc.NewServer()

	greetpb.RegisterGreetServiceServer(grpcServer, &GreetServer{})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}

func listenTo(address string) net.Listener {
	fmt.Println("greetings!")
	// listen on the default port for gRPC
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln("unable to connect to port")
	}
	return listener
}
