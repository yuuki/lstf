// +build linux

package netutil

import (
	"fmt"

	"github.com/elastic/gosigar/sys/linux"
)

// NetlinkConnections returns connection stats.
func NetlinkConnections() ([]*linux.InetDiagMsg, error) {
	req := linux.NewInetDiagReq()
	msgs, err := linux.NetlinkInetDiag(req)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

// FilterByLocalListeningPorts filters ConnectionStat slice by the local listening ports.
func FilterByLocalListeningPorts(conns []*linux.InetDiagMsg) ([]string, error) {
	ports := []string{}
	for _, conn := range conns {
		if linux.TCPState(conn.State) != linux.TCP_LISTEN {
			continue
		}
		sip := conn.SrcIP().String()
		if sip == "0.0.0.0" || sip == "127.0.0.1" || sip == "::" {
			ports = append(ports, fmt.Sprintf("%d", conn.SrcPort()))
		}
	}
	return ports, nil
}

// LodalListeningPorts returns the local listening ports.
func LocalListeningPorts() ([]string, error) {
	msgs, err := NetlinkConnections()
	if err != nil {
		return nil, err
	}
	ports := make([]string, 0, len(msgs))
	for _, diag := range msgs {
		if linux.TCPState(diag.State) != linux.TCP_LISTEN {
			continue
		}
		ports = append(ports, fmt.Sprintf("%v", diag.SrcPort()))
	}
	return ports, nil
}
