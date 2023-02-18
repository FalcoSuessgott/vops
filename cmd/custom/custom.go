package custom

import (
	"fmt"
	"strings"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/exec"
	"github.com/FalcoSuessgott/vops/pkg/flags"
	"github.com/spf13/cobra"
)

type customOptions struct {
	Command    string
	Cluster    string
	AllCluster bool
	List       bool
}

// NewCustomCmd vops custom command.
func NewCustomCmd(cfg string) *cobra.Command {
	var c *config.Config

	o := &customOptions{}

	cmd := &cobra.Command{
		Use:           "custom",
		Aliases:       []string{"c"},
		Short:         "run any custom command for a vault cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			fmt.Println("[ Custom Command ]")
			fmt.Printf("using %s\n", cfg)

			c, err = config.ParseConfig(cfg)
			if err != nil {
				return err
			}

			if len(c.CustomCmds) == 0 {
				return fmt.Errorf("a least one custom command is required")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if o.List {
				fmt.Println("\n[ Available Commands ]")
				for name, cmd := range c.CustomCmds {
					fmt.Printf("\"%s\": \"%s\"\n", name, cmd)
				}

				fmt.Println("\nrun any available command with \"vops custom -x \"<command name>\" -c \"<cluster-name>\".")

				return nil
			}

			if o.AllCluster {
				for _, cluster := range c.Cluster {
					if err := o.runCustomCommand(cluster, c.CustomCmds); err != nil {
						return err
					}
				}

				return nil
			}

			cluster, err := c.GetCluster(o.Cluster)
			if err != nil {
				return err
			}

			if err := o.runCustomCommand(*cluster, c.CustomCmds); err != nil {
				return err
			}

			return nil
		},
	}

	flags.AllClusterFlag(cmd, o.AllCluster)
	flags.ClusterFlag(cmd, o.Cluster)

	cmd.Flags().BoolVarP(&o.List, "list", "l", o.List, "list all available custom commands")
	cmd.Flags().StringVarP(&o.Command, "command", "x", o.Cluster, "the name of command to run")

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
