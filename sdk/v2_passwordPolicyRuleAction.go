// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

// PasswordPolicyRuleAction represents a password policy rule action.
// The Requirement field was added to support selfServicePasswordReset requirement configuration.
// The logic has been updated by AI to support new fields, review carefully
type PasswordPolicyRuleAction struct {
	Access      string                         `json:"access,omitempty"`
	Requirement *PasswordPolicyRuleRequirement `json:"requirement,omitempty"` // Added to support SSPR requirement fields. The logic has been updated by AI to support new fields, review carefully
}

// PasswordPolicyRuleRequirement holds the SSPR requirement configuration including
// primary methods, step-up authentication, and access control mode.
// The logic has been updated by AI to support new fields, review carefully
type PasswordPolicyRuleRequirement struct {
	Primary       *PasswordPolicyRuleRequirementPrimary `json:"primary,omitempty"`
	StepUp        *PasswordPolicyRuleRequirementStepUp  `json:"stepUp,omitempty"`
	AccessControl string                                `json:"accessControl,omitempty"`
}

// PasswordPolicyRuleRequirementPrimary holds the list of primary authentication methods
// allowed for self-service password reset (e.g., otp, push, sms, email, voice).
// The logic has been updated by AI to support new fields, review carefully
type PasswordPolicyRuleRequirementPrimary struct {
	Methods []string `json:"methods,omitempty"`
}

// PasswordPolicyRuleRequirementStepUp holds the step-up authentication requirement
// for self-service password reset.
// The logic has been updated by AI to support new fields, review carefully
type PasswordPolicyRuleRequirementStepUp struct {
	Required bool `json:"required"`
}

func NewPasswordPolicyRuleAction() *PasswordPolicyRuleAction {
	return &PasswordPolicyRuleAction{}
}

func (a *PasswordPolicyRuleAction) IsPolicyInstance() bool {
	return true
}
