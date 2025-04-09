package kafka

import (
	"fmt"

	"github.com/IBM/sarama"

	"google.golang.org/protobuf/proto"
)

var config = sarama.NewConfig()

var producer sarama.AsyncProducer
var producerErr error

func init() {
	config.ClientID = "auth-api-gateway"
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = false

	producer, producerErr = sarama.NewAsyncProducer([]string{"localhost:9092"}, config)
	if producerErr != nil {
		panic(producerErr)
	}
}

func Publish(topic string, protoMessage proto.Message, subscriberId string, requestId string) {
	protoBytes, err := proto.Marshal(protoMessage)
	if err != nil {
		fmt.Println("Failed to marshal protobuf message:", err)
		return
	}
	producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(protoBytes),
		Headers: []sarama.RecordHeader{
			{Key: []byte("subscriberId"), Value: []byte(subscriberId)},
			{Key: []byte("requestId"), Value: []byte(requestId)},
		},
	}
}
