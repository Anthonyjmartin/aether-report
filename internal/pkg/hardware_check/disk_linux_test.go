package hardware_check

import (
	"encoding/json"
	"reflect"
	"strings"
	"syscall"
	"testing"
)

const mockDiskFile string = `rootfs / rootfs rw 0 0
sysfs /sys sysfs rw,seclabel,nosuid,nodev,noexec,relatime 0 0
/dev/mapper/centos-root / xfs rw,seclabel,relatime,attr2,inode64,noquota 0 0`

var fakeStatfs = syscall.Statfs_t{
	Bsize:  4096,
	Blocks: 114434612,
	Bfree:  78317375,
	Bavail: 72486975,
	Files:  29138944,
	Ffree:  28437091,
}
var fakeDiskSlice []DiskDetails

func Test_getPercent(t *testing.T) {
	type args struct {
		total uint64
		avail uint64
	}
	tests := []struct {
		name         string
		args         args
		wantSPercent string
		wantIPercent int
	}{
		{
			name: "20% check",
			args: struct {
				total uint64
				avail uint64
			}{total: 500, avail: 400},
			wantSPercent: "20%",
			wantIPercent: 20,
		},
		{
			name: "100% check",
			args: struct {
				total uint64
				avail uint64
			}{total: 500, avail: 0},
			wantSPercent: "100%",
			wantIPercent: 100,
		},
		{
			name: "0% check",
			args: struct {
				total uint64
				avail uint64
			}{total: 500, avail: 500},
			wantSPercent: "0%",
			wantIPercent: 0,
		},
		{
			name: "-% inode check",
			args: struct {
				total uint64
				avail uint64
			}{total: 0, avail: 0},
			wantSPercent: "-%",
			wantIPercent: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSPercent, gotIPercent := getPercent(tt.args.total, tt.args.avail)
			if gotSPercent != tt.wantSPercent {
				t.Errorf("getPercent() gotSPercent = %v, want %v", gotSPercent, tt.wantSPercent)
			}
			if gotIPercent != tt.wantIPercent {
				t.Errorf("getPercent() gotIPercent = %v, want %v", gotIPercent, tt.wantIPercent)
			}
		})
	}
}

func Test_checkAlert(t *testing.T) {
	type args struct {
		blocks      uint64
		blocksAvail uint64
		inodes      uint64
		inodesFree  uint64
		blockSize   int64
	}
	tests := []struct {
		name           string
		args           args
		wantBlockAlert string
		wantInodeAlert string
	}{
		{
			name: "Block: ok, Inode: ok",
			args: struct {
				blocks      uint64
				blocksAvail uint64
				inodes      uint64
				inodesFree  uint64
				blockSize   int64
			}{blocks: 1073741824, blocksAvail: 1073741824, inodes: 20, inodesFree: 20, blockSize: 4096},
			wantBlockAlert: "ok",
			wantInodeAlert: "ok",
		},
		{
			name: "Block: warn (for percent), Inode: ok",
			args: struct {
				blocks      uint64
				blocksAvail uint64
				inodes      uint64
				inodesFree  uint64
				blockSize   int64
			}{blocks: 100, blocksAvail: 9, inodes: 20, inodesFree: 20, blockSize: 4096},
			wantBlockAlert: "warn",
			wantInodeAlert: "ok",
		},
		{
			name: "Block: alert (for percent), Inode: ok",
			args: struct {
				blocks      uint64
				blocksAvail uint64
				inodes      uint64
				inodesFree  uint64
				blockSize   int64
			}{blocks: 100, blocksAvail: 4, inodes: 20, inodesFree: 20, blockSize: 4096},
			wantBlockAlert: "alert",
			wantInodeAlert: "ok",
		},
		{
			name: "Block: ok (for more than 20GB), Inode: ok",
			args: struct {
				blocks      uint64
				blocksAvail uint64
				inodes      uint64
				inodesFree  uint64
				blockSize   int64
			}{blocks: 4978176000, blocksAvail: 448035840, inodes: 20, inodesFree: 20, blockSize: 4096},
			wantBlockAlert: "ok",
			wantInodeAlert: "ok",
		},
		{
			name: "Block: ok (for more than 20GB), Inode: alert",
			args: struct {
				blocks      uint64
				blocksAvail uint64
				inodes      uint64
				inodesFree  uint64
				blockSize   int64
			}{blocks: 4978176000, blocksAvail: 448035840, inodes: 100, inodesFree: 3, blockSize: 4096},
			wantBlockAlert: "ok",
			wantInodeAlert: "alert",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBlockAlert, gotInodeAlert := checkAlert(tt.args.blocks, tt.args.blocksAvail, tt.args.inodes, tt.args.inodesFree, tt.args.blockSize)
			if gotBlockAlert != tt.wantBlockAlert {
				t.Errorf("checkAlert() gotBlockAlert = %v, want %v", gotBlockAlert, tt.wantBlockAlert)
			}
			if gotInodeAlert != tt.wantInodeAlert {
				t.Errorf("checkAlert() gotInodeAlert = %v, want %v", gotInodeAlert, tt.wantInodeAlert)
			}
		})
	}
}

