package snapshot

import (
	"bytes"
	"fmt"
	"path"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/flags"
	"github.com/FalcoSuessgott/vops/pkg/fs"
	"github.com/FalcoSuessgott/vops/pkg/utils"
	"github.com/FalcoSuessgott/vops/pkg/vault"
	"github.com/spf13/cobra"
)

type snapshotOptions struct {
	Cluster    string
	AllCluster bool
}

// NewSnapshotCmd snapshot command.
func NewSnapshotCmd(cfg string) *cobra.Command {
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
		newSnapSavewCmd(cfg),
		newSnapRestoreCmd(cfg),
	)

	return cmd
}

func newSnapSavewCmd(cfg string) *cobra.Command {
	var c *config.Config

	o := &snapshotOptions{}

	cmd := &cobra.Command{
		Use:           "save",
		Aliases:       []string{"s"},
		Short:         "save a snapshot of a vault cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			fmt.Println("[ Snapshot Save ]")
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
					if err := saveSnapshot(cluster); err != nil {
						return err
					}
				}

				return nil
			}

			cluster, err := c.GetCluster(o.Cluster)
			if err != nil {
				return err
			}

			return saveSnapshot(*cluster)
		},
	}

	flags.AllClusterFlag(cmd, o.AllCluster)
	flags.ClusterFlag(cmd, o.Cluster)

	return cmd
}

func newSnapRestoreCmd(cfg string) *cobra.Command {
	var c *config.Config

	o := &snapshotOptions{}

	cmd := &cobra.Command{
		Use:           "restore",
		Aliases:       []string{"r"},
		Short:         "restore a snapshot of a vault cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			fmt.Println("[ Snapshot Save ]")
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
					if err := restoreSnapshot(cluster); err != nil {
						return err
					}
				}

				return nil
			}

			cluster, err := c.GetCluster(o.Cluster)
			if err != nil {
				return err
			}

			return restoreSnapshot(*cluster)
		},
	}

	flags.AllClusterFlag(cmd, o.AllCluster)
	flags.ClusterFlag(cmd, o.Cluster)

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
