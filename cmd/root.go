package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/utils"
	"github.com/spf13/cobra"
)

const cfgFileEnvVar = "VOPS_CONFIG"

var (
	cfgFile    = "vops.yaml"
	cfg        *config.Config
	cluster    string
	allCluster bool
)

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
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			fmt.Println(utils.PrintHeader(cmd.Use, cfgFile))

			cfg, err = config.ParseConfig(cfgFile)
			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.PersistentFlags().StringVar(&cfgFile, "config", cfgFile, "path to the vops configfile")
	cmd.PersistentFlags().StringVarP(&cluster, "cluster", "c", cluster, "name of the vault cluster")
	cmd.PersistentFlags().BoolVarP(&allCluster, "all-cluster", "A", allCluster, "perform action for all cluster defined in the vops configuration file")

	cmd.AddCommand(
		initCmd(),
		unsealCmd(),
		sealCmd(),
		rekeyCmd(),
		adhocCmd(),
		generateRootCmd(),
		versionCmd(v),
		snapshotCmd(),
		customCmd(),
		configCmd(writer),
		manCmd().Cmd,
		uiCmd(),
		loginCmd(),
		tokenCmd(),
		completionCmd(),
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
