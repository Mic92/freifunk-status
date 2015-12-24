package sysinfo

import (
	"bytes"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Mic92/freifunk-status/bmxd"
	"github.com/Mic92/freifunk-status/ddmesh"
	"github.com/Mic92/freifunk-status/nvram"
	"github.com/Mic92/freifunk-status/proc"
)

const (
	NVRAM_CONF = "/etc/nvram.conf"
)

func system() System {
	var system System
	system.Uname = uname()
	system.Date = time.Now().Format(time.UnixDate)
	system.Nameserver = make([]string, 0)
	uptime, err := proc.Uptime()
	if err != nil {
		log.Printf("failed to get uptime: %s", err)
	}
	system.Uptime = uptime
	bmxStatus, err := bmxd.Status()
	if err != nil {
		log.Printf("failed to get bmx status: %s", err)
	}
	system.Bmxd = bmxStatus
	return system
}

func backbone(config map[string]string) []Accept {
	accepts := make([]Accept, 0)
	for key, value := range config {
		if strings.HasPrefix(key, "backbone_accept_") {
			split := strings.SplitN(value, ":", 2)
			if len(split) != 2 {
				continue
			}
			accepts = append(accepts, Accept{split[0], split[0]})
		}
		if strings.HasPrefix(key, "backbone_range_") {
			split := strings.SplitN(value, ":", 2)
			if len(split) != 2 {
				continue
			}
			accepts = append(accepts, Accept{split[0], split[1]})
		}
	}
	return accepts
}

func common(config map[string]string, node ddmesh.Node) Common {
	var c Common
	c.Node = strconv.Itoa(node.Id)
	c.City = config["city"]
	c.Ip = node.IPv4.String()
	c.Domain = ddmesh.DOMAIN
	return c
}

func contact(config map[string]string) Contact {
	var c Contact
	c.Name = config["contact_name"]
	c.Location = config["contact_location"]
	c.Email = config["contact_email"]
	c.Note = config["contact_note"]
	return c
}

func gps(config map[string]string) Gps {
	var g Gps
	g.Latitude = config["gps_latitude"]
	g.Longitude = config["gps_longitude"]
	g.Altitude = config["gps_altitude"]
	return g
}

func connections(node ddmesh.Node) ([]Connection, error) {
	conns, err := proc.Tcp4()
	if err != nil {
		log.Printf("failed to get open tcp connections: %s\n", err)
	} else {
		conns = []proc.Connection{}
	}

	conns2, err := proc.Udp4()
	if err != nil {
		log.Printf("failed to get open udp connections: %s\n", err)
	} else {
		conns2 = []proc.Connection{}
	}
	conns = append(conns, conns2...)
	filtered := make([]Connection, 0)
	for _, conn := range conns {
		if bytes.Equal(conn.Ip, node.IPv4) ||
			len(conn.Ip) == 4 && conn.Ip[0] == 169 && conn.Ip[1] == 254 {
			c := Connection{}
			c.Local.Ip = conn.Ip.String()
			c.Local.Port = strconv.Itoa(int(conn.Port))
			c.Foreign.Ip = conn.ForeignIp.String()
			c.Foreign.Port = strconv.Itoa(int(conn.ForeignPort))
			filtered = append(filtered)
		}
	}
	return filtered, nil
}

func New() Sysinfo {
	info := Sysinfo{}
	info.Version = "8"
	info.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)

	config, err := nvram.Open(NVRAM_CONF)
	id := 0
	if err == nil {
		id, _ = strconv.Atoi(config["ddmesh_node"])
	} else {
		log.Printf("Failed to get nvram.conf: %s: %s\n", NVRAM_CONF, err)
		config = make(map[string]string, 0)
	}
	node := ddmesh.NewNode(id)
	info.Data.Common = common(config, node)
	info.Data.Contact = contact(config)
	info.Data.Gps = gps(config)
	info.Data.Backbone.Accept = backbone(config)
	info.Data.System = system()

	stats, err := statistics(config)
	if err != nil {
		log.Printf("Failed to get statistics: %s\n", err)
	}
	info.Data.Statistics = stats

	links, err := bmxd.Links()
	if err != nil {
		log.Printf("Failed to get bmxd links: %s\n", err)
	}
	info.Data.Bmxd.RoutingTables.Route.Link = links

	links, globalHna, err := bmxd.Hna()
	if err != nil {
		log.Printf("Failed to get bmxd hna: %s\n", err)
	}
	info.Data.Bmxd.RoutingTables.Hna.Link = links
	info.Data.Bmxd.RoutingTables.Hna.Global = globalHna

	status, err := bmxd.Gateways()
	if err != nil {
		log.Printf("Failed to get bmxd gateways: %s\n", err)
	}
	info.Data.Bmxd.Gateways = status

	bmxdInfo, err := bmxd.Info()
	if err == nil {
		info.Data.Bmxd.Info = bmxdInfo
	} else {
		info.Data.Bmxd.Info = []string{}
	}

	conns, err := connections(node)
	if err != nil {
		log.Printf("Failed to get open connections: %s\n", err)
	}
	info.Data.Connections = conns

	return info
}
