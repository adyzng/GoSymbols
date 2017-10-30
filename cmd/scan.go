package cmd

import (
	"github.com/adyzng/GoSymbols/config"
	"github.com/adyzng/GoSymbols/symbol"
	"github.com/urfave/cli"

	log "gopkg.in/clog.v1"
)

// Admin ...
var Admin = cli.Command{
	Name:        "scan",
	Usage:       "Scan the exist symbol store",
	Description: "Scan the exist symbol store, and generate the config.",
	Action:      runAdmin,
	Flags: []cli.Flag{
		stringFlag("path, p", `D:\SymbolServer`, "Exist symbol store path"),
	},
}

func runAdmin(c *cli.Context) error {
	path := config.Destination
	if c.IsSet("path") {
		path = c.String("path")
		log.Trace("[App] Scan path %s", path)
	}

	log.Info("[App] Scan store %s", path)
	ss := symbol.GetServer()

	ss.WalkBuilders(func(b symbol.Builder) error {
		log.Info("[App] Load branch %s.", b.Name())
		return nil
	})

	return ss.ScanStore(path)
}
