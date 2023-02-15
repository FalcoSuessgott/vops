package generateroot

import (
	"fmt"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/vault"
	"github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
)

type generateRootOptions struct {
	Cluster    string
	AllCluster bool
}

func newDefaultGenerateRootOptions() *generateRootOptions {
	return &generateRootOptions{}
}

// NewGenerateRootCmd vops rekey command.
func NewGenerateRootCmd(cfg string) *cobra.Command {
	o := newDefaultGenerateRootOptions()

	cmd := &cobra.Command{
		Use:           "generate-root",
		Aliases:       []string{"gr"},
		Short:         "generates a new root token for a single or all cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return config.ValidateConfig(cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("[ Generate Root Token ]")
			fmt.Printf("using %s\n", cfg)

			config, err := config.ParseConfig(cfg)
			if err != nil {
				return err
			}

			if o.AllCluster {
				for _, cluster := range config.Cluster {
					if err := generateRoot(cluster); err != nil {
						return err
					}
				}

				return nil
			}

			cluster, err := config.GetCluster(o.Cluster)
			if err != nil {
				return err
			}

			if err := generateRoot(*cluster); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&o.Cluster, "cluster", "c", o.Cluster, "name of the vault cluster to initialize")
	cmd.Flags().BoolVarP(&o.AllCluster, "all-cluster", "A", o.AllCluster, "initialize all cluster defined in the vops configuration file")

	return cmd
}

func generateRoot(cluster config.Cluster) error {
	var token *api.GenerateRootStatusResponse

	fmt.Printf("\n[ %s ]\n", cluster.Name)

	if err := cluster.ApplyEnvironmentVariables(); err != nil {
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
