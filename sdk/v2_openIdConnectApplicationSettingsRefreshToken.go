package sdk

import "encoding/json"

type OpenIdConnectApplicationSettingsRefreshToken struct {
	Leeway       int64  `json:"-"`
	LeewayPtr    *int64 `json:"leeway"`
	RotationType string `json:"rotation_type,omitempty"`
}

func (a *OpenIdConnectApplicationSettingsRefreshToken) MarshalJSON() ([]byte, error) {
	type Alias OpenIdConnectApplicationSettingsRefreshToken
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.Leeway != 0 {
		result.LeewayPtr = Int64Ptr(a.Leeway)
	}
	return json.Marshal(&result)
}

func (a *OpenIdConnectApplicationSettingsRefreshToken) UnmarshalJSON(data []byte) error {
	type Alias OpenIdConnectApplicationSettingsRefreshToken

	result := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.LeewayPtr != nil {
		a.Leeway = *result.LeewayPtr
		a.LeewayPtr = result.LeewayPtr
	}
	return nil
}
