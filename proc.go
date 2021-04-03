package peerinfo

import (
	"fmt"
	"os"
	"strconv"
)

func findProcFd(title string) (*Process, error) {
	// list all files in /proc
	plist, err := os.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("failed to list /proc: %w", err)
	}

	// for each proc
	for _, proc := range plist {
		name := proc.Name()
		// check if name is only digits
		pid, err := strconv.ParseUint(name, 10, 64)
		if err != nil {
			// not numeric: not what we're looking for
			continue
		}

		// let's list fds
		fdlist, err := os.ReadDir("/proc/" + name + "/fd")
		if err != nil {
			// mmmh...
			continue
		}

		for _, fd := range fdlist {
			fdName := fd.Name()
			fdInt, err := strconv.ParseUint(fdName, 10, 64)
			if err != nil {
				// mmmmh...
				continue
			}

			lnk, err := os.Readlink("/proc/" + name + "/fd/" + fdName)
			if err != nil {
				// mmmmmmmmmh...
				continue
			}

			if lnk == title {
				// this is the thing we were looking for
				return &Process{Pid: int(pid), fd: int(fdInt)}, nil
			}
		}
	}

	// not found
	return nil, fmt.Errorf("could not find fd titled %s", title)
}
