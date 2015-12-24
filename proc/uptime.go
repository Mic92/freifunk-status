package proc

import (
	"fmt"
	"syscall"
	"time"
)

func Uptime() (string, error) {
	var info syscall.Sysinfo_t
	err := syscall.Sysinfo(&info)
	if err != nil {
		return "", err
	}
	days := info.Uptime / (60 * 60 * 24)

	uptime := ""
	if days != 0 {
		s := ""
		if days > 1 {
			s = "s"
		}
		uptime += fmt.Sprintf("%d day%s, ", days, s)
	}

	minutes := info.Uptime / 60
	hours := minutes / 60
	hours %= 24
	minutes %= 60

	uptime += fmt.Sprintf("%2d:%02d", hours, minutes)

	h, m, s := time.Now().Clock()
	ret := fmt.Sprintf(" %d:%d:%d up %s  load average: %.2f, %.2f, %.2f",
		h, m, s,
		uptime,
		float64(info.Loads[0])/65536,
		float64(info.Loads[1])/65536,
		float64(info.Loads[2])/65536)
	return ret, nil
}
