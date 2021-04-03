package peerinfo

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Process identifies a specific running process on the machine
type Process struct {
	Pid int
	fd  int
}

// Exe will return the path of the running executable
func (c *Process) Exe() (string, error) {
	// readlink of /proc/*/exe
	return os.Readlink(fmt.Sprintf("/proc/%d/exe", c.Pid))
}

// Cmdline will return an array of arguments used to run this executable
func (c *Process) Cmdline() ([]string, error) {
	// read /proc/*/cmdline, split on \0
	data, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", c.Pid))
	if err != nil {
		return nil, err
	}

	if data[len(data)-1] == 0 {
		data = data[:len(data)-1]
	}

	return strings.Split(string(data), "\x00"), nil
}
