// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PolicyRuleConditions struct {
	App                   *AppAndInstancePolicyRuleCondition             `json:"app,omitempty"`
	Apps                  *AppInstancePolicyRuleCondition                `json:"apps,omitempty"`
	AuthContext           *PolicyRuleAuthContextCondition                `json:"authContext,omitempty"`
	AuthProvider          *PasswordPolicyAuthenticationProviderCondition `json:"authProvider,omitempty"`
	BeforeScheduledAction *BeforeScheduledActionPolicyRuleCondition      `json:"beforeScheduledAction,omitempty"`
	Clients               *ClientPolicyCondition                         `json:"clients,omitempty"`
	Context               *ContextPolicyRuleCondition                    `json:"context,omitempty"`
	Device                *DevicePolicyRuleCondition                     `json:"device,omitempty"`
	GrantTypes            *GrantTypePolicyRuleCondition                  `json:"grantTypes,omitempty"`
	Groups                *GroupPolicyRuleCondition                      `json:"groups,omitempty"`
	IdentityProvider      *IdentityProviderPolicyRuleCondition           `json:"identityProvider,omitempty"`
	MdmEnrollment         *MDMEnrollmentPolicyRuleCondition              `json:"mdmEnrollment,omitempty"`
	Network               *PolicyNetworkCondition                        `json:"network,omitempty"`
	People                *PolicyPeopleCondition                         `json:"people,omitempty"`
	Platform              *PlatformPolicyRuleCondition                   `json:"platform,omitempty"`
	Risk                  *RiskPolicyRuleCondition                       `json:"risk,omitempty"`
	RiskScore             *RiskScorePolicyRuleCondition                  `json:"riskScore,omitempty"`
	Scopes                *OAuth2ScopesMediationPolicyRuleCondition      `json:"scopes,omitempty"`
	UserIdentifier        *UserIdentifierPolicyRuleCondition             `json:"userIdentifier,omitempty"`
	UserStatus            *UserStatusPolicyRuleCondition                 `json:"userStatus,omitempty"`
	Users                 *UserPolicyRuleCondition                       `json:"users,omitempty"`
}

func NewPolicyRuleConditions() *PolicyRuleConditions {
	return &PolicyRuleConditions{}
}

func (a *PolicyRuleConditions) IsPolicyInstance() bool {
	return true
}
