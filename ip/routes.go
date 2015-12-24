package ip

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Route struct {
	Net     string
	Gateway string
	Dev     string
	Proto   string
	Scope   string
	Src     string
	Metric  string
}

func parseRoute(line string) (*Route, error) {
	words := strings.Split(line, " ")
	var r Route
	if len(words) < 3 {
		return nil, fmt.Errorf("Invalid line in routing table: %s", line)
	}
	r.Net = words[0]

	for idx, word := range words[1:] {
		switch word {
		case "via":
			r.Gateway = words[idx+2]
		case "dev":
			r.Dev = words[idx+2]
		case "proto":
			r.Proto = words[idx+2]
		case "scope":
			r.Scope = words[idx+2]
		case "src":
			r.Src = words[idx+2]
		case "metric":
			r.Metric = words[idx+2]
		}
	}
	return &r, nil
}

func Routes(table string, ipv6 bool) ([]Route, error) {
	var version string
	if ipv6 {
		version = "-6"
	} else {
		version = "-4"
	}
	args := []string{"ip", version, "route", "list", "table", table}
	cmd := exec.Command(args[0], args[1:]...)
	out := bytes.Buffer{}
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to execute: %s: %s", args)
	}
	var routes []Route
	for _, line := range strings.Split(out.String(), "\n") {
		if line == "" {
			continue
		}
		if route, err := parseRoute(line); err == nil {
			routes = append(routes, *route)
		} else {
			return routes, err
		}
	}
	return routes, nil
}
