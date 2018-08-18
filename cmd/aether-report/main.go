package main

import (
	flag "github.com/ogier/pflag"
	"gitlab.com/anthony.j.martin/aether-report/internal/pkg/hardware_check"
)

var (
	diskOutputFmt     string
	diskHumanRead     bool
	diskDisplayInodes bool
	outputToFile      string
)

func main() {
	flag.Parse()
	hardware_check.RunDiskInfo(diskOutputFmt, outputToFile, diskHumanRead, diskDisplayInodes)
}

func init() {
	flag.StringVarP(&diskOutputFmt, "output", "o", "text", "Output format.")
	flag.BoolVarP(&diskHumanRead, "hr", "h", false, "Display disk storage as human-readable.")
	flag.BoolVarP(&diskDisplayInodes, "inode", "i", false, "Display disk Inode information")
	flag.StringVarP(&outputToFile, "outfile", "f", "", "Output to file.")
}
