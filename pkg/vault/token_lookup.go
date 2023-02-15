package vault

import "fmt"

// TokenLookup unseals a vault node, returns true if the node is unsealed returns the attached policies.
func (v *Vault) TokenLookup() ([]interface{}, error) {
	token := v.Client.Auth().Token()

	data, err := token.LookupSelf()
	if err != nil {
		return nil, err
	}

	//nolint: forcetypeassert
	if v, ok := data.Data["policies"]; ok {
		return v.([]interface{}), nil
	}

	return nil, fmt.Errorf("couldnt lookup token policies")
}
