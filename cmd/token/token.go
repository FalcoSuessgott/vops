package token

import (
	"fmt"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

type tokenOptions struct {
	Cluster string
}

func newDefaultTokenOptions() *tokenOptions {
	return &tokenOptions{}
}

// NewTokenCmd login command.
func NewTokenCmd(cfg string) *cobra.Command {
	o := newDefaultTokenOptions()

	cmd := &cobra.Command{
		Use:           "token",
		Short:         "copy the token from the token exec command to your clipboard buffer",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return config.ValidateConfig(cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("[ Token ]")
			fmt.Printf("using %s\n", cfg)

			config, err := config.ParseConfig(cfg)
			if err != nil {
				return err
			}

			cluster, err := config.GetCluster(o.Cluster)
			if err != nil {
				return err
			}

			fmt.Printf("\n[ %s ]\n", cluster.Name)
			fmt.Printf("copying token for cluster %s\n", cluster.Name)

			cluster.ExtraEnv["VAULT_ADDR"] = cluster.Addr
			cluster.ExtraEnv["VAULT_TOKEN"] = cluster.Token

			if err := cluster.ApplyEnvironmentVariables(cluster.ExtraEnv); err != nil {
				return err
			}

			if err := cluster.RunTokenExecCommand(); err != nil {
				return err
			}

			if err := clipboard.WriteAll(cluster.Token); err != nil {
				return err
			}

			fmt.Printf("token for cluster %s copied to clipboard buffer.\n", cluster.Name)

			return nil
		},
	}

	return cmd
}
