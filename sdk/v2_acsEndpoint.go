// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import "encoding/json"

type AcsEndpoint struct {
	Index    int64  `json:"-"`
	IndexPtr *int64 `json:"index,omitempty"`
	Url      string `json:"url,omitempty"`
}

func (a *AcsEndpoint) MarshalJSON() ([]byte, error) {
	type Alias AcsEndpoint
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.Index != 0 {
		result.IndexPtr = Int64Ptr(a.Index)
	}
	return json.Marshal(&result)
}

func (a *AcsEndpoint) UnmarshalJSON(data []byte) error {
	type Alias AcsEndpoint

	result := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.IndexPtr != nil {
		a.Index = *result.IndexPtr
		a.IndexPtr = result.IndexPtr
	}
	return nil
}
