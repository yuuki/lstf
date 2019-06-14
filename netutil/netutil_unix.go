// +build freebsd darwin

package netutil

import (
	"fmt"

	gnet "github.com/shirou/gopsutil/net"
	"golang.org/x/xerrors"
)

// FilterByLocalListeningPorts filters ConnectionStat slice by the local listening ports.
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
func FilterByLocalListeningPorts(conns []gnet.ConnectionStat) ([]string, error) {
	ports := []string{}
	for _, conn := range conns {
		if conn.Status != "LISTEN" {
			continue
		}
		if conn.Laddr.IP == "0.0.0.0" || conn.Laddr.IP == "127.0.0.1" || conn.Laddr.IP == "::" {
			ports = append(ports, fmt.Sprintf("%d", conn.Laddr.Port))
		}
	}
	return ports, nil
}

// LocalListeningPorts returns the local listening ports.
func LocalListeningPorts() ([]string, error) {
	conns, err := gnet.Connections("tcp")
	if err != nil {
		return nil, xerrors.Errorf("gopsutil/net.Connections() failed: %v", err)
	}
	return FilterByLocalListeningPorts(conns)
}
