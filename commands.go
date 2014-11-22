package main

import (
	"github.com/mitchellh/cli"
	"github.com/srspnda/consul-check/command"
	"os"
)

// Commands is the mapping of all the available Consul check commands.
var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{Writer: os.Stdout}

	Commands = map[string]cli.CommandFactory{
		"disk": func() (cli.Command, error) {
			return &command.DiskCommand{
				Ui: ui,
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Version: Version,
				Ui:      ui,
			}, nil
		},
	}
}
