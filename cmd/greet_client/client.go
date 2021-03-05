package main

import (
	"context"
	"log"

	greetpb "github.com/stephenjlovell/grpc-go-example/api/go/pkg/greetpb"
	"google.golang.org/grpc"
)

func main() {
	cc := connect()
	defer cc.Close()

	client := greetpb.NewGreetServiceClient(cc)
	response := sendRequest(client)
	log.Printf("received response: %v\n", response.GetResponse())
}

func connect() *grpc.ClientConn {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v\n", err)
	}
	return cc
}

func sendRequest(client greetpb.GreetServiceClient) *greetpb.GreetResponse {
	log.Println("executing RPC call...")
	request := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Steve",
			LastName:  "Lovell",
		},
	}
	response, err := client.Greet(context.Background(), request)
	if err != nil {
		log.Fatalf("failed to receive response: %v\n", err)
	}
	return response
}
