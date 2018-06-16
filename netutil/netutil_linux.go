// +build linux

package netutil

import (
	"encoding/binary"
	"errors"
	"fmt"
	"syscall"

	"github.com/vishvananda/netlink/nl"
)

// LodalListeningPorts returns the local listening ports.
func LocalListeningPorts() ([]string, error) {
	s, err := nl.Subscribe(syscall.NETLINK_INET_DIAG)
	if err != nil {
		return nil, err
	}
	defer s.Close()
	req := nl.NewNetlinkRequest(TCPDIAG_GETSOCK, syscall.NLM_F_REQUEST|syscall.NLM_F_DUMP)
	req.AddData(&inetDiagReq{
		Family: syscall.AF_INET,
		States: 1 << TCP_LISTEN,
	})
	s.Send(req)
	msgs, err := s.Receive()
	if err != nil {
		return nil, err
	}
	if len(msgs) == 0 {
		return nil, errors.New("no message nor error from netlink")
	}

	var diagMsgs []*inetDiagMsg
done:
	for _, msg := range msgs {
		if msg.Header.Type == syscall.NLMSG_DONE {
			break done
		}
		if msg.Header.Type == syscall.NLMSG_ERROR {
			errval := native.Uint32(msg.Data[:4])
			return nil, fmt.Errorf("netlink error: %d", -errval)
		}
		diagMsg, err := deserializeInetDiag(msg.Data)
		if err != nil {
			return nil, err
		}
		diagMsgs = append(diagMsgs, diagMsg)
	}

	ports := make([]string, 0, len(diagMsgs))
	for _, diagMsg := range diagMsgs {
		ports = append(ports,
			fmt.Sprintf("%d", binary.BigEndian.Uint16(diagMsg.ID.SrcPort[0:2])))
	}
	return ports, nil
}
