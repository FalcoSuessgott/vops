package cmd

import (
	"fmt"
	"os"

	mcoral "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
)

// ManCmd manpage command.
type ManCmd struct {
	Cmd *cobra.Command
}

func manCmd() *ManCmd {
	root := &ManCmd{}

	c := &cobra.Command{
		Use:                   "man",
		Short:                 "Generates GoReleaser's command line manpages",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Hidden:                true,
		Args:                  cobra.NoArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			manPage, err := mcoral.NewManPage(1, root.Cmd.Root())
			if err != nil {
				return err
			}

			_, err = fmt.Fprint(os.Stdout, manPage.Build(roff.NewDocument()))

			return err
		},
	}

	root.Cmd = c

	return root
}
