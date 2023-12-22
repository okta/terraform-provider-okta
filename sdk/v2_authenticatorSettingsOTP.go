// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import "encoding/json"

type AuthenticatorSettingsOTP struct {
	AcceptableAdjacentIntervals int    `json:"acceptableAdjacentIntervals"`
	Algorithm                   string `json:"algorithm"`
	Encoding                    string `json:"encoding"`
	PassCodeLength              int    `json:"passCodeLength"`
	Protocol                    string `json:"protocol"`
	TimeIntervalInSeconds       int    `json:"timeIntervalInSeconds"`
}

func (a *AuthenticatorSettingsOTP) MarshalJSON() ([]byte, error) {
	type Alias AuthenticatorSettingsOTP
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	return json.Marshal(&result)
}

func (a *AuthenticatorSettingsOTP) UnmarshalJSON(data []byte) error {
	type Alias AuthenticatorSettingsOTP

	result := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	return nil
}
