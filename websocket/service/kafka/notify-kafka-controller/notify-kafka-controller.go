package notifyKafkaController

import (
	"errors"
	"websocket/proto"
	"websocket/service/kafka/consumer"
	"websocket/socket"

	googleProto "google.golang.org/protobuf/proto"
)

type StatusMessage struct {
	SubscriptionId string `json:"subscriptionId"`
	RequestId      string `json:"requestId"`
	Message        string `json:"message"`
	Status         string `json:"status"`
}

func Register(pool *socket.WebSocketPool) {
	createHandler := func(status string) func(message googleProto.Message, subscriberId string, requestId string) error {
		return func(message googleProto.Message, subscriberId string, requestId string) error {
			statusMessage, ok := message.(*proto.StatusMessage)
			if !ok {
				return errors.New("error casting message to StatusMessage")
			}
			if pool.HasConnection(subscriberId) {
				return pool.SendMessage(subscriberId, StatusMessage{
					SubscriptionId: subscriberId,
					RequestId:      requestId,
					Message:        statusMessage.Message,
					Status:         status,
				})
			}
			return errors.New("no subscriber connected to this service")
		}
	}

	handlerSuccess := createHandler("success")
	handlerFailure := createHandler("failure")

	go func() {
		consumer.StartKafkaConsumer([]string{"localhost:9092"}, "v1.account.register.status.failure", "websocket", &proto.StatusMessage{}, handlerFailure)
	}()
	go func() {
		consumer.StartKafkaConsumer([]string{"localhost:9092"}, "v1.account.register.status.success", "websocket", &proto.StatusMessage{}, handlerSuccess)
	}()
	go func() {
		consumer.StartKafkaConsumer([]string{"localhost:9092"}, "v1.status.failure", "websocket", &proto.StatusMessage{}, handlerFailure)
	}()
	go func() {
		consumer.StartKafkaConsumer([]string{"localhost:9092"}, "v1.status.success", "websocket", &proto.StatusMessage{}, handlerSuccess)
	}()
}
