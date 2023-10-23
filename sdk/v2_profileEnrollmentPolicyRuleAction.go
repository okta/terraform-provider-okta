package sdk

type ProfileEnrollmentPolicyRuleAction struct {
	Access                     string                                            `json:"access,omitempty"`
	ActivationRequirements     *ProfileEnrollmentPolicyRuleActivationRequirement `json:"activationRequirements,omitempty"`
	PreRegistrationInlineHooks []*PreRegistrationInlineHook                      `json:"preRegistrationInlineHooks,omitempty"`
	ProfileAttributes          []*ProfileEnrollmentPolicyRuleProfileAttribute    `json:"profileAttributes,omitempty"`
	ProgressiveProfilingAction string                                            `json:"progressiveProfilingAction,omitempty"`
	TargetGroupIds             []string                                          `json:"targetGroupIds,omitempty"`
	UiSchemaId                 string                                            `json:"uiSchemaId,omitempty"`
	UnknownUserAction          string                                            `json:"unknownUserAction,omitempty"`
}

func NewProfileEnrollmentPolicyRuleAction() *ProfileEnrollmentPolicyRuleAction {
	return &ProfileEnrollmentPolicyRuleAction{}
}

func (a *ProfileEnrollmentPolicyRuleAction) IsPolicyInstance() bool {
	return true
}
