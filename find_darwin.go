package peerinfo

import (
	"errors"
	"net"
)

func Find(loc, rem net.Addr) (*Process, error) {
	return nil, errors.New("not supported on darwin")
}
