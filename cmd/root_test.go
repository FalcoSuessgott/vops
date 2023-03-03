package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/FalcoSuessgott/vops/pkg/fs"
	"github.com/stretchr/testify/require"
)

func TestE2E(t *testing.T) {
	testCases := []struct {
		name            string
		command         []string
		err             bool
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
			command: []string{"init", "-A"},
			err:     false,
		},
		{
			name:    "unseal",
			command: []string{"unseal", "-A"},
			err:     false,
		},
		{
			name:    "seal",
			command: []string{"seal", "-A"},
			err:     false,
		},
		{
			name:    "unseal",
			command: []string{"unseal", "-A"},
			err:     false,
		},
		{
			name:    "snapsht save",
			command: []string{"snapshot", "save", "-A"},
			err:     false,
		},
		{
			name:    "login",
			command: []string{"login", "-A"},
			err:     false,
		},
		{
			name:    "rekey",
			command: []string{"rekey", "-A"},
			err:     false,
		},
		{
			name:    "custom list",
			command: []string{"custom", "-l"},
			err:     false,
		},
		// {
		// 	name:    "custom",
		// 	command: []string{"custom", "--command", "status", "-A"},
		// 	err:     false,
		// },
		// {
		// 	name:    "custom error",
		// 	command: []string{"custom", "-x", "invalid", "-A"},
		// 	err:     true,
		// },
		{
			name:    "adhoc",
			command: []string{"adhoc", "-x", "vault status", "-A"},
			err:     false,
		},
		{
			name:    "adhoc error",
			command: []string{"adhoc"},
			err:     true,
		},
		{
			name:    "generate-root",
			command: []string{"generate-root", "-A"},
			err:     false,
		},
	}

	b := bytes.NewBufferString("")

	// 
	os.Unsetenv("VOPS_CONFIG")
	for _, tc := range testCases {
		fmt.Println(tc.command)

		cmd := NewRootCmd("", b)
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
