// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type OktaSignOnPolicyRuleSignonActions struct {
	Access                  string                                    `json:"access,omitempty"`
	FactorLifetime          int64                                     `json:"-"`
	FactorLifetimePtr       *int64                                    `json:"factorLifetime,omitempty"`
	FactorPromptMode        string                                    `json:"factorPromptMode,omitempty"`
	RememberDeviceByDefault *bool                                     `json:"rememberDeviceByDefault,omitempty"`
	RequireFactor           *bool                                     `json:"requireFactor,omitempty"`
	Session                 *OktaSignOnPolicyRuleSignonSessionActions `json:"session,omitempty"`
}
