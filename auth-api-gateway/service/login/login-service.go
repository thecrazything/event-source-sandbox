package login

import (
	"projects/auth-api-gateway/proto"
	"projects/auth-api-gateway/service/kafka"
)

func Login(username string, password string, loginRequestId string, subscriberId string, requestId string) {
	credentials := proto.LoginUser{
		LoginRequestId: loginRequestId,
		Username: username,
		Password: password,
	}
	kafka.Publish("v1.account.login", &credentials, subscriberId, requestId)
}