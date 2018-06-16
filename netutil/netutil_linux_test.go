// +build linux

package netutil

import "testing"

func TestNetlinkConnections(t *testing.T) {
	conns, err := NetlinkConnections()
	if err != nil {
		t.Fatalf("should not raise error: %v", err)
	}
	if len(conns) == 0 {
		t.Error("NetlinkConnections() should not be len == 0")
	}
}
