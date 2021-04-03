package peerinfo

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func parseHexIp(ip string) (net.IP, int, error) {
	// parse 0100007F:A25E (127.0.0.1)
	// parse 00000000000000000000000001000000:0277 (::1:631) as well

	info := strings.Split(ip, ":")
	if len(info) != 2 {
		// bad
		return nil, -1, fmt.Errorf("could not parse ip %s", ip)
	}

	// parse hex values
	ipB, err := hex.DecodeString(info[0])
	if err != nil {
		return nil, -1, fmt.Errorf("failed to parse ip of %s: %w", ip, err)
	}

	port, err := strconv.ParseUint(info[1], 16, 16) // 16bits port number
	if err != nil {
		return nil, -1, fmt.Errorf("failed to parse port of %s: %w", port, err)
	}

	for i := len(ipB) - 4; i >= 0; i -= 4 {
		quickRevBytes(ipB[i : i+4])
	}

	return net.IP(ipB), int(port), nil
}

func quickRevBytes(b []byte) {
	b[0], b[1], b[2], b[3] = b[3], b[2], b[1], b[0]
}
