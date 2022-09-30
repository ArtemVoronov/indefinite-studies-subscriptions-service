package services

import (
	"log"
	"sync"

	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/app"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/services/auth"
	kafkaService "github.com/ArtemVoronov/indefinite-studies-utils/pkg/services/kafka"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/utils"
)

type Services struct {
	auth          *auth.AuthGRPCService
	kafkaProducer *kafkaService.KafkaProducerService
}

var once sync.Once
var instance *Services

func Instance() *Services {
	once.Do(func() {
		if instance == nil {
			instance = createServices()
		}
	})
	return instance
}

func createServices() *Services {
	authcreds, err := app.LoadTLSCredentialsForClient(utils.EnvVar("AUTH_SERVICE_CLIENT_TLS_CERT_PATH"))
	if err != nil {
		log.Fatalf("unable to load TLS credentials: %s", err)
	}

	kafkaProducer, err := kafkaService.CreateKafkaProducerService(utils.EnvVar("KAFKA_HOST") + ":" + utils.EnvVar("KAFKA_PORT"))
	if err != nil {
		log.Fatalf("unable to create kafka producer: %s", err)
	}

	return &Services{
		auth:          auth.CreateAuthGRPCService(utils.EnvVar("AUTH_SERVICE_GRPC_HOST")+":"+utils.EnvVar("AUTH_SERVICE_GRPC_PORT"), &authcreds),
		kafkaProducer: kafkaProducer,
	}
}

func (s *Services) Shutdown() {
	s.auth.Shutdown()
	s.kafkaProducer.Shutdown()
}

func (s *Services) Auth() *auth.AuthGRPCService {
	return s.auth
}

func (s *Services) KafkaProducer() *kafkaService.KafkaProducerService {
	return s.kafkaProducer
}
