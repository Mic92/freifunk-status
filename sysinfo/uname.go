package sysinfo

import (
	"fmt"
	"syscall"
)

func fieldToString(bs [65]int8) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		if v == 0 {
			return string(b[:i])
		}
		if v < 0 {
			b[i] = byte(256 + int(v))
		} else {
			b[i] = byte(v)
		}
	}
	return string(b)
}

func uname() string {
	var n syscall.Utsname
	// never fails
	syscall.Uname(&n)
	return fmt.Sprintf("%s %s %s %s %s %s",
		fieldToString(n.Sysname),
		fieldToString(n.Nodename),
		fieldToString(n.Release),
		fieldToString(n.Version),
		fieldToString(n.Machine),
		fieldToString(n.Domainname))
}
