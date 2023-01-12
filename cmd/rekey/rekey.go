package rekey

import (
	"fmt"

	"github.com/FalcoSuessgott/vops/pkg/config"
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

func newDefaultRekeyOptions() *rekeyOptions {
	return &rekeyOptions{}
}

// NewRekeyCmd vops rekey command.
func NewRekeyCmd(cfg string) *cobra.Command {
	o := newDefaultRekeyOptions()

	cmd := &cobra.Command{
		Use:           "rekey",
		Aliases:       []string{"rk"},
		Short:         "rekey a single or all vault cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("[ Rekeying ]")
			fmt.Printf("using %s\n", cfg)

			config, err := config.ParseConfig(cfg)
			if err != nil {
				return err
			}

			// AllCluster
			if o.AllCluster {
				for _, cluster := range config.Cluster {
					if err := o.rekeyCluster(cluster); err != nil {
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

			if err := o.rekeyCluster(*cluster); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&o.Cluster, "cluster", "c", o.Cluster, "name of the vault cluster to initialize")
	cmd.Flags().IntVarP(&o.Shares, "shares", "s", o.Shares, "Number of keyshares")
	cmd.Flags().IntVarP(&o.Threshold, "threshold", "t", o.Threshold, "Number of required keys to unseal vault")
	cmd.Flags().BoolVarP(&o.AllCluster, "all-cluster", "A", o.AllCluster, "unseal all cluster defined in the vops configuration file")

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

	if err := cluster.ApplyEnvironmentVariables(); err != nil {
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
