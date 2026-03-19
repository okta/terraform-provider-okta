// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type ProvisioningGroups struct {
	Action              string   `json:"action,omitempty"`
	Assignments         []string `json:"assignments,omitempty"`
	Filter              []string `json:"filter,omitempty"`
	SourceAttributeName string   `json:"sourceAttributeName,omitempty"`
}
