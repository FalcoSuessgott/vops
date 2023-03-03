package cmd

import (
	"fmt"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/vault"
	"github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
)

func generateRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "generate-root",
		Aliases:       []string{"gr"},
		Short:         "generate a new root token for a vault cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if allCluster {
				for _, cluster := range cfg.Cluster {
					if err := generateRoot(cluster); err != nil {
						return err
					}
				}

				return nil
			}

			cluster, err := cfg.GetCluster(cluster)
			if err != nil {
				return err
			}

			if err := generateRoot(*cluster); err != nil {
				return err
			}

			return nil
		},
	}

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

	keyFile, err := cluster.GetKeyFile()
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

	keys := keyFile.Keys

	if cluster.Keys.Autounseal {
		keys = keyFile.RecoveryKeys
	}

	for _, key := range keys {
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
