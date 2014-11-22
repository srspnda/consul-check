package command

import (
	"flag"
	"fmt"
	"github.com/cloudfoundry/gosigar"
	"github.com/mitchellh/cli"
	"os"
	"strings"
)

const outputFormat = "%-15s %4s %4s %5s %4s %-15s\n"

func formatSize(size uint64) string {
	return sigar.FormatSize(size * 1024)
}

// DiskCommand is a Command implementation that checks local partitions for
// usage levels
type DiskCommand struct {
	Ui cli.Ui
}

// Run is the function mapped to the DiskCommand implementation.
// This is invoked upon calling `consul-check disk ...`
func (d *DiskCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("disk", flag.ContinueOnError)
	cmdFlags.Usage = func() { d.Ui.Output(d.Help()) }

	var warning, failure int
	cmdFlags.IntVar(&warning, "warning", 85, "disk check is warning")
	cmdFlags.IntVar(&failure, "failure", 95, "disk check is failure")

	// Set exitCode to 0 as default exit code. Set to 1 if something is
	// greater than or equal to warning. Set to 2 if something is greater
	// than or equal to critical.
	var exitCode = 0

	fsList := sigar.FileSystemList{}
	fsList.Get()

	fmt.Fprintf(
		os.Stdout,
		outputFormat,
		"Filesystem",
		"Size",
		"Used",
		"Avail",
		"Use%",
		"Mounted on",
	)

	for _, fs := range fsList.List {
		dirName := fs.DirName

		usage := sigar.FileSystemUsage{}
		usage.Get(dirName)

		used := int(usage.UsePercent())

		fmt.Fprintf(
			os.Stdout,
			outputFormat,
			fs.DevName,
			formatSize(usage.Total),
			formatSize(usage.Used),
			formatSize(usage.Avail),
			sigar.FormatPercent(usage.UsePercent()),
			dirName,
		)

		switch exitCode {
		case 0:
			if used >= critical {
				exitCode = 2
			} else if used >= warning {
				exitCode = 1
			}
		case 1:
			if used >= critical {
				exitCode = 2
			}
		}
	}

	return exitCode
}

// Help returns a string that is the usage for the DiskCommand
func (d *DiskCommand) Help() string {
	helpText := `
Usage: consul-check disk <options>

	Check local partitions for usage levels

Options:

  -warning=<level>   Percent usage for a disk check to be at level warning.
                  This will return to Consul an exit code of 1.
  -failure=<level>   Percent usage for a disk check to be at level failure.
                  This will return to Consul an exit code of 2.
`

	return strings.TrimSpace(helpText)
}

// Synopsis of the DiskCommand implementation.
func (d *DiskCommand) Synopsis() string {
	return "Checks the local system for disk usage"
}
