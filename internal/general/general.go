package general

import (
	"bufio"
	"os"
	"path/filepath"
	"fmt"
	"strings"
)

// Human readable storage sizes
const (
	B  = 1
	KB = B << 10
	MB = KB << 10
	GB = MB << 10
	TB = GB << 10
)

//StringInSlice return true if string is in slice else return false.
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a || strings.HasPrefix(a, b) {
			return true
		}
	}
	return false
}

// HumanReadableFormat takes a float64 size in bytes and reuturns a human readable version as a string.
func HumanReadableFormat(rawSize float64) (humanReadString string) {
	switch {
	case rawSize >= TB:
		humanReadString = fmt.Sprintf("%.2fTB", rawSize/float64(TB))
	case rawSize >= GB:
		humanReadString = fmt.Sprintf("%.2fGB", rawSize/float64(GB))
	case rawSize >= MB:
		humanReadString = fmt.Sprintf("%.2fMB", rawSize/float64(MB))
	case rawSize >= KB:
		humanReadString = fmt.Sprintf("%.2fKB", rawSize/float64(KB))
	default:
		humanReadString = fmt.Sprintf("%.2fB", rawSize)
	}
	return humanReadString
}

// FileReadLines reads the lines of a provided filepath and returns lines as a slice of strings.
func FileReadLines(filepath string) (fileStrings []string, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileStrings = append(fileStrings, scanner.Text())
	}
	return
}