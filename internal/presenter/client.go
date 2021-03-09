package presenter

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/lushc/hacker-news-scraper/internal/datastore"
	pb "github.com/lushc/hacker-news-scraper/protobufs"
)

const (
	serverHostnameEnv = "SERVER_HOSTNAME"
	serverPortEnv     = "SERVER_PORT"
	timeout           = 30 * time.Second
)

var (
	errServerHostnameEnv = fmt.Errorf("missing env var %s", serverHostnameEnv)
	errServerPortEnv     = fmt.Errorf("missing env var %s", serverPortEnv)
)

// ItemStreamClient is an alias for the protobuf interfaces which stream item messages
type ItemStreamClient interface {
	Recv() (*pb.Item, error)
	grpc.ClientStream
}

// ItemStreamer is a wrapper function to simplify getting an ItemStreamClient
type ItemStreamer func(cancelCtx context.Context) (ItemStreamClient, error)

// Client wraps the gRPC client's streaming functions in channels
type Client struct {
	client pb.APIClient
	conn   *grpc.ClientConn
}

// NewClient will connect to the gRPC server in blocking mode
func NewClient() (*Client, error) {
	// TODO: viper config instead
	hostname, ok := os.LookupEnv(serverHostnameEnv)
	if !ok {
		return nil, errServerHostnameEnv
	}

	port, ok := os.LookupEnv(serverPortEnv)
	if !ok {
		return nil, errServerPortEnv
	}

	// TODO: use TLS
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", hostname, port), grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	return &Client{
		client: pb.NewAPIClient(conn),
		conn:   conn,
	}, nil
}

// Close will close the underlying connection to the gRPC server
func (c Client) Close() {
	c.conn.Close()
}

// WrapListAll will wrap the gRPC APIClient's ListAll so that it can be used as a Source
func (c Client) WrapListAll(ctx context.Context, in *emptypb.Empty) Source {
	return wrapStream(ctx, func(cancelCtx context.Context) (ItemStreamClient, error) {
		return c.client.ListAll(cancelCtx, in)
	})
}

// WrapListType will wrap the gRPC APIClient's ListType so that it can be used as a Source
func (c Client) WrapListType(ctx context.Context, in *pb.TypeRequest) Source {
	return wrapStream(ctx, func(cancelCtx context.Context) (ItemStreamClient, error) {
		return c.client.ListType(cancelCtx, in)
	})
}

// wrapStream will wrap an ItemStreamer so that when invoked, the stream is created and its responses are pushed into a channel
func wrapStream(ctx context.Context, streamer ItemStreamer) Source {
	return func(items chan<- datastore.Item, errs chan<- error, ephemeral bool) {
		if ephemeral {
			defer close(items)
			defer close(errs)
		}

		cancelCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		stream, err := streamer(cancelCtx)
		if err != nil {
			errs <- fmt.Errorf("creating stream: %w", err)
			return
		}

		for {
			item, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					errs <- fmt.Errorf("receiving stream: %w", err)
				}
				return
			}

			items <- *datastore.Ptoi(item)
		}
	}
}
