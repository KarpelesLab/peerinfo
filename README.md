[![GoDoc](https://godoc.org/github.com/KarpelesLab/peerinfo?status.svg)](https://godoc.org/github.com/KarpelesLab/peerinfo)

# peerinfo

Grab information about the process of a connected TCP socket.

This allows doing stuff we typically do with unix sockets with programs that
do not allow unix sockets.

It's all qemu's fault by the way, see https://bugs.launchpad.net/bugs/1903470

## Usage

Calling `peerinfo.Find` will return details on the process owning the given
socket address on the local machine. This is done by reading files in `/proc`
so this library is linux only.

```go
import "github.com/KarpelesLab/peerinfo"

func xxx(c net.Conn) {
	// remote and local are reversed as we want to remote side of the socket
	ci, err := peerinfo.Find(c.RemoteAddr(), c.LocalAddr())
	
	// check err, etc...
}
```

## With Unix sockets?

Go notoriously lacks a way to get the peer's information on a unix socket. As
it turns out, an issue was opened in 2010 about this (get `SO_PEERCRED`
information from an open socket), but was hijacked into sending FDs and
credentials over unix sockets and the original issue was cast aside for [the
issue to be finally closed in 2017 without a proper resolution](https://github.com/golang/go/issues/1101#issuecomment-339649510).

Because of this, this package will also provide a method to obtain peer
credentials based on a `net.UnixConn`. Maybe one day Go will have a proper API
to perform this.

You can either directly query info with `peerinfo.UnixPeer()` if you know it
will be a unix conn, or let `peerinfo` detect the type of connection.

```go
import "github.com/KarpelesLab/peerinfo"

func xxx(c net.Conn) {
	// just pass the connection and it'll be automatic
	ci, err := peerinfo.FindPeer(c)
	if err != nil {
		...
	}

	// can use ci.Pid, ci.GetUid(), ci.GetGid()
```

