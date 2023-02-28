package unseal

import (
	"fmt"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/flags"
	"github.com/FalcoSuessgott/vops/pkg/vault"
	"github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
)

type unsealOptions struct {
	Cluster    string
	Node       string
	AllNodes   bool
	AllCluster bool
}

// NewUnsealCmd vops unseal command.
// nolint: gocognit, cyclop
func NewUnsealCmd(cfg string) *cobra.Command {
	var c *config.Config

	o := &unsealOptions{
		AllNodes: true,
	}

	cmd := &cobra.Command{
		Use:           "unseal",
		Aliases:       []string{"u"},
		Short:         "unseal a single node or a all nodes of a vault cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			fmt.Println("[ Unseal ]")
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
					if err := unsealCluster(cluster); err != nil {
						return err
					}
				}

				return nil
			}

			if o.AllNodes {
				cluster, err := c.GetCluster(o.Cluster)
				if err != nil {
					return err
				}

				if err := unsealCluster(*cluster); err != nil {
					return err
				}

				return nil
			}

			cluster, err := c.GetCluster(o.Cluster)
			if err != nil {
				return err
			}

			fmt.Printf("\n[ %s ]\n", cluster.Name)

			if err := cluster.ApplyEnvironmentVariables(cluster.ExtraEnv); err != nil {
				return err
			}

			keys, err := cluster.GetKeyFile()
			if err != nil {
				return err
			}

			fmt.Printf("using keyfile \"%s\"\n", cluster.Keys.Path)

			for _, n := range cluster.Nodes {
				if n == o.Node {
					if err := unsealNode(n, keys); err != nil {
						return err
					}
				}
			}

			return fmt.Errorf("invalid node \"%s\" for cluster \"%s\"", o.Node, cluster.Name)
		},
	}

	flags.AllClusterFlag(cmd, o.AllCluster)
	flags.ClusterFlag(cmd, o.Cluster)

	cmd.Flags().StringVarP(&o.Node, "node", "n", o.Node, "unseal a single vault node")
	cmd.Flags().BoolVarP(&o.AllNodes, "all", "a", o.AllNodes, "unseal all nodes of a cluster")

	return cmd
}

func unsealCluster(cluster config.Cluster) error {
	fmt.Printf("\n[ %s ]\n", cluster.Name)

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

	for _, node := range cluster.Nodes {
		if err := unsealNode(node, keys); err != nil {
			return err
		}
	}

	fmt.Printf("cluster \"%s\" unsealed\n", cluster.Name)

	return nil
}

func unsealNode(node string, keys *api.InitResponse) error {
	client, err := vault.NewClient(node)
	if err != nil {
		return err
	}

	for _, k := range keys.Keys {
		isUnsealed, err := client.Unseal(k)
		if err != nil {
			return err
		}

		fmt.Printf("unsealing node \"%s\"\n", node)

		if isUnsealed {
			break
		}
	}

	return nil
}
