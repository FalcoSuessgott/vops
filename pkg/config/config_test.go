package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	testCases := []struct {
		name   string
		path   string
		expCfg *Config
		err    bool
	}{
		{
			name: "single cluster",
			path: "testdata/config_1.yaml",
			expCfg: &Config{
				CustomCmds: map[string]interface{}{
					"list-peers": "vault operator raft list-peers",
				},
				Cluster: []Cluster{
					{
						Name:         "cluster-1",
						Addr:         "https://test.vault.de",
						TokenExecCmd: "vault login",
						Keys: &KeyConfig{
							Autounseal: false,
							Path:       "file.json",
							Shares:     5,
							Threshold:  3,
						},
						SnapshotDir: "snapshot/",
						Nodes: []string{
							"vault-server-01", "vault-server-02", "vault-server-03",
						},
						Env: map[string]interface{}{},
						ExtraEnv: map[string]interface{}{
							"VAULT_SKIP_VERIFY": true,
						},
					},
				},
			},
			err: false,
		},
	}

	for _, tc := range testCases {
		cfg, err := ParseConfig(tc.path)
		if tc.err {
			require.Error(t, err, tc.name)
		}

		for i := range cfg.Cluster {
			cfg.Cluster[i].Env = map[string]interface{}{}
		}

		fmt.Println(tc.name, cfg)
		assert.Equal(t, tc.expCfg, cfg, tc.name)
	}
}
