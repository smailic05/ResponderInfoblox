package svc

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/my-responder/pkg/dapr"
	"github.com/my-responder/pkg/pb"
	"github.com/spf13/viper"
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
	mtx         sync.RWMutex
}

func (s *server) GetRequestsFromServer() int {
	s.mtx.RLock()
	tmp := s.Requests
	s.mtx.RUnlock()
	return int(tmp)
}

func (s *server) IncRequests() {
	s.mtx.Lock()
	s.Requests++
	s.mtx.Unlock()
}

// GetVersion returns the current version of the service
func (*server) GetVersion(context.Context, *empty.Empty) (*pb.VersionResponse, error) {
	return &pb.VersionResponse{Version: version}, nil
}

func (s *server) GetDescription(ctx context.Context, req *pb.GetDescriptionRequest) (*pb.GetDescriptionResponse, error) {
	if req.GetService() != 2 {
		id := s.GetRequestsFromServer()
		data := dapr.Message{Id: id}
		s.pubsub.Publish("info", data)
		resp := <-s.pubsub.Buffer
		for id != resp.Id {
			s.pubsub.Buffer <- resp
			resp = <-s.pubsub.Buffer
		}
		return &pb.GetDescriptionResponse{Description: resp.Data}, nil
	}
	s.IncRequests()
	return &pb.GetDescriptionResponse{Description: s.Description}, nil
}

func (s *server) UpdateDescription(ctx context.Context, req *pb.UpdateDescriptionRequest) (*pb.UpdateDescriptionResponse, error) {
	if req.GetDescription() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Description can't be empty")
	}
	if req.GetService() != 2 {
		id := s.GetRequestsFromServer()
		data := dapr.Message{Id: id, Data: req.Description}
		s.pubsub.Publish("update-info", data)
		resp := <-s.pubsub.Buffer
		for id != resp.Id {
			s.pubsub.Buffer <- resp
			resp = <-s.pubsub.Buffer
		}
		return &pb.UpdateDescriptionResponse{Description: resp.Data}, nil
	}
	s.IncRequests()
	s.Description = req.Description
	return &pb.UpdateDescriptionResponse{Description: s.Description}, nil
}

func (s *server) GetUptime(ctx context.Context, req *pb.GetUptimeRequest) (*pb.GetUptimeResponse, error) {
	if req.GetService() != 2 {
		id := s.GetRequestsFromServer()
		data := dapr.Message{Id: id}
		s.pubsub.Publish("uptime", data)
		resp := <-s.pubsub.Buffer
		for id != resp.Id {
			s.pubsub.Buffer <- resp
			resp = <-s.pubsub.Buffer
		}
		i, err := strconv.Atoi(resp.Data)
		if err != nil {
			return nil, err
		}
		return &pb.GetUptimeResponse{Uptime: int64(i)}, nil
	}
	id := s.GetRequestsFromServer()
	data := dapr.Message{Id: id}
	s.pubsub.Publish("get-mode", data)
	resp := <-s.pubsub.Buffer
	for id != resp.Id {
		s.pubsub.Buffer <- resp
		resp = <-s.pubsub.Buffer
	}
	i, err := strconv.Atoi(resp.Data)
	if err != nil {
		return nil, err
	}
	s.IncRequests()
	if i == 0 {
		return &pb.GetUptimeResponse{}, status.Error(codes.Unavailable, "Service is Unavailable")
	}
	uptime := time.Now().Unix() - s.Timestamp.Unix()
	return &pb.GetUptimeResponse{Uptime: uptime}, nil
}

func (s *server) GetRequests(ctx context.Context, req *pb.GetRequestsRequest) (*pb.GetRequestsResponse, error) {
	if req.GetService() != 2 {
		id := s.GetRequestsFromServer()
		data := dapr.Message{Id: id}
		s.pubsub.Publish("requests", data)
		resp := <-s.pubsub.Buffer
		for id != resp.Id {
			s.pubsub.Buffer <- resp
			resp = <-s.pubsub.Buffer
		}
		i, err := strconv.Atoi(resp.Data)
		if err != nil {
			return nil, err
		}
		return &pb.GetRequestsResponse{Requests: int64(i)}, nil
	}
	s.IncRequests()
	return &pb.GetRequestsResponse{Requests: int64(s.GetRequestsFromServer())}, nil
}
func (s *server) GetMode(ctx context.Context, req *pb.GetModeRequest) (*pb.GetModeResponse, error) {
	id := s.GetRequestsFromServer()
	data := dapr.Message{Id: int(id)}
	s.pubsub.Publish("get-mode", data)
	resp := <-s.pubsub.Buffer
	for id != resp.Id {
		s.pubsub.Buffer <- resp
		resp = <-s.pubsub.Buffer
	}
	i, err := strconv.Atoi(resp.Data)
	if err != nil {
		return nil, err
	}
	s.IncRequests()
	return &pb.GetModeResponse{Mode: int64(i)}, nil
}

func (s *server) SetMode(ctx context.Context, req *pb.SetModeRequest) (*pb.SetModeResponse, error) {
	id := s.GetRequestsFromServer()
	data := dapr.Message{Id: id}
	s.pubsub.Publish("set-mode", data)
	resp := <-s.pubsub.Buffer
	for id != resp.Id {
		s.pubsub.Buffer <- resp
		resp = <-s.pubsub.Buffer
	}
	i, err := strconv.Atoi(resp.Data)
	if err != nil {
		return nil, err
	}
	s.IncRequests()
	return &pb.SetModeResponse{Mode: int64(i)}, nil
}

func (s *server) Restart(ctx context.Context, req *pb.RestartRequest) (*pb.RestartResponse, error) {
	if req.Service != 2 {
		s.pubsub.Publish("restart", dapr.Message{})
		return &pb.RestartResponse{}, nil
	}
	s.Description = viper.GetString("app.id")
	s.Timestamp = time.Now()
	s.Requests = 0
	return &pb.RestartResponse{}, nil
}

// NewBasicServer returns an instance of the default server interface
func NewBasicServer(dapr *dapr.PubSub) (pb.MyResponderServer, error) {
	return &server{Description: "Responder", Timestamp: time.Now(), pubsub: dapr}, nil
}
