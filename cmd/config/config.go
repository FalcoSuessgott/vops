package config

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/utils"
	"github.com/spf13/cobra"
)

// NewConfigCmd config command.
func NewConfigCmd(cfg string) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "config",
		Aliases:       []string{"c", "cfg"},
		Short:         "config",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return config.ValidateConfig(cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		newConfigExampleCmd(),
		newConfigValidateCmd(cfg),
	)

	return cmd
}

func newConfigExampleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "example",
		Short:         "prints an example configuration",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			exampleCfg := &config.Config{
				CustomCmds: map[string]interface{}{
					"list-peers": "vault operator raft list-peers",
					"status":     "vault status",
				},
				Cluster: []config.Cluster{
					{
						Name:         "cluster-1",
						Addr:         "http://127.0.0.1:8200",
						TokenExecCmd: "jq -r '.root_token' {{ .Keys.Path }}",
						Keys: &config.KeyConfig{
							Threshold: 1,
							Shares:    1,
							Path:      "{{ .Name }}.json",
						},
						SnapshotDir: "{{ .Name }}/",
						Nodes: []string{
							"{{ .Addr }}",
						},
						ExtraEnv: map[string]interface{}{
							"VAULT_TLS_SKIP_VERIFY": true,
						},
					},
				},
			}

			fmt.Println(string(utils.ToYAML(&exampleCfg)))

			return nil
		},
	}

	return cmd
}

func newConfigValidateCmd(cfg string) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "validate",
		Aliases:       []string{"v", "val"},
		Short:         "validates a vops configuration file",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := config.ParseConfig(cfg)
			if err != nil {
				return err
			}

			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.TabIndent)
			fmt.Fprintln(w, "Name\tAddr\tTokenExecCmd\tToken Policies\tNodes\tKey\tSnapshotDir")

			for _, cluster := range config.Cluster {
				fmt.Fprintln(w, cluster)
			}

			w.Flush()

			return nil
		},
	}

	return cmd
}
