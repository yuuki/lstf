// +build linux

package netutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNetlinkConnections(t *testing.T) {
	conns, err := NetlinkConnections()
	if err != nil {
		t.Fatalf("should not raise error: %v", err)
	}
	if len(conns) == 0 {
		t.Error("NetlinkConnections() should not be len == 0")
	}
}

func TestParseProcStat(t *testing.T) {
	cur, _ := os.Getwd()
	root := filepath.Join(cur, "../testdata")
	pid := 10000

	pname, err := parseProcStat(root, pid)
	if err != nil {
		t.Fatal(err)
	}

	if pname != "nginx" {
		t.Errorf("process name '%s' should be 'nginx'", pname)
	}
}
