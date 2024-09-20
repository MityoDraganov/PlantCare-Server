package connectionManager

import (
	"PlantCare/websocket/wsTypes"
	"sync"
)

// ConnectionManager holds the WebSocket Connections
type ConnectionManager struct {
	Mu          sync.RWMutex
	Connections map[string]*wsTypes.Connection
}

// Exported global instance of ConnectionManager
var ConnManager = &ConnectionManager{
	Connections: make(map[string]*wsTypes.Connection),
	Mu: sync.RWMutex{},
}

func (cm *ConnectionManager) AddConnection(potID string, conn wsTypes.Connection) {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	cm.Connections[potID] = &conn
}

// RemoveConnection removes a WebSocket connection
func (cm *ConnectionManager) RemoveConnection(potID string) {
	cm.Mu.Lock()
	defer cm.Mu.Unlock()
	delete(cm.Connections, potID)
}

// GetConnection retrieves a connection by pot ID
func (cm *ConnectionManager) GetConnection(potID string) (*wsTypes.Connection, bool) {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()
	conn, ok := cm.Connections[potID]
	return conn, ok
}

// GetAllConnections returns all active Connections
func (cm *ConnectionManager) GetAllConnections() map[string]*wsTypes.Connection {
	cm.Mu.RLock()
	defer cm.Mu.RUnlock()
	conns := make(map[string]*wsTypes.Connection)
	for k, v := range cm.Connections {
		conns[k] = v
	}
	return conns
}
