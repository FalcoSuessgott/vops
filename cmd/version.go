package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func versionCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "version",
		Short:         "print vops version",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "vops %s\n", version)

			return nil
		},
	}

	return cmd
}
