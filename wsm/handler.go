package wsm

import (
	"errors"
	"fmt"
	"sync"

	"rlp-member-service/log"

	"github.com/gorilla/websocket"
)

type WebSocketManager struct {
	clients map[*websocket.Conn]string
	manset  map[string]map[*websocket.Conn]struct{} // Use a map instead of a slice to improve lookup and deletion efficiency
	mu      sync.RWMutex                            // Read-write lock
}

// Global WebSocketManager
var wsManager = &WebSocketManager{
	clients: make(map[*websocket.Conn]string),
	manset:  make(map[string]map[*websocket.Conn]struct{}),
}

// Get the WebSocketManager singleton
func RetrieveWsManager() *WebSocketManager {
	return wsManager
}

func (wm *WebSocketManager) Stat() (allConn, auditAllConn int) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	allConn = len(wm.clients)
	total := 0
	for _, connections := range wm.manset {
		total += len(connections)
	}
	auditAllConn = total
	return allConn, auditAllConn
}

// Add WebSocket connection
func (wm *WebSocketManager) AddClient(chain, ca string, ws *websocket.Conn) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	key := buildKey(chain, ca)

	// If the key does not exist, initialize the map
	if _, exists := wm.manset[key]; !exists {
		wm.manset[key] = make(map[*websocket.Conn]struct{})
	}
	wm.manset[key][ws] = struct{}{} // Use a map for deduplication
	wm.clients[ws] = key
}

// Delete WebSocket connection
func (wm *WebSocketManager) RemoveClient(ws *websocket.Conn) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	// Get key
	key, exists := wm.clients[ws]
	if !exists {
		return
	}

	// Delete the WebSocket connection from manset[key]
	if connections, ok := wm.manset[key]; ok {
		delete(connections, ws)
		if len(connections) == 0 {
			delete(wm.manset, key) // Clean up empty keys
		}
	}

	// Delete the clients mapping
	delete(wm.clients, ws)
}

// Broadcast message to the specified chain and CA
func (wm *WebSocketManager) Broadcast(chain, ca string, msg string) {
	wm.mu.RLock() // Only a read lock is needed
	defer wm.mu.RUnlock()
	key := buildKey(chain, ca)

	// Get all WebSocket connections under this key
	connections, exists := wm.manset[key]
	if !exists {
		return
	}

	// Send message
	for ws := range connections {
		if err := ws.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			log.Error("Failed to send message:", err)
			ws.Close()
			wm.RemoveClient(ws) // Remove after closing the connection
		}
	}
}

// Send an individual message to the specified WebSocket connection
func (wm *WebSocketManager) SendToClient(ws *websocket.Conn, msg string) error {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	// Check if the connection exists
	if _, exists := wm.clients[ws]; !exists {
		return errors.New("WebSocket client not found")
	}
	return ws.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (wm *WebSocketManager) BroadcastToAll(msg []byte) error {
	wm.mu.RLock()
	connections := make([]*websocket.Conn, 0, len(wm.clients))
	for ws := range wm.clients {
		connections = append(connections, ws) // Copy connection
	}
	wm.mu.RUnlock()

	var failedConnections []*websocket.Conn

	for _, ws := range connections {
		err := ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			failedConnections = append(failedConnections, ws) // Log failed connections
		}
	}

	if len(failedConnections) > 0 {
		wm.mu.Lock()
		for _, ws := range failedConnections {
			delete(wm.clients, ws) // Remove invalid WebSocket connections from clients
			ws.Close()
		}
		wm.mu.Unlock()
	}

	if len(failedConnections) > 0 {
		return fmt.Errorf("failed to send message to %d clients", len(failedConnections))
	}

	return nil
}

// Generate  key
func buildKey(chain, ca string) string {
	return fmt.Sprintf("%s-%s", chain, ca)
}
