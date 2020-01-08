package event

import (
	"sync"
)

// Manager ...
type Manager struct {
	sync.RWMutex

	listeners map[Type][]*Listener
}

// NewManager ...
func NewManager() *Manager {
	return &Manager{
		listeners: make(map[Type][]*Listener),
	}
}

// Register ...
func (m *Manager) Register(lis *Listener) {
	m.Lock()
	defer m.Unlock()
	ml := m.listeners[lis.Type()]
	ml = append(ml, lis)
}

// Remove ...
func (m *Manager) Remove(lis *Listener) {
	m.Lock()
	defer m.Unlock()
}

// Notify ...
func (m *Manager) Notify(e *Event) {

}
