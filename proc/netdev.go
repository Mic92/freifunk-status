package proc

import (
	"bufio"
	"fmt"
	"os"
)

const NETDEV_PATH = "/proc/net/dev"

type IfStat struct {
	Name         string
	RxBytes      uint
	RxPackets    uint
	RxErrs       uint
	RxDrop       uint
	RxFifo       uint
	RxFrame      uint
	RxCompressed uint
	RxMulticast  uint
	TxBytes      uint
	TxPackets    uint
	TxErrs       uint
	TxDrop       uint
	TxFifo       uint
	TxFrame      uint
	TxCompressed uint
	TxMulticast  uint
}

func Netdev() (map[string]IfStat, error) {
	f, err := os.Open(NETDEV_PATH)
	if err != nil {
		e := fmt.Errorf("Error opening %s: %s", NETDEV_PATH, err)
		return nil, e
	}
	scanner := bufio.NewScanner(f)
	stats := make(map[string]IfStat)
	line := 0
	for scanner.Scan() {
		line += 1
		if line < 3 {
			continue
		}
		var s IfStat
		fmt.Sscan(scanner.Text(),
			&s.Name,
			&s.RxBytes,
			&s.RxPackets,
			&s.RxErrs,
			&s.RxDrop,
			&s.RxFifo,
			&s.RxFrame,
			&s.RxCompressed,
			&s.RxMulticast,
			&s.TxBytes,
			&s.TxPackets,
			&s.TxErrs,
			&s.TxDrop,
			&s.TxFifo,
			&s.TxFrame,
			&s.TxCompressed,
			&s.TxMulticast)

		s.Name = s.Name[:len(s.Name)-1]
		stats[s.Name] = s
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse nvram.conf: %s", err)
	}
	return stats, nil
}
