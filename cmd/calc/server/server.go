package main

import (
	"log"

	calcpb "github.com/stephenjlovell/grpc-go-example/api/go/pkg/calcpb"
	calc "github.com/stephenjlovell/grpc-go-example/internal/calc"
	"github.com/stephenjlovell/grpc-go-example/internal/shared"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener := shared.ListenTo()
	grpcServer := grpc.NewServer(shared.GetCreds())

	calcpb.RegisterCalcServiceServer(grpcServer, &calc.Server{})
	// to use reflection from evans CLI:
	// $ evans -p 50052 -r -t --cacert=./ssl/ca.crt --host=localhost
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
