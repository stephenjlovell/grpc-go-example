package blog

import (
	"context"
	"log"

	"github.com/stephenjlovell/grpc-go-example/api/go/pkg/blogpb"
	"github.com/stephenjlovell/grpc-go-example/internal/blog/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	log.Printf("created post %s by author %s", post.GetTitle(), post.GetAuthorId())
	return &blogpb.CreatePostResponse{
		Post: &blogpb.Post{
			Id:       id,
			AuthorId: post.GetAuthorId(),
			Title:    post.GetTitle(),
			Content:  post.GetContent(),
		},
	}, nil
}

func (s *Server) GetPost(ctx context.Context, req *blogpb.GetPostRequest) (*blogpb.GetPostResponse, error) {
	blogId := req.GetPostId()
	id, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot parse ID: %v", err)
	}
	data := new(db.Post)
	filter := bson.D{{"_id", id}}
	findErr := db.GetCollection("posts").FindOne(ctx, filter).Decode(data)
	if findErr != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot find post with id %v: %v", id, findErr)
	}

	return &blogpb.GetPostResponse{
		Post: toPost(data),
	}, nil
}

func (s *Server) UpdatePost(ctx context.Context, req *blogpb.UpdatePostRequest) (*blogpb.UpdatePostResponse, error) {
	post := req.GetPost()
	id, err := primitive.ObjectIDFromHex(post.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot parse ID: %v", err)
	}
	data := &db.Post{
		AuthorID: post.GetAuthorId(),
		Title:    post.GetTitle(),
		Content:  post.GetContent(),
	}
	filter := bson.D{{"_id", id}}

	res := db.GetCollection("posts").FindOneAndUpdate(ctx, filter, data)
	if res.Err() == mongo.ErrNoDocuments {
		return nil, status.Errorf(codes.InvalidArgument, "no post matches ID %v: %v", post.GetId(), res.Err())
	}

	return &blogpb.UpdatePostResponse{
		Post: toPost(data),
	}, nil
}

func toPost(data *db.Post) *blogpb.Post {
	return &blogpb.Post{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorID,
		Content:  data.Content,
		Title:    data.Title,
	}
}
