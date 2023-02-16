package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/FalcoSuessgott/vops/cmd/config"
	"github.com/FalcoSuessgott/vops/cmd/custom"
	"github.com/FalcoSuessgott/vops/cmd/generateroot"
	"github.com/FalcoSuessgott/vops/cmd/initialize"
	"github.com/FalcoSuessgott/vops/cmd/manpage"
	"github.com/FalcoSuessgott/vops/cmd/rekey"
	"github.com/FalcoSuessgott/vops/cmd/seal"
	"github.com/FalcoSuessgott/vops/cmd/snapshot"
	"github.com/FalcoSuessgott/vops/cmd/unseal"
	"github.com/FalcoSuessgott/vops/cmd/version"
	"github.com/spf13/cobra"
)

const cfgFileEnvVar = "VOPS_CONFIG"

var cfgFile = "vops.yaml"

// NewRootCmd vops root command.
func NewRootCmd(v string, writer io.Writer) *cobra.Command {
	if v, ok := os.LookupEnv(cfgFileEnvVar); ok {
		cfgFile = v
	}

	cmd := &cobra.Command{
		Use:           "vops",
		Short:         "A HashiCorp Vault cluster operations tool",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.PersistentFlags().StringVar(&cfgFile, "config", cfgFile, "path to config file")

	// sub commands
	cmd.AddCommand(
		initialize.NewInitCmd(cfgFile),
		unseal.NewUnsealCmd(cfgFile),
		seal.NewSealCmd(cfgFile),
		rekey.NewRekeyCmd(cfgFile),
		generateroot.NewGenerateRootCmd(cfgFile),
		version.NewVersionCmd(v),
		config.NewConfigCmd(cfgFile),
		snapshot.NewSnapshotCmd(cfgFile),
		custom.NewCustomCmd(cfgFile),
		config.NewConfigCmd(cfgFile),
		manpage.NewManCmd().Cmd,
	)

	return cmd
}

// Execute invokes the command.
func Execute(version string) error {
	if err := NewRootCmd(version, os.Stdout).Execute(); err != nil {
		return fmt.Errorf("[ERROR] %w", err)
	}

	return nil
}
