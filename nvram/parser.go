package nvram

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func Open(f string) (map[string]string, error) {
	file, err := os.Open(f) // For read access.
	if err != nil {
		return nil, err
	}
	return Read(io.Reader(file))
}

func Read(f io.Reader) (map[string]string, error) {
	scanner := bufio.NewScanner(f)
	config := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		splits := strings.SplitN(line, "=", 2)
		if len(splits) == 2 {
			config[splits[0]] = splits[1]
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse nvram.conf: %s", err)
	}
	return config, nil
}
