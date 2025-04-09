package session

import (
	"errors"
	"fmt"
	"projects/auth-api-gateway/proto"
	"projects/auth-api-gateway/service/kafka"
	"projects/auth-api-gateway/service/kafka/consumer"
	"projects/auth-api-gateway/service/redis"

	googleProto "google.golang.org/protobuf/proto"
)

type Session struct {
	SessionKey string
}

type SessionStatus int

const (
	WAITING SessionStatus = iota
	VALID
	INVALID
)

func FindRegisteredSession(requestId string) (SessionStatus, Session) {
	client := redis.GetRedisClient()
	value, err := client.Get("v1.account.session." + requestId)
	if err != nil {
		fmt.Printf("Error getting session: %v\n", err)
		return WAITING, Session{}
	}
	return VALID, Session{
		SessionKey: value,
	}
}

func ClearRegisteredSession(requestId string) {
	client := redis.GetRedisClient()
	client.Delete("v1.account.session.status." + requestId)
}

func ListenForNewSessions() {
	handler := func(message googleProto.Message, subscriberId string, requestId string) error {
		// Handle the message
		loginResponse, ok := message.(*proto.LoginResponse)
		if !ok {
			fmt.Printf("Error casting message to LoginResponse: %v %v \n", message, ok)
			return errors.New("error casting message to LoginResponse")
		}

		responseTopic := "v1.status.success"
		client := redis.GetRedisClient()
		if loginResponse.Status == "SUCCESS" {
			println("Login successful, setting session in redis")
			err := client.Set("v1.account.session."+loginResponse.LoginRequestId, loginResponse.Token, 0)
			if err != nil {
				fmt.Printf("Error setting session in redis: %v\n", err)
				return err
			}
		} else {
			responseTopic = "v1.status.failure"
		}

		println("")
		kafka.Publish(responseTopic, &proto.StatusMessage{Message: loginResponse.Status}, subscriberId, requestId)

		return nil
	}

	consumer.StartKafkaConsumer([]string{"localhost:9092"}, "v1.account.login.result", "auth-api-gateway", &proto.LoginResponse{}, handler)
}
