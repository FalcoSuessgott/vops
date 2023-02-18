package token

import (
	"fmt"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/flags"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

type tokenOptions struct {
	Cluster string
}

// NewTokenCmd token command.
func NewTokenCmd(cfg string) *cobra.Command {
	var c *config.Config

	o := &tokenOptions{}

	cmd := &cobra.Command{
		Use:           "token",
		Short:         "copy the token from the token exec command to your clipboard buffer",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			fmt.Println("[ Copy Token ]")
			fmt.Printf("using %s\n", cfg)

			c, err = config.ParseConfig(cfg)
			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := c.GetCluster(o.Cluster)
			if err != nil {
				return err
			}

			fmt.Printf("\n[ %s ]\n", cluster.Name)
			fmt.Printf("copying token for cluster %s\n", cluster.Name)

			if cluster.TokenExecCmd == "" {
				return fmt.Errorf("no token exec command defined")
			}

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

	flags.ClusterFlag(cmd, o.Cluster)

	return cmd
}
