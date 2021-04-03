package peerinfo

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func findSocket(file string, locIp net.IP, locPort int, remIp net.IP, remPort int) (int, error) {
	// /proc/net/tcp

	f, err := os.Open(file)
	if err != nil {
		return -1, fmt.Errorf("failed to open %s: %w", file, err)
	}
	defer f.Close()

	// first line:   sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
	r := bufio.NewReader(f)
	r.ReadString('\n')

	for {
		lin, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// simply not found
				return -1, nil
			}
			return -1, fmt.Errorf("while reading %s: %w", file, err)
		}
		linA := strings.Fields(lin) // split
		if len(linA) < 10 {
			// shouldn't happen
			return -1, fmt.Errorf("failed to parse %s: invalid line", file)
		}
		// decode line's IP & port
		//  18: 0100007F:A25E 0100007F:107D 01 00000000:00000000 00:00000000 00000000  1000        0 2629411 1 000000008148de13 21 0 0 10 -1
		//  54: 0100007F:107D 0100007F:A25E 01 00000000:00000000 02:000000D1 00000000  1000        0 2637058 2 000000003976dda9 20 4 28 10 -1
		// 1 = local addr
		// 2 = remote addr
		// 9 = socket

		linLocIp, linLocPort, err := parseHexIp(linA[1])
		if err != nil {
			return -1, fmt.Errorf("failed to read %s: %w", file, err)
		}
		linRemIp, linRemPort, err := parseHexIp(linA[2])
		if err != nil {
			return -1, fmt.Errorf("failed to read %s: %w", file, err)
		}

		if !locIp.Equal(linLocIp) || !remIp.Equal(linRemIp) {
			continue
		}
		if locPort != linLocPort || remPort != linRemPort {
			continue
		}

		// we found our line, parse the socket
		sock, err := strconv.ParseInt(linA[9], 10, 64)
		if err != nil {
			return -1, fmt.Errorf("failed to parse socket in %s: %w", file, err)
		}

		return int(sock), nil
	}
}
