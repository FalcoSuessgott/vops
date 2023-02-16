package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/FalcoSuessgott/vops/pkg/fs"
	"github.com/FalcoSuessgott/vops/pkg/utils"
)

const (
	clusterEnvVar = "VOPS_CLUSTER"
)

var (
	errNoClusterDefined     = errors.New("no cluster defined")
	errNoSuchClusterDefined = errors.New("no such cluster defined")
)

// Config holds the config file parameters.
type Config struct {
	Cluster    []Cluster              `json:"Cluster" yaml:"Cluster,omitempty"`
	CustomCmds map[string]interface{} `json:"CustomCmds" yaml:"CustomCmds,omitempty"`
}

// ParseConfig reads and parses a vops config file.
func ParseConfig(path string) (*Config, error) {
	cfg := &Config{}

	out := fs.ReadFile(path)

	utils.FromYAML(out, &cfg)

	if len(cfg.Cluster) == 0 {
		return nil, errNoClusterDefined
	}

	for i, c := range cfg.Cluster {
		c.Env = utils.GetEnvs()

		renderedCluster, err := c.RenderConfig()
		if err != nil {
			return nil, fmt.Errorf("error while rendering config for cluster %s: %w", c.Name, err)
		}

		cfg.Cluster[i] = *renderedCluster
	}

	return cfg, nil
}

// ValidateConfig validates config.
func ValidateConfig(path string) error {
	cfg := &Config{}

	utils.FromYAML(fs.ReadFile(path), &cfg)

	return nil
}

// GetCluster returns the vault struct matching the name.
func (c *Config) GetCluster(name string) (*Cluster, error) {
	if v, ok := os.LookupEnv(clusterEnvVar); ok {
		name = v
	}

	if name == "" && len(c.Cluster) == 1 {
		return &c.Cluster[0], nil
	}

	for _, v := range c.Cluster {
		if v.Name == name {
			return &v, nil
		}
	}

	return nil, errNoSuchClusterDefined
}
