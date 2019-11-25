package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func handleErr(e error, doPanic bool) {
	if e != nil && e != io.EOF {
		if doPanic {
			panic(e)
		}
		log.SetOutput(os.Stderr)
		fmt.Fprintf(os.Stderr, "Error: %v\n", e)
	}
}

func nullByteToSpace(b []byte) []byte {
	for i, v := range b {
		if v == 0x0 {
			b[i] = 0x20 // ASCII space character
		}
	}
	return b
}

// isTargetProcess returns true if given pid is referring to executable
// identified by target, otherwise it returns false.
func isTargetProcess(pid int, target string) bool {
	var fileInfo os.FileInfo
	var exe, lnk string
	var err error

	lnk = fmt.Sprintf("/proc/%d/exe", pid)
	if fileInfo, err = os.Lstat(lnk); err != nil {
		return false
	}
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		exe, err = os.Readlink(lnk)
		if err != nil {
			return false
		}
	}
	return strings.Compare(target, exe) == 0
}
