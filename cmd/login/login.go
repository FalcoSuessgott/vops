package login

import (
	"fmt"
	"strings"

	"github.com/FalcoSuessgott/vops/pkg/config"
	"github.com/FalcoSuessgott/vops/pkg/exec"
	"github.com/spf13/cobra"
)

type loginOptions struct {
	Cluster string
}

func newDefaultLoginOptions() *loginOptions {
	return &loginOptions{}
}

// NewLoginCmd login command.
func NewLoginCmd(cfg string) *cobra.Command {
	o := newDefaultLoginOptions()

	cmd := &cobra.Command{
		Use:           "login",
		Short:         "perform a vault login command for the specified cluster",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return config.ValidateConfig(cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("[ Login ]")
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
			fmt.Printf("performing a vault login to %s\n", cluster.Name)

			cluster.ExtraEnv["VAULT_ADDR"] = cluster.Addr
			cluster.ExtraEnv["VAULT_TOKEN"] = cluster.Token

			if err := cluster.ApplyEnvironmentVariables(cluster.ExtraEnv); err != nil {
				return err
			}

			if err := cluster.RunTokenExecCommand(); err != nil {
				return err
			}

			loginCmd := fmt.Sprintf("vault login %s", cluster.Token)

			fmt.Println("executed token exec command")
			fmt.Println()
			fmt.Printf("$> vault login $(%s)", cluster.TokenExecCmd)
			fmt.Println()

			out, err := exec.Run(strings.Split(loginCmd, " "))
			if err != nil {
				return err
			}

			fmt.Println(string(out))

			return nil
		},
	}

	return cmd
}
