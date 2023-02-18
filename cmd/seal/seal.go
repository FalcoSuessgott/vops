package seal

import (
	"fmt"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/flags"
	"github.com/FalcoSuessgott/vops/pkg/vault"
	"github.com/spf13/cobra"
)

type sealOptions struct {
	Cluster    string
	AllCluster bool
}

// NewSealCmd vops seal command.
func NewSealCmd(cfg string) *cobra.Command {
	var c *config.Config

	o := &sealOptions{}

	cmd := &cobra.Command{
		Use:           "seal",
		Aliases:       []string{"s"},
		Short:         "seal a cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			fmt.Println("[ Seal ]")
			fmt.Printf("using %s\n", cfg)

			c, err = config.ParseConfig(cfg)
			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.AllCluster {
				for _, cluster := range c.Cluster {
					if err := sealCluster(cluster); err != nil {
						return err
					}
				}

				return nil
			}

			cluster, err := c.GetCluster(o.Cluster)
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

	flags.AllClusterFlag(cmd, o.AllCluster)
	flags.ClusterFlag(cmd, o.Cluster)

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
