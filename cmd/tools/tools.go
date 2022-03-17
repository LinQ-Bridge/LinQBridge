package tools

import "github.com/urfave/cli"

var CMD = cli.Command{
	Name:    "tool",
	Aliases: []string{"t"},
	Usage:   "linq Tool",
	Subcommands: []cli.Command{
		ConfigCMD,
		DeployCMD,
		GenesisCMD,
		NodekeyCMD,
	},
}
