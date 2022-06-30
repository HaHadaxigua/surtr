package cmd

import (
	"github.com/HaHadaxigua/surtr/internal/http"

	"github.com/urfave/cli/v2"
)

// NewHttpCommand start http server
func NewHttpCommand() *cli.Command {
	return &cli.Command{
		Name:        "http",
		Description: "start http server",
		Usage:       "Surtr http",
		Aliases:     []string{"h"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "conf",
				Aliases:     []string{"c"},
				DefaultText: "config http service",
				Usage:       "surtr http -conf xxx.yaml",
				Required:    false,
			},
		},
		Action: func(c *cli.Context) error {
			scv, err := http.New(c.String("c"))
			if err != nil {
				return err
			}
			scv.Start(c.Context)
			return nil
		},
	}
}
