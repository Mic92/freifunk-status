package bmxd

import (
	"fmt"
	"regexp"
	"strings"
)

type Gateway struct {
	Ip string `json:"ip"`
}

type GatewayStatus struct {
	Selected  string    `json:"selected"`
	Preferred string    `json:"preferred"`
	Gateways  []Gateway `json:"gateways"`
}

var preferred = regexp.MustCompile(`preferred gateway: ([0-9.]+)`)

func Gateways() (GatewayStatus, error) {
	var status GatewayStatus
	out, err := run("-c", "--gateways")
	if err != nil {
		return status, err
	}
	lines := strings.Split(out, "\n")
	if len(lines) == 0 {
		return status, fmt.Errorf("Got not output for `bmxd -c --gateways`")
	}
	m := preferred.FindStringSubmatch(lines[0])
	if m != nil && len(m) > 1 {
		status.Preferred = m[1]
	}
	lines = lines[1:]
	gateways := make([]Gateway, 0)
	for _, line := range lines {
		var ip string
		fmt.Sscan(line, &ip)
		gateways = append(gateways, Gateway{ip})
	}
	status.Gateways = gateways
	return status, nil
}
