// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type Provisioning struct {
	Action        string                  `json:"action,omitempty"`
	Conditions    *ProvisioningConditions `json:"conditions,omitempty"`
	Groups        *ProvisioningGroups     `json:"groups,omitempty"`
	ProfileMaster *bool                   `json:"profileMaster,omitempty"`
}
