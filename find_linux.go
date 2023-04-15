package peerinfo

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func Find(loc, rem net.Addr) (*Process, error) {
	var locIp, remIp net.IP
	var locPort, remPort int

	switch a := loc.(type) {
	case *net.TCPAddr:
		locIp = a.IP
		locPort = a.Port
	case *net.UDPAddr:
		locIp = a.IP
		locPort = a.Port
	default:
		return nil, fmt.Errorf("unsupported address type %T", loc)
	}

	switch a := rem.(type) {
	case *net.TCPAddr:
		remIp = a.IP
		remPort = a.Port
	case *net.UDPAddr:
		remIp = a.IP
		remPort = a.Port
	default:
		return nil, fmt.Errorf("unsupported address type %T", rem)
	}

	for _, fn := range []string{"/proc/net/tcp", "/proc/net/tcp6", "/proc/net/udp", "/proc/net/udp6"} {
		sockId, err := findSocket(fn, locIp, locPort, remIp, remPort)
		if err != nil {
			return nil, err
		}

		if sockId != -1 {
			// found
			return findProcFd(fmt.Sprintf("socket:[%d]", sockId))
		}
	}

	// not found
	return nil, fmt.Errorf("error finding socket: %w", os.ErrNotExist)
}

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
				res := &Process{Pid: int(pid), Uid: -1, Gid: -1, fd: int(fdInt)}
				return res, nil
			}
		}
	}

	// not found
	return nil, fmt.Errorf("could not find fd titled %s", title)
}
