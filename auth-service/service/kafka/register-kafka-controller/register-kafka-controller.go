package registerKafkaController

import (
	"errors"
	"projects/auth-service/proto"
	"projects/auth-service/service/kafka/consumer"
	"projects/auth-service/service/register"

	googleProto "google.golang.org/protobuf/proto"
)

func Register() {
	handler := func(message googleProto.Message, subsriberId string, requestId string) error {
		// Handle the message
		registerMessage, ok := message.(*proto.RegisterUser)
		if !ok {
			return errors.New("error casting message to RegisterUser")
		}
		println("Subscriber ID:", subsriberId)
		println("Request ID:", requestId)
		println(registerMessage.Username)
		
		register.RegisterUser(registerMessage.Username, registerMessage.Password, subsriberId, requestId)

		return nil
	}

	consumer.StartKafkaConsumer([]string{"localhost:9092"}, "v1.account.register", "auth-service", &proto.RegisterUser{}, handler)
}
