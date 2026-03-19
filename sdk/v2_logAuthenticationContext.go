// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import "encoding/json"

type LogAuthenticationContext struct {
	AuthenticationProvider string     `json:"authenticationProvider,omitempty"`
	AuthenticationStep     int64      `json:"-"`
	AuthenticationStepPtr  *int64     `json:"authenticationStep,omitempty"`
	CredentialProvider     string     `json:"credentialProvider,omitempty"`
	CredentialType         string     `json:"credentialType,omitempty"`
	ExternalSessionId      string     `json:"externalSessionId,omitempty"`
	Interface              string     `json:"interface,omitempty"`
	Issuer                 *LogIssuer `json:"issuer,omitempty"`
}

func (a *LogAuthenticationContext) MarshalJSON() ([]byte, error) {
	type Alias LogAuthenticationContext
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.AuthenticationStep != 0 {
		result.AuthenticationStepPtr = Int64Ptr(a.AuthenticationStep)
	}
	return json.Marshal(&result)
}

func (a *LogAuthenticationContext) UnmarshalJSON(data []byte) error {
	type Alias LogAuthenticationContext

	result := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.AuthenticationStepPtr != nil {
		a.AuthenticationStep = *result.AuthenticationStepPtr
		a.AuthenticationStepPtr = result.AuthenticationStepPtr
	}
	return nil
}
