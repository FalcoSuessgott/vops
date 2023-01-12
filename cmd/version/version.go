package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCmd version subcommand.
func NewVersionCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "version",
		Short:         "print vops version",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "vops %s\n", version)

			return nil
		},
	}

	return cmd
}
