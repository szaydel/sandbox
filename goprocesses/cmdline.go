package main

import (
	"fmt"
	"os"
	"strings"
)

// CommandLine is a representation of a process' arguments and name
// separated into a slice of argument strings and a name string.
type CommandLine struct {
	args        []string
	programName string
}

// ProgramName is the name of the program for the given process, often
// known as argv[0].
func (c CommandLine) ProgramName() string {
	return c.programName
}

// Args is the set of arguments, sans what would be argv[0], i.e. ProgramName.
// It is effectively argv[1...N]
func (c CommandLine) Args() []string {
	return c.args
}

func cmdLineArgs(pid int) *CommandLine {
	var buf = make([]byte, 256)
	var f *os.File
	var err error
	var n int
	var programName string
	var restOfCmdline []string

	var cmdLinePath = fmt.Sprintf("/proc/%d/cmdline", pid)
	if f, err = os.Open(cmdLinePath); err != nil {
		handleErr(err, false)
	}
	defer f.Close()
	if n, err = f.Read(buf); err != nil {
		handleErr(err, false)
	}
	cmdLineSlice := strings.Split(string(nullByteToSpace(buf[:n])), " ")
	programName = cmdLineSlice[0]

	if len(cmdLineSlice) > 1 {
		restOfCmdline = cmdLineSlice[1:]
	} else {
		restOfCmdline = []string{}
	}
	return &CommandLine{
		programName: programName,
		args:        restOfCmdline,
	}
}
