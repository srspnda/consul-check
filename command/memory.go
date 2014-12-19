package command

import (
	"flag"
	"fmt"
	"github.com/cloudfoundry/gosigar"
	"github.com/mitchellh/cli"
	"os"
	"strings"
)

func format(size uint64) uint64 {
	return size / 1024 / 1024
}

func getPercentUtil(used, free uint64) int {
	return int(float64(used) / (float64(used + free)) * 100)
}

// MemoryCommand is a Command implementation that checks local virtual memory
// utilization levels
type MemoryCommand struct {
	Ui cli.Ui
}

// Run is the function mapped to the MemoryCommand implementation.
// This is invoked upon calling `consul-check memory ...`
func (m *MemoryCommand) Run(args []string) int {
	var warnLevel, critLevel int

	cmdFlags := flag.NewFlagSet("memory", flag.ContinueOnError)
	cmdFlags.Usage = func() { m.Ui.Output(m.Help()) }
	cmdFlags.IntVar(&warnLevel, "warn", 95, "memory util is level WARN")
	cmdFlags.IntVar(&critLevel, "crit", 99, "memory util is level CRIT")

	if err := cmdFlags.Parse(args); err != nil {
		return 0
	}

	mem := sigar.Mem{}
	mem.Get()
	swap := sigar.Swap{}
	swap.Get()

	fmt.Fprintf(
		os.Stdout,
		"%18s %10s %10s\n",
		"total", "used", "free",
	)

	fmt.Fprintf(
		os.Stdout,
		"Mem:    %10d %10d %10d\n",
		format(mem.Total), format(mem.Used), format(mem.Free),
	)

	fmt.Fprintf(
		os.Stdout,
		"-/+ buffers/cache: %10d %10d\n",
		format(mem.ActualUsed), format(mem.ActualFree),
	)

	fmt.Fprintf(
		os.Stdout,
		"Swap:   %10d %10d %10d\n",
		format(swap.Total), format(swap.Used), format(swap.Free),
	)

	exitCode := 0

	var memUsed, memFree uint64
	memUsed = format(mem.ActualUsed)
	memFree = format(mem.ActualFree)

	switch l := getPercentUtil(memUsed, memFree); {
	case l > critLevel:
		exitCode = 2
	case l > warnLevel:
		exitCode = 1
	}

	return exitCode
}

// Help returns a string that is the usage for the MemoryCommand
func (m *MemoryCommand) Help() string {
	helpText := `
Usage: consul-check memory <options>

	Check memory utilization of the local system

Options:

  -warn=95   Percent memory utilization for check to be WARN
  -crit=99   Percent memory utilization for check to be CRIT
`

	return strings.TrimSpace(helpText)
}

// Synopsis of the MemoryCommand implementation.
func (m *MemoryCommand) Synopsis() string {
	return "Checks the local memory utilization"
}
