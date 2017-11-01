package symbol

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/adyzng/GoSymbols/config"
	log "gopkg.in/clog.v1"
)

var (
	symSvr *sserver
	once   sync.Once
)

const (
	symConfig = "symbols.json"
)

// sserver ...
//
type sserver struct {
	lck      sync.RWMutex
	builders map[string]Builder
}

// GetServer return single instance of sserver
//
func GetServer() *sserver {
	once.Do(func() {
		symSvr = &sserver{
			builders: make(map[string]Builder, 1),
		}
	})
	return symSvr
}

// Scan exist symbol store
func (ss *sserver) ScanStore(path string) error {
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		log.Error(2, "[SS] Enum symbol store %s failed: %v.", path, err)
		return err
	}

	func() {
		ss.lck.Lock()
		defer ss.lck.Unlock()
		for _, f := range fs {
			if !f.IsDir() {
				continue
			}
			b := NewBranch(f.Name(), f.Name())
			if b.CanBrowse() || b.CanUpdate() {
				ss.builders[strings.ToLower(f.Name())] = b
				log.Info("[SS] Load branch %s.", b.Name())
			}
		}
	}()
	return ss.SaveBranchs("")
}

// Modify branch
func (ss *sserver) Modify(branch *Branch) Builder {
	ss.lck.Lock()
	defer ss.lck.Unlock()

	lower := strings.ToLower(branch.StoreName)
	if b, ok := ss.builders[lower]; ok {
		nb := NewBranch2(branch)
		if nb.CanUpdate() || nb.CanBrowse() {
			bb := b.GetBranch()
			*bb = *nb.GetBranch()
			/*
				ob.BuildName = nb.BuildName
				ob.StoreName = nb.StoreName
				ob.BuildPath = nb.BuildPath
				ob.StorePath = nb.StorePath
			*/
			return b
		}
	}
	return nil
}

// Get reture given branch,  if not exist return nil
func (ss *sserver) Get(storeName string) Builder {
	ss.lck.RLock()
	defer ss.lck.RUnlock()

	lower := strings.ToLower(storeName)
	b, _ := ss.builders[lower]
	return b
}

// AddBranch if already exist, do nothing.
func (ss *sserver) Add(buildName, storeName string) Builder {
	ss.lck.Lock()
	defer ss.lck.Unlock()

	// exist one
	if b, ok := ss.builders[storeName]; ok {
		return b
	}

	// new one
	b := NewBranch(buildName, storeName)
	if b.CanBrowse() || b.CanUpdate() {
		ss.builders[strings.ToLower(storeName)] = b
		return b
	}

	return nil
}

// DeleteBranch remove given branch
func (ss *sserver) Delete(storeName string) Builder {
	ss.lck.Lock()
	defer ss.lck.Unlock()

	lower := strings.ToLower(storeName)
	if b, ok := ss.builders[lower]; ok {
		delete(ss.builders, lower)
		return b
	}
	return nil
}

// WalkBuilders walk all exist builders, the handler should be return asap.
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

// LoadBranchs scan local symbol store for exist branchs.
func (ss *sserver) LoadBranchs() error {
	fpath := filepath.Join(config.AppPath, symConfig)
	fd, err := os.OpenFile(fpath, os.O_RDONLY, 666)
	if err != nil {
		log.Error(2, "[SS] Read symbols config file %s failed: %v.", fpath, err)
		return err
	}

	var arr []*Branch
	dec := json.NewDecoder(fd)
	if err := dec.Decode(&arr); err != nil {
		return err
	}

	ss.lck.Lock()
	defer ss.lck.Unlock()

	for _, b := range arr {
		ss.builders[strings.ToLower(b.StoreName)] = NewBranch2(b)
	}
	return nil
}

// SaveBranchs ...
func (ss *sserver) SaveBranchs(path string) error {
	if path == "" {
		path = config.AppPath
	}

	fpath := filepath.Join(path, symConfig)
	fd, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 666)
	if err != nil {
		log.Error(2, "[SS] Open file %s failed: %v.", fpath, err)
		return err
	}

	ss.lck.Lock()
	defer func() {
		fd.Close()
		ss.lck.Unlock()
	}()

	arr := make([]*Branch, 0, len(ss.builders))
	for _, b := range ss.builders {
		arr = append(arr, b.GetBranch())
	}

	enc := json.NewEncoder(fd)
	enc.SetIndent("", "\t")
	return enc.Encode(arr)
}

// Run ...
func (ss *sserver) Run(done <-chan struct{}) {
	var wg sync.WaitGroup
	log.Info("[SS] Symbol server start ...")

	ticker := time.NewTicker(time.Hour * 2)
	defer ticker.Stop()

	if err := ss.LoadBranchs(); err != nil {
		log.Error(2, "[SS] Load branchs failed: %v.", err)
		return
	}

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

		if err := ss.SaveBranchs(""); err != nil {
			log.Error(2, "[SS] Save branchs list failed: %v.", err)
		}

		select {
		case <-done:
			log.Warn("[SS] Receive stop signal.")
			break LOOP
		case <-ticker.C:
			wg.Wait()
			break
		}
	}
	if err := ss.SaveBranchs(""); err != nil {
		log.Error(2, "[SS] Save branchs list failed: %v.", err)
	}

	wg.Wait()
	log.Info("[SS] Symbol server stop.")
}
