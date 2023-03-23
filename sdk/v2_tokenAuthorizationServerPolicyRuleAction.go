package sdk

import "encoding/json"

type TokenAuthorizationServerPolicyRuleAction struct {
	AccessTokenLifetimeMinutes     int64                                               `json:"-"`
	AccessTokenLifetimeMinutesPtr  *int64                                              `json:"accessTokenLifetimeMinutes"`
	InlineHook                     *TokenAuthorizationServerPolicyRuleActionInlineHook `json:"inlineHook,omitempty"`
	RefreshTokenLifetimeMinutes    int64                                               `json:"-"`
	RefreshTokenLifetimeMinutesPtr *int64                                              `json:"refreshTokenLifetimeMinutes"`
	RefreshTokenWindowMinutes      int64                                               `json:"-"`
	RefreshTokenWindowMinutesPtr   *int64                                              `json:"refreshTokenWindowMinutes"`
}

func (a *TokenAuthorizationServerPolicyRuleAction) MarshalJSON() ([]byte, error) {
	type Alias TokenAuthorizationServerPolicyRuleAction
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.AccessTokenLifetimeMinutes != 0 {
		result.AccessTokenLifetimeMinutesPtr = Int64Ptr(a.AccessTokenLifetimeMinutes)
	}
	if a.RefreshTokenLifetimeMinutes != 0 {
		result.RefreshTokenLifetimeMinutesPtr = Int64Ptr(a.RefreshTokenLifetimeMinutes)
	}
	if a.RefreshTokenWindowMinutes != 0 {
		result.RefreshTokenWindowMinutesPtr = Int64Ptr(a.RefreshTokenWindowMinutes)
	}
	return json.Marshal(&result)
}

func (a *TokenAuthorizationServerPolicyRuleAction) UnmarshalJSON(data []byte) error {
	type Alias TokenAuthorizationServerPolicyRuleAction

	result := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.AccessTokenLifetimeMinutesPtr != nil {
		a.AccessTokenLifetimeMinutes = *result.AccessTokenLifetimeMinutesPtr
		a.AccessTokenLifetimeMinutesPtr = result.AccessTokenLifetimeMinutesPtr
	}
	if result.RefreshTokenLifetimeMinutesPtr != nil {
		a.RefreshTokenLifetimeMinutes = *result.RefreshTokenLifetimeMinutesPtr
		a.RefreshTokenLifetimeMinutesPtr = result.RefreshTokenLifetimeMinutesPtr
	}
	if result.RefreshTokenWindowMinutesPtr != nil {
		a.RefreshTokenWindowMinutes = *result.RefreshTokenWindowMinutesPtr
		a.RefreshTokenWindowMinutesPtr = result.RefreshTokenWindowMinutesPtr
	}
	return nil
}
