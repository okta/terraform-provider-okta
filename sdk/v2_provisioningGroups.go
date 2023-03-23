package sdk

type ProvisioningGroups struct {
	Action              string   `json:"action,omitempty"`
	Assignments         []string `json:"assignments,omitempty"`
	Filter              []string `json:"filter,omitempty"`
	SourceAttributeName string   `json:"sourceAttributeName,omitempty"`
}
