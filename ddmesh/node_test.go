package ddmesh_test

import (
	"github.com/Mic92/freifunk-status/ddmesh"
	"testing"
)

func TestNewNode(t *testing.T) {
	n := ddmesh.NewNode(51004)
	if n.IPv4.String() != "10.201.200.5" {
		t.Errorf("incorrect ipv4: %s", n.IPv4.String())
	}
	if n.IPv6.String() != "fd11:11ae:7466:c73c::" {
		t.Errorf("incorrect ipv6: %s", n.IPv6.String())
	}
}
