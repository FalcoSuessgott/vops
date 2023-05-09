package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

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

// nolint: cyclop
func rekeyCluster(cluster config.Cluster) error {
	fmt.Printf("\n[ %s ]\n", cluster.Name)
	fmt.Printf("performing a rekey for %s with %d shares and a threshold of %d\n", cluster.Name, cluster.Keys.Shares, cluster.Keys.Threshold)

	if err := cluster.ApplyEnvironmentVariables(cluster.ExtraEnv); err != nil {
		return err
	}

	if cluster.Keys == nil || cluster.Keys.Path == "" {
		return fmt.Errorf("a key file containing unseal/recovery keys for that cluster is required")
	}

	keyFile, err := cluster.GetKeyFile()
	if err != nil {
		return err
	}

	fmt.Printf("using keyfile \"%s\"\n", cluster.Keys.Path)

	v, err := vault.NewClient(cluster.Addr)
	if err != nil {
		return err
	}

	rekeyInit, err := v.RekeyInit(cluster.Keys.Shares, cluster.Keys.Threshold, cluster.Keys.Autounseal)
	if err != nil {
		return err
	}

	fmt.Println("initialized rekey process")

	var newKeys *api.RekeyUpdateResponse

	keys := keyFile.Keys

	if cluster.Keys.Autounseal {
		keys = keyFile.RecoveryKeys
	}

	for i, key := range keys {
		resp, err := v.RekeyUpdate(key, rekeyInit.Nonce, cluster.Keys.Autounseal)
		if err != nil {
			return err
		}

		fmt.Printf("[%02d/%02d] successfully entered key\n", i+1, cluster.Keys.Threshold)

		if resp.Complete {
			fmt.Println("rekeying successfully completed")

			newKeys = resp

			break
		}
	}

	fileName := strings.TrimSuffix(cluster.Keys.Path, filepath.Ext(cluster.Keys.Path))
	newName := fmt.Sprintf("%s_%s%s", fileName, utils.GetCurrentDate(), filepath.Ext(cluster.Keys.Path))

	fs.RenameFile(cluster.Keys.Path, newName)

	fmt.Printf(
		"renamed keyfile \"%s\" for cluster \"%s\" to \"%s\".\n"+
			"Hint: snapshots depend on the unseal/recovery keys from the moment the snapshot has been created.\n"+
			"This way you always have the matching unseal/recovery keys for the specific snapshot if needed ready.\n",
		cluster.Keys.Path, cluster.Name, newName,
	)

	newKeyfile := &api.InitResponse{
		RootToken: keyFile.RootToken,
	}

	if cluster.Keys.Autounseal {
		newKeyfile.RecoveryKeys = newKeys.Keys
		newKeyfile.RecoveryKeysB64 = newKeys.KeysB64
	} else {
		newKeyfile.Keys = newKeys.Keys
		newKeyfile.KeysB64 = newKeys.KeysB64
	}

	if err := fs.WriteToFile(utils.ToJSON(newKeyfile), cluster.Keys.Path); err != nil {
		return err
	}

	return nil
}
