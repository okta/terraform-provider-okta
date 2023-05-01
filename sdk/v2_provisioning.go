package sdk

type Provisioning struct {
	Action        string                  `json:"action,omitempty"`
	Conditions    *ProvisioningConditions `json:"conditions,omitempty"`
	Groups        *ProvisioningGroups     `json:"groups,omitempty"`
	ProfileMaster *bool                   `json:"profileMaster,omitempty"`
}
