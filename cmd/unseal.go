package cmd

import (
	"fmt"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/vault"
	"github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
)

type unsealOptions struct {
	Node     string
	AllNodes bool
}

// nolint: gocognit, cyclop
func unsealCmd() *cobra.Command {
	o := &unsealOptions{
		AllNodes: true,
	}

	cmd := &cobra.Command{
		Use:           "unseal",
		Aliases:       []string{"u"},
		Short:         "unseal a single node or a all nodes of a vault cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if allCluster {
				for _, cluster := range cfg.Cluster {
					if err := unsealCluster(cluster); err != nil {
						return err
					}
				}

				return nil
			}

			if o.AllNodes {
				cluster, err := cfg.GetCluster(cluster)
				if err != nil {
					return err
				}

				if err := unsealCluster(*cluster); err != nil {
					return err
				}

				return nil
			}

			cluster, err := cfg.GetCluster(cluster)
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
