package cmd

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/utils"
	"github.com/spf13/cobra"
)

func configCmd(w io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "config",
		Aliases:       []string{"c", "cfg"},
		Short:         "validate the vops configuration file",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		configExampleCmd(w),
		configValidateCmd(),
	)

	return cmd
}

func configExampleCmd(w io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "example",
		Short:         "prints an example configuration",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
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

			fmt.Fprintln(w, string(utils.ToYAML(&exampleCfg)))

			return nil
		},
	}

	return cmd
}

func configValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "validate",
		Aliases:       []string{"v", "val"},
		Short:         "validates a vops configuration file",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println()

			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.TabIndent)

			for _, cluster := range cfg.Cluster {
				fmt.Fprintln(w, cluster)
			}

			w.Flush()

			return nil
		},
	}

	return cmd
}
