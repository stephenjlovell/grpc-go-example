package main

import (
	"context"
	"fmt"
	"log"

	calcpb "github.com/stephenjlovell/grpc-go-example/api/go/pkg/calcpb"
	calc "github.com/stephenjlovell/grpc-go-example/internal/calc"
	"google.golang.org/grpc"
)

func main() {
	cc := connect()
	defer cc.Close()

	client := calcpb.NewCalcServiceClient(cc)
	response := sendRequest(client)
	fmt.Printf("The answer is: %v\n", response.GetResult())
}

func connect() *grpc.ClientConn {
	cc, err := grpc.Dial(calc.LISTEN_ADDRESS, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v\n", err)
	}
	return cc
}

func sendRequest(client calcpb.CalcServiceClient) *calcpb.CalcResponse {
	fmt.Println("executing RPC call...")
	request := &calcpb.CalcRequest{
		Operation: calcpb.Operations_ADD,
		Operands:  []int64{1, 1, 2, 3, 5, 8},
	}
	response, err := client.Calculate(context.Background(), request)
	if err != nil {
		log.Fatalf("failed to receive response: %v\n", err)
	}
	return response

}
