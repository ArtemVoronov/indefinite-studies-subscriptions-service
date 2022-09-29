package services

import (
	"log"
	"sync"

	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/app"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/services/auth"
	kafkaService "github.com/ArtemVoronov/indefinite-studies-utils/pkg/services/kafka"
	"github.com/ArtemVoronov/indefinite-studies-utils/pkg/utils"
	// "github.com/confluentinc/confluent-kafka-go/kafka"
)

// TODO: use kafka consumer for implementation of subscriptions or move its using to other services

type Services struct {
	auth          *auth.AuthGRPCService
	kafkaProducer *kafkaService.KafkaProducerService
	// kafkaConsumer *kafkaService.KafkaConsumerService

	// quit              chan struct{}
	// kafkaMessagesChan chan *kafka.Message
	// kafkaErrorsChan   chan error
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

	// kafkaConsumer, err := kafkaService.CreateKafkaConsumerService(utils.EnvVar("KAFKA_HOST")+":"+utils.EnvVar("KAFKA_PORT"), utils.EnvVar("KAFKA_GROUP_ID"))
	// if err != nil {
	// 	log.Fatalf("unable to create kafka consumer: %s", err)
	// }

	// quit := make(chan struct{})

	// kafkaMessagesChan, kafkaErrorsChan := kafkaConsumer.PollTopics(quit, "myTopic", 10_000)
	// // kafkaMessagesChan, kafkaErrorsChan := kafkaConsumer.SubscribeTopics(quit, "myTopic", "^hello", 10*time.Second)

	// go func() {
	// 	for e := range kafkaMessagesChan {
	// 		fmt.Printf("================Message on %s: %s\n", e.TopicPartition, string(e.Value))
	// 	}
	// }()
	// go func() {
	// 	for e := range kafkaErrorsChan {
	// 		fmt.Println(e)
	// 	}
	// }()

	return &Services{
		auth:          auth.CreateAuthGRPCService(utils.EnvVar("AUTH_SERVICE_GRPC_HOST")+":"+utils.EnvVar("AUTH_SERVICE_GRPC_PORT"), &authcreds),
		kafkaProducer: kafkaProducer,
		// kafkaConsumer: kafkaConsumer,

		// quit:              quit,
		// kafkaMessagesChan: kafkaMessagesChan,
		// kafkaErrorsChan:   kafkaErrorsChan,
	}
}

func (s *Services) Shutdown() {
	// defer close(s.quit)
	// defer close(s.kafkaMessagesChan)
	// defer close(s.kafkaErrorsChan)

	s.auth.Shutdown()
	s.kafkaProducer.Shutdown()
	// s.kafkaConsumer.Shutdown()
}

func (s *Services) Auth() *auth.AuthGRPCService {
	return s.auth
}

func (s *Services) KafkaProducer() *kafkaService.KafkaProducerService {
	return s.kafkaProducer
}

// func (s *Services) KafkaConsumer() *kafkaService.KafkaConsumerService {
// 	return s.kafkaConsumer
// }
