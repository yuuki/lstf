package netutil

import (
	"net"
	"strings"
)

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
