package cmd

import (
	"fmt"
	"strings"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/exec"
	"github.com/spf13/cobra"
)

type customOptions struct {
	Command string
	List    bool
}

func customCmd() *cobra.Command {
	o := &customOptions{}

	cmd := &cobra.Command{
		Use:           "custom",
		Aliases:       []string{"c"},
		Short:         "run any custom command for a vault cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.List {
				fmt.Println("\n[ Available Commands ]")
				for name, cmd := range cfg.CustomCmds {
					fmt.Printf("\"%s\": \"%s\"\n", name, cmd)
				}

				fmt.Println("\nrun any available command with \"vops custom -x \"<command name>\" -c \"<cluster-name>\".")

				return nil
			}

			if allCluster {
				for _, cluster := range cfg.Cluster {
					if err := o.runCustomCommand(cluster, cfg.CustomCmds); err != nil {
						return err
					}
				}

				return nil
			}

			cluster, err := cfg.GetCluster(cluster)
			if err != nil {
				return err
			}

			if err := o.runCustomCommand(*cluster, cfg.CustomCmds); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&o.List, "list", "l", o.List, "list all available custom commands")

	return cmd
}

func (o *customOptions) runCustomCommand(cluster config.Cluster, cmds map[string]interface{}) error {
	cmd, ok := cmds[o.Command]
	if !ok {
		return fmt.Errorf("invalid command")
	}

	//nolint: forcetypeassert
	parts := strings.Split(cmd.(string), " ")

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
	fmt.Printf("$> %s", cmd)
	fmt.Println()

	out, err := exec.Run(parts)
	if err != nil {
		return err
	}

	fmt.Println(string(out))

	return nil
}
