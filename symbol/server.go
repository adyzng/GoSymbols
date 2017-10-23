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
	ss.lck.Lock()
	defer ss.lck.Unlock()

	ss.builders = append(ss.builders, branch)
	return branch
}

// Delete builder
//
func (ss *sserver) DeleteBuilder(storeName string) Builder {
	ss.lck.Lock()
	defer ss.lck.Unlock()

	for i, b := range ss.builders {
		if b.Name() == storeName {
			if i == len(ss.builders)-1 {
				ss.builders = ss.builders[:i]
			} else {
				bs := ss.builders[:]
				ss.builders = append(bs[:i], bs[i+1:]...)
			}
			b.Delete()
			return b
		}
	}
	return nil
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
				if err := b.Init(); err != nil {
					log.Warn("[SS] Failed to load branch %s: %v.", b.Name(), err)
				} else {
					log.Info("[SS] Load branch %s.", b.Name())
				}
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
	ss.WalkBuilders(func(bu Builder) error {
		wg.Add(1)
		log.Info("[SS] Parse branch %s.", bu.Name())
		go func() {
			defer wg.Done()
			bu.ParseBuilds(nil)
		}()
		return nil
	})
	wg.Wait()

LOOP:
	for {
		ss.WalkBuilders(func(bu Builder) error {
			if bu.CanUpdate() {
				go func() {
					wg.Add(1)
					defer wg.Done()
					log.Trace("[SS] Trigger branch %s.", bu.Name())
					bu.AddBuild("")
				}()
			} else {
				log.Trace("[SS] Can't update branch %s.", bu.Name())
			}
			return nil
		})

		select {
		case <-done:
			log.Warn("[SS] Receive stop signal.")
			break LOOP
		case <-ticker.C:
			wg.Wait()
			break
		}
	}
	ss.WalkBuilders(func(bu Builder) error {
		if err := bu.Persist(); err != nil {
			log.Error(2, "[SS] Save branch failed %v.", err)
		}
		return nil
	})

	wg.Wait()
	log.Info("[SS] Symbol server stop.")
}
