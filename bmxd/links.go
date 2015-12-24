package bmxd

import (
	"fmt"

	"github.com/Mic92/freifunk-status/ip"
)

type Link struct {
	Target    string `json:"target"`
	Interface string `json:"interface"`
}

func Links() ([]Link, error) {
	routes, err := ip.Routes("bat_route", false)
	if err != nil {
		return nil, fmt.Errorf("failed to list bat_route table: %s", err)
	}
	links := make([]Link, 0)
	for _, route := range routes {
		if route.Scope == "link" {
			l := Link{route.Net, route.Dev}
			links = append(links, l)
		}
	}
	return links, nil
}
