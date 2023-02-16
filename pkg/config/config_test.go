package config

import (
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
							Path:      "file.json",
							Shares:    5,
							Threshold: 3,
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
		{
			name: "render config",
			path: "testdata/config_2.yaml",
			expCfg: &Config{
				Cluster: []Cluster{
					{
						Name:         "cluster-1",
						Addr:         "https://test.vault.de",
						TokenExecCmd: "vault login https://test.vault.de",
						Nodes: []string{
							"https://test.vault.de",
						},
						Env: map[string]interface{}{},
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

		// ignore envs
		for i := range cfg.Cluster {
			cfg.Cluster[i].Env = map[string]interface{}{}
		}

		assert.Equal(t, tc.expCfg, cfg, tc.name)
	}
}