func Test_convertSize(t *testing.T) {
	type args struct {
		Blocks uint64
		Bsize  int64
	}
	tests := []struct {
		name             string
		args             args
		wantSizeAsString string
	}{
		{
			name: "Return Bytes",
			args: args{
				Blocks: 0,
				Bsize:  4096,
			},
			wantSizeAsString: "0.00B",
		},
		{
			name: "Return KiloBytes",
			args: args{
				Blocks: 1,
				Bsize:  4096,
			},
			wantSizeAsString: "4.00KB",
		},
		{
			name: "Return MegaBytes",
			args: args{
				Blocks: 1024,
				Bsize:  4096,
			},
			wantSizeAsString: "4.00MB",
		},
		{
			name: "Return GigaBytes",
			args: args{
				Blocks: 1048576,
				Bsize:  4096,
			},
			wantSizeAsString: "4.00GB",
		},
		{
			name: "Return TeraBytes",
			args: args{
				Blocks: 1073741824,
				Bsize:  4096,
			},
			wantSizeAsString: "4.00TB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSizeAsString := convertSize(tt.args.Blocks, tt.args.Bsize); gotSizeAsString != tt.wantSizeAsString {
				t.Errorf("convertSize() = %v, want %v", gotSizeAsString, tt.wantSizeAsString)
			}
		})
	}
}

func Test_getDiskInfo(t *testing.T) {
	type args struct {
		mount        string
		fakeDiskInfo syscall.Statfs_t
	}
	tests := []struct {
		name    string
		args    args
		wantFs  syscall.Statfs_t
		wantErr bool
	}{
		{
			name: "Return fake data",
			args: args{
				mount:        "ignore",
				fakeDiskInfo: fakeStatfs,
			},
			wantFs: syscall.Statfs_t{
				Bsize:  4096,
				Blocks: 114434612,
				Bfree:  78317375,
				Bavail: 72486975,
				Files:  29138944,
				Ffree:  28437091,
			},
			wantErr: false,
		},
		{
			name: "Return error",
			args: args{
				mount:        "ignore",
				fakeDiskInfo: syscall.Statfs_t{},
			},
			wantFs:  syscall.Statfs_t{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFs, err := getDiskInfo(tt.args.mount, tt.args.fakeDiskInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDiskInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFs, tt.wantFs) {
				t.Errorf("getDiskInfo() = %v, want %v", gotFs, tt.wantFs)
			}
		})
	}
}

