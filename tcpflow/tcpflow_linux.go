// +build linux

package tcpflow

import (
	"fmt"

	"github.com/elastic/gosigar/sys/linux"

	"github.com/yuuki/lstf/netutil"
)

// GetHostFlows gets host flows by Linux netlink API.
func GetHostFlows() (HostFlows, error) {
	conns, err := netutil.NetlinkConnections()
	if err != nil {
		return nil, err
	}
	ports, err := netutil.FilterByLocalListeningPorts(conns)
	if err != nil {
		return nil, err
	}
	flows := HostFlows{}
	for _, conn := range conns {
		if linux.TCPState(conn.State) == linux.TCP_LISTEN {
			continue
		}
		lport, rport := fmt.Sprintf("%d", conn.SrcPort()), fmt.Sprintf("%d", conn.DstPort())
		if contains(ports, lport) {
			flows.insert(&HostFlow{
				Direction: FlowPassive,
				Local:     &AddrPort{Addr: conn.SrcIP().String(), Port: lport},
				Peer:      &AddrPort{Addr: conn.DstIP().String(), Port: "many"},
			})
		} else {
			flows.insert(&HostFlow{
				Direction: FlowActive,
				Local:     &AddrPort{Addr: conn.SrcIP().String(), Port: "many"},
				Peer:      &AddrPort{Addr: conn.DstIP().String(), Port: rport},
			})
		}
	}
	return flows, nil
}
