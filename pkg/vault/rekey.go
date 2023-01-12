package vault

import "github.com/hashicorp/vault/api"

// RekeyInit inits a rekeying of a vault server.
func (v *Vault) RekeyInit(shares, threshold int, recoveryKeys bool) (*api.RekeyStatusResponse, error) {
	var err error

	opts := &api.RekeyInitRequest{
		SecretShares:    shares,
		SecretThreshold: threshold,
	}

	if recoveryKeys {
		resp, err := v.Client.Sys().RekeyRecoveryKeyInit(opts)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}

	resp, err := v.Client.Sys().RekeyInit(opts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// RekeyUpdate rekeys a vault server.
func (v *Vault) RekeyUpdate(key, nonce string) (*api.RekeyUpdateResponse, error) {
	resp, err := v.Client.Sys().RekeyUpdate(key, nonce)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
