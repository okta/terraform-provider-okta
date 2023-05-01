package sdk

type UserTypeCondition struct {
	Exclude []string `json:"exclude,omitempty"`
	Include []string `json:"include,omitempty"`
}

func NewUserTypeCondition() *UserTypeCondition {
	return &UserTypeCondition{}
}

func (a *UserTypeCondition) IsPolicyInstance() bool {
	return true
}
