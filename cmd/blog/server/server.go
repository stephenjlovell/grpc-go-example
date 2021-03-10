package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	blogpb "github.com/stephenjlovell/grpc-go-example/api/go/pkg/blogpb"
	"github.com/stephenjlovell/grpc-go-example/internal/blog/db"
	blogServer "github.com/stephenjlovell/grpc-go-example/internal/blog/server"
	"github.com/stephenjlovell/grpc-go-example/internal/shared"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	listener := shared.ListenTo()
	grpcServer := grpc.NewServer(shared.GetCreds())

	blogpb.RegisterBlogServiceServer(grpcServer, &blogServer.Server{})
	// to use reflection from evans CLI:
	// $ evans -p 50053 -r -t --cacert=./ssl/ca.crt --host=localhost
	reflection.Register(grpcServer)

	go func() {
		log.Println("starting blog server... beep boop...")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v\n", err)
		}
	}()

	sigKill := make(chan os.Signal, 1)
	signal.Notify(sigKill, os.Interrupt) // relay incoming signals to sigKill

	<-sigKill // block until OS signal

	log.Printf("\ngracefully shutting down server...\n")

	fmt.Printf("finishing running grpc requests...")
	grpcServer.GracefulStop()
	fmt.Printf("done.\n")

	fmt.Printf("closing listener...")
	listener.Close()
	fmt.Printf("done.\n")

	db.GracefulDisconnect()

}
