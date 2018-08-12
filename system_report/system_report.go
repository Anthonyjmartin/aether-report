package main

import (
	"fmt"
	"gitlab.com/anthony.j.martin/golang/system_programs/hardware_check"
	"github.com/shirou/gopsutil/host"
	"syscall"
	"encoding/json"
	"log"
)

const (
	B  = 1
	KB = B << 10
	MB = KB << 10
	GB = MB << 10
	TB = GB << 10
)

func main() {
	mounts := hardware_check.ListMounts()
	for i := range mounts {
		disk := hardware_check.DiskUsage(mounts[i])
		fmt.Printf("Listing %s:\n", mounts[i])
		fmt.Printf("  ├─All: %.2f GB\n", float64(disk.All)/float64(GB))
		fmt.Printf("  ├─Used: %.2f GB\n", float64(disk.Used)/float64(GB))
		fmt.Printf("  ├─Free: %.2f GB\n", float64(disk.Free)/float64(GB))
		fmt.Printf("  └─Percent: %.2f%%\n\n", disk.Percent)
	}

	s, _ := host.Info()
	fmt.Printf("Hostname:\t%v\n"+
		"Uptime:\t%v\n"+
		"BootTime:\t%v\n"+
		"Procs:\t%v\n"+
		"OS:\t%v\n"+
		"Platform:\t%v\n"+
		"PlatformFamily:\t%v\n"+
		"PlatformVersion:\t%v\n"+
		"KervelVersion:\t%v\n"+
		"VirtualizationSystem:\t%v\n"+
		"VirtualizationRole:\t%v\n"+
		"UUID:\t%v\n\n",
		s.Hostname, s.Uptime, s.BootTime, s.Procs, s.OS, s.Platform, s.PlatformFamily, s.PlatformVersion,
		s.KernelVersion, s.VirtualizationSystem, s.VirtualizationRole, s.HostID)

	//docerkinfo, err := software_check.DockerDetails()
	//if err == nil {
	//	for i := 0; i < len(docerkinfo); i++ {
	//		fmt.Println(docerkinfo[i])
	//	}
	//} else {
	//	fmt.Println(err)
	//}

	fs := syscall.Statfs_t{}
	err := syscall.Statfs("/", &fs)
	if err != nil {
		fmt.Println(err)
	}
	data, err := json.Marshal(fs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)
	fmt.Printf("Bsize:\t%v\n"+
		"Blocks:\t%v\n"+
		"Bfree:\t%v\n"+
		"Bavail:\t%v\n"+
		"Inode Total:\t%v\n"+
		"Inodes Available:\t%v\n"+
		"Fsid:\t%v\n"+
		"Type:\t%v\n"+
		"Flags:\t%v\n",
		fs.Bsize, fs.Blocks, fs.Bfree, fs.Bavail, fs.Files, fs.Ffree, fs.Fsid, fs.Type, fs.Flags)
}
