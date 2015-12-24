package bmxd

import (
	"fmt"

	"github.com/Mic92/freifunk-status/ip"
)

type GlobalHna struct {
	Target    string `json:"target"`
	Via       string `json:"via"`
	Interface string `json:"interface"`
}

func Hna() ([]Link, []GlobalHna, error) {
	routes, err := ip.Routes("bat_hna", false)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list bat_hna table: %s", err)
	}
	links := make([]Link, 0)
	globals := make([]GlobalHna, 0)
	for _, route := range routes {
		if route.Gateway == "" {
			links = append(links, Link{route.Net, route.Dev})
		} else {
			global := GlobalHna{route.Net, route.Gateway, route.Dev}
			globals = append(globals, global)
		}
	}
	return links, globals, nil
}
