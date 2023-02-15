package initialize

import (
	"fmt"

	"github.com/FalcoSuessgott/vops/pkg/config"
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

func newDefaultInitOptions() *initOptions {
	return &initOptions{}
}

// NewInitCmd vops init command.
func NewInitCmd(cfg string) *cobra.Command {
	o := newDefaultInitOptions()

	cmd := &cobra.Command{
		Use:           "init",
		Aliases:       []string{"i"},
		Short:         "initialize a vault cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("[ Intialization ]")
			fmt.Printf("using %s\n", cfg)

			config, err := config.ParseConfig(cfg)
			if err != nil {
				return err
			}

			if o.AllCluster {
				for _, cluster := range config.Cluster {
					if err := o.initializeCluster(cluster); err != nil {
						return err
					}
				}

				return nil
			}

			cluster, err := config.GetCluster(o.Cluster)
			if err != nil {
				return err
			}

			if err := o.initializeCluster(*cluster); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&o.Cluster, "cluster", "c", o.Cluster, "name of a cluster specified in the vops configuration file")
	cmd.Flags().BoolVarP(&o.Status, "status", "S", o.Status, "print the initialization status of a cluster")
	cmd.Flags().IntVarP(&o.Shares, "shares", "s", o.Shares, "Number of keyshares")
	cmd.Flags().IntVarP(&o.Threshold, "threshold", "t", o.Threshold, "Number of required keys to unseal vault")
	cmd.Flags().BoolVarP(&o.AllCluster, "all-cluster", "A", o.AllCluster, "initialize all cluster defined in the vops configuration file")

	return cmd
}

func (o *initOptions) initializeCluster(cluster config.Cluster) error {
	if o.Shares > 0 {
		cluster.Keys.Shares = o.Shares
	}

	if o.Threshold > 0 {
		cluster.Keys.Threshold = o.Threshold
	}

	fmt.Printf("\n[ %s ]\n", cluster.Name)
	fmt.Printf("attempting intialization of cluster \"%s\" with %d shares and a threshold of %d\n", cluster.Name, o.Shares, o.Threshold)

	if err := cluster.ApplyEnvironmentVariables(); err != nil {
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
