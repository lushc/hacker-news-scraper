package api

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/lushc/hacker-news-scraper/internal/datastore"
	pb "github.com/lushc/hacker-news-scraper/protobufs"
)

var (
	enumTypes = map[pb.TypeRequest_Type]datastore.ItemType{
		pb.TypeRequest_JOB:   datastore.Job,
		pb.TypeRequest_STORY: datastore.Story,
	}
)

type server struct {
	pb.UnimplementedAPIServer
}

func (s server) ListAll(empty *emptypb.Empty, stream pb.API_ListAllServer) error {
	panic("implement me")
}

func (s server) ListType(request *pb.TypeRequest, stream pb.API_ListTypeServer) error {
	fmt.Println(enumTypes[*request.Type.Enum()])
	panic("implement me")
}

func StartServer(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s := grpc.NewServer()
	pb.RegisterAPIServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
