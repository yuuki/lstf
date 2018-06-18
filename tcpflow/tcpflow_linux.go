// +build linux

package tcpflow

import (
	"fmt"

	"github.com/elastic/gosigar/sys/linux"

	"github.com/yuuki/lstf/netutil"
)

// GetHostFlows gets host flows by netlink, and try to get by procfs if it fails.
func GetHostFlows() (HostFlows, error) {
	flows, err := GetHostFlowsByNetlink()
	if err != nil {
		// fallback to procfs
		return GetHostFlowsByProcfs()
	}
	return flows, err
}

// GetHostFlowsByNetlink gets host flows by Linux netlink API.
func GetHostFlowsByNetlink() (HostFlows, error) {
	conns, err := netutil.NetlinkConnections()
	if err != nil {
		return nil, err
	}
	ports, err := netutil.NetlinkFilterByLocalListeningPorts(conns)
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
