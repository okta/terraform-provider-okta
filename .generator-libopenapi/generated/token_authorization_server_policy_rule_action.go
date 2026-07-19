// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// TokenAuthorizationServerPolicyRuleAction represents the TokenAuthorizationServerPolicyRuleAction schema
type TokenAuthorizationServerPolicyRuleAction struct {
	// Timeframe when the refresh token is valid. The minimum is 10 minutes. The maximum is five years (2,628,000 minutes).
	RefreshTokenWindowMinutes int `json:"refreshTokenWindowMinutes,omitempty"`
	// Lifetime of the access token in minutes. The minimum is five minutes. The maximum is one day.
	AccessTokenLifetimeMinutes int `json:"accessTokenLifetimeMinutes,omitempty"`
	InlineHook interface{} `json:"inlineHook,omitempty"`
	// Lifetime of the refresh token is the minimum access token lifetime.
	RefreshTokenLifetimeMinutes int `json:"refreshTokenLifetimeMinutes,omitempty"`
}
