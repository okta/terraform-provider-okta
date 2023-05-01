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
