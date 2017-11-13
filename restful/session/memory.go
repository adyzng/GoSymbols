package session

import (
	"errors"
	"sync"
	"time"

	"github.com/adyzng/GoSymbols/restful/uuid"
)

// MemoryStore store sessin in memory
//
type MemoryStore struct {
	mx   sync.RWMutex
	sess map[string]*sessData
}
type sessData struct {
	expireAt int64
	data     interface{}
}

// NewMemStore ...
//
func NewMemStore() SessStore {
	return &MemoryStore{
		sess: make(map[string]*sessData),
	}
}

// Get session
func (m *MemoryStore) Get(id string) interface{} {
	m.mx.RLock()
	defer m.mx.RUnlock()
	if d, ok := m.sess[id]; ok {
		return d.data
	}
	return nil
}

// Set session data
func (m *MemoryStore) Set(id string, data interface{}) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	if sd, ok := m.sess[id]; ok {
		sd.expireAt = time.Now().Add(SessTimeout).Unix()
		sd.data = data
		return nil
	}
	return errors.New("session not exist")
}

// Delete session
func (m *MemoryStore) Delete(id string) interface{} {
	m.mx.Lock()
	defer m.mx.Unlock()
	if data, ok := m.sess[id]; ok {
		delete(m.sess, id)
		return data.data
	}
	return nil
}

// Create new session
func (m *MemoryStore) Create(data interface{}) string {
	id := uuid.NewUUID()
	m.mx.Lock()
	defer m.mx.Unlock()
	m.sess[id] = &sessData{
		data:     data,
		expireAt: time.Now().Add(SessTimeout).Unix(),
	}
	return id
}

// Udpate all sessions, if timeout, delete it
// update `max` items at most each time.
func (m *MemoryStore) Udpate(max int) int {
	m.mx.Lock()
	defer m.mx.Unlock()

	total := max
	nowUnix := time.Now().Unix()

	for key, val := range m.sess {
		if val.expireAt < nowUnix {
			delete(m.sess, key)
		}
		if total--; total == 0 {
			break
		}
	}
	return max
}
