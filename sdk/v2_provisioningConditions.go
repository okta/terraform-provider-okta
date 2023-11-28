// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type ProvisioningConditions struct {
	Deprovisioned *ProvisioningDeprovisionedCondition `json:"deprovisioned,omitempty"`
	Suspended     *ProvisioningSuspendedCondition     `json:"suspended,omitempty"`
}
