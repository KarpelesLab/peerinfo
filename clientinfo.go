package peerinfo

import (
	"fmt"
	"net"
	"os"
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
