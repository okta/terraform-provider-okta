package sdk

import "encoding/json"

type IdentityProviderCredentialsTrust struct {
	Audience                   string `json:"audience,omitempty"`
	Issuer                     string `json:"issuer,omitempty"`
	Kid                        string `json:"kid,omitempty"`
	Revocation                 string `json:"revocation,omitempty"`
	RevocationCacheLifetime    int64  `json:"-"`
	RevocationCacheLifetimePtr *int64 `json:"revocationCacheLifetime,omitempty"`
}

func (a *IdentityProviderCredentialsTrust) MarshalJSON() ([]byte, error) {
	type Alias IdentityProviderCredentialsTrust
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.RevocationCacheLifetime != 0 {
		result.RevocationCacheLifetimePtr = Int64Ptr(a.RevocationCacheLifetime)
	}
	return json.Marshal(&result)
}

func (a *IdentityProviderCredentialsTrust) UnmarshalJSON(data []byte) error {
	type Alias IdentityProviderCredentialsTrust

	result := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.RevocationCacheLifetimePtr != nil {
		a.RevocationCacheLifetime = *result.RevocationCacheLifetimePtr
		a.RevocationCacheLifetimePtr = result.RevocationCacheLifetimePtr
	}
	return nil
}
