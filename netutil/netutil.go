package netutil

import (
	"fmt"
	"net"
	"os"
	"strings"

	gnet "github.com/shirou/gopsutil/net"
)

// LocalListeningPorts returns the local listening ports.
// eg. [199, 111, 46131, 53, 8953, 25, 2812, 80, 8081, 22]
// -----------------------------------------------------------------------------------
// [y_uuki@host ~]$ netstat -tln
// Active Internet connections (only servers)
// Proto Recv-Q Send-Q Local Address               Foreign Address             State
// tcp        0      0 0.0.0.0:199                 0.0.0.0:*                   LISTEN
// tcp        0      0 0.0.0.0:111                 0.0.0.0:*                   LISTEN
// tcp        0      0 0.0.0.0:46131               0.0.0.0:*                   LISTEN
// tcp        0      0 127.0.0.1:53                0.0.0.0:*                   LISTEN
// tcp        0      0 127.0.0.1:8953              0.0.0.0:*                   LISTEN
// tcp        0      0 127.0.0.1:25                0.0.0.0:*                   LISTEN
// tcp        0      0 0.0.0.0:2812                0.0.0.0:*                   LISTEN
// tcp        0      0 :::80                       :::*                        LISTEN
// tcp        0      0 :::8081                     :::*                        LISTEN
// tcp        0      0 :::22                       :::*                        LISTEN
func LocalListeningPorts() ([]string, error) {
	conns, err := gnet.Connections("tcp")
	if err != nil {
		return nil, err
	}
	ports := []string{}
	for _, conn := range conns {
		if conn.Laddr.IP == "0.0.0.0" || conn.Laddr.IP == "127.0.0.1" || conn.Laddr.IP == "::" {
			ports = append(ports, fmt.Sprintf("%d", conn.Laddr.Port))
		}
	}
	return ports, nil
}

// ResolveAddr lookup first hostname from IP Address.
func ResolveAddr(addr string) string {
	hostnames, _ := net.LookupAddr(addr)
	if len(hostnames) > 0 {
		return strings.TrimSuffix(hostnames[0], ".")
	}
	return addr
}

// LocalIPAddrs gets the string slice of localhost IPaddrs.
func LocalIPAddrs() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	addrStrings := make([]string, 0, len(addrs))
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				addrStrings = append(addrStrings, ipnet.IP.String())
			}
		}
	}
	return addrStrings, nil
}

var (
	// IPConntrackPath are ip_conntrack path.
	IPConntrackPath = "/proc/net/ip_conntrack" // old kernel
	// NFConntrackPath are nf_conntrack path.
	NFConntrackPath = "/proc/net/nf_conntrack" // new kernel
)

// FindConntrackPath returns the conntrack proc path if exists.
func FindConntrackPath() string {
	for _, path := range []string{IPConntrackPath, NFConntrackPath} {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}
