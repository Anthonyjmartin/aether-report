package hardwarechecks

// PhysicalMemoryStat represents the physical memory of a system.
type PhysicalMemoryStat struct {
	Total          uint64  `json:"total"`
	Free           uint64  `json:"free"`
	Available      uint64  `json:"available"`
	Used           uint64  `json:"used"`
	UsedPercent    float64 `json:"usedprecent"`
	Buffers        uint64  `json:"buffers"`
	Cached         uint64  `json:"cached"`
	SwapCached     uint64  `json:"swapcahced"`
	Active         uint64  `json:"active"`
	Inactive       uint64  `json:"inactive"`
	SwapTotal      uint64  `json:"swaptotal"`
	SwapFree       uint64  `json:"swapfree"`
	Dirty          uint64  `json:"dirty"`
	Writeback      uint64  `json:"writeback"`
	AnonPages      uint64  `json:"anonpages"`
	Mapped         uint64  `json:"mapped"`
	Shmem          uint64  `json:"shmem"`
	Slab           uint64  `json:"slab"`
	SReclaimable   uint64  `json:"sreclaimable"`
	SUnreclaimable uint64  `json:"sunreclaimable"`
	KernelStack    uint64  `json:"kernelstack"`
	PageTables     uint64  `json:"pagetables"`
	WritebackTmp   uint64  `json:"writebacktmp"`
	CommitLimit    uint64  `json:"commitlimit"`
	CommittedAS    uint64  `json:"committedas"`
	VmallocTotal   uint64  `json:"vmalloctotal"`
	VmallocUsed    uint64  `json:"vmallocused"`
	VmallocChunk   uint64  `json:"vmallocchunk"`
	AnonHugePages  uint64  `json:"anonhugepages"`
	HugePagesTotal uint64  `json:"hugepagestotal"`
	HugePagesFree  uint64  `json:"hugepagesfree"`
	HugePagesRsvd  uint64  `json:"hugepagesrsvd"`
	HugePagesSurp  uint64  `json:"hugepagessurp"`
	Hugepagesize   uint64  `json:"hugepagesize"`
}
