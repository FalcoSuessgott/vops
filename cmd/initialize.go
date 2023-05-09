package cmd

import (
	"fmt"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/fs"
	"github.com/FalcoSuessgott/vops/pkg/utils"
	"github.com/FalcoSuessgott/vops/pkg/vault"
	"github.com/spf13/cobra"
)

var status bool

func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "init",
		Aliases:       []string{"i"},
		Short:         "initialize a vault cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if allCluster {
				for _, cluster := range cfg.Cluster {
					if err := initializeCluster(cluster); err != nil {
						return err
					}
				}

				return nil
			}

			cluster, err := cfg.GetCluster(cluster)
			if err != nil {
				return err
			}

			if err := initializeCluster(*cluster); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&status, "status", "S", status, "print the initialization status of a cluster")

	return cmd
}

func initializeCluster(cluster config.Cluster) error {
	fmt.Printf("\n[ %s ]\n", cluster.Name)
	fmt.Printf("attempting intialization of cluster \"%s\" with %d shares and a threshold of %d\n", cluster.Name, cluster.Keys.Shares, cluster.Keys.Threshold)

	if cluster.Keys == nil || cluster.Keys.Path == "" {
		return fmt.Errorf("a keyfile location is required")
	}

	if err := cluster.ApplyEnvironmentVariables(cluster.ExtraEnv); err != nil {
		return err
	}

	v, err := vault.NewClient(cluster.Addr)
	if err != nil {
		return err
	}

	isInitialized, err := v.IsInitialized()
	if err != nil {
		return err
	}

	if status || isInitialized {
		fmt.Printf("%s is already initialized.\n", cluster.Name)

		return nil
	}

	resp, err := v.Initialize(cluster.Keys.Shares, cluster.Keys.Threshold, cluster.Keys.Autounseal)
	if err != nil {
		return err
	}

	if err := fs.WriteToFile(utils.ToJSON(resp), cluster.Keys.Path); err != nil {
		return err
	}

	fmt.Printf("successfully initialized %s and wrote keys to %s.\n", cluster.Name, cluster.Keys.Path)

	return nil
}
