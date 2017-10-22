package symbol

import (
	"io/ioutil"
	"sync"
	"time"

	"github.com/adyzng/goPDB/config"
	log "gopkg.in/clog.v1"
)

var (
	symSvr = &sserver{}
	once   sync.Once
)

// sserver ...
//
type sserver struct {
	lck      sync.RWMutex
	builders []Builder
}

// GetServer return single instance of sserver
//
func GetServer() *sserver {
	return symSvr
}

// GetBuilderList get branch list
//
func (ss *sserver) GetBuilderList() []Builder {
	return ss.builders[:]
}

// GetBuilder reture given build if not exist, return nil
//
func (ss *sserver) GetBuilder(storeName string) Builder {
	ss.lck.RLock()
	defer ss.lck.RUnlock()

	for _, v := range ss.builders {
		if v.Name() == storeName {
			return v
		}
	}
	return nil
}

// AddBuilder if already exist, do nothing.
//
func (ss *sserver) AddBuilder(buildName, storeName string) Builder {
	if b := ss.GetBuilder(storeName); b != nil {
		return b
	}

	branch := NewBranch(buildName, storeName)
	if err := branch.Init(); err != nil {
		if err = branch.SetSubpath("", ""); err != nil {
			return nil
		}
	}

	ss.lck.Lock()
	defer ss.lck.Unlock()

	ss.builders = append(ss.builders, branch)
	return branch
}

// WalkBuilders walk all exist builders, the handler should be return asap.
//
func (ss *sserver) WalkBuilders(handler func(branch Builder) error) error {
	var err error
	if handler == nil {
		return nil
	}
	ss.lck.RLock()
	defer ss.lck.RUnlock()

	for _, b := range ss.builders {
		if err = handler(b); err != nil {
			break
		}
	}
	return err
}

// LoadBuilder scan local symbol store for exist branchs.
//
func (ss *sserver) LoadBuilders() error {
	root := config.Destination
	fi, err := ioutil.ReadDir(root)
	if err != nil {
		log.Error(2, "[SS] Enum symbol store %s failed: %v.", root, err)
		return err
	}
	for _, f := range fi {
		if f.IsDir() {
			b := ss.AddBuilder(f.Name(), f.Name())
			if b != nil {
				log.Info("[SS] Load branch %s.", b.Name())
			}
		}
	}
	return nil
}

// Run ...
//
func (ss *sserver) Run(done <-chan struct{}) {
	var wg sync.WaitGroup
	log.Info("[SS] Symbol server start ...")

	ticker := time.NewTicker(time.Hour * 2)
	defer ticker.Stop()

	ss.LoadBuilders()
	log.Info("[SS] There are %d branches.", len(ss.builders))

LOOP:
	for {
		ss.lck.RLock()
		//log.Info("[SS] Branchs Count %d.", len(ss.builders))

		for idx, b := range ss.builders {
			if !b.CanUpdate() {
				continue
			}
			wg.Add(1)
			go func(build Builder, idx int) {
				defer wg.Done()
				log.Info("[SS] Trigger %d: %v.", idx, build.Name())
				build.Add("")
			}(b, idx)
		}
		ss.lck.RUnlock()

		select {
		case <-done:
			log.Warn("[SS] Receive stop signal.")
			break LOOP
		case <-ticker.C:
			break
		}
	}

	wg.Wait()
	log.Info("[SS] Symbol server stop.")
}
