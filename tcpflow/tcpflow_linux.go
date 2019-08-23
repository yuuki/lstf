// +build linux

package tcpflow

import (
	"fmt"

	"github.com/elastic/gosigar/sys/linux"
	"golang.org/x/xerrors"

	"github.com/yuuki/lstf/netutil"
)

// GetHostFlows gets host flows by netlink, and try to get by procfs if it fails.
func GetHostFlows(opt *GetHostFlowsOption) (HostFlows, error) {
	flows, err := GetHostFlowsByNetlink(opt)
	if err != nil {
		var netlinkErr *netutil.NetlinkError
		if xerrors.As(err, &netlinkErr) {
			// fallback to procfs
			return GetHostFlowsByProcfs()
		}
		return nil, err
	}
	return flows, nil
}

// GetHostFlowsByNetlink gets host flows by Linux netlink API.
func GetHostFlowsByNetlink(opt *GetHostFlowsOption) (HostFlows, error) {
	var userEnts netutil.UserEnts
	if opt.Processes {
		var err error
		userEnts, err = netutil.BuildUserEntries()
		if err != nil {
			return nil, err
		}
	}
	conns, err := netutil.NetlinkConnections()
	if err != nil {
		return nil, err
	}
	lconns, err := netutil.NetlinkFilterByLocalListeningPorts(conns)
	if err != nil {
		return nil, err
	}

	ports := make([]string, 0, len(lconns))
	lportEnt := make(netutil.UserEntByLport, len(lconns))
	for _, lconn := range lconns {
		sport := fmt.Sprintf("%d", lconn.SrcPort())
		ports = append(ports, sport)
		if userEnts != nil {
			lportEnt[sport] = userEnts[lconn.Inode]
		}
	}

	flows := HostFlows{}
	for _, conn := range conns {
		switch linux.TCPState(conn.State) {
		case linux.TCP_LISTEN:
			continue
		case linux.TCP_SYN_SENT:
			continue
		case linux.TCP_SYN_RECV:
			continue
		}

		switch opt.Filter {
		case FilterAll:
		case FilterPublic:
			if netutil.IsPrivateIP(conn.DstIP()) {
				continue
			}
		case FilterPrivate:
			if !netutil.IsPrivateIP(conn.DstIP()) {
				continue
			}
		}

		var ent *netutil.UserEnt
		// inode 0 means that it provides no process information
		if userEnts != nil && conn.Inode != 0 {
			ent = userEnts[conn.Inode]
		}

		lport, rport := fmt.Sprintf("%d", conn.SrcPort()), fmt.Sprintf("%d", conn.DstPort())
		if contains(ports, lport) {
			// passive open
			if ent == nil {
				ent = lportEnt[lport]
			}
			hf := &HostFlow{
				Direction: FlowPassive,
				Local:     &AddrPort{Addr: conn.SrcIP().String(), Port: lport},
				Peer:      &AddrPort{Addr: conn.DstIP().String(), Port: "many"},
			}
			if ent != nil {
				hf.Process = &Process{
					Name: ent.Pname(),
					Pgid: ent.Pgrp(),
				}
			}
			flows.insert(hf)
		} else {
			// active open
			hf := &HostFlow{
				Direction: FlowActive,
				Local:     &AddrPort{Addr: conn.SrcIP().String(), Port: "many"},
				Peer:      &AddrPort{Addr: conn.DstIP().String(), Port: rport},
			}
			if ent != nil {
				hf.Process = &Process{
					Name: ent.Pname(),
					Pgid: ent.Pgrp(),
				}
			}
			flows.insert(hf)
		}
	}

	if !opt.Numeric {
		for _, flow := range flows {
			flow.setLookupedName()
		}
	}
	return flows, nil
}

// GetHostFlowsByProcfs gets host flows from procfs.
func GetHostFlowsByProcfs() (HostFlows, error) {
	conns, err := netutil.ProcfsConnections()
	if err != nil {
		return nil, err
	}
	ports, err := netutil.FilterByLocalListeningPorts(conns)
	if err != nil {
		return nil, err
	}
	flows := HostFlows{}
	for _, conn := range conns {
		switch conn.Status {
		case linux.TCP_LISTEN:
			continue
		case linux.TCP_SYN_SENT:
			continue
		case linux.TCP_SYN_RECV:
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
