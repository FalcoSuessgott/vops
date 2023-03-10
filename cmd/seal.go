package cmd

import (
	"fmt"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/vault"
	"github.com/spf13/cobra"
)

func sealCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "seal",
		Aliases:       []string{"s"},
		Short:         "seal a cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if allCluster {
				for _, cluster := range cfg.Cluster {
					if err := sealCluster(cluster); err != nil {
						return err
					}
				}

				return nil
			}

			cluster, err := cfg.GetCluster(cluster)
			if err != nil {
				return err
			}

			if err := sealCluster(*cluster); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func sealCluster(cluster config.Cluster) error {
	fmt.Printf("\n[ %s ]\n", cluster.Name)

	if cluster.TokenExecCmd == "" {
		return fmt.Errorf("no token exec command defined")
	}

	if err := cluster.ApplyEnvironmentVariables(cluster.ExtraEnv); err != nil {
		return err
	}

	if err := cluster.RunTokenExecCommand(); err != nil {
		return err
	}

	fmt.Println("executed token exec command")

	c, err := vault.NewTokenClient(cluster.Addr, cluster.Token)
	if err != nil {
		return err
	}

	if err := c.Seal(); err != nil {
		return err
	}

	fmt.Printf("cluster \"%s\" sealed\n", cluster.Name)

	return nil
}
