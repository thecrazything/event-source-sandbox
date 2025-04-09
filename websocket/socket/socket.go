package socket

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebSocketPool struct {
	connections map[string]*websocket.Conn
	mu          sync.Mutex
}

func NewWebSocketPool() *WebSocketPool {
	return &WebSocketPool{
		connections: make(map[string]*websocket.Conn),
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (pool *WebSocketPool) handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}

	id := uuid.New().String()

	pool.mu.Lock()
	pool.connections[id] = conn
	pool.mu.Unlock()

	// Send the ID to the client
	initialMessage := map[string]string{"subscriptionId": id}
	if err := conn.WriteJSON(initialMessage); err != nil {
		pool.removeConnection(id)
		return
	}

	go pool.listenForMessages(id, conn)
}

func (pool *WebSocketPool) listenForMessages(id string, conn *websocket.Conn) {
	defer pool.removeConnection(id)

	// Read a single message and terminate the connection
	_, _, err := conn.ReadMessage()
	if err == nil {
		conn.Close()
	}
}

func (pool *WebSocketPool) HasConnection(id string) bool {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	_, exists := pool.connections[id]
	return exists
}

func (pool *WebSocketPool) SendMessage(id string, message interface{}) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	conn, exists := pool.connections[id]
	if !exists {
		return fmt.Errorf("connection with ID %s not found", id)
	}

	return conn.WriteJSON(message)
}

func (pool *WebSocketPool) removeConnection(id string) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if conn, exists := pool.connections[id]; exists {
		conn.Close()
		delete(pool.connections, id)
	}
}

func StartWebSocketServer(pool *WebSocketPool, addr string) error {
	http.HandleFunc("/ws", pool.handleConnection)
	fmt.Printf("WebSocket server started at %s\n", addr)
	http.ListenAndServe(addr, nil)
	return nil
}
