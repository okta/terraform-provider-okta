package sdk

type UserCondition struct {
	Exclude []string `json:"exclude,omitempty"`
	Include []string `json:"include,omitempty"`
}

func NewUserCondition() *UserCondition {
	return &UserCondition{}
}

func (a *UserCondition) IsPolicyInstance() bool {
	return true
}
