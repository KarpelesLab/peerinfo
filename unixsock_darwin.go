package peerinfo

import (
	"net"

	"golang.org/x/sys/unix"
)

func UnixPeer(c *net.UnixConn) (*Process, error) {
	raw, err := c.SyscallConn()
	if err != nil {
		return nil, err
	}

	res := &Process{fd: -1, Pid: -1, Uid: -1, Gid: -1}

	err2 := raw.Control(func(fd uintptr) {
		res.Pid, err = unix.GetsockoptInt(int(fd), unix.SOL_LOCAL, unix.LOCAL_PEERPID)
		if err != nil {
			return
		}
		cred, err := unix.GetsockoptXucred(int(fd), unix.SOL_LOCAL, unix.LOCAL_PEERCRED)
		if err == nil {
			res.Uid = int(cred.Uid)
			if cred.Ngroups > 0 {
				// TODO get all groups?
				res.Gid = int(cred.Groups[0])
			}
		}
	})
	if err2 != nil {
		return res, err2
	}
	if err != nil {
		return res, err
	}

	return res, nil
}
