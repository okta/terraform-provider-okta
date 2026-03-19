// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import "encoding/json"

type PasswordPolicyRecoveryEmailRecoveryToken struct {
	TokenLifetimeMinutes    int64  `json:"-"`
	TokenLifetimeMinutesPtr *int64 `json:"tokenLifetimeMinutes"`
}

func NewPasswordPolicyRecoveryEmailRecoveryToken() *PasswordPolicyRecoveryEmailRecoveryToken {
	return &PasswordPolicyRecoveryEmailRecoveryToken{
		TokenLifetimeMinutes:    10080,
		TokenLifetimeMinutesPtr: Int64Ptr(10080),
	}
}

func (a *PasswordPolicyRecoveryEmailRecoveryToken) IsPolicyInstance() bool {
	return true
}

func (a *PasswordPolicyRecoveryEmailRecoveryToken) MarshalJSON() ([]byte, error) {
	type Alias PasswordPolicyRecoveryEmailRecoveryToken
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.TokenLifetimeMinutes != 0 {
		result.TokenLifetimeMinutesPtr = Int64Ptr(a.TokenLifetimeMinutes)
	}
	return json.Marshal(&result)
}

func (a *PasswordPolicyRecoveryEmailRecoveryToken) UnmarshalJSON(data []byte) error {
	type Alias PasswordPolicyRecoveryEmailRecoveryToken

	result := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.TokenLifetimeMinutesPtr != nil {
		a.TokenLifetimeMinutes = *result.TokenLifetimeMinutesPtr
		a.TokenLifetimeMinutesPtr = result.TokenLifetimeMinutesPtr
	}
	return nil
}
