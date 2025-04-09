package register

import (
	"projects/auth-service/service/kafka/producer"
	"projects/auth-service/proto"
)

func RegisterUser(username string, password string, subscriberId string, requestId string) {
	// Simulate user registration logic
	println("Registering user:", username)
	println("Password:", password)
	println("Subscriber ID:", subscriberId)
	println("Request ID:", requestId)

	// Here you would typically save the user to a database
	// For this example, we'll just print the details
	sendSuccess(username, subscriberId, requestId)
}

func sendError(message string, subscriberId string, requestId string) {
	sendStatus("v1.account.register.status.failure", message, subscriberId, requestId)
}

func sendSuccess(message string, subscriberId string, requestId string) {
	sendStatus("v1.account.register.status.success", message, subscriberId, requestId)
}

func sendStatus(topic string, message string, subscriberId string, requestId string) {
	producer.Publish(topic, &proto.StatusMessage{Message: message}, subscriberId, requestId)
}