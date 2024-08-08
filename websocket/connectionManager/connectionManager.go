package connectionManager

import (
	"sync"
	"PlantCare/websocket/wsTypes"
)

type ConnectionManager struct {
	mu          sync.RWMutex
	connections map[string]*wsTypes.Connection
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*wsTypes.Connection),
	}
}

// AddConnection adds a connection to the manager with the given pot ID.
func (cm *ConnectionManager) AddConnection(potID string, conn *wsTypes.Connection) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.connections[potID] = conn
}

// RemoveConnection removes a connection from the manager.
func (cm *ConnectionManager) RemoveConnection(potID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.connections, potID)
}

// GetConnection retrieves a connection by pot ID.
func (cm *ConnectionManager) GetConnection(potID string) (*wsTypes.Connection, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	conn, ok := cm.connections[potID]
	return conn, ok
}

// GetAllConnections returns all connections.
func (cm *ConnectionManager) GetAllConnections() map[string]*wsTypes.Connection {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	conns := make(map[string]*wsTypes.Connection)
	for k, v := range cm.connections {
		conns[k] = v
	}
	return conns
}
