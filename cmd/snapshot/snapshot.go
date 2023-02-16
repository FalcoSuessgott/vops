package snapshot

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

type snapshotOptions struct {
	Cluster    string
	Backup     bool
	Restore    bool
	AllCluster bool
}

func newDefaultSnapshotOptions() *snapshotOptions {
	return &snapshotOptions{
		Backup: true,
	}
}

// NewSnapshotCmd snapshot command.
func NewSnapshotCmd(cfg string) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "snapshot",
		Aliases:       []string{"ss"},
		Short:         "creates or restorees a snapshot from a single or all vault cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return config.ValidateConfig(cfg)
		},
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
	o := newDefaultSnapshotOptions()

	cmd := &cobra.Command{
		Use:           "save",
		Aliases:       []string{"s"},
		Short:         "saves a snapshot of a single or all vault cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("[ Snapshot Save ]")
			fmt.Printf("using %s\n", cfg)

			config, err := config.ParseConfig(cfg)
			if err != nil {
				return err
			}

			if o.AllCluster {
				for _, cluster := range config.Cluster {
					if err := saveSnapshot(cluster); err != nil {
						return err
					}
				}

				return nil
			}

			cluster, err := config.GetCluster(o.Cluster)
			if err != nil {
				return err
			}

			return saveSnapshot(*cluster)
		},
	}

	cmd.Flags().StringVarP(&o.Cluster, "cluster", "c", o.Cluster, "name of the vault cluster to initialize")
	cmd.Flags().BoolVarP(&o.AllCluster, "all-cluster", "A", o.AllCluster, "unseal all cluster defined in the vops configuration file")

	return cmd
}

func newSnapRestoreCmd(cfg string) *cobra.Command {
	o := newDefaultSnapshotOptions()

	cmd := &cobra.Command{
		Use:           "restore",
		Aliases:       []string{"r"},
		Short:         "restore a snapshot of a single or all vault cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("[ Snapshot Restore ]")
			fmt.Printf("using %s\n", cfg)

			config, err := config.ParseConfig(cfg)
			if err != nil {
				return err
			}

			if o.AllCluster {
				for _, cluster := range config.Cluster {
					if err := restoreSnapshot(cluster); err != nil {
						return err
					}
				}

				return nil
			}

			cluster, err := config.GetCluster(o.Cluster)
			if err != nil {
				return err
			}

			return restoreSnapshot(*cluster)
		},
	}

	cmd.Flags().StringVarP(&o.Cluster, "cluster", "c", o.Cluster, "name of the vault cluster to initialize")
	cmd.Flags().BoolVarP(&o.AllCluster, "all-cluster", "A", o.AllCluster, "unseal all cluster defined in the vops configuration file")

	return cmd
}

func saveSnapshot(cluster config.Cluster) error {
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
