package ddmesh

import (
	"net"
	"strconv"
)

type Node struct {
	Id       int
	IPv4     net.IP
	IPv6     net.IP
	Hostname string
}

const DOMAIN = "freifunk-dresden.de"

func NewNode(id int) Node {
	upper := byte(id / 255)
	lower := byte(id%255 + 1)
	ip4 := net.IPv4(10, 201, upper, lower)

	high := byte((id & ^0xFF) >> 8)
	low := byte(id & 0xFF)
	ip6 := net.IP{0xfd, 0x11, 0x11, 0xae, 0x74, 0x66, high, low, 0, 0, 0, 0, 0, 0, 0, 0}
	return Node{id, ip4, ip6, "r" + strconv.Itoa(id)}
}
