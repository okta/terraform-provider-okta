// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
