// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type UserCondition struct {
	// Exclude deliberately has no omitempty. A policy's default rule reports
	// "people": {"users": {"exclude": []}} and Okta rejects any update whose
	// conditions object is not structurally identical ("Cannot modify the conditions
	// object because it is read-only"). omitempty drops empty slices whether nil or
	// not, so with it there is no way to echo that empty array back. See GH-2788.
	Exclude []string `json:"exclude"`
	Include []string `json:"include,omitempty"`
}

func NewUserCondition() *UserCondition {
	return &UserCondition{}
}

func (a *UserCondition) IsPolicyInstance() bool {
	return true
}
