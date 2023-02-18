package rekey

import (
	"fmt"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/flags"
	"github.com/FalcoSuessgott/vops/pkg/fs"
	"github.com/FalcoSuessgott/vops/pkg/utils"
	"github.com/FalcoSuessgott/vops/pkg/vault"
	"github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
)

type rekeyOptions struct {
	Cluster    string
	Shares     int
	Threshold  int
	AllCluster bool
}

// NewRekeyCmd vops rekey command.
func NewRekeyCmd(cfg string) *cobra.Command {
	var c *config.Config

	o := &rekeyOptions{}

	cmd := &cobra.Command{
		Use:           "rekey",
		Aliases:       []string{"rk"},
		Short:         "rekey a cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			fmt.Println("[ Rekey ]")
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
					if err := o.rekeyCluster(cluster); err != nil {
						return err
					}
				}

				return nil
			}

			cluster, err := c.GetCluster(o.Cluster)
			if err != nil {
				return err
			}

			if err := o.rekeyCluster(*cluster); err != nil {
				return err
			}

			return nil
		},
	}

	flags.AllClusterFlag(cmd, o.AllCluster)
	flags.ClusterFlag(cmd, o.Cluster)

	cmd.Flags().IntVarP(&o.Shares, "shares", "s", o.Shares, "Number of keyshares")
	cmd.Flags().IntVarP(&o.Threshold, "threshold", "t", o.Threshold, "Number of required keys to unseal vault")

	return cmd
}

//nolint: cyclop
func (o *rekeyOptions) rekeyCluster(cluster config.Cluster) error {
	if o.Shares > 0 {
		cluster.Keys.Shares = o.Shares
	}

	if o.Threshold > 0 {
		cluster.Keys.Threshold = o.Threshold
	}

	fmt.Printf("\n[ %s ]\n", cluster.Name)
	fmt.Printf("performing a rekey for %s with %d shares and a threshold of %d\n", cluster.Name, cluster.Keys.Shares, cluster.Keys.Threshold)

	if cluster.Keys.Path == "" {
		return fmt.Errorf("a key file containing unseal/recovery keys for that cluster is required")
	}

	if err := cluster.ApplyEnvironmentVariables(cluster.ExtraEnv); err != nil {
		return err
	}

	keys, err := cluster.GetKeyFile()
	if err != nil {
		return err
	}

	fmt.Printf("using keyfile \"%s\"\n", cluster.Keys.Path)

	v, err := vault.NewClient(cluster.Addr)
	if err != nil {
		return err
	}

	rekeyInit, err := v.RekeyInit(cluster.Keys.Shares, cluster.Keys.Threshold, false)
	if err != nil {
		return err
	}

	var newKeys *api.RekeyUpdateResponse

	for _, key := range keys.Keys {
		resp, err := v.RekeyUpdate(key, rekeyInit.Nonce)
		if err != nil {
			return err
		}

		if resp.Complete {
			fmt.Println("rekeying successfully completed")

			newKeys = resp

			break
		}
	}

	newName := fmt.Sprintf("%s_%s", cluster.Keys.Path, utils.GetCurrentDate())

	fs.RenameFile(cluster.Keys.Path, newName)

	fmt.Printf(
		"renamed keyfile \"%s\" for cluster \"%s\" to \"%s\""+
			"(snapshots depend on the unseal/recovery keys from the moment the snapshot has been created. "+
			"This way you always have the matching unseal/recovery keys ready.\n",
		cluster.Keys.Path, cluster.Name, newName,
	)

	if err := fs.WriteToFile(utils.ToJSON(newKeys), cluster.Keys.Path); err != nil {
		return err
	}

	return nil
}
