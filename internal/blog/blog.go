package blog

import "github.com/stephenjlovell/grpc-go-example/api/go/pkg/blogpb"

const (
	ListenAddress = "localhost:50053"
)

// Server is a placeholder for where our server logic would reside.
type Server struct {
	// this is awkward but necessary to provide guarantees about our interface to calcpb.RegisterCalcServiceServer
	blogpb.UnimplementedBlogServiceServer
}
