package main

import (
	"fmt"
	"log"

	greetpb "github.com/stephenjlovell/grpc-go-example/api/go/pkg"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v\n", err)
	}

	defer cc.Close()

	client := greetpb.NewGreetServiceClient(cc)
	fmt.Printf("created client: %v\n", client)
}
