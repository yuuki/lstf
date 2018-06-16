// +build linux

package netutil

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/pkg/errors"
	"github.com/vishvananda/netlink/nl"
)

const (
	// TCPF_ALL is a flag to request all sockets in any TCP state.
	TCPF_ALL = ^uint32(0)

	// TCPDIAG_GETSOCK is the netlink message type for requesting TCP diag data.
	// https://github.com/torvalds/linux/blob/v4.0/include/uapi/linux/inet_diag.h#L7
	TCPDIAG_GETSOCK = 18
)

const (
	// https://github.com/torvalds/linux/blob/5924bbecd0267d87c24110cbe2041b5075173a25/include/net/tcp_states.h#L16
	TCP_ESTABLISHED uint8 = iota + 1
	TCP_SYN_SENT
	TCP_SYN_RECV
	TCP_FIN_WAIT1
	TCP_FIN_WAIT2
	TCP_TIME_WAIT
	TCP_CLOSE
	TCP_CLOSE_WAIT
	TCP_LAST_ACK
	TCP_LISTEN
	TCP_CLOSING
)

var (
	native       = nl.NativeEndian()
	networkOrder = binary.BigEndian
)

// inetDiagSockID contains the socket identity.
// https://github.com/torvalds/linux/blob/v4.0/include/uapi/linux/inet_diag.h#L13
type inetDiagSockID struct {
	SrcPort   [2]byte  // Source port (big-endian).
	DstPort   [2]byte  // Destination port (big-endian).
	SrcIP     [16]byte // Source IP
	DstIP     [16]byte // Destination IP
	Interface uint32
	Cookie    [2]uint32
}

var sizeofInetDiagMsg = int(unsafe.Sizeof(inetDiagMsg{}))

// inetDiagMsg contains a socket identifier and netstat information.
type inetDiagMsg struct {
	Family  uint8 // Address family
	State   uint8 // TCP state
	Timer   uint8
	Retrans uint8
	ID      inetDiagSockID
	Expires uint32
	RQueue  uint32
	WQueue  uint32
	UID     uint32
	INode   uint32
}

var sizeofInetDiagReq = int(unsafe.Sizeof(inetDiagReq{}))

// inetDiagReq represents request diagnostic data.
// https://github.com/torvalds/linux/blob/v4.0/include/uapi/linux/inet_diag.h#L37
type inetDiagReq struct {
	Family uint8
	SrcLen uint8
	DstLen uint8
	Ext    uint8
	ID     inetDiagSockID
	States uint32 // States to dump.
	DBs    uint32 // Tables to dump.
}

func (r *inetDiagReq) Serialize() []byte {
	buf := bytes.NewBuffer(make([]byte, sizeofInetDiagReq))
	buf.Reset()
	if err := binary.Write(buf, native, r); err != nil {
		// This never returns an error.
		panic(err)
	}
	return buf.Bytes()
}

func (r *inetDiagReq) Len() int { return sizeofInetDiagReq }

func deserializeInetDiag(b []byte) (*inetDiagMsg, error) {
	if len(b) < sizeofInetDiagMsg {
		return nil, fmt.Errorf("socket data short read (%d); want %d", len(b), sizeofInetDiagMsg)
	}
	r := bytes.NewReader(b)
	inetDiagMsg := &inetDiagMsg{}
	err := binary.Read(r, native, inetDiagMsg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal inet_diag_msg")
	}
	return inetDiagMsg, nil
}
