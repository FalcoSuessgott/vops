package generateroot

import (
	"fmt"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/flags"
	"github.com/FalcoSuessgott/vops/pkg/vault"
	"github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
)

type generateRootOptions struct {
	Cluster    string
	AllCluster bool
}

// NewGenerateRootCmd vops generate root command.
func NewGenerateRootCmd(cfg string) *cobra.Command {
	var c *config.Config

	o := &generateRootOptions{}

	cmd := &cobra.Command{
		Use:           "generate-root",
		Aliases:       []string{"gr"},
		Short:         "generate a new root token for a vault cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			fmt.Println("[ Generate Root Token ]")
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
					if err := generateRoot(cluster); err != nil {
						return err
					}
				}

				return nil
			}

			cluster, err := c.GetCluster(o.Cluster)
			if err != nil {
				return err
			}

			if err := generateRoot(*cluster); err != nil {
				return err
			}

			return nil
		},
	}

	flags.AllClusterFlag(cmd, o.AllCluster)
	flags.ClusterFlag(cmd, o.Cluster)

	return cmd
}

//nolint: cyclop
func generateRoot(cluster config.Cluster) error {
	var token *api.GenerateRootStatusResponse

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

	v, err := vault.NewClient(cluster.Addr)
	if err != nil {
		return err
	}

	otp, err := v.GenerateOTP()
	if err != nil {
		return api.ErrIncompleteSnapshot
	}

	fmt.Println("generated on OTP for root token creation")

	regenRoot, err := v.GenerateRootInit(otp)
	if err != nil {
		return err
	}

	fmt.Println("started root token generation process")

	for _, key := range keys.Keys {
		resp, err := v.GenerateRootUpdate(key, regenRoot.Nonce)
		if err != nil {
			return err
		}

		if resp.Complete {
			fmt.Println("root token generation completed")

			token = resp

			break
		}
	}

	rootToken, err := v.DecodeRootToken(token.EncodedRootToken, otp)
	if err != nil {
		return err
	}

	fmt.Printf(
		"new root token: \"%s\" "+
			"(make sure to update your token exec commands in your vops configfile if necessary.)\n", rootToken)

	return nil
}
