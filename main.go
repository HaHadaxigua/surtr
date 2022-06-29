package main

import (
	"fmt"
	"os"

	. "github.com/HaHadaxigua/surtr/cmd"
	. "github.com/HaHadaxigua/surtr/global"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	BuildTime   = ""
	BuildNumber = ""
	GitCommit   = ""
	Version     = "1.0.0"
)

func main() {
	logrus.Infof("WELCOME %s", Surtr)
	if err := newApp().Run(os.Args); err != nil {
		logrus.Fatal("failed to init cli")
	}
}

func newApp() *cli.App {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer,
			"version:    %s\n"+
				"Git Commit: %s\n"+
				"Build Time: %s\n"+
				"Build:      %s\n",
			c.App.Version, GitCommit, BuildTime, BuildNumber)
	}
	return &cli.App{
		Name:                 Surtr,
		Version:              Version,
		Usage:                "Setup your environment",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			NewHttpCommand(),
		},
	}
}
