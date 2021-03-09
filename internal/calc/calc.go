package calc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/stephenjlovell/grpc-go-example/api/go/pkg/calcpb"
)

const (
	LISTEN_ADDRESS = "0.0.0.0:50052"
)

// Server is a placeholder for where our server logic would reside.
type Server struct {
	// this is awkward but necessary to provide guarantees about our interface to calcpb.RegisterCalcServiceServer
	calcpb.UnimplementedCalcServiceServer
}

// GetPrimes decomposes req.Value into a stream of primes
func (s *Server) GetPrimes(req *calcpb.PrimeRequest, stream calcpb.CalcService_GetPrimesServer) error {
	v := req.GetValue()
	sent := []uint32{}
	if v < 2 {
		return nil
	}
	for i := uint32(2); v > 1; {
		if v%i == 0 {
			err := stream.Send(&calcpb.PrimeResponse{
				Value:  i,
				JobUid: req.GetJobUid(),
			})
			if err != nil {
				log.Printf("[%s] WARNING: unable to stream response to client: %v", req.GetJobUid(), err)
				return err
			}
			sent = append(sent, i)
			v /= i
		} else {
			i++
		}
	}
	// log.Printf("[%s] %s:%v => %8f", uid, op, operands, result)
	log.Printf("[%s] PRIMES: %v => %v", req.GetJobUid(), req.GetValue(), sent)
	return nil
}

// Calculate generates a response to the rpc call
func (s *Server) Calculate(ctx context.Context, pb *calcpb.CalcRequest) (*calcpb.CalcResponse, error) {
	result, err := calculateResult(pb)
	if err != nil {
		log.Printf("[%s] WARNING: %v", pb.GetJobUid(), err)
		return nil, err
	}
	logResult(result, pb)
	return &calcpb.CalcResponse{
		Result: result,
		JobUid: pb.GetJobUid(),
	}, nil
}

func logResult(result float64, pb *calcpb.CalcRequest) {
	op := strings.ToLower(pb.GetOperation().String())
	operands := pb.GetOperands()
	uid := pb.GetJobUid()
	log.Printf("[%s] %s: %v => %8f", uid, op, operands, result)
}

func calculateResult(pb *calcpb.CalcRequest) (float64, error) {
	operands := pb.GetOperands()
	if len(operands) < 2 {
		return 0, fmt.Errorf("invalid operation: too few operands")
	}
	switch pb.GetOperation() {
	case calcpb.Operations_ADD:
		return add(operands)
	case calcpb.Operations_SUBTRACT:
		return subtract(operands)
	case calcpb.Operations_MULTIPLY:
		return multiply(operands)
	case calcpb.Operations_DIVIDE:
		return divide(operands)
	default:
		return 0, fmt.Errorf("invalid operation requested: %v", pb.GetOperation())
	}
}

func add(operands []int64) (float64, error) {
	var result float64
	for _, n := range operands {
		result += float64(n)
	}
	return result, nil
}

func subtract(operands []int64) (float64, error) {
	result := float64(operands[0])
	for _, n := range operands[1:] {
		result -= float64(n)
	}
	return result, nil
}

func multiply(operands []int64) (float64, error) {
	result := float64(operands[0])
	for _, n := range operands[1:] {
		result *= float64(n)
	}
	return result, nil
}

func divide(operands []int64) (float64, error) {
	result := float64(operands[0])
	for _, n := range operands[1:] {
		if n == 0 {
			return 0, fmt.Errorf("cannot divide by zero")
		}
		result /= float64(n)
	}
	return result, nil
}
