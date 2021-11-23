package dapr

import (
	"context"
	"errors"
	"fmt"

	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/grpc"
	"github.com/sirupsen/logrus"
)

type PubSub struct {
	client         client.Client
	Logger         *logrus.Logger
	TopicSubscribe string
	Name           string
	Buffer         chan string
}

func InitPubsub(topic string, pubsubName string, appPort int, grpcPort int, log *logrus.Logger) (*PubSub, error) {
	var err error
	init := false
	ps := &PubSub{
		Logger:         log,
		TopicSubscribe: topic,
		Name:           pubsubName,
		Buffer:         make(chan string, 100),
	}

	if pubsubName != "" && topic != "" && grpcPort >= 1 {
		ps.initSubscriber(appPort)
		init = true
	}

	if appPort >= 1 {
		if ps.client, err = ps.initPublisher(grpcPort); err != nil {
			return nil, err
		}
		init = true
	}

	if init {
		return ps, nil
	}
	return nil, fmt.Errorf("pubsub disabled")
}

func (p *PubSub) initPublisher(port int) (client.Client, error) {
	clientPub, err := client.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to open atlas pubsub connection: %v", err)
	}
	return clientPub, nil
}

func (p *PubSub) initSubscriber(appPort int) {
	s, err := daprd.NewService(fmt.Sprintf(":%d", appPort))

	if err != nil {
		p.Logger.Fatalf("failed to start the server: %v", err)
	}

	subscription := &common.Subscription{
		PubsubName: p.Name,
		Topic:      p.TopicSubscribe,
	}
	if err := s.AddTopicEventHandler(subscription, p.eventHandler); err != nil {
		p.Logger.Fatalf("error adding handler: %v", err)
	}

	// start the server to handle incoming events
	go func(service common.Service) {
		if err := service.Start(); err != nil {
			p.Logger.Fatalf("server error: %v", err)
		}
	}(s)
}

func (p *PubSub) eventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	p.Buffer <- fmt.Sprintf("%s", e.Data)
	return false, nil
}

func (p *PubSub) Publish(topic string, msg []byte) error {
	if p.client == nil {
		return errors.New("client is not initialized")
	}

	err := p.client.PublishEvent(context.Background(), p.Name, topic, msg)
	return err
}
