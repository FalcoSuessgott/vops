package unseal

import (
	"fmt"

	"github.com/FalcoSuessgott/vops/pkg/config"
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

func newDefaultUnsealOptions() *unsealOptions {
	return &unsealOptions{
		AllNodes: true,
	}
}

// NewUnsealCmd vops unseal command.
// nolint: gocognit, cyclop
func NewUnsealCmd(cfg string) *cobra.Command {
	o := newDefaultUnsealOptions()

	cmd := &cobra.Command{
		Use:           "unseal",
		Aliases:       []string{"u"},
		Short:         "unseal a single node or a single cluster or all cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return config.ValidateConfig(cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("[ Unseal ]")
			fmt.Printf("using %s\n", cfg)

			config, err := config.ParseConfig(cfg)
			if err != nil {
				return err
			}

			if o.Node != "" {
				o.AllNodes = false
			}

			// All Cluster
			if o.AllCluster {
				for _, cluster := range config.Cluster {
					if err := unsealCluster(cluster); err != nil {
						return err
					}
				}

				return nil
			}

			// All Nodes
			if o.AllNodes {
				cluster, err := config.GetCluster(o.Cluster)
				if err != nil {
					return err
				}

				if err := unsealCluster(*cluster); err != nil {
					return err
				}

				return nil
			}

			// Single Node
			cluster, err := config.GetCluster(o.Cluster)
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

	cmd.Flags().StringVarP(&o.Cluster, "cluster", "c", o.Cluster, "name of the vault cluster to initialize")
	cmd.Flags().StringVarP(&o.Node, "node", "n", o.Node, "unseal a single vault node")
	cmd.Flags().BoolVarP(&o.AllNodes, "all", "a", o.AllNodes, "unseal all nodes of a cluster")
	cmd.Flags().BoolVarP(&o.AllCluster, "all-cluster", "A", o.AllCluster, "unseal all cluster defined in the vops configuration file")

	return cmd
}

func unsealCluster(cluster config.Cluster) error {
	fmt.Printf("\n[ %s ]\n", cluster.Name)

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
