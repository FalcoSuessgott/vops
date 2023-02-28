package ui

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/flags"
	"github.com/spf13/cobra"
)

type uiOptions struct {
	Cluster string
}

// NewUICmd ui command.
func NewUICmd(cfg string) *cobra.Command {
	var c *config.Config

	o := &uiOptions{}

	cmd := &cobra.Command{
		Use:           "ui",
		Short:         "open the UI of the specified vault cluster address in your browser",
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
			cluster, err := c.GetCluster(o.Cluster)
			if err != nil {
				return err
			}

			fmt.Printf("\n[ %s ]\n", cluster.Name)
			fmt.Printf("opening %s\n", cluster.Addr)

			//nolint: gosec
			switch runtime.GOOS {
			case "linux":
				err = exec.Command("xdg-open", cluster.Addr).Start()
			case "windows":
				err = exec.Command("rundll32", "url.dll,FileProtocolHandler", cluster.Addr).Start()
			case "darwin":
				err = exec.Command("open", cluster.Addr).Start()
			default:
				err = fmt.Errorf("unsupported platform")
			}

			if err != nil {
				return err
			}

			return nil
		},
	}

	flags.ClusterFlag(cmd, o.Cluster)

	return cmd
}
