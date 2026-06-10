// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type AccessPolicyRuleConditions struct {
	App                   *AppAndInstancePolicyRuleCondition             `json:"app,omitempty"`
	Apps                  *AppInstancePolicyRuleCondition                `json:"apps,omitempty"`
	AuthContext           *PolicyRuleAuthContextCondition                `json:"authContext,omitempty"`
	AuthProvider          *PasswordPolicyAuthenticationProviderCondition `json:"authProvider,omitempty"`
	BeforeScheduledAction *BeforeScheduledActionPolicyRuleCondition      `json:"beforeScheduledAction,omitempty"`
	Clients               *ClientPolicyCondition                         `json:"clients,omitempty"`
	Context               *ContextPolicyRuleCondition                    `json:"context,omitempty"`
	Device                *DeviceAccessPolicyRuleCondition               `json:"device,omitempty"`
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
	ElCondition           *AccessPolicyRuleCustomCondition               `json:"elCondition,omitempty"`
	UserType              *UserTypeCondition                             `json:"userType,omitempty"`
	Office365Client       *Office365ClientCondition                      `json:"office365Client,omitempty"`
}

func NewAccessPolicyRuleConditions() *AccessPolicyRuleConditions {
	return &AccessPolicyRuleConditions{}
}

func (a *AccessPolicyRuleConditions) IsPolicyInstance() bool {
	return true
}
