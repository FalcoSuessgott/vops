package cmd

import (
	"bytes"
	"fmt"
	"path"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/fs"
	"github.com/FalcoSuessgott/vops/pkg/utils"
	"github.com/FalcoSuessgott/vops/pkg/vault"
	"github.com/spf13/cobra"
)

func snapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "snapshot",
		Aliases:       []string{"ss"},
		Short:         "save or restore a snapshot of a vault cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		snapSavewCmd(),
		snapRestoreCmd(),
	)

	return cmd
}

func snapSavewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "save",
		Aliases:       []string{"s"},
		Short:         "save a snapshot of a vault cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if allCluster {
				for _, cluster := range cfg.Cluster {
					if err := saveSnapshot(cluster); err != nil {
						return err
					}
				}

				return nil
			}

			cluster, err := cfg.GetCluster(cluster)
			if err != nil {
				return err
			}

			return saveSnapshot(*cluster)
		},
	}

	return cmd
}

func snapRestoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "restore",
		Aliases:       []string{"r"},
		Short:         "restore a snapshot of a vault cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if allCluster {
				for _, cluster := range cfg.Cluster {
					if err := restoreSnapshot(cluster); err != nil {
						return err
					}
				}

				return nil
			}

			cluster, err := cfg.GetCluster(cluster)
			if err != nil {
				return err
			}

			return restoreSnapshot(*cluster)
		},
	}

	return cmd
}

func saveSnapshot(cluster config.Cluster) error {
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

	v, err := vault.NewTokenClient(cluster.Addr, cluster.Token)
	if err != nil {
		return err
	}

	w, err := v.SnapshotBackup()
	if err != nil {
		return err
	}

	fs.CreateDirIfNotExist(cluster.SnapshotDir)

	snapshotName := path.Join(cluster.SnapshotDir, utils.GetCurrentDate())

	if fs.WriteToFile(w.Bytes(), snapshotName) != nil {
		return err
	}

	fmt.Printf("created snapshot file \"%s\" for cluster \"%s\"\n", snapshotName, cluster.Name)

	return nil
}

func restoreSnapshot(cluster config.Cluster) error {
	fmt.Printf("\n[ %s ]\n", cluster.Name)

	if err := cluster.ApplyEnvironmentVariables(cluster.ExtraEnv); err != nil {
		return err
	}

	if err := cluster.RunTokenExecCommand(); err != nil {
		return err
	}

	fmt.Println("executed token exec command")

	v, err := vault.NewTokenClient(cluster.Addr, cluster.Token)
	if err != nil {
		return err
	}

	var b bytes.Reader

	if err := v.SnapshotRestore(&b, true); err != nil {
		return err
	}

	fmt.Printf("restrored snapshot for %s\n", cluster.Name)

	return nil
}
