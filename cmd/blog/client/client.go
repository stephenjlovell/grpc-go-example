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
	id := doCreatePost(client)
	doGetPost(client, id)
	doUpdatePost(client, id)
}

func doUpdatePost(client blogpb.BlogServiceClient, id string) {
	req := &blogpb.UpdatePostRequest{
		Post: &blogpb.Post{
			Id:       id,
			AuthorId: "stephenjlovell",
			Title:    "my first MediumClone post (edited)",
			Content:  "No more lorem, this is the real stuff!",
		},
	}

	res, err := client.UpdatePost(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to update blog post: %v", err)
	}

	log.Printf("post updated: %v", res.GetPost())
}

func doCreatePost(client blogpb.BlogServiceClient) string {
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
		return ""
	}
	log.Printf("Post with id %s created: %v", res.GetPost().GetId(), res.GetPost())
	return res.GetPost().GetId()
}

func doGetPost(client blogpb.BlogServiceClient, id string) {
	log.Println("reading blog entry")
	req := &blogpb.GetPostRequest{
		PostId: id,
	}
	res, err := client.GetPost(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to read from blog: %v", err)
	}
	log.Printf("retrieved post: %v", res.GetPost())
}
