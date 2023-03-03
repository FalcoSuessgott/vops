package cmd

import (
	"fmt"
	"strings"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/exec"
	"github.com/spf13/cobra"
)

var adhocCommand string

func adhocCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "adhoc",
		Short:         "run any command",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if adhocCommand == "" {
				return fmt.Errorf("no command specified. Use --command flag")
			}

			if allCluster {
				for _, cluster := range cfg.Cluster {
					if err := runAdhocCommand(cluster); err != nil {
						return err
					}
				}

				return nil
			}

			cluster, err := cfg.GetCluster(cluster)
			if err != nil {
				return err
			}

			if err := runAdhocCommand(*cluster); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&adhocCommand, "command", "x", adhocCommand, "the command to run")

	return cmd
}

func runAdhocCommand(cluster config.Cluster) error {
	parts := strings.Split(adhocCommand, " ")

	fmt.Printf("\n[ %s ]\n", cluster.Name)

	cluster.ExtraEnv["VAULT_ADDR"] = cluster.Addr
	cluster.ExtraEnv["VAULT_TOKEN"] = cluster.Token

	if err := cluster.ApplyEnvironmentVariables(cluster.ExtraEnv); err != nil {
		return err
	}

	if err := cluster.RunTokenExecCommand(); err != nil {
		return err
	}

	fmt.Println("token exec command successful")
	fmt.Println()
	fmt.Printf("$> %s", adhocCommand)
	fmt.Println()

	out, err := exec.Run(parts)
	if err != nil {
		return err
	}

	fmt.Println(string(out))

	return nil
}
