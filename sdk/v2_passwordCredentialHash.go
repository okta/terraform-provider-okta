package sdk

import "encoding/json"

type PasswordCredentialHash struct {
	Algorithm     string `json:"algorithm,omitempty"`
	Salt          string `json:"salt,omitempty"`
	SaltOrder     string `json:"saltOrder,omitempty"`
	Value         string `json:"value,omitempty"`
	WorkFactor    int64  `json:"-"`
	WorkFactorPtr *int64 `json:"workFactor,omitempty"`
}

func (a *PasswordCredentialHash) MarshalJSON() ([]byte, error) {
	type Alias PasswordCredentialHash
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.WorkFactor != 0 {
		result.WorkFactorPtr = Int64Ptr(a.WorkFactor)
	}
	return json.Marshal(&result)
}

func (a *PasswordCredentialHash) UnmarshalJSON(data []byte) error {
	type Alias PasswordCredentialHash

	result := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.WorkFactorPtr != nil {
		a.WorkFactor = *result.WorkFactorPtr
		a.WorkFactorPtr = result.WorkFactorPtr
	}
	return nil
}
