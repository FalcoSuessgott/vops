package seal

import (
	"fmt"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/vault"
	"github.com/spf13/cobra"
)

type sealOptions struct {
	Cluster    string
	AllCluster bool
}

func newDefaultsealOptions() *sealOptions {
	return &sealOptions{
		AllCluster: true,
	}
}

// NewSealCmd vops seal command.
func NewSealCmd(cfg string) *cobra.Command {
	o := newDefaultsealOptions()

	cmd := &cobra.Command{
		Use:           "seal",
		Aliases:       []string{"s"},
		Short:         "seals a single or all cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("[ Seal ]")
			fmt.Printf("using %s\n", cfg)

			config, err := config.ParseConfig(cfg)
			if err != nil {
				return err
			}

			// All Cluster
			if o.AllCluster {
				for _, cluster := range config.Cluster {
					if err := sealCluster(cluster); err != nil {
						return err
					}
				}

				return nil
			}

			// Single Node
			cluster, err := config.GetCluster(o.Cluster)
			if err != nil {
				return err
			}

			if err := sealCluster(*cluster); err != nil {
				return err
			}

			fmt.Printf("cluster \"%s\" sealed\n", cluster.Name)

			return nil
		},
	}

	cmd.Flags().StringVarP(&o.Cluster, "cluster", "c", o.Cluster, "name of the vault cluster to initialize")
	cmd.Flags().BoolVarP(&o.AllCluster, "all-cluster", "A", o.AllCluster, "unseal all cluster defined in the vops configuration file")

	return cmd
}

func sealCluster(cluster config.Cluster) error {
	fmt.Printf("\n[ %s ]\n", cluster.Name)

	if err := cluster.ApplyEnvironmentVariables(); err != nil {
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
