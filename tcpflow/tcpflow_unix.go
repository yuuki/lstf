// +build darwin freebsd

package tcpflow

import (
	"fmt"

	gnet "github.com/shirou/gopsutil/net"
	"github.com/yuuki/lstf/netutil"
)

// GetHostFlows gets host flows.
func GetHostFlows() (HostFlows, error) {
	conns, err := gnet.Connections("tcp")
	if err != nil {
		return nil, err
	}
	ports, err := netutil.FilterByLocalListeningPorts(conns)
	if err != nil {
		return nil, err
	}
	flows := HostFlows{}
	for _, conn := range conns {
		if conn.Status == "LISTEN" {
			continue
		}
		lport := fmt.Sprintf("%d", conn.Laddr.Port)
		rport := fmt.Sprintf("%d", conn.Raddr.Port)
		if contains(ports, lport) {
			flows.insert(&HostFlow{
				Direction: FlowPassive,
				Local:     &AddrPort{Addr: conn.Laddr.IP, Port: lport},
				Peer:      &AddrPort{Addr: conn.Raddr.IP, Port: "many"},
			})
		} else {
			flows.insert(&HostFlow{
				Direction: FlowActive,
				Local:     &AddrPort{Addr: conn.Laddr.IP, Port: "many"},
				Peer:      &AddrPort{Addr: conn.Raddr.IP, Port: rport},
			})
		}
	}
	return flows, nil
}
