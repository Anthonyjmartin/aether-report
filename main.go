package main

import (
	"fmt"
	"github.com/urfave/cli"
	"gitlab.com/anthony.j.martin/aether-report/internal/pkg/hardwarechecks"
	"os"
	"syscall"
	"time"
)

var (
	diskOutputFmt     string
	diskHumanRead     bool
	diskDisplayInodes bool
	version           string
)

func init() {
	cli.VersionFlag = cli.BoolFlag{Name: "version, V"}

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "version=%s\n", c.App.Version)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "aether-report"
	app.Version = version
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		{
			Name:  "Anthony Martin",
			Email: "anthony.j.martin142@gmail.com",
		},
	}
	app.Usage = "Collect and report system information."
	app.HideHelp = false
	app.HideVersion = false

	app.Commands = []cli.Command{
		{
			Name:     "disk",
			Category: "Hardware Checks",
			Usage:    "Runs report of disk.",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "humanread, H", Destination: &diskHumanRead, Usage: "Display disk storage as human-readable."},
				cli.BoolFlag{Name: "inode, i", Destination: &diskDisplayInodes, Usage: "Display inode information."},
				cli.StringFlag{Name: "output, o", Value: "text", Destination: &diskOutputFmt, Usage: "Chose output `FORMAT` <text|json>."},
			},
			Action: func(c *cli.Context) error {
				hardwarechecks.RunDiskInfo(diskOutputFmt, diskHumanRead, diskDisplayInodes, "/etc/mtab", syscall.Statfs_t{})
				return nil
			},
		},
		{
			Name:  "all",
			Usage: "Runs report on all checks.",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "humanread, H", Destination: &diskHumanRead, Usage: "Display disk storage as human-readable."},
				cli.StringFlag{Name: "output, o", Value: "text", Destination: &diskOutputFmt, Usage: "Chose output `FORMAT` [(text)|json]."},
			},
			Action: func(c *cli.Context) error {
				hardwarechecks.RunDiskInfo(diskOutputFmt, diskHumanRead, diskDisplayInodes, "/etc/mtab", syscall.Statfs_t{})
				return nil
			},
		},
		{
			Name:  "report",
			Usage: "Runs all checks and reports back to configured aether-station system.",
			Action: func(c *cli.Context) error {
				fmt.Println("This is a place holder as aether-station is still a WIP.")
				return nil
			},
		},
	}

	_ = app.Run(os.Args)
}
