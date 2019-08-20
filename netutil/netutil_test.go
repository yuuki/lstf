package netutil

import (
	"net"
	"testing"
)

func TestLocalIPAddrss(t *testing.T) {
	addrs, err := LocalIPAddrs()
	if err != nil {
		t.Fatalf("should not raise error: %v", err)
	}
	if len(addrs) == 0 {
		t.Error("localIPAddrs() should not be len == 0")
	}
}

func TestLocalListeningPorts(t *testing.T) {
	ports, err := LocalListeningPorts()
	if err != nil {
		t.Fatalf("should not raise error: %v", err)
	}
	if len(ports) == 0 {
		t.Error("localIPAddrs() should not be len == 0")
	}
}

func TestIsPrivateIP(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		{"192.168.10.111", true},
		{"172.16.10.111", true},
		{"10.1.10.111", true},
		{"192.0.2.111", false},
	}
	for _, tt := range tests {
		in := net.ParseIP(tt.in)
		if IsPrivateIP(in) != tt.out {
			t.Errorf("IsPrivateIP(%v) got: %v", in, tt.out)
		}
	}
}
