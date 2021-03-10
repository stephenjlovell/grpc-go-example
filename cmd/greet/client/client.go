package main

import (
	"context"
	"io"
	"log"
	"time"

	greetpb "github.com/stephenjlovell/grpc-go-example/api/go/pkg/greetpb"
	clientLib "github.com/stephenjlovell/grpc-go-example/internal/shared/client"
	"google.golang.org/grpc/status"
)

const (
	ListenAddress = "localhost:50051"
)

func main() {
	cc := clientLib.Connect()
	defer cc.Close()
	client := greetpb.NewGreetServiceClient(cc)
	doUnaryRequest(client)
	doServerStreaming(client)
	doClientStreaming(client)
	doBiStreaming(client)
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

func doUnaryRequest(client greetpb.GreetServiceClient) {
	log.Println("executing unary RPC calls...")
	unaryRequestWithTimeout(client, 300*time.Millisecond)
	unaryRequestWithTimeout(client, 100*time.Millisecond) // will time out
}

func unaryRequestWithTimeout(client greetpb.GreetServiceClient, timeout time.Duration) {
	request := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Steve",
			LastName:  "Lovell",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	response, err := client.Greet(ctx, request)
	if err != nil {
		grpcErr, ok := status.FromError(err)
		if ok {
			log.Printf("WARNING: %v", grpcErr.Message())
		} else {
			log.Fatalf("failed to receive response: %v\n", err)
		}
		return
	}
	log.Printf("received unary response: %v\n", response.GetResponse())
}
