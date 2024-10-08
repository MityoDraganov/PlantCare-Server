package otaManager

import (
	"sync"
)

type PendingOTA struct {
	PotID      string   // Pot ID that the OTA is pending for
	DriverURLs []string // List of driver URLs for the OTA update
}

type PendingOTAManager struct {
	mu          sync.RWMutex
	pendingOTAs map[string]PendingOTA // Track potID and its corresponding PendingOTA info
}

// Global instance of the PendingOTAManager
var OTAManager = &PendingOTAManager{
	pendingOTAs: make(map[string]PendingOTA),
}

// IsOTAPending checks if there's a pending OTA update for a given potID
func (p *PendingOTAManager) IsOTAPending(potID string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	_, exists := p.pendingOTAs[potID]
	return exists
}

// AddOTAPending adds a pending OTA for a given potID with associated driver URLs
func (p *PendingOTAManager) AddOTAPending(potID string, driverURLs []string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pendingOTAs[potID] = PendingOTA{
		PotID:      potID,
		DriverURLs: driverURLs,
	}
}

// RemoveOTAPending removes the pending OTA status for the given potID
func (p *PendingOTAManager) RemoveOTAPending(potID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.pendingOTAs, potID)
}

// GetPendingOTA retrieves the pending OTA details (including potID and driver URLs) for a given potID
func (p *PendingOTAManager) GetPendingOTA(potID string) (PendingOTA, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	ota, exists := p.pendingOTAs[potID]
	return ota, exists
}
