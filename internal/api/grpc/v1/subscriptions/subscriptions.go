package subscriptions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ArtemVoronov/indefinite-studies-subscriptions-service/internal/services"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/log"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/services/kafka"
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
		return nil, fmt.Errorf("unable to add event: %w", err)
	}
	return &subscriptions.PutEventReply{}, nil
}

func (s *SubscriptionsServiceServer) PutSendEmailEvent(ctx context.Context, in *subscriptions.PutSendEmailEventRequest) (*subscriptions.PutSendEmailEventReply, error) {
	dto := kafka.SendEmailEvent{
		Sender:    in.GetSender(),
		Recepient: in.GetRecepient(),
		Subject:   in.GetSubject(),
		Body:      in.GetBody(),
	}
	data, err := json.Marshal(dto)
	if err != nil {
		return nil, fmt.Errorf("unable to add SEND_EMAIL event: %w", err)
	}

	err = services.Instance().KafkaProducer().CreateMessage(kafka.EVENT_TYPE_SEND_EMAIL, string(data))
	if err != nil {
		return nil, fmt.Errorf("unable to add SEND_EMAIL event: %w", err)
	}
	// TODO clean
	log.Info(fmt.Sprintf("put event: %v", dto))
	return &subscriptions.PutSendEmailEventReply{}, nil
}
