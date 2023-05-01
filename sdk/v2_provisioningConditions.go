package sdk

type ProvisioningConditions struct {
	Deprovisioned *ProvisioningDeprovisionedCondition `json:"deprovisioned,omitempty"`
	Suspended     *ProvisioningSuspendedCondition     `json:"suspended,omitempty"`
}
