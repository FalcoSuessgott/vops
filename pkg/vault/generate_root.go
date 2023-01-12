package vault

import (
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/roottoken"
)

// GenerateRootInit initialized the regeneration of a root token.
func (v *Vault) GenerateRootInit(otp string) (*api.GenerateRootStatusResponse, error) {
	resp, err := v.Client.Sys().GenerateRootInit(otp, "")
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GenerateRootUpdate enters a key to a generate root process.
func (v *Vault) GenerateRootUpdate(key, nonce string) (*api.GenerateRootStatusResponse, error) {
	resp, err := v.Client.Sys().GenerateRootUpdate(key, nonce)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DecodeRootToken encoded a root token.
func (v *Vault) DecodeRootToken(encodedToken, otp string) (string, error) {
	t, err := roottoken.DecodeToken(encodedToken, otp, len(otp))
	if err != nil {
		return "", err
	}

	return t, nil
}

// GenerateOTP generates a otp.
func (v *Vault) GenerateOTP() (string, error) {
	resp, err := v.Client.Sys().GenerateRootStatus()
	if err != nil {
		return "", err
	}

	t, err := roottoken.GenerateOTP(resp.OTPLength)
	if err != nil {
		return "", err
	}

	return t, nil
}
