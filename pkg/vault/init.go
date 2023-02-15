package vault

import (
	"github.com/hashicorp/vault/api"
)

// Initialize initializes a vault server.
func (v *Vault) Initialize(shares, threshold int, recoveryKeys bool) (*api.InitResponse, error) {
	opts := &api.InitRequest{}

	if recoveryKeys {
		opts.RecoveryShares = shares
		opts.RecoveryThreshold = threshold
	} else {
		opts.SecretShares = shares
		opts.SecretThreshold = threshold
	}

	resp, err := v.Client.Sys().Init(opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// IsInitialized returns true if vault is already initialized.
func (v *Vault) IsInitialized() (bool, error) {
	return v.Client.Sys().InitStatus()
}

// IsSealed returns true if vault is already initialized.
func (v *Vault) IsSealed() (bool, error) {
	resp, err := v.Client.Sys().Health()
	if err != nil {
		return false, err
	}

	return resp.Sealed, nil
}
