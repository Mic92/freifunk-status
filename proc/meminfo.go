package proc

import (
	"bufio"
	"fmt"
	"os"
)

const MEMINFO_PATH = "/proc/meminfo"

func Meminfo() (map[string]string, error) {
	f, err := os.Open(MEMINFO_PATH)
	if err != nil {
		e := fmt.Errorf("Error opening %s: %s", MEMINFO_PATH, err)
		return nil, e
	}
	scanner := bufio.NewScanner(f)
	stats := make(map[string]string)
	line := 0
	for scanner.Scan() {
		line += 1
		if line < 3 {
			continue
		}
		var name, value string
		fmt.Sscan(scanner.Text(), &name, &value)
		name = name[:len(name)-1]
		stats[name] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse nvram.conf: %s", err)
	}
	return stats, nil
}
