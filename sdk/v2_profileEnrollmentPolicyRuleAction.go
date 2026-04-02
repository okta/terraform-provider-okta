// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
	EnrollAuthenticatorTypes   []string                                          `json:"enrollAuthenticatorTypes,omitempty"`
}

func NewProfileEnrollmentPolicyRuleAction() *ProfileEnrollmentPolicyRuleAction {
	return &ProfileEnrollmentPolicyRuleAction{}
}

func (a *ProfileEnrollmentPolicyRuleAction) IsPolicyInstance() bool {
	return true
}
