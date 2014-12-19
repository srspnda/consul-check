package command

import (
	"flag"
	"fmt"
	"github.com/cloudfoundry/gosigar"
	"github.com/mitchellh/cli"
	"os"
	"runtime"
	"strings"
	"time"
)

// CpuloadCommand is a Command implementation that checks the cpu load of
// the local system
type CpuloadCommand struct {
	Ui cli.Ui
}

// Run is the function mapped to the CpuloadCommand implementation.
// This is invoked upon calling `consul-check cpuload ...`
func (c *CpuloadCommand) Run(args []string) int {
	numCpus := float64(runtime.NumCPU())

	var warnLevel, critLevel float64
	cmdFlags := flag.NewFlagSet("cpuload", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	cmdFlags.Float64Var(&warnLevel, "warn", numCpus*1.25, "CPU load is level WARN")
	cmdFlags.Float64Var(&critLevel, "crit", numCpus*2.5, "CPU load is level CRIT")

	if err := cmdFlags.Parse(args); err != nil {
		return 0
	}

	concreteSigar := sigar.ConcreteSigar{}

	uptime := sigar.Uptime{}
	uptime.Get()

	avg, err := concreteSigar.GetLoadAverage()
	if err != nil {
		fmt.Printf("Failed to get load average")
		return 0
	}

	fmt.Fprintf(
		os.Stdout,
		" %s up %s load average: %.2f, %.2f, %.2f\n",
		time.Now().Format("15:04:05"),
		uptime.Format(),
		avg.One,
		avg.Five,
		avg.Fifteen,
	)

	exitCode := 0

	switch l := avg.One; {
	case l > critLevel:
		exitCode = 2
	case l > warnLevel:
		exitCode = 1
	}

	return exitCode
}

// Help returns a string that is the usage for the CpuloadCommand
func (c *CpuloadCommand) Help() string {
	helpText := `
Usage: consul-check cpuload <options>

	Check CPU load of the local system

Options:

  -warn=<level>   CPU load for check to be WARN
  -crit=<level>   CPU load for check to be CRIT
`

	return strings.TrimSpace(helpText)
}

// Synopsis of the CpuloadCommand implementation.
func (c *CpuloadCommand) Synopsis() string {
	return "Checks the CPU load of the local system"
}
