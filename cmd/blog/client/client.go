package main

import (
	"context"
	"log"

	"github.com/stephenjlovell/grpc-go-example/api/go/pkg/blogpb"
	"github.com/stephenjlovell/grpc-go-example/internal/shared"
)

func main() {
	cc := shared.Connect()
	defer cc.Close()
	client := blogpb.NewBlogServiceClient(cc)

	req := &blogpb.CreatePostRequest{
		Post: &blogpb.Post{
			AuthorId: "stephenjlovell",
			Title:    "my first MediumClone post",
			Content:  "Lorem ipsum dolor sedet...",
		},
	}

	res, err := client.CreatePost(context.Background(), req)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}
	log.Printf("Post created: %v", res.GetPost())
}
