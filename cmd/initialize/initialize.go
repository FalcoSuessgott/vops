package initialize

import (
	"fmt"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/flags"
	"github.com/FalcoSuessgott/vops/pkg/fs"
	"github.com/FalcoSuessgott/vops/pkg/utils"
	"github.com/FalcoSuessgott/vops/pkg/vault"
	"github.com/spf13/cobra"
)

type initOptions struct {
	Shares     int
	Threshold  int
	Cluster    string
	Status     bool
	AllCluster bool
}

// NewInitCmd vops init command.
func NewInitCmd(cfg string) *cobra.Command {
	var c *config.Config

	o := &initOptions{}

	cmd := &cobra.Command{
		Use:           "init",
		Aliases:       []string{"i"},
		Short:         "initialize a vault cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			fmt.Println("[ Initialize ]")
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
					if err := o.initializeCluster(cluster); err != nil {
						return err
					}
				}

				return nil
			}

			cluster, err := c.GetCluster(o.Cluster)
			if err != nil {
				return err
			}

			if err := o.initializeCluster(*cluster); err != nil {
				return err
			}

			return nil
		},
	}

	flags.AllClusterFlag(cmd, o.AllCluster)
	flags.ClusterFlag(cmd, o.Cluster)

	cmd.Flags().BoolVarP(&o.Status, "status", "S", o.Status, "print the initialization status of a cluster")
	cmd.Flags().IntVarP(&o.Shares, "shares", "s", o.Shares, "Number of keyshares")
	cmd.Flags().IntVarP(&o.Threshold, "threshold", "t", o.Threshold, "Number of required keys to unseal vault")

	return cmd
}

//nolint: cyclop
func (o *initOptions) initializeCluster(cluster config.Cluster) error {
	if o.Shares > 0 {
		cluster.Keys.Shares = o.Shares
	}

	if o.Threshold > 0 {
		cluster.Keys.Threshold = o.Threshold
	}

	fmt.Printf("\n[ %s ]\n", cluster.Name)
	fmt.Printf("attempting intialization of cluster \"%s\" with %d shares and a threshold of %d\n", cluster.Name, cluster.Keys.Shares, cluster.Keys.Threshold)

	if cluster.Keys.Path == "" {
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

	if o.Status || isInitialized {
		fmt.Printf("%s is already initialized.\n", cluster.Name)

		return nil
	}

	resp, err := v.Initialize(cluster.Keys.Shares, cluster.Keys.Threshold, false)
	if err != nil {
		return err
	}

	if err := fs.WriteToFile(utils.ToJSON(resp), cluster.Keys.Path); err != nil {
		return err
	}

	fmt.Printf("successfully initialized %s and wrote keys to %s.\n", cluster.Name, cluster.Keys.Path)

	return nil
}
