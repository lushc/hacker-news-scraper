package api

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/lushc/hacker-news-scraper/internal/datastore"
	pb "github.com/lushc/hacker-news-scraper/protobufs"
)

const (
	serverPortEnv = "SERVER_PORT"
)

var (
	errServerPortEnv = fmt.Errorf("missing env var %s", serverPortEnv)
)

type Server struct {
	port   int
	reader datastore.Reader
	pb.UnimplementedAPIServer
}

func NewServer(reader datastore.Reader) (*Server, error) {
	// TODO: viper config instead
	portEnv, ok := os.LookupEnv(serverPortEnv)
	if !ok {
		log.Fatal(errServerPortEnv)
	}

	port, err := strconv.Atoi(portEnv)
	if err != nil {
		log.Fatal(err)
	}

	return &Server{
		port:   port,
		reader: reader,
	}, nil
}

func (s Server) ListAll(empty *emptypb.Empty, stream pb.API_ListAllServer) error {
	items, err := s.reader.All(stream.Context())
	if err != nil {
		return fmt.Errorf("ListAll read: %w", err)
	}

	for _, item := range items {
		if err := stream.Send(datastore.Itop(*item)); err != nil {
			return fmt.Errorf("ListAll send: %w", err)
		}
	}

	return nil
}

func (s Server) ListType(request *pb.TypeRequest, stream pb.API_ListTypeServer) error {
	items, err := s.reader.ByItemType(stream.Context(), datastore.EnumTypes[request.Type])
	if err != nil {
		return fmt.Errorf("ListType read: %w", err)
	}

	for _, item := range items {
		if err := stream.Send(datastore.Itop(*item)); err != nil {
			return fmt.Errorf("ListType send: %w", err)
		}
	}

	return nil
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	gs := grpc.NewServer()
	pb.RegisterAPIServer(gs, s)
	if err := gs.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
