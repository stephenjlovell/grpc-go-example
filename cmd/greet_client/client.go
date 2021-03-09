package main

import (
	"context"
	"io"
	"log"

	greetpb "github.com/stephenjlovell/grpc-go-example/api/go/pkg/greetpb"
	"google.golang.org/grpc"
)

func main() {
	cc := connect()
	defer cc.Close()
	client := greetpb.NewGreetServiceClient(cc)
	doUnaryRequest(client)
	doServerStreaming(client)
	doClientStreaming(client)
}

func connect() *grpc.ClientConn {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v\n", err)
	}
	return cc
}

func doClientStreaming(client greetpb.GreetServiceClient) {

}

func doServerStreaming(client greetpb.GreetServiceClient) {
	// make streaming api request
	responseStream := requestServerStreaming(client)
	for {
		msg, err := responseStream.Recv()
		if err == io.EOF {
			// no more responses will be sent
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v\n", err)
		}
		log.Printf("received streamed response: %s\n", msg.GetResponse())
	}
}

func requestServerStreaming(client greetpb.GreetServiceClient) greetpb.GreetService_GreetManyTimesClient {
	log.Println("sending a stream of RPC calls...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Steve",
			LastName:  "Lovell",
		},
	}
	responseStream, err := client.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to receive response: %v\n", err)
	}
	return responseStream
}

// GreetManyTimes(ctx context.Context, in *GreetManyTimesRequest, opts ...grpc.CallOption) (GreetService_GreetManyTimesClient, error)

func doUnaryRequest(client greetpb.GreetServiceClient) {
	log.Println("executing single RPC call...")
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
	log.Printf("received unary response: %v\n", response.GetResponse())
}
