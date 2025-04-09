package main

import (
	notifyKafkaController "websocket/service/kafka/notify-kafka-controller"
	"websocket/socket"
)

func main() {
	var pool = socket.NewWebSocketPool()
	notifyKafkaController.Register(pool)
	socket.StartWebSocketServer(pool, "localhost:8081")
}
