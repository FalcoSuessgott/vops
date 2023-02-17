package ui

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/spf13/cobra"
)

type uiOptions struct {
	Cluster string
}

func newDefaultUIptions() *uiOptions {
	return &uiOptions{}
}

// NewUICmd ui command.
func NewUICmd(cfg string) *cobra.Command {
	o := newDefaultUIptions()

	cmd := &cobra.Command{
		Use:           "ui",
		Short:         "open the specified vault cluster address in your browser",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return config.ValidateConfig(cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("[ UI ]")
			fmt.Printf("using %s\n", cfg)

			config, err := config.ParseConfig(cfg)
			if err != nil {
				return err
			}

			cluster, err := config.GetCluster(o.Cluster)
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

	return cmd
}
