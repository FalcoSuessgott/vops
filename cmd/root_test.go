package cmd

import (
	"bytes"
	"io"
	"testing"

	"github.com/FalcoSuessgott/vops/pkg/fs"
	"github.com/stretchr/testify/require"
)

func TestE2E(t *testing.T) {
	testCases := []struct {
		name    string
		command []string
		err     bool
	}{
		{
			name:    "example config",
			command: []string{"config", "example", "--config", "vops.yml"},
			err:     false,
		},
		{
			name:    "config val",
			command: []string{"config", "validate"},
			err:     false,
		},
		{
			name:    "initialize",
			command: []string{"init", "-c", "cluster-1"},
			err:     false,
		},
		{
			name:    "unseal",
			command: []string{"unseal", "-c", "cluster-1"},
			err:     false,
		},
		{
			name:    "seal",
			command: []string{"seal", "-c", "cluster-1"},
			err:     false,
		},
		{
			name:    "unseal",
			command: []string{"unseal", "-c", "cluster-1"},
			err:     false,
		},
		{
			name:    "snapsht save",
			command: []string{"snapshot", "save", "-c", "cluster-1"},
			err:     false,
		},
		{
			name:    "login",
			command: []string{"login", "-c", "cluster-1"},
			err:     false,
		},
		{
			name:    "rekey",
			command: []string{"rekey", "-c", "cluster-1"},
			err:     false,
		},
		{
			name:    "custom",
			command: []string{"custom", "-x", "status", "cluster-1"},
			err:     false,
		},
		{
			name:    "generate-root",
			command: []string{"generate-root", "-c", "cluster-1"},
			err:     false,
		},
	}

	b := bytes.NewBufferString("")
	cmd := NewRootCmd("", b)

	for _, tc := range testCases {
		cmd.SetArgs(tc.command)

		err := cmd.Execute()
		if tc.err {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}

		// write example config file
		if tc.name == "example config" {
			out, _ := io.ReadAll(b)

			// nolint
			fs.WriteToFile(out, "vops.yml")
		}
	}
}
