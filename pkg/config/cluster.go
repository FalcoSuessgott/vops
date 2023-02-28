package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/FalcoSuessgott/vops/pkg/exec"
	"github.com/FalcoSuessgott/vops/pkg/fs"
	"github.com/FalcoSuessgott/vops/pkg/template"
	"github.com/FalcoSuessgott/vops/pkg/utils"
	"github.com/FalcoSuessgott/vops/pkg/vault"
	"github.com/hashicorp/vault/api"
)

// Cluster struct of a single vops vault cluster.
type Cluster struct {
	Name         string                 `json:"Name" yaml:"Name,omitempty"`
	Addr         string                 `json:"Addr" yaml:"Addr,omitempty"`
	Token        string                 `json:"-" yaml:"-,omitempty"`
	TokenExecCmd string                 `json:"TokenExecCmd" yaml:"TokenExecCmd,omitempty"`
	Keys         *KeyConfig             `json:"Keys" yaml:"Keys"`
	SnapshotDir  string                 `json:"SnapshotDirectory,omitempty" yaml:"SnapshotDirectory,omitempty"`
	Nodes        []string               `json:"Nodes,omitempty" yaml:"Nodes,omitempty"`
	Env          map[string]interface{} `json:"Env" yaml:"Env,omitempty"`
	ExtraEnv     map[string]interface{} `json:"ExtraEnv,omitempty" yaml:"ExtraEnv,omitempty"`
}

// KeyConfig keyconfig parameters.
type KeyConfig struct {
	Path      string `json:"Path,omitempty" yaml:"Path,omitempty"`
	Shares    int    `json:"Shares,omitempty" yaml:"Shares,omitempty"`
	Threshold int    `json:"Threshold,omitempty" yaml:"Threshold,omitempty"`
}

// RunTokenExecCommand executes the token command.
func (c *Cluster) RunTokenExecCommand() error {
	out, err := exec.Run(strings.Split(c.TokenExecCmd, " "))
	if err != nil {
		return fmt.Errorf("error while executing token command: %w", err)
	}

	c.Token = strings.TrimSuffix(string(out), "\n")

	return nil
}

// GetKeyFile reads the defined keyfile and returns it.
func (c *Cluster) GetKeyFile() (*api.InitResponse, error) {
	resp := &api.InitResponse{}

	keyfile := fs.ReadFile(c.Keys.Path)

	if err := utils.FromJSON(keyfile, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// ApplyEnvironmentVariables applies the environment variables specfied for a single vault.
func (c *Cluster) ApplyEnvironmentVariables(envs map[string]interface{}) error {
	for k, v := range envs {
		if err := os.Setenv(k, fmt.Sprintf("%v", v)); err != nil {
			return err
		}

		fmt.Printf("applying %s\n", k)
	}

	return nil
}

// RenderConfig renders the config until all templates are replaced.
func (c *Cluster) RenderConfig() (*Cluster, error) {
	d := utils.ToJSON(c)
	m := map[string]interface{}{}

	if err := json.Unmarshal(d, &m); err != nil {
		return nil, err
	}

	data, err := template.Render(d, m)
	if err != nil {
		return nil, err
	}

	renderedConfig := &Cluster{}

	if err := utils.FromJSON(data.Bytes(), renderedConfig); err != nil {
		return nil, fmt.Errorf("cannot render values to vops config")
	}

	return renderedConfig, nil
}

func (c Cluster) String() string {
	policies := []string{}

	if err := c.RunTokenExecCommand(); err == nil {
		client, err := vault.NewTokenClient(c.Addr, c.Token)
		if err == nil {
			pols, err := client.TokenLookup()
			if err == nil {
				policies = append(policies, fmt.Sprintf("%v", pols))
			}
		}
	}

	return fmt.Sprintf(
		"Name:\t%s\n"+
			"Address:\t%s\n"+
			"TokenExecCmd:\t%s\n"+
			"TokenExecCmd Policies:\t%s\n"+
			"Nodes:\t[%s]\n"+
			"Key Config:\t{Path: %s, Shares: %d, Threshold: %d}\n"+
			"Snapshot Directory:\t%s\n",
		c.Name,
		c.Addr,
		c.TokenExecCmd,
		strings.Join(policies, ","),
		strings.Join(c.Nodes, ","),
		c.Keys.Path,
		c.Keys.Shares,
		c.Keys.Threshold,
		c.SnapshotDir,
	)
}
