package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adyzng/GoSymbols/cmd"
	"github.com/adyzng/GoSymbols/config"
	"github.com/urfave/cli"

	log "gopkg.in/clog.v1"
)

func init() {
	fpath, _ := filepath.Abs(config.LogPath)
	if err := os.MkdirAll(filepath.Dir(fpath), 666); err != nil {
		fmt.Printf("[App] Create log folder failed: %v.", err)
	}

	log.New(log.FILE, log.FileConfig{
		Level:      log.TRACE,
		Filename:   filepath.Join(fpath, "app.log"),
		BufferSize: 2048,
		FileRotationConfig: log.FileRotationConfig{
			Rotate:  true,
			MaxDays: 30,
			MaxSize: 50 * (1 << 20),
		},
	})
}

const APP_VER = "0.0.0.1"

func main() {
	app := cli.NewApp()
	app.Name = config.AppName
	app.Usage = "A self-service symbol store"
	app.Version = APP_VER
	app.Commands = []cli.Command{
		cmd.Web,
		cmd.Admin,
		cmd.AddBuild,
	}

	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
	log.Shutdown()
}
