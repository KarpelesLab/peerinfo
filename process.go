package peerinfo

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// Process identifies a specific running process on the machine
type Process struct {
	Pid int
	Uid int
	Gid int
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

func (p *Process) GetUid() int {
	if p.Uid != -1 {
		return p.Uid
	}
	p.readOwner()
	return p.Uid
}

func (p *Process) GetGid() int {
	if p.Gid != -1 {
		return p.Gid
	}
	p.readOwner()
	return p.Gid
}

func (p *Process) readOwner() {
	if st, err := os.Stat(filepath.Join("/proc", strconv.Itoa(p.Pid))); err == nil {
		p.setOwner(st)
	}
}

func (p *Process) setOwner(s fs.FileInfo) {
	switch v := s.Sys().(type) {
	case *syscall.Stat_t:
		p.Uid = int(v.Uid)
		p.Gid = int(v.Gid)
	default:
		p.Uid = -1
		p.Gid = -1
	}
}
