// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import "encoding/json"

type AuthenticatorSettings struct {
	AllowedFor                string          `json:"allowedFor,omitempty"`
	AppInstanceId             string          `json:"appInstanceId,omitempty"`
	ChannelBinding            *ChannelBinding `json:"channelBinding,omitempty"`
	Compliance                *Compliance     `json:"compliance,omitempty"`
	TokenLifetimeInMinutes    int64           `json:"-"`
	TokenLifetimeInMinutesPtr *int64          `json:"tokenLifetimeInMinutes,omitempty"`
	UserVerification          string          `json:"userVerification,omitempty"`
	EnrollmentSecurityLevel   string          `json:"enrollmentSecurityLevel,omitempty"`
	UserVerificationMethods   []string        `json:"userVerificationMethods,omitempty"`
}

func (a *AuthenticatorSettings) MarshalJSON() ([]byte, error) {
	type Alias AuthenticatorSettings
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.TokenLifetimeInMinutes != 0 {
		result.TokenLifetimeInMinutesPtr = Int64Ptr(a.TokenLifetimeInMinutes)
	}
	return json.Marshal(&result)
}

func (a *AuthenticatorSettings) UnmarshalJSON(data []byte) error {
	type Alias AuthenticatorSettings

	result := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.TokenLifetimeInMinutesPtr != nil {
		a.TokenLifetimeInMinutes = *result.TokenLifetimeInMinutesPtr
		a.TokenLifetimeInMinutesPtr = result.TokenLifetimeInMinutesPtr
	}
	return nil
}
