package main

import (
	"context"
	"io"
	"log"
	"time"

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
	doBiStreaming(client)
}

func connect() *grpc.ClientConn {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v\n", err)
	}
	return cc
}

func doBiStreaming(client greetpb.GreetServiceClient) {
	log.Println("starting a bidirectional streaming RPC")
	stream, err := client.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream: %v\n", err)
	}

	reqs := []*greetpb.GreetEveryoneRequest{
		newGreetEveryoneRequest("Steve"),
		newGreetEveryoneRequest("Jon"),
		newGreetEveryoneRequest("April"),
		newGreetEveryoneRequest("Julie"),
		newGreetEveryoneRequest("Mark"),
	}

	waitc := make(chan struct{})
	// send
	go func() {
		for _, req := range reqs {
			if sendErr := stream.Send(req); sendErr != nil {
				log.Fatalf("failed to send request: %v", sendErr)
			}
			time.Sleep(100 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// receive
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				break
			}
			if err != nil {
				close(waitc)
				log.Fatalf("failed to receive response: %v", err)
			}
			log.Printf("received response: %v", res.GetResponse())
		}
	}()

	<-waitc
}

func doClientStreaming(client greetpb.GreetServiceClient) {
	log.Println("starting a client streaming RPC")
	stream, err := client.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream: %v\n", err)
	}

	reqs := []*greetpb.LongGreetRequest{
		newLongGreetRequest("Steve"),
		newLongGreetRequest("Jon"),
		newLongGreetRequest("April"),
		newLongGreetRequest("Julie"),
	}
	for _, req := range reqs {
		if err := stream.Send(req); err != nil {
			log.Fatalf("failed to send request: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}
	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to receive response: %v", err)
	}
	log.Printf("received client streaming response: %s", response.GetResponse())
}

func newLongGreetRequest(firstName string) *greetpb.LongGreetRequest {
	return &greetpb.LongGreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: firstName,
		},
	}
}

func newGreetEveryoneRequest(firstName string) *greetpb.GreetEveryoneRequest {
	return &greetpb.GreetEveryoneRequest{
		Greeting: &greetpb.Greeting{
			FirstName: firstName,
		},
	}
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
