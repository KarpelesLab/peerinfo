[![GoDoc](https://godoc.org/github.com/KarpelesLab/peerinfo?status.svg)](https://godoc.org/github.com/KarpelesLab/peerinfo)

# peerinfo

Grab information about the process of a connected TCP socket.

This allows doing stuff we typically do with unix sockets with programs that
do not allow unix sockets.

It's all qemu's fault by the way, see https://bugs.launchpad.net/bugs/1903470
