package command

import (
	"bytes"
	"fmt"
	"github.com/mitchellh/cli"
)

// VersionCommand is a Command implementation that prints the version of
// consul-check
type VersionCommand struct {
	Version string
	Ui      cli.Ui
}

// Help returns a string that is the usage for the VersionCommand
func (v *VersionCommand) Help() string {
	return ""
}

// Run is a function mapped to the VersionCommand implementation.
// This is invoked upon calling `consul-check <-v|--version|version>`.
func (v *VersionCommand) Run(_ []string) int {
	var versionString bytes.Buffer
	fmt.Fprintf(&versionString, "consul-check v%s", v.Version)

	v.Ui.Output(versionString.String())
	return 0
}

// Synopsis of the VersionCommand implementation.
func (v *VersionCommand) Synopsis() string {
	return "Prints the consul-check version"
}
