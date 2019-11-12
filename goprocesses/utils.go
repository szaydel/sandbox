package main

import (
	"fmt"
	"io"
	"os"
)

func handleErr(e error, doPanic bool) {
	if e != nil && e != io.EOF {
		if doPanic {
			panic(e)
		}
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
