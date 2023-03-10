package vault

import "github.com/hashicorp/vault/api"

// RekeyInit inits a rekeying of a vault server.
func (v *Vault) RekeyInit(shares, threshold int, recoveryKeys bool) (*api.RekeyStatusResponse, error) {
	var fn func(config *api.RekeyInitRequest) (*api.RekeyStatusResponse, error)

	opts := &api.RekeyInitRequest{
		SecretShares:    shares,
		SecretThreshold: threshold,
	}

	if recoveryKeys {
		fn = v.Client.Sys().RekeyRecoveryKeyInit
	} else {
		fn = v.Client.Sys().RekeyInit
	}

	resp, err := fn(opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// RekeyUpdate rekeys a vault server.
func (v *Vault) RekeyUpdate(key, nonce string, recoveryKeys bool) (*api.RekeyUpdateResponse, error) {
	var fn func(shard, nonce string) (*api.RekeyUpdateResponse, error)

	if recoveryKeys {
		fn = v.Client.Sys().RekeyRecoveryKeyUpdate
	} else {
		fn = v.Client.Sys().RekeyUpdate
	}

	resp, err := fn(key, nonce)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
