package hardwarechecks

import (
	"strconv"
	"strings"
)

func PhysicalMemoryStatDetails(memInfo []string) (*PhysicalMemoryStat, error) {
	memParsed := &PhysicalMemoryStat
	for i range memInfo {
		line := memInfo[i]
		fields := strings.Split(line, ":")
		if len(fields) !=2 {
			continue
		}

		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])
		value = strings.Replace(value, "kb", "", -1)

		vnum, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return memParsed, err
		}

		switch key {
		case "MemTotal":
			memParsed.MemTotal = vnum * 1024
		case "MemFree":
			memParsed.MemFree = vnum * 1024
		case "MemAvailable":
			memParsed.MemAvailable = vnum * 1024
		}
	}
}