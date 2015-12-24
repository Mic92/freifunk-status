package sysinfo

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/Mic92/freifunk-status/proc"
)

func fmtTraffic(s proc.IfStat) string {
	return fmt.Sprintf("%d,%d", s.RxBytes, s.TxBytes)
}

func netStatistics(stats map[string]interface{},
	config map[string]string,
	ifstats map[string]proc.IfStat) {
	for ifname, ifstat := range ifstats {
		if ifname == config["ifname"] ||
			ifname == "lo" ||
			ifname == "icvpn" ||
			ifname == "bat" {
			continue
		}
		traffic := fmtTraffic(ifstat)
		stats["traffic_"+ifname] = traffic
	}
	stats["accepted_user_count"] = "0"
	stats["dhcp_count"] = "0"
	stats["dhcp_lease"] = "0"
	stats["traffic_adhoc"] = ""
	stats["traffic_ap"] = ""
	stats["traffic_wan"] = ""
	stats["traffic_ovpn"] = ""
	stats["traffic_icvpn"] = ""

	if ifstat, ok := ifstats[config["ifname"]]; ok {
		stats["traffic_wan"] = fmtTraffic(ifstat)
	}
	if ifstat, ok := ifstats["vpn0"]; ok {
		stats["traffic_ovpn"] = fmtTraffic(ifstat)
	}
	if ifstat, ok := ifstats["icvpn"]; ok {
		stats["traffic_icvpn"] = fmtTraffic(ifstat)
	}
}

func gatewayUsage() []map[string]string {
	content, err := ioutil.ReadFile("/var/statistic/gateway_usage")
	if err != nil {
		return make([]map[string]string, 0)
	}
	lines := strings.Split(string(content), "\n")
	usage := make([]map[string]string, 0)
	for _, line := range lines {
		ary := strings.SplitN(line, ":", 2)
		if len(ary) != 2 {
			continue
		}
		entry := make(map[string]string, 1)
		entry[ary[0]] = ary[1]
		usage = append(usage, entry)
	}
	return usage
}

func statistics(config map[string]string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	ifstats, err := proc.Netdev()
	if err != nil {
		log.Printf("Failed to get Netdev statistics: %s\n", err)
	}
	netStatistics(stats, config, ifstats)
	meminfo, err := proc.Meminfo()
	if err != nil {
		log.Printf("Failed to get Memory statistics: %s\n", err)
	}
	for field, value := range meminfo {
		stats["meminfo_"+field] = value
	}
	if s, err := ioutil.ReadFile("/proc/loadavg"); err == nil {
		stats["cpu_load"] = string(s[:len(s)-1])
	} else {
		log.Printf("Failed to read /proc/loadavg: %s", err)
	}
	if s, err := ioutil.ReadFile("/proc/stat"); err == nil {
		stat := strings.SplitN(string(s), "\n", 2)[0]
		stats["cpu_stat"] = stat[5:len(stat)]
	} else {
		log.Printf("Failed to read /proc/stat: %s", err)
	}
	stats["gateway_usage"] = gatewayUsage()
	return stats, nil
}
