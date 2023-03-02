package cmd

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

func tokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "token",
		Short:         "copy the token from the token exec command to your clipboard buffer",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := cfg.GetCluster(cluster)
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

	return cmd
}
