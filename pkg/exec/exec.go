package exec

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Run runs the given command and returns the output.
//nolint: gosec
func Run(cmd []string) ([]byte, error) {
	var stdout, stderr bytes.Buffer

	c := exec.Command("bash", "-c", strings.Join(cmd, " "))

	c.Stdout = &stdout
	c.Stderr = &stderr

	if c.Run() != nil {
		return nil, fmt.Errorf("error while running command: %v", stderr.String())
	}

	return stdout.Bytes(), nil
}
