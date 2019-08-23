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

	FilterAll     = "all"
	FilterPublic  = "public"
	FilterPrivate = "private"
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
	Name string `json:"name"`
	Addr string `json:"addr"`
	Port string `json:"port"`
}

// String returns the string representation of the AddrPort.
func (a *AddrPort) String() string {
	if a.Name == "" {
		return net.JoinHostPort(a.Addr, a.Port)
	}
	return net.JoinHostPort(a.Name, a.Port)
}

// PortInt returnts integer representation.
func (a *AddrPort) PortInt() int {
	if a.Port == "many" {
		return 0
	}
	i, _ := strconv.Atoi(a.Port)
	return i
}

// Process represents a OS process.
type Process struct {
	Name string `json:"name"`
	Pgid int    `json:"pgid"`
}

// HostFlow represents a `host flow`.
type HostFlow struct {
	Direction   FlowDirection `json:"direction"`
	Local       *AddrPort     `json:"local"`
	Peer        *AddrPort     `json:"peer"`
	Connections int64         `json:"connections"`
	Process     *Process      `json:"process,omitempty"`
}

// String returns the string representation of HostFlow.
func (f *HostFlow) String() string {
	var entStr string
	if f.Process != nil {
		entStr = fmt.Sprintf("\t(\"%s\",pgid=%d)", f.Process.Name, f.Process.Pgid)
	}
	switch f.Direction {
	case FlowActive:
		return fmt.Sprintf("%s\t-->\t%s\t%d%s", f.Local, f.Peer, f.Connections, entStr)
	case FlowPassive:
		return fmt.Sprintf("%s\t<--\t%s\t%d%s", f.Local, f.Peer, f.Connections, entStr)
	}
	return ""
}

// UniqKey returns the unique identifier key for connections flow.
func (f *HostFlow) UniqKey() string {
	return fmt.Sprintf("%d-%s-%s", f.Direction, f.Local, f.Peer)
}

// setLookupedName replaces f.Addr into lookuped name.
func (f *HostFlow) setLookupedName() {
	f.Local.Name = netutil.ResolveAddr(f.Local.Addr)
	f.Peer.Name = netutil.ResolveAddr(f.Peer.Addr)
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
	} else {
		if hf[key].Process == nil {
			hf[key].Process = flow.Process
		}
	}
	hf[key].Connections++
}

func contains(strs []string, s string) bool {
	for _, str := range strs {
		if str == s {
			return true
		}
	}
	return false
}

// GetHostFlowsOption represens an option for func GetHostFlows().
type GetHostFlowsOption struct {
	Numeric   bool
	Processes bool
	Filter    string
}
