// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type LogTransaction struct {
	Detail interface{} `json:"detail,omitempty"`
	Id     string      `json:"id,omitempty"`
	Type   string      `json:"type,omitempty"`
}
