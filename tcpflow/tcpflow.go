package tcpflow

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/yuuki/lstf/netutil"
)

// FlowDirection are bitmask that represents both Active or Passive.
type FlowDirection int

const (
	// FlowUnknown are unknown flow.
	FlowUnknown FlowDirection = 1 << iota
	// FlowActive are 'active open'.
	FlowActive
	// FlowPassive are 'passive open'
	FlowPassive
)

// String returns string representation.
func (c FlowDirection) String() string {
	switch c {
	case FlowActive:
		return "active"
	case FlowPassive:
		return "passive"
	case FlowUnknown:
		return "unknown"
	}
	return ""
}

// MarshalJSON returns human readable `mode` format.
func (c FlowDirection) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// AddrPort are <addr>:<port>
type AddrPort struct {
	Addr string `json:"addr"`
	Port string `json:"port"`
}

// String returns the string representation of the AddrPort.
func (a *AddrPort) String() string {
	return net.JoinHostPort(a.Addr, a.Port)
}

// PortInt returnts integer representation.
func (a *AddrPort) PortInt() int {
	if a.Port == "many" {
		return 0
	}
	i, _ := strconv.Atoi(a.Port)
	return i
}

// HostFlow represents a `host flow`.
type HostFlow struct {
	Direction   FlowDirection `json:"direction"`
	Local       *AddrPort     `json:"local"`
	Peer        *AddrPort     `json:"peer"`
	Connections int64         `json:"connections"`
}

// String returns the string representation of HostFlow.
func (f *HostFlow) String() string {
	switch f.Direction {
	case FlowActive:
		return fmt.Sprintf("%s\t --> \t%s \t%d", f.Local, f.Peer, f.Connections)
	case FlowPassive:
		return fmt.Sprintf("%s\t <-- \t%s \t%d", f.Local, f.Peer, f.Connections)
	}
	return ""
}

// UniqKey returns the unique key for connections aggregation
func (f *HostFlow) UniqKey() string {
	return fmt.Sprintf("%d-%s-%s", f.Direction, f.Local, f.Peer)
}

// ReplaceLookupedName replaces f.Addr into lookuped name.
func (f *HostFlow) ReplaceLookupedName() {
	f.Peer.Addr = netutil.ResolveAddr(f.Peer.Addr)
}

// HostFlows represents a group of host flow by unique key.
type HostFlows map[string]*HostFlow

// MarshalJSON converts map into list.
func (hf HostFlows) MarshalJSON() ([]byte, error) {
	list := make([]HostFlow, 0, len(hf))
	for _, f := range hf {
		list = append(list, *f)
	}
	return json.Marshal(list)
}

func (hf HostFlows) insert(flow *HostFlow) {
	key := flow.UniqKey()
	if _, ok := hf[key]; !ok {
		hf[key] = flow
	}
	hf[key].Connections++
	return
}

func contains(strs []string, s string) bool {
	for _, str := range strs {
		if str == s {
			return true
		}
	}
	return false
}