func Test_parseDiskFile(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    []ValidDisks
		wantErr bool
	}{
		{
			name: "Fake file with 1 correct disk.",
			args: args{
				file: mockDiskFile,
			},
			want: []ValidDisks{
				{
					Partition: "/dev/mapper/centos-root",
					Mount:     "/",
					Type:      "xfs",
					Options:   "rw,seclabel,relatime,attr2,inode64,noquota",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDiskFile(strings.NewReader(tt.args.file))
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDiskFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDiskFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMounts(t *testing.T) {
	type args struct {
		diskFile     string
		testDiskInfo syscall.Statfs_t
	}
	tests := []struct {
		name    string
		args    args
		want    []DiskDetails
		wantErr bool
	}{
		{
			name: "Return first error",
			args: args{
				diskFile:     "",
				testDiskInfo: syscall.Statfs_t{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Returns correct details.",
			args: args{
				diskFile:     mockDiskFile,
				testDiskInfo: fakeStatfs,
			},
			want: []DiskDetails{
				{
					Name:          "/",
					Partition:     "/dev/mapper/centos-root",
					PartitionType: "xfs",
					ReadOnly:      false,
					Blocks: DiskBlocks{
						Blocks:   114434612,
						Bsize:    4096,
						Bfree:    78317375,
						Bavail:   72486975,
						Bused:    41947637,
						Bpercent: "36%",
						Balert:   "ok",
					},
					Inodes: DiskInodes{
						Inodes:   29138944,
						Ifree:    28437091,
						Iused:    701853,
						Ipercent: "2%",
						Ialert:   "ok",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getMounts(tt.args.diskFile, tt.args.testDiskInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("getMounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getMounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_textOutput(t *testing.T) {
	type args struct {
		humanRead    bool
		inode        bool
		diskFile     string
		testDiskInfo syscall.Statfs_t
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Return error",
			args: args{
				humanRead:    false,
				inode:        false,
				diskFile:     "",
				testDiskInfo: syscall.Statfs_t{},
			},
			wantErr: true,
		},
		{
			name: "Return inode",
			args: args{
				humanRead:    false,
				inode:        true,
				diskFile:     mockDiskFile,
				testDiskInfo: fakeStatfs,
			},
			wantErr: false,
		},
		{
			name: "Return blocks",
			args: args{
				humanRead:    false,
				inode:        false,
				diskFile:     mockDiskFile,
				testDiskInfo: fakeStatfs,
			},
			wantErr: false,
		},
		{
			name: "Return human readable",
			args: args{
				humanRead:    true,
				inode:        false,
				diskFile:     mockDiskFile,
				testDiskInfo: fakeStatfs,
			},
			wantErr: false,
		},
		{
			name: "Return -i and -h join error",
			args: args{
				humanRead:    true,
				inode:        true,
				diskFile:     mockDiskFile,
				testDiskInfo: fakeStatfs,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := textOutput(tt.args.humanRead, tt.args.inode, tt.args.diskFile, tt.args.testDiskInfo); (err != nil) != tt.wantErr {
				t.Errorf("textOutput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunDiskInfo(t *testing.T) {
	type args struct {
		outputFmt    string
		humanRead    bool
		inode        bool
		diskFile     string
		fakeDiskInfo syscall.Statfs_t
	}
	fakeDiskSlice = append(fakeDiskSlice, DiskDetails{
		Name:          "/",
		Partition:     "/dev/mapper/centos-root",
		PartitionType: "xfs",
		ReadOnly:      false,
		Blocks: DiskBlocks{
			Blocks:   114434612,
			Bsize:    4096,
			Bfree:    78317375,
			Bavail:   72486975,
			Bused:    41947637,
			Bpercent: "36%",
			Balert:   "ok",
		},
		Inodes: DiskInodes{
			Inodes:   29138944,
			Ifree:    28437091,
			Iused:    701853,
			Ipercent: "2%",
			Ialert:   "ok",
		},
	},
	)
	jsonCheckData, _ := getMounts(mockDiskFile, fakeStatfs)
	convertedJson, _ := json.Marshal(jsonCheckData)

	tests := []struct {
		name           string
		args           args
		wantJsonReturn []byte
		wantErr        bool
	}{
		{
			name: "Text output",
			args: args{
				outputFmt:    "text",
				humanRead:    false,
				inode:        false,
				diskFile:     mockDiskFile,
				fakeDiskInfo: fakeStatfs,
			},
			wantJsonReturn: nil,
			wantErr:        false,
		},
		{
			name: "json output",
			args: args{
				outputFmt:    "json",
				humanRead:    false,
				inode:        false,
				diskFile:     mockDiskFile,
				fakeDiskInfo: fakeStatfs,
			},
			wantJsonReturn: convertedJson,
			wantErr:        false,
		},
		{
			name: "Get error from both -h and -i",
			args: args{
				outputFmt:    "text",
				humanRead:    true,
				inode:        true,
				diskFile:     mockDiskFile,
				fakeDiskInfo: fakeStatfs,
			},
			wantJsonReturn: nil,
			wantErr:        true,
		},
		{
			name: "json error",
			args: args{
				outputFmt:    "json",
				humanRead:    false,
				inode:        false,
				diskFile:     "",
				fakeDiskInfo: syscall.Statfs_t{},
			},
			wantJsonReturn: nil,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotJsonReturn, err, _ := RunDiskInfo(tt.args.outputFmt, tt.args.humanRead, tt.args.inode, tt.args.diskFile, tt.args.fakeDiskInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunDiskInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotJsonReturn, tt.wantJsonReturn) {
				t.Errorf("RunDiskInfo() = %v, want %v", gotJsonReturn, tt.wantJsonReturn)
			}
		})
	}
}
