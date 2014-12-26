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
	var warnLevel, critLevel int

	cmdFlags := flag.NewFlagSet("disk", flag.ContinueOnError)
	cmdFlags.Usage = func() { d.Ui.Output(d.Help()) }
	cmdFlags.IntVar(&warnLevel, "warn", 85, "disk usage is level WARN")
	cmdFlags.IntVar(&critLevel, "crit", 95, "disk usage is level CRIT")

	if err := cmdFlags.Parse(args); err != nil {
		return 0
	}

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

	fsList := sigar.FileSystemList{}
	fsList.Get()

	exitCode := 0

	for _, fs := range fsList.List {
		dirName := fs.DirName

		fsUsage := sigar.FileSystemUsage{}
		fsUsage.Get(dirName)

		usedPercent := fsUsage.UsePercent()

		fmt.Fprintf(
			os.Stdout,
			outputFormat,
			fs.DevName,
			formatSize(fsUsage.Total),
			formatSize(fsUsage.Used),
			formatSize(fsUsage.Avail),
			sigar.FormatPercent(usedPercent),
			dirName,
		)

		switch level := int(usedPercent); {
		case level > critLevel:
			exitCode = 2
		case level > warnLevel && exitCode < 2:
			exitCode = 1
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

  -warn=<level>   Percent disk usage for check to be WARN
  -crit=<level>   Percent disk usage for check to be CRIT
`

	return strings.TrimSpace(helpText)
}

// Synopsis of the DiskCommand implementation.
func (d *DiskCommand) Synopsis() string {
	return "Checks the local system for disk usage"
}
