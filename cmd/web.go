package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/adyzng/GoSymbols/config"
	"github.com/adyzng/GoSymbols/route"
	"github.com/adyzng/GoSymbols/symbol"
	"github.com/urfave/cli"

	log "gopkg.in/clog.v1"
)

// Web server subcommand
//
var Web = cli.Command{
	Name:        "serve",
	Usage:       "Start symbol store server",
	Description: "Symbol server will take care of symbols, and serve the web portal.",
	Action:      runWeb,
	Flags: []cli.Flag{
		stringFlag("port, p", "3000", "Given port number to prevent conflict"),
		stringFlag("config, c", "config/server.ini", "Custom configuration file path"),
	},
}

func runWeb(c *cli.Context) error {
	if c.IsSet("config") {
		config.LoadConfig(c.String("config"))
	}
	if c.IsSet("port") {
		config.Port = c.Uint("port")
	}

	done := make(chan struct{}, 1)
	serv := http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Address, config.Port),
		Handler:      route.NewRouter(),
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	log.Info("[App] Start %s ...", config.AppName)
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		log.Info("[App] Listening %s", serv.Addr)
		serv.ListenAndServe()
	}()
	go func() {
		defer wg.Done()
		symbol.GetServer().Run(done)
	}()
	go func() {
		defer wg.Done()
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, os.Kill)

		<-sigs
		close(done)
		close(sigs)
		log.Info("[App] Receive terminate signal!")
		serv.Shutdown(context.Background())
	}()

	wg.Wait()
	log.Info("[App] Stopped %s.", config.AppName)
	return nil
}
