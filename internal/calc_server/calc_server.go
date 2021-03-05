package calc

import (
	"context"
	"fmt"

	"github.com/stephenjlovell/grpc-go-example/api/go/pkg/calcpb"
)

// CalcServer is a placeholder for where our server logic would reside.
type CalcServer struct {
	// this is awkward but necessary to provide guarantees about our interface to calcpb.RegisterCalcServiceServer
	calcpb.UnimplementedCalcServiceServer
}

// Calculate generates a response to the rpc call
func (s *CalcServer) Calculate(ctx context.Context, pb *calcpb.CalcRequest) (*calcpb.CalcResponse, error) {
	result, err := calculateResult(pb)
	if err != nil {
		return nil, err
	}
	return &calcpb.CalcResponse{
		Result: result,
	}, nil
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
