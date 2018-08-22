package hardware_check

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gitlab.com/anthony.j.martin/aether-report/internal/pkg/util_funcs"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"text/tabwriter"
)

type DiskBlocks struct {
	Blocks   uint64 `json:"blocks"`
	Bsize    int64  `json:"block_size"`
	Bfree    uint64 `json:"blocks_free"`
	Bavail   uint64 `json:"blocks_avail"`
	Bused    uint64 `json:"blocks_used"`
	Bpercent string `json:"blocks_percent"`
	Balert   string `json:"blocks_alert"`
}

type DiskInodes struct {
	Inodes   uint64 `json:"inodes"`
	Ifree    uint64 `json:"inodes_free"`
	Iused    uint64 `json:"inodes_used"`
	Ipercent string `json:"inodes_percent"`
	Ialert   string `json:"inodes_alert"`
}
type DiskDetails struct {
	Name          string     `json:"name"`
	Partition     string     `json:"partition"`
	PartitionType string     `json:"partition_type"`
	ReadOnly      bool       `json:"read_only"`
	Blocks        DiskBlocks `json:"disk_blocks"`
	Inodes        DiskInodes `json:"disk_inodes"`
}

const (
	B  = 1
	KB = B << 10
	MB = KB << 10
	GB = MB << 10
	TB = GB << 10
)

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
	"squashfs",
	"debugfs",
	"efivarfs",
	"cgroup",
	"mqueue",
	"hugetlbfs",
	"fuse",
	"config",
	"configfs",
	"pstore",
	"securityfs",
	"nsfs",
	"selinuxfs",
	"tracefs",
	"overlay",
}

// Get list of mounted filesystems
func getMounts() []DiskDetails {
	var diskDetails []DiskDetails
	path := "/etc/mtab"
	inFile, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer inFile.Close()
	skipMountRegex, err := regexp.Compile("^/(proc|snap)/") // We do not want reports on virtual filesystems.
	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		data := strings.Fields(scanner.Text())
		if !skipMountRegex.MatchString(data[1]) && !util_funcs.StringInSlice(data[2], excludedFsTypes) {
			ro := true
			opts := strings.Split(data[3], ",")
			for _, opt := range opts {
				if opt == "rw" {
					ro = false
					break
				}
			}
			fs := syscall.Statfs_t{}
			err := syscall.Statfs(data[1], &fs)
			if err != nil {
				//fmt.Println(err)  // Uncomment when debugging.
				continue
			}

			bPercent, _ := getPercent(fs.Blocks, fs.Bavail)
			iPercent, _ := getPercent(fs.Files, fs.Ffree)
			blockAlert, inodeAlert := checkAlert(fs)
			mount := DiskDetails{
				Name:          data[1],
				Partition:     data[0],
				PartitionType: data[2],
				ReadOnly:      ro,
				Blocks: DiskBlocks{
					Blocks:   fs.Blocks,
					Bsize:    fs.Bsize,
					Bfree:    fs.Bfree,
					Bavail:   fs.Bavail,
					Bused:    fs.Blocks - fs.Bavail,
					Bpercent: bPercent,
					Balert:   blockAlert,
				},
				Inodes: DiskInodes{
					Inodes:   fs.Files,
					Ifree:    fs.Ffree,
					Iused:    fs.Files - fs.Ffree,
					Ipercent: iPercent,
					Ialert:   inodeAlert,
				},
			}
			diskDetails = append(diskDetails, mount)
		}
	}
	return diskDetails
}

// Calculate percent for blocks and inodes.
func getPercent(total uint64, avail uint64) (sPercent string, iPercent int) {
	iPercent = int(float64(total-avail) / float64(total) * 100)
	if iPercent >= 0 {
		sPercent = strconv.Itoa(iPercent) + "%"
	} else {
		sPercent = "-%"
	}
	return
}

