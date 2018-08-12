package hardware_check

import (
	"syscall"
	"os"
	"bufio"
	"strings"
	)

type DiskStatus struct {
	All     uint64 `json:"all"`
	Used    uint64 `json:"used"`
	Free    uint64 `json:"free"`
	Percent float64 `json:"percent"`
}

var excludedFsTypes = []string{
	"autofs",
	"usbfs",
	"rootfs",
	"proc",
	"sysfs",
	"devtmpfs",
	"devpts",
	"tmpfs",
	"binfmt_misc",
	"rpc_pipefs",
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Get list of mounted filesystems
func ListMounts() (mounts []string) {
	path := "/proc/mounts"
	inFile, err := os.Open(path)
	defer inFile.Close()
	check(err)
	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		data := strings.Fields(scanner.Text())
		if strings.HasPrefix(data[0], "/") && !strings.HasPrefix(data[1], "/proc/") && !stringInSlice(data[2], excludedFsTypes) {
			mounts = append(mounts, data[1])
		}
	}
	return
}

// disk usage of path/disk
func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bavail * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	disk.Percent = float64(disk.Used)/float64(disk.All)*100
	return
}
