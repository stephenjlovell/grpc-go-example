package main

import (
	clientLib "github.com/stephenjlovell/grpc-go-example/internal/shared/client"
)

func main() {
	cc := clientLib.Connect()
	defer cc.Close()
	// client := blogpb.NewBlogServiceClient(cc)

	// TODO
}
