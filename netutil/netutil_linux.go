// +build linux

package netutil

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/elastic/gosigar/sys/linux"

	gnet "github.com/shirou/gopsutil/net"
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

// NetlinkFilterByLocalListeningPorts filters ConnectionStat slice by the local listening ports.
func NetlinkFilterByLocalListeningPorts(conns []*linux.InetDiagMsg) ([]string, error) {
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

// NetlinkLodalListeningPorts returns the local listening ports.
func NetlinkLocalListeningPorts() ([]string, error) {
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

const (
	tcpProcFilename = "/proc/net/tcp"
)

// Addr is <addr>:<port>.
type Addr struct {
	IP   string `json:"ip"`
	Port uint32 `json:"port"`
}

// ConnectionStat represents staticstics for a connection.
type ConnectionStat struct {
	Laddr  Addr
	Raddr  Addr
	Status linux.TCPState
}

// ProcfsConnections returns connection stats.
// ref. https://github.com/shirou/gopsutil/blob/c23bcca55e77b8389d84b09db8c5ac2b472070ef/net/net_linux.go#L656
func ProcfsConnections() ([]*ConnectionStat, error) {
	body, err := ioutil.ReadFile(tcpProcFilename)
	if err != nil {
		return nil, err
	}

	lines := bytes.Split(body, []byte("\n"))
	conns := make([]*ConnectionStat, 0, len(lines)-1)
	for _, line := range lines[1:] {
		l := strings.Fields(string(line))
		if len(l) < 10 {
			continue
		}
		laddr := l[1]
		raddr := l[2]
		status, err := strconv.ParseUint(l[3], 16, 8)
		if err != nil {
			log.Printf("decode error: %v", err)
		}
		la, err := decodeAddress(laddr)
		if err != nil {
			continue
		}
		ra, err := decodeAddress(raddr)
		if err != nil {
			continue
		}

		conns = append(conns, &ConnectionStat{
			Laddr:  la,
			Raddr:  ra,
			Status: linux.TCPState(status),
		})
	}

	return conns, nil
}

// decodeAddress decode addresse represents addr in proc/net/*
// ex:
// "0500000A:0016" -> "10.0.0.5", 22
// "0085002452100113070057A13F025401:0035" -> "2400:8500:1301:1052:a157:7:154:23f", 53
// ref. https://github.com/shirou/gopsutil/blob/c23bcca55e77b8389d84b09db8c5ac2b472070ef/net/net_linux.go#L600
func decodeAddress(src string) (Addr, error) {
	t := strings.Split(src, ":")
	if len(t) != 2 {
		return Addr{}, fmt.Errorf("does not contain port, %s", src)
	}
	addr := t[0]
	port, err := strconv.ParseInt("0x"+t[1], 0, 64)
	if err != nil {
		return Addr{}, fmt.Errorf("invalid port, %s", src)
	}
	decoded, err := hex.DecodeString(addr)
	if err != nil {
		return Addr{}, fmt.Errorf("decode error, %s", err)
	}
	var ip net.IP
	// Assumes this is little_endian
	ip = net.IP(gnet.Reverse(decoded))
	return Addr{
		IP:   ip.String(),
		Port: uint32(port),
	}, nil
}

// FilterByLocalListeningPorts filters ConnectionStat slice by the local listening ports.
func FilterByLocalListeningPorts(conns []*ConnectionStat) ([]string, error) {
	ports := []string{}
	for _, conn := range conns {
		if conn.Status != linux.TCP_LISTEN {
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
	conns, err := ProcfsConnections()
	if err != nil {
		return nil, err
	}
	return FilterByLocalListeningPorts(conns)
}
