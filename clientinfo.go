package peerinfo

import (
	"net"
)

func FindPeer(n net.Conn) (*Process, error) {
	switch v := n.(type) {
	case *net.TCPConn:
		return Find(v.RemoteAddr(), v.LocalAddr())
	case *net.UDPConn:
		return Find(v.RemoteAddr(), v.LocalAddr())
	case *net.UnixConn:
		return UnixPeer(v)
	default:
		return nil, ErrUnsupportedConnectionType
	}
}
