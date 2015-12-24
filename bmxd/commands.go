package bmxd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func run(args ...string) (string, error) {
	cmd := exec.Command("bmxd", args...)
	out := bytes.Buffer{}
	stderr := bytes.Buffer{}
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		cmd := strings.Join(args, " ")
		return "", fmt.Errorf("failed to run bmxd:\n$ bmxd %s: %s", cmd, err)
	}
	return strings.TrimSpace(string(out.Bytes())), nil
}

func Status() (string, error) {
	return run("-c", "--status")
}

func Info() ([]string, error) {
	out, err := run("-ci")
	if err != nil {
		return nil, err
	}
	return strings.Split(out, "\n"), nil
}
