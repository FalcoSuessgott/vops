package custom

import (
	"fmt"
	"strings"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/exec"
	"github.com/spf13/cobra"
)

type customOptions struct {
	Command    string
	Cluster    string
	AllCluster bool
	List       bool
}

func newCustomOptions() *customOptions {
	return &customOptions{}
}

// NewCustomCmd vops init command.
func NewCustomCmd(cfg string) *cobra.Command {
	o := newCustomOptions()

	cmd := &cobra.Command{
		Use:           "custom",
		Aliases:       []string{"c"},
		Short:         "run a custom command for a single or al vault cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return config.ValidateConfig(cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("[ Custom ]")
			fmt.Printf("using %s\n", cfg)

			config, err := config.ParseConfig(cfg)
			if err != nil {
				return err
			}

			if o.List {
				fmt.Println("\n[ Available Commands ]")
				for name, cmd := range config.CustomCmds {
					fmt.Printf("\"%s\": \"%s\"\n", name, cmd)
				}

				fmt.Println("\nrun any available command with \"vops custom -x \"<command name>\" -c \"<cluster-name>\".")

				return nil
			}

			if o.AllCluster {
				for _, cluster := range config.Cluster {
					if err := o.runCustomCommand(cluster, config.CustomCmds); err != nil {
						return err
					}
				}

				return nil
			}

			cluster, err := config.GetCluster(o.Cluster)
			if err != nil {
				return err
			}

			if err := o.runCustomCommand(*cluster, config.CustomCmds); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&o.List, "list", "l", o.List, "list available custom commands")
	cmd.Flags().StringVarP(&o.Command, "command", "x", o.Cluster, "the name of command defined in vops.yaml to run")
	cmd.Flags().StringVarP(&o.Cluster, "cluster", "c", o.Cluster, "name of a cluster specified in the vops configuration file")
	cmd.Flags().BoolVarP(&o.AllCluster, "all-cluster", "A", o.AllCluster, "initialize all cluster defined in the vops configuration file")

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
