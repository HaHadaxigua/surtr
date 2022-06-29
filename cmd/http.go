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
		Action: func(c *cli.Context) error {
			http.New().Start(c.Context)
			return nil
		},
	}
}
