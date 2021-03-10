package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	greetpb "github.com/stephenjlovell/grpc-go-example/api/go/pkg/greetpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	ListenAddress = "localhost:50051"
	CertFile      = "ssl/server.crt"
	KeyFile       = "ssl/server.pem"
)

// GreetServer is our server implentation.
type GreetServer struct {
	// this is awkward but necessary to provide guarantees about our interface to greetpb.RegisterGreetServiceServer
	greetpb.UnimplementedGreetServiceServer
}

// Greet generates a response to the rpc call after sleeping 200ms
func (s *GreetServer) Greet(ctx context.Context, pb *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	firstName := pb.GetGreeting().GetFirstName()
	lastName := pb.GetGreeting().GetLastName()
	response := "Hello " + firstName + " " + lastName + "!"

	for i := 0; i < 20; i++ {
		time.Sleep(10 * time.Millisecond) // pretend to do work in small increments
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("deadline exceeded")
			return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
		}
	}
	log.Println(response)
	return &greetpb.GreetResponse{
		Response: response,
	}, nil
}

// GreetManyTimes implements an example of server streaming
func (s *GreetServer) GreetManyTimes(pb *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	firstName := pb.GetGreeting().GetFirstName()
	lastName := pb.GetGreeting().GetLastName()

	for i := 1; i <= 10; i++ {
		str := strconv.Itoa(i) + ": Hello " + firstName + " " + lastName + "!"
		resp := &greetpb.GreetManyTimesResponse{
			Response: str,
		}
		if err := stream.Send(resp); err != nil {
			return fmt.Errorf("server failed to send response: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

// LongGreet implements an example of client streaming
func (s *GreetServer) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	str := "Welcome"
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Response: str + "!",
			})
		}
		if err != nil {
			return fmt.Errorf("server failed to read from stream: %v", err)
		}
		firstName := req.Greeting.GetFirstName()
		str += ", " + firstName
	}
}

// GreetEveryone implements bi-directional streaming
func (s *GreetServer) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("server failed to read from stream: %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		response := "Hello " + firstName + "!"
		log.Println(response)
		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Response: response,
		})
		if sendErr != nil {
			return fmt.Errorf("server failed to send to stream: %v", sendErr)
		}
	}
}

func main() {
	listener := listenTo(ListenAddress)

	grpcServer := grpc.NewServer(getCreds())

	greetpb.RegisterGreetServiceServer(grpcServer, &GreetServer{})
	// to use reflection from evans CLI:
	// $ evans -p 50051 -r -t --cacert=./ssl/ca.crt --host=localhost
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
	fmt.Println("greetings!")
	// listen on the default port for gRPC
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln("unable to connect to port")
	}
	return listener
}
