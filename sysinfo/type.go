package sysinfo

import (
	"github.com/Mic92/freifunk-status/bmxd"
)

type Firmware struct {
	Version     string `json:"version"`
	Id          string `json:"DISTRIB_ID"`
	Release     string `json:"DISTRIB_RELEASE"`
	Revision    string `json:"DISTRIB_REVISION"`
	Codename    string `json:"DISTRIB_CODENAME"`
	Target      string `json:"DISTRIB_TARGET"`
	Description string `json:"DISTRIB_DESCRIPTION"`
}

type System struct {
	Uptime     string   `json:"uptime"`
	Uname      string   `json:"uname"`
	Nameserver []string `json:"nameserver"`
	Date       string   `json:"date"`
	Board      string   `json:"board"`
	Model      string   `json:"model"`
	Model2     string   `json:"model2"`
	Cpuinfo    string   `json:"cpuinfo"`
	Bmxd       string   `json:"bmxd"`
}

type Common struct {
	City   string `json:"city"`
	Node   string `json:"node"`
	Domain string `json:"domain"`
	Ip     string `json:"ip"`
}

type Gps struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Altitude  string `json:"altitude"`
}
type Contact struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Email    string `json:"email"`
	Note     string `json:"note"`
}

type Bmxd struct {
	RoutingTables struct {
		Route struct {
			Link []bmxd.Link `json:"link"`
		} `json:"route"`
		Hna struct {
			Link   []bmxd.Link      `json:"link"`
			Global []bmxd.GlobalHna `json:"global"`
		} `json:"hna"`
	}
	Gateways bmxd.GatewayStatus `json:"gateways"`
	Info     []string           `json:"info"`
}

type Accept struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

type Connection struct {
	Local struct {
		Ip   string `json:"ip"`
		Port string `json:"port"`
	} `json:"local"`
	Foreign struct {
		Ip   string `json:"ip"`
		Port string `json:"port"`
	} `json:"foreign"`
}

type Sysinfo struct {
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
	Data      struct {
		Firmware Firmware `json:"firmware"`
		System   System   `json:"system"`
		Common   Common   `json:"common"`
		Gps      Gps      `json:"gps"`
		Contact  Contact  `json:"contact"`
		Backbone struct {
			Accept []Accept `json:"accept"`
		} `json:"backbone"`
		Statistics     map[string]interface{} `json:"statistics"`
		Bmxd           Bmxd                   `json:"bmxd"`
		InternetTunnel string                 `json:"internet_tunnel"`
		Connections    []Connection           `json:"connection"`
	} `json:"data"`
}
