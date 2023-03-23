package sdk

type ScheduledUserLifecycleAction struct {
	Status string `json:"status,omitempty"`
}

func NewScheduledUserLifecycleAction() *ScheduledUserLifecycleAction {
	return &ScheduledUserLifecycleAction{}
}

func (a *ScheduledUserLifecycleAction) IsPolicyInstance() bool {
	return true
}
