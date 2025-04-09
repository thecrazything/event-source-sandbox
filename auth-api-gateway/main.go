package main

import (
	"projects/auth-api-gateway/api"
	"projects/auth-api-gateway/service/session"
)

func main() {
	go func() {
		session.ListenForNewSessions()
	}()
	api.Start("localhost:8080")
}
