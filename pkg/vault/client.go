package vault

import (
	"github.com/hashicorp/vault/api"
)

// Vault represents a vault struct used for reading and writing secrets.
type Vault struct {
	Client *api.Client
}

func NewClient(addr string) (*Vault, error) {
	cfg := &api.Config{
		Address: addr,
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Vault{client}, nil
}

func NewTokenClient(addr, token string) (*Vault, error) {
	cfg := &api.Config{
		Address: addr,
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	client.SetToken(token)

	return &Vault{client}, nil
}
