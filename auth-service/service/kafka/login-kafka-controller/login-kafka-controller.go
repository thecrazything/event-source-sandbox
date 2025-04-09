package login

import (
	"errors"
	"projects/auth-service/proto"
	"projects/auth-service/service/kafka/consumer"
	"projects/auth-service/service/kafka/producer"
	"projects/auth-service/service/login"

	googleProto "google.golang.org/protobuf/proto"
)

func Register() {
	handler := func(message googleProto.Message, subsriberId string, requestId string) error {
		// Handle the message
		loginMessage, ok := message.(*proto.LoginUser)
		if !ok {
			return errors.New("error casting message to LoginUser")
		}

		token, err := login.LoginUser(loginMessage.Username, loginMessage.Password)
		response := "SUCCESS"
		if err != nil {
			println("Error logging in user:", err.Error())
			response = "FAILURE"
		}
		producer.Publish("v1.account.login.result", &proto.LoginResponse{LoginRequestId: loginMessage.LoginRequestId, Status: response, Token: token}, subsriberId, requestId)
		return nil
	}

	consumer.StartKafkaConsumer([]string{"localhost:9092"}, "v1.account.login", "auth-service", &proto.LoginUser{}, handler)
}
