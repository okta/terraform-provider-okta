// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PreRegistrationInlineHook struct {
	InlineHookId string `json:"inlineHookId,omitempty"`
}

func NewPreRegistrationInlineHook() *PreRegistrationInlineHook {
	return &PreRegistrationInlineHook{}
}

func (a *PreRegistrationInlineHook) IsPolicyInstance() bool {
	return true
}
