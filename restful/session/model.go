package session

import (
	"sync"
	"time"

	log "gopkg.in/clog.v1"
)

// StorType ...
type StorType int

const (
	SessTimeout  = time.Hour * 24
	CookieSessID = "session_id"
	CookieMaxAge = time.Hour * 24 / time.Second
)
const (
	_          StorType = iota
	MemStore            // memory store
	RedisStore          // redis store
)

var (
	sessMgr *SessManager
	once    sync.Once
)

// SessStore interface
//
type SessStore interface {
	// Get session data
	Get(id string) interface{}
	// Set session data
	Set(id string, data interface{}) error
	// Update sessions if timeout, delete it
	Udpate(max int) int

	// Delete session
	Delete(id string) interface{}
	// Create an new session
	Create(data interface{}) string
}

// SessManager manage all sessions
//
type SessManager struct {
	SessStore
}

// GetManager create an session manager
//
func GetManager(typ ...StorType) *SessManager {
	once.Do(func() {
		typ = append(typ, MemStore)
		switch typ[0] {
		case MemStore:
			sessMgr = &SessManager{
				SessStore: NewMemStore(),
			}
			go sessMgr.scan()
		case RedisStore:
			panic("redis session store isn't implement yet")
		default:
			panic("unknown session storage type")
		}
	})
	return sessMgr
}

func (s *SessManager) scan() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	log.Trace("[Session] Session updater starting ...")
	for {
		select {
		case <-ticker.C:
			s.Udpate(1000)
		}
	}
}
