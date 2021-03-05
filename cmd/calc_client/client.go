package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	calcpb "github.com/stephenjlovell/grpc-go-example/api/go/pkg/calcpb"
	calc "github.com/stephenjlovell/grpc-go-example/internal/calc"
	"google.golang.org/grpc"
)

func main() {
	// create 4 example clients with their own connection:
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go startClient(i, &wg)
	}
	wg.Wait()
	log.Println("all clients finished working")
}

func startClient(i int, wg *sync.WaitGroup) {
	log.Printf("client %d: starting\n", i)
	cc := connect()
	defer cc.Close()
	defer wg.Done()

	client := calcpb.NewCalcServiceClient(cc)

	for j := 0; j < 10; j++ {
		// sleep for 101-300 ms
		time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)

		response, err := sendRequest(client)
		if err != nil {
			log.Printf("client %d job %d WARNING: %v\n", i, j, err)
			continue
		}
		log.Printf("client %d: job %d result: %v\n", i, j, response.GetResult())
	}
	log.Printf("client %d finished working\n", i)
}

func connect() *grpc.ClientConn {
	cc, err := grpc.Dial(calc.LISTEN_ADDRESS, grpc.WithInsecure())
	if err != nil {
		// blow everything up if the server won't speak to us
		log.Fatalf("Failed to connect to server: %v\n", err)
	}
	return cc
}

func sendRequest(client calcpb.CalcServiceClient) (*calcpb.CalcResponse, error) {
	request := makeRandomRequest()
	response, err := client.Calculate(context.Background(), request)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	return response, nil
}

func makeRandomRequest() *calcpb.CalcRequest {
	operation := rand.Intn(5)
	operands := make([]int64, rand.Intn(5)+2)
	for i := range operands {
		operands[i] = int64(rand.Intn(99) + 1)
	}
	return &calcpb.CalcRequest{
		Operation: calcpb.Operations(operation),
		Operands:  operands,
	}
}
