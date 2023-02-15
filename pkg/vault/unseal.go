package vault

// Unseal unseals a vault node, returns true if the node is unsealed.
func (v *Vault) Unseal(key string) (bool, error) {
	resp, err := v.Client.Sys().Unseal(key)
	if err != nil {
		return false, err
	}

	return !resp.Sealed, nil
}

// Seal seals a cluster.
func (v *Vault) Seal() error {
	err := v.Client.Sys().Seal()
	if err != nil {
		return err
	}

	return nil
}
