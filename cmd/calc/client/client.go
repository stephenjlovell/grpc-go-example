package main

import (
	"context"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	calcpb "github.com/stephenjlovell/grpc-go-example/api/go/pkg/calcpb"
	calc "github.com/stephenjlovell/grpc-go-example/internal/calc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

const (
	CA_CERT_FILE = "ssl/ca.crt"
)

func main() {
	doSquareRoot()
	doAverage()
	streamPrimes()
	doArithmetic()
}

func doSquareRoot() {
	cc := connect()
	defer cc.Close()
	client := calcpb.NewCalcServiceClient(cc)
	squareRootOf(9, client)  // works
	squareRootOf(-1, client) // throws INVALID_ARGUMENT
}

func squareRootOf(v int64, client calcpb.CalcServiceClient) {
	requestID := uuid.NewString()
	res, err := client.SquareRoot(context.Background(), &calcpb.SquareRootRequest{
		Value:  v,
		JobUid: requestID,
	})
	if err != nil {
		grpcErr, ok := status.FromError(err)
		if ok { // it's a GRPC error
			log.Printf("[%s] WARNING: %s", requestID, grpcErr.Message())
		} else { // it's something else
			log.Fatalf("[%s] %v", requestID, err)
		}
	} else {
		logResult(requestID, "square root", v, res.GetResult())
	}
}

func logResult(requestID, op string, in, out interface{}) {
	log.Printf("[%s] %s: %v => %v", requestID, op, in, out)
}

// client streaming example
func doAverage() {
	cc := connect()
	defer cc.Close()
	client := calcpb.NewCalcServiceClient(cc)

	stream, err := client.GetAverage(context.Background())
	requestID := uuid.NewString()
	if err != nil {
		log.Printf("[%s] WARNING: could not open stream: %v", requestID, err)
		return
	}
	seq := fibonacciOfLength(25)
	for _, v := range seq {
		err := stream.Send(&calcpb.AverageRequest{
			JobUid: requestID,
			Value:  v,
		})
		if err != nil {
			log.Printf("[%s] WARNING: request failed to send: %v", requestID, err)
			return
		}
	}
	// await results
	res, err := stream.CloseAndRecv()
	logResult(requestID, "average", seq, res.GetResult())
}

// 1, 1, 2, 3, 5, 8... till your int64 overfloweth
func fibonacciOfLength(n int) []int64 {
	seq := make([]int64, n, n)
	seq[0] = 1
	seq[1] = 1
	for i := 2; i < n; i++ {
		seq[i] = seq[i-1] + seq[i-2]
	}
	return seq
}

func streamPrimes() {
	cc := connect()
	defer cc.Close()
	client := calcpb.NewCalcServiceClient(cc)
	getPrimesFor(client, 2)
	getPrimesFor(client, 42)
	getPrimesFor(client, 12345678)
}

func getPrimesFor(client calcpb.CalcServiceClient, val uint32) {
	req := &calcpb.PrimeRequest{
		Value:  val,
		JobUid: uuid.NewString(),
	}
	stream, err := client.GetPrimes(context.Background(), req)
	if err != nil {
		log.Printf("[%s] WARNING: request failed: %v", req.GetJobUid(), err)
		return
	}
	results := []uint32{}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("[%s] WARNING: request failed: %v", req.GetJobUid(), err)
			break
		}
		results = append(results, resp.GetValue())
	}
	logResult(req.GetJobUid(), "primes", val, results)
}

func doArithmetic() {
	// create 4 worker threads sharing a connection
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
		response, err := client.Calculate(context.Background(), r)
		if err != nil {
			log.Printf("[%s] client:%d WARNING: %v", r.GetJobUid(), i, err)
			continue
		}
		log.Printf("[%s] client:%d result: %v", response.GetJobUid(), i, response.GetResult())
	}
	log.Printf("client:%d finished working", i)
}

func connect() *grpc.ClientConn {
	creds, sslErr := credentials.NewClientTLSFromFile(CA_CERT_FILE, "")
	if sslErr != nil {
		log.Fatalf("Failed to load CA trust certificate: %v", sslErr)
	}
	cc, err := grpc.Dial(calc.ListenAddress, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v\n", err)
	}
	return cc
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
