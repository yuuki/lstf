// +build darwin freebsd

package tcpflow

import (
	"fmt"
	"net"

	gnet "github.com/shirou/gopsutil/net"
	"github.com/yuuki/lstf/netutil"
	"golang.org/x/xerrors"
)

// GetHostFlows gets host flows.
// TODO: implement processes option
func GetHostFlows(opt *GetHostFlowsOption) (HostFlows, error) {
	conns, err := gnet.Connections("tcp")
	if err != nil {
		return nil, xerrors.Errorf("gopsutil/net.Connections(): %v", err)
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

		switch opt.Filter {
		case FilterAll:
		case FilterPublic:
			if netutil.IsPrivateIP(net.ParseIP(conn.Raddr.IP)) {
				continue
			}
		case FilterPrivate:
			if !netutil.IsPrivateIP(net.ParseIP(conn.Raddr.IP)) {
				continue
			}
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
	if !opt.Numeric {
		for _, flow := range flows {
			flow.setLookupedName()
		}
	}
	return flows, nil
}
