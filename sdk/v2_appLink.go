// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import "encoding/json"

type AppLink struct {
	AppAssignmentId  string `json:"appAssignmentId,omitempty"`
	AppInstanceId    string `json:"appInstanceId,omitempty"`
	AppName          string `json:"appName,omitempty"`
	CredentialsSetup *bool  `json:"credentialsSetup,omitempty"`
	Hidden           *bool  `json:"hidden,omitempty"`
	Id               string `json:"id,omitempty"`
	Label            string `json:"label,omitempty"`
	LinkUrl          string `json:"linkUrl,omitempty"`
	LogoUrl          string `json:"logoUrl,omitempty"`
	SortOrder        int64  `json:"-"`
	SortOrderPtr     *int64 `json:"sortOrder,omitempty"`
}

func (a *AppLink) MarshalJSON() ([]byte, error) {
	type Alias AppLink
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.SortOrder != 0 {
		result.SortOrderPtr = Int64Ptr(a.SortOrder)
	}
	return json.Marshal(&result)
}

func (a *AppLink) UnmarshalJSON(data []byte) error {
	type Alias AppLink

	result := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.SortOrderPtr != nil {
		a.SortOrder = *result.SortOrderPtr
		a.SortOrderPtr = result.SortOrderPtr
	}
	return nil
}
