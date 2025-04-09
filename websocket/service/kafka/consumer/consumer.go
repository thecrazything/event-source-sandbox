package consumer

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
)

type ProtoMessageHandler func(message proto.Message, subscriberId string, requestId string) error

func StartKafkaConsumer(brokers []string, topic string, groupID string, protoMessage proto.Message, handler ProtoMessageHandler) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.AutoCommit.Enable = true

	log.Printf("Starting Kafka consumer for topic: %s, groupID: %s", topic, groupID)

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	defer consumerGroup.Close()

	ctx := context.Background()
	for {
		err := consumerGroup.Consume(ctx, []string{topic}, &consumerGroupHandler{
			protoMessage: protoMessage,
			handler:      handler,
		})
		if err != nil {
			log.Printf("Error consuming messages: %v", err)
		}
	}
}

type consumerGroupHandler struct {
	protoMessage proto.Message
	handler      ProtoMessageHandler
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	log.Printf("Consuming messages from topic: %s", claim.Topic())
	for message := range claim.Messages() {
		log.Printf("Received message: %s", message.Value)
		protoMessageInstance := proto.Clone(h.protoMessage)
		err := proto.Unmarshal(message.Value, protoMessageInstance)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		requestId, subscriberId := getMetaHeaders(message)
		err = h.handler(protoMessageInstance, subscriberId, requestId)
		if err != nil {
			log.Printf("Error handling message: %v", err)
		} else {
			session.MarkMessage(message, "")
		}
	}
	return nil
}

func getMetaHeaders(message *sarama.ConsumerMessage) (string, string) {
	requestId, subscriberId := "", ""
	for _, header := range message.Headers {
		switch string(header.Key) {
		case "requestId":
			requestId = string(header.Value)
		case "subscriberId":
			subscriberId = string(header.Value)
		}
	}
	return requestId, subscriberId
}
