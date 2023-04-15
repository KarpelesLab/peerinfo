package peerinfo

import (
	"net"
	"syscall"
)

func UnixPeer(c *net.UnixConn) (*Process, error) {
	raw, err := c.SyscallConn()
	if err != nil {
		return nil, err
	}

	var cred *syscall.Ucred
	err2 := raw.Control(func(fd uintptr) {
		cred, err = syscall.GetsockoptUcred(int(fd), syscall.SOL_SOCKET, syscall.SO_PEERCRED)
	})
	if err2 != nil {
		return nil, err2
	}
	if err != nil {
		return nil, err
	}

	res := &Process{
		fd:  -1,
		Pid: int(cred.Pid),
		Uid: int(cred.Uid),
		Gid: int(cred.Gid),
	}

	return res, nil
}
