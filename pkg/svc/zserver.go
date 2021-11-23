package svc

import (
	"context"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/my-responder/pkg/dapr"
	"github.com/my-responder/pkg/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// version is the current version of the service
	version = "0.0.1"
)

// Default implementation of the MyResponder server interface
type server struct {
	pb.UnimplementedMyResponderServer
	Description string
	Timestamp   time.Time
	Requests    int64
	pubsub      *dapr.PubSub
}

// GetVersion returns the current version of the service
func (server) GetVersion(context.Context, *empty.Empty) (*pb.VersionResponse, error) {
	return &pb.VersionResponse{Version: version}, nil
}

func (s *server) GetDescription(ctx context.Context, req *pb.GetDescriptionRequest) (*pb.GetDescriptionResponse, error) {
	if req.GetService() != 2 {
		data := []byte("ping")
		s.pubsub.Publish("info", data)
		return &pb.GetDescriptionResponse{Description: <-s.pubsub.Buffer}, nil
	}
	s.Requests += 1
	return &pb.GetDescriptionResponse{Description: s.Description}, nil
}

func (s *server) UpdateDescription(ctx context.Context, req *pb.UpdateDescriptionRequest) (*pb.UpdateDescriptionResponse, error) {
	if req.GetDescription() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Description can't be empty")
	}
	if req.GetService() != 2 {
		data := []byte(req.Description)
		s.pubsub.Publish("update-info", data)
		return &pb.UpdateDescriptionResponse{Description: <-s.pubsub.Buffer}, nil
	}
	s.Requests += 1
	s.Description = req.Description
	return &pb.UpdateDescriptionResponse{Description: s.Description}, nil
}

func (s *server) GetUptime(ctx context.Context, req *pb.GetUptimeRequest) (*pb.GetUptimeResponse, error) {
	if req.GetService() != 2 {
		data := []byte("ping")
		s.pubsub.Publish("uptime", data)
		i, err := strconv.Atoi(<-s.pubsub.Buffer)
		if err != nil {
			return nil, err
		}
		return &pb.GetUptimeResponse{Uptime: int64(i)}, nil
	}
	s.Requests += 1
	uptime := time.Now().Unix() - s.Timestamp.Unix()
	return &pb.GetUptimeResponse{Uptime: uptime}, nil
}

func (s *server) GetRequests(ctx context.Context, req *pb.GetRequestsRequest) (*pb.GetRequestsResponse, error) {
	if req.GetService() != 2 {
		data := []byte("ping")
		s.pubsub.Publish("requests", data)
		i, err := strconv.Atoi(<-s.pubsub.Buffer)
		if err != nil {
			return nil, err
		}
		return &pb.GetRequestsResponse{Requests: int64(i)}, nil
	}
	s.Requests += 1
	return &pb.GetRequestsResponse{Requests: s.Requests}, nil
}

// NewBasicServer returns an instance of the default server interface
func NewBasicServer(dapr *dapr.PubSub) (pb.MyResponderServer, error) {
	return &server{Description: "Responder", Timestamp: time.Now(), pubsub: dapr}, nil
}
