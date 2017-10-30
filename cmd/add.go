package cmd

import (
	"errors"

	"github.com/adyzng/GoSymbols/symbol"
	"github.com/urfave/cli"

	log "gopkg.in/clog.v1"
)

// AddBuild ...
var AddBuild = cli.Command{
	Name:        "add",
	Usage:       "Add given build for specified branch.",
	Description: "Add the given build of specified branch exist in symbol store.",
	Action:      addBuild,
	Flags: []cli.Flag{
		stringFlag("branch, b", "", "The branch name in the symbol store."),
		stringFlag("version, v", "", "The build version, empty version for the latest build."),
	},
}

func addBuild(c *cli.Context) error {
	bname := ""
	if c.IsSet("branch") {
		bname = c.String("branch")
	}
	if bname == "" {
		return errors.New("empty branch name")
	}

	build := ""
	if c.IsSet("build") {
		build = c.String("build")
	}

	log.Info("[App] Add build %s for branch %s", build, bname)
	ss := symbol.GetServer()
	if err := ss.LoadBranchs(); err != nil {
		return err
	}

	builder := ss.Get(bname)
	if builder == nil {
		log.Warn("[App] Branch %s not exist.", bname)
		return errors.New("branch not exist")
	}

	return builder.AddBuild(build)
}
