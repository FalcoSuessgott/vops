package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

func uiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "ui",
		Short:         "open the UI of the specified vault cluster address in your browser",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cluster, err := cfg.GetCluster(cluster)
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
