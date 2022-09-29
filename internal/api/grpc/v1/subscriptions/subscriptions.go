package subscriptions

import (
	"context"
	"fmt"

	"github.com/ArtemVoronov/indefinite-studies-subscriptions-service/internal/services"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/services/subscriptions"
	"google.golang.org/grpc"
)

type SubscriptionsServiceServer struct {
	subscriptions.UnimplementedSubscriptionsServiceServer
}

func RegisterServiceServer(s *grpc.Server) {
	subscriptions.RegisterSubscriptionsServiceServer(s, &SubscriptionsServiceServer{})
}

func (s *SubscriptionsServiceServer) PutEvent(ctx context.Context, in *subscriptions.PutEventRequest) (*subscriptions.PutEventReply, error) {
	err := services.Instance().KafkaProducer().CreateMessage(in.GetEventType(), in.GetEventBody())
	if err != nil {
		return nil, fmt.Errorf("unable to add event: %s", err)
	}
	return &subscriptions.PutEventReply{}, nil
}

func (s *SubscriptionsServiceServer) GetEvent(ctx context.Context, in *subscriptions.GetEventRequest) (*subscriptions.GetEventReply, error) {
	// TODO
	// return &subscriptions.GetEventReply{}, nil
	return nil, fmt.Errorf("not implemented")
}