// Determine if storage needs to have a warning, alert, or is ok.
func checkAlert(fs syscall.Statfs_t) (blockAlert string, inodeAlert string) {
	_, bPercent := getPercent(fs.Blocks, fs.Bavail)
	bAvail := fs.Bsize * int64(fs.Bavail)
	_, iPercent := getPercent(fs.Files, fs.Ffree)

	switch {
	case bPercent < 90 || bAvail >= 20*GB && bPercent < 95:
		blockAlert = "ok"
	case bPercent >= 90 && bPercent < 95:
		blockAlert = "warn"
	default:
		blockAlert = "alert"
	}

	if iPercent >= 95 {
		inodeAlert = "alert"
	} else {
		inodeAlert = "ok"
	}
	return
}

// Convert DiskDetails fields to human-readable formats.
func convertSize(Blocks uint64, Bsize int64) (sizeAsString string) {
	sizedBlocks := float64(Blocks * uint64(Bsize))
	switch {
	case sizedBlocks >= TB:
		sizeAsString = fmt.Sprintf("%.2fTB", sizedBlocks/float64(TB))
	case sizedBlocks >= GB:
		sizeAsString = fmt.Sprintf("%.2fGB", sizedBlocks/float64(GB))
	case sizedBlocks >= MB:
		sizeAsString = fmt.Sprintf("%.2fMB", sizedBlocks/float64(MB))
	case sizedBlocks >= KB:
		sizeAsString = fmt.Sprintf("%.2fKB", sizedBlocks/float64(KB))
	default:
		sizeAsString = fmt.Sprintf("%.2fB", sizedBlocks)
	}
	return sizeAsString
}

// Output data for "text" format.
func textOutput(humanRead bool, inode bool) error {
	fmt.Println("#####   Disk Usage Stats   #####")

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 2, 4, 0, ' ', 0)
	defer w.Flush()

	diskDetails := getMounts()

	switch {
	case inode:
		fmt.Fprintln(w, "Filesystem  \tInodes  \tUsed  \tAvail  \tUse%  \tMount")
	case humanRead:
		fmt.Fprintln(w, "Filesystem  \tSize  \tUsed  \tAvail  \tUse%  \tMount")
	default:
		fmt.Fprintln(w, "Filesystem  \tBlocks  \tUsed  \tAvail  \tUse%  \tMount\t BlockSize")
	}

	for i := range diskDetails {
		if inode {
			diskI := diskDetails[i].Inodes
			fmt.Fprintf(w, "%s  \t%d  \t%d  \t%d  \t%s  \t%s  \n", diskDetails[i].Partition, diskI.Inodes, diskI.Iused, diskI.Ifree, diskI.Ipercent, diskDetails[i].Name)
		} else {

			diskB := diskDetails[i].Blocks
			if humanRead {
				fmt.Fprintf(w, "%s  \t%s  \t%s  \t%s  \t%s  \t%s  \n", diskDetails[i].Partition, convertSize(diskB.Blocks, diskB.Bsize),
					convertSize(diskB.Bused, diskB.Bsize), convertSize(diskB.Bavail, diskB.Bsize),
					diskB.Bpercent, diskDetails[i].Name)
			} else {
				fmt.Fprintf(w, "%s  \t%d  \t%d  \t%d  \t%s  \t%s  \t%d  \n", diskDetails[i].Partition, diskB.Blocks, diskB.Bused, diskB.Bavail, diskB.Bpercent, diskDetails[i].Name, diskB.Bsize)
			}
		}
	}
	return nil
}

// Output data for "json" format.
func jsonOutput() ([]byte, error) {
	diskDetails := getMounts()
	return json.Marshal(diskDetails)
}

// Process data based on passed variables.
func RunDiskInfo(outputFmt string, humanRead bool, inode bool) (textReturn error, jsonReturn []byte, err error) {
	if humanRead && inode {
		fmt.Fprintln(os.Stderr, "\nError: Cannot use both -h and -i  flags.\n\nRun 'aether-report COMMAND --help' for more information on a command.")
		return
	}
	if outputFmt == "text" {
		textReturn = textOutput(humanRead, inode)
		return
	} else if outputFmt == "json" {
		jsonReturn, err = jsonOutput()
		fmt.Println(string(jsonReturn))
	}
	return
}
