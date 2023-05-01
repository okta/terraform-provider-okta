package sdk

import "encoding/json"

type LogSecurityContext struct {
	AsNumber    int64  `json:"-"`
	AsNumberPtr *int64 `json:"asNumber,omitempty"`
	AsOrg       string `json:"asOrg,omitempty"`
	Domain      string `json:"domain,omitempty"`
	IsProxy     *bool  `json:"isProxy,omitempty"`
	Isp         string `json:"isp,omitempty"`
}

func (a *LogSecurityContext) MarshalJSON() ([]byte, error) {
	type Alias LogSecurityContext
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.AsNumber != 0 {
		result.AsNumberPtr = Int64Ptr(a.AsNumber)
	}
	return json.Marshal(&result)
}

func (a *LogSecurityContext) UnmarshalJSON(data []byte) error {
	type Alias LogSecurityContext

	result := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.AsNumberPtr != nil {
		a.AsNumber = *result.AsNumberPtr
		a.AsNumberPtr = result.AsNumberPtr
	}
	return nil
}
