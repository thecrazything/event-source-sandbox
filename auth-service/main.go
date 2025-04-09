package main

import (
	loginKafkaController "projects/auth-service/service/kafka/login-kafka-controller"
	registerKafkaController "projects/auth-service/service/kafka/register-kafka-controller"
)

func main() {
	go func() {
		registerKafkaController.Register()
	}()
	loginKafkaController.Register()
}
