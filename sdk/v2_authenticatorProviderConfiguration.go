// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import "encoding/json"

type AuthenticatorProviderConfiguration struct {
	AuthPort         int64                                            `json:"-"`
	AuthPortPtr      *int64                                           `json:"authPort,omitempty"`
	Host             string                                           `json:"host,omitempty"`
	HostName         string                                           `json:"hostName,omitempty"`
	InstanceId       string                                           `json:"instanceId,omitempty"`
	IntegrationKey   string                                           `json:"integrationKey,omitempty"`
	SecretKey        string                                           `json:"secretKey,omitempty"`
	SharedSecret     string                                           `json:"sharedSecret,omitempty"`
	UserNameTemplate *AuthenticatorProviderConfigurationUserNamePlate `json:"userNameTemplate,omitempty"`
	APNS             *APNS                                            `json:"apns,omitempty"`
	FCM              *FCM                                             `json:"fcm,omitempty"`
}

func (a *AuthenticatorProviderConfiguration) MarshalJSON() ([]byte, error) {
	type Alias AuthenticatorProviderConfiguration
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.AuthPort != 0 {
		result.AuthPortPtr = Int64Ptr(a.AuthPort)
	}
	return json.Marshal(&result)
}

func (a *AuthenticatorProviderConfiguration) UnmarshalJSON(data []byte) error {
	type Alias AuthenticatorProviderConfiguration

	result := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.AuthPortPtr != nil {
		a.AuthPort = *result.AuthPortPtr
		a.AuthPortPtr = result.AuthPortPtr
	}
	return nil
}
