package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/adyzng/goPDB/config"
	"github.com/adyzng/goPDB/symbol"

	log "gopkg.in/clog.v1"
)

func init() {
	log.New(log.FILE, log.FileConfig{
		Level:      log.INFO,
		Filename:   "log/gPDB.log",
		BufferSize: 2048,
		FileRotationConfig: log.FileRotationConfig{
			Daily:   true,
			Rotate:  true,
			MaxDays: 30,
			MaxSize: 50 * (1 << 20),
		},
	})
}

func main() {
	log.Info("[App] Start %s ...", config.AppName)

	done := make(chan struct{}, 1)
	serv := http.Server{
		Addr:         config.Address,
		Handler:      NewRouter(),
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		log.Info("[App] Listening %s", serv.Addr)
		serv.ListenAndServe()
		wg.Done()
	}()
	go func() {
		symbol.GetServer().Run(done)
		wg.Done()
	}()
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, os.Kill)

		<-sigs
		close(done)
		close(sigs)
		log.Info("[App] Receive terminate signal!")

		serv.Shutdown(context.Background())
		wg.Done()
	}()

	wg.Wait()
	log.Info("[App] Stop %s.", config.AppName)
	log.Shutdown()
}
