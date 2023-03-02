package cmd

import (
	"fmt"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/fs"
	"github.com/FalcoSuessgott/vops/pkg/utils"
	"github.com/FalcoSuessgott/vops/pkg/vault"
	"github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
)

func rekeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "rekey",
		Aliases:       []string{"rk"},
		Short:         "rekey a cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if allCluster {
				for _, cluster := range cfg.Cluster {
					if err := rekeyCluster(cluster); err != nil {
						return err
					}
				}

				return nil
			}

			cluster, err := cfg.GetCluster(cluster)
			if err != nil {
				return err
			}

			if err := rekeyCluster(*cluster); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func rekeyCluster(cluster config.Cluster) error {
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
