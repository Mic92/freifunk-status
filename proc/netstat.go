package proc

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	PROC_TCP  = "/proc/net/tcp"
	PROC_UDP  = "/proc/net/udp"
	PROC_TCP6 = "/proc/net/tcp6"
	PROC_UDP6 = "/proc/net/udp6"
)

var STATE = map[string]string{
	"01": "ESTABLISHED",
	"02": "SYN_SENT",
	"03": "SYN_RECV",
	"04": "FIN_WAIT1",
	"05": "FIN_WAIT2",
	"06": "TIME_WAIT",
	"07": "CLOSE",
	"08": "CLOSE_WAIT",
	"09": "LAST_ACK",
	"0A": "LISTEN",
	"0B": "CLOSING",
}

type Connection struct {
	Ip          net.IP
	Port        uint
	ForeignIp   net.IP
	ForeignPort uint
	State       string
}

func hexToBytes(chars []byte) []byte {
	bytes := make([]byte, len(chars)/2)
	for idx, c := range chars {
		switch {
		case '0' <= c && c <= '9':
			c -= '0'
		case 'a' <= c && c <= 'f':
			c -= 'a' + 10
		case 'A' <= c && c <= 'F':
			c -= 'A' + 10
		}
		if idx%2 == 1 {
			bytes[(idx-1)/2] += c << 4
		} else {
			bytes[idx/2] += c
		}
	}
	return bytes
}

func parseAddr(addr string) (ip net.IP, port uint) {
	ip_port := strings.Split(addr, ":")
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, hexToBytes([]byte(ip_port[0])))
	ip = net.IP(buf.Bytes())
	p := hexToBytes([]byte(ip_port[1]))
	port = uint(p[0])<<8 + uint(p[1])
	return
}

func netstat(path string) ([]Connection, error) {
	f, err := os.Open(path)
	if err != nil {
		e := fmt.Errorf("Error opening %s: %s", path, err)
		return nil, e
	}
	scanner := bufio.NewScanner(f)

	var connections []Connection

	n := 0
	for scanner.Scan() {
		n += 1
		if n < 2 {
			continue
		}
		var ignored, laddr, raddr, state string
		_, err := fmt.Sscan(scanner.Text(), &ignored, &laddr, &raddr)
		if err != nil {
			continue
		}
		// local ip and port
		// TODO parsing still incorrect
		ip, port := parseAddr(laddr)
		fip, fport := parseAddr(raddr)
		// fmt.Printf("laddr: %s:%d/raddr: %s:%d\n", ip.String(), port, fip.String(), fport)

		c := Connection{ip, port, fip, fport, STATE[state]}
		connections = append(connections, c)
	}

	return connections, nil
}

func Tcp4() ([]Connection, error) {
	return netstat(PROC_TCP)
}
func Tcp6() ([]Connection, error) {
	return netstat(PROC_TCP6)
}
func Udp4() ([]Connection, error) {
	return netstat(PROC_UDP)
}
func Udp6() ([]Connection, error) {
	return netstat(PROC_UDP6)
}
