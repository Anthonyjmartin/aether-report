package main

import (
	"fmt"
	flag "github.com/ogier/pflag"
	"gitlab.com/anthony.j.martin/aether-report/internal/pkg/hardware_check"
	"os"
)

var (
	diskOutputFmt     string
	diskHumanRead     bool
	diskDisplayInodes bool
	allOutputFmt      string
	outputToFile      string
)

func main() {
	diskCommand := flag.NewFlagSet("disk", flag.ExitOnError)
	allCommand := flag.NewFlagSet("all", flag.ExitOnError)
	diskCommand.StringVarP(&diskOutputFmt, "output", "o", "text", "Output format.")
	diskCommand.BoolVarP(&diskHumanRead, "humanread", "h", false, "Display disk storage as human-readable.")
	diskCommand.BoolVarP(&diskDisplayInodes, "inode", "i", false, "Display disk Inode information")
	allCommand.StringVarP(&allOutputFmt, "output", "o", "text", "Output format.")
	flag.StringVarP(&outputToFile, "outfile", "f", "", "Relative or full path to file for output.")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `
Usage: aether-report [OPTIONS] COMMAND [CMDOPTIONS]

Outputs or sends system information. (eg. Disk, CPU, OS info)

Commands:
  all    Run all hardware and software checks.
  disk   Run disk hardware check.

Options:`)
		flag.PrintDefaults()
		fmt.Println()
		fmt.Fprintln(os.Stderr, "Run 'aether-report COMMAND --help' for more information on a command.")
		fmt.Println()
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "disk":
		diskCommand.Parse(os.Args[2:])
		hardware_check.RunDiskInfo(diskOutputFmt, outputToFile, diskHumanRead, diskDisplayInodes)
	case "all":
		allCommand.Parse(os.Args[2:])
		hardware_check.RunDiskInfo(allOutputFmt, outputToFile, diskHumanRead, diskDisplayInodes)
	default:
		flag.Usage()
		os.Exit(1)
	}
}
