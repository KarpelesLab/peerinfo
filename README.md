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
