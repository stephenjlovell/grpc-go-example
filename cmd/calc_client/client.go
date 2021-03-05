package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
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
	log.Printf("client %d: starting", i)
	cc := connect()
	defer cc.Close()
	defer wg.Done()

	client := calcpb.NewCalcServiceClient(cc)

	for j := 0; j < 10; j++ {
		// sleep for 101-300 ms
		time.Sleep(time.Duration(rand.Intn(300)+200) * time.Millisecond)

		r := makeRandomRequest()
		response, err := sendRequest(client, r)
		if err != nil {
			log.Printf("[%s] client:%d WARNING: %v", r.GetJobUid(), i, err)
			continue
		}
		log.Printf("[%s] client:%d result: %v", response.GetJobUid(), i, response.GetResult())
	}
	log.Printf("client:%d finished working", i)
}

func connect() *grpc.ClientConn {
	cc, err := grpc.Dial(calc.LISTEN_ADDRESS, grpc.WithInsecure())
	if err != nil {
		// blow everything up if the server won't speak to us
		log.Fatalf("Failed to connect to server: %v\n", err)
	}
	return cc
}

func sendRequest(client calcpb.CalcServiceClient, r *calcpb.CalcRequest) (*calcpb.CalcResponse, error) {
	response, err := client.Calculate(context.Background(), r)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	return response, nil
}

// NOTE: this intentionally generates occasional malformed requests in order to make sure the
// server can gracefully handle them.
func makeRandomRequest() *calcpb.CalcRequest {
	operands := make([]int64, rand.Intn(10))
	for i := range operands {
		operands[i] = int64(rand.Intn(100))
	}
	return &calcpb.CalcRequest{
		Operation: calcpb.Operations(rand.Intn(5)),
		Operands:  operands,
		JobUid:    uuid.NewString(),
	}
}
