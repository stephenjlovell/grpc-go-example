package blog

import (
	"context"

	"github.com/stephenjlovell/grpc-go-example/api/go/pkg/blogpb"
	"github.com/stephenjlovell/grpc-go-example/internal/blog/db"
)

const (
	ListenAddress = "localhost:50053"
)

// Server is a placeholder for where our server logic would reside.
type Server struct {
	// this is awkward but necessary to provide guarantees about our interface to calcpb.RegisterCalcServiceServer
	blogpb.UnimplementedBlogServiceServer
}

// CreatePost handles unary post creation
func (s *Server) CreatePost(ctx context.Context, req *blogpb.CreatePostRequest) (*blogpb.CreatePostResponse, error) {
	post := req.GetPost()

	id, err := db.GetCollection("posts").SaveOne(ctx, db.Post{
		AuthorID: post.GetAuthorId(),
		Title:    post.GetTitle(),
		Content:  post.GetContent(),
	})
	if err != nil {
		return nil, err
	}

	return &blogpb.CreatePostResponse{
		Post: &blogpb.Post{
			Id:       id,
			AuthorId: post.GetAuthorId(),
			Title:    post.GetTitle(),
			Content:  post.GetContent(),
		},
	}, nil
}
