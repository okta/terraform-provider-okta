// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"encoding/json"
	"time"
)

type SocialAuthToken struct {
	ExpiresAt       *time.Time `json:"expiresAt,omitempty"`
	Id              string     `json:"id,omitempty"`
	Scopes          []string   `json:"scopes,omitempty"`
	Token           string     `json:"token,omitempty"`
	TokenAuthScheme string     `json:"tokenAuthScheme,omitempty"`
	TokenType       string     `json:"tokenType,omitempty"`
}

func (a *SocialAuthToken) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	var token map[string]interface{}
	err := json.Unmarshal(data, &token)
	if err != nil {
		return err
	}
	if ea, found := token["expiresAt"]; found {
		if expiresAt, err := time.Parse(time.RFC3339, ea.(string)); err == nil {
			a.ExpiresAt = &expiresAt
		}
	}
	a.Id, _ = token["id"].(string)
	if scopes, found := token["scopes"]; found {
		_scopes := scopes.([]interface{})
		a.Scopes = make([]string, len(_scopes))
		for i, scope := range _scopes {
			a.Scopes[i] = scope.(string)
		}
	}
	a.Token, _ = token["token"].(string)
	a.TokenAuthScheme, _ = token["tokenAuthScheme"].(string)
	a.TokenType, _ = token["tokenType"].(string)

	return nil
}
