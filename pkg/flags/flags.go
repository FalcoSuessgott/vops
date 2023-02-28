package flags

import (
	"github.com/spf13/cobra"
)

// AllClusterFlag returns the cobra flag for all clusters.
func AllClusterFlag(cmd *cobra.Command, flag bool) {
	cmd.Flags().BoolVarP(&flag, "all-cluster", "A", flag, "perform action for all cluster defined in the vops configuration file")
}

// ClusterFlag returns the cobra flag for a single cluster.
func ClusterFlag(cmd *cobra.Command, flag string) {
	cmd.Flags().StringVarP(&flag, "cluster", "c", flag, "name of the vault cluster")
}
