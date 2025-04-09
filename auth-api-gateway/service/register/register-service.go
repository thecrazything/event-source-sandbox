package register

import (
	"projects/auth-api-gateway/proto"
	"projects/auth-api-gateway/service/kafka"
)

func Register(username string, password string, subscriberId string, requestId string) {
	account := proto.RegisterUser{
		Username: username,
		Password: password,
	}

	// TODO validation

	kafka.Publish("v1.account.register", &account, subscriberId, requestId)
}
