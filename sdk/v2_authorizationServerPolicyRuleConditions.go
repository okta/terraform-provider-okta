package sdk

type AuthorizationServerPolicyRuleConditions struct {
	Clients    *ClientPolicyCondition                    `json:"clients,omitempty"`
	GrantTypes *GrantTypePolicyRuleCondition             `json:"grantTypes,omitempty"`
	People     *PolicyPeopleCondition                    `json:"people,omitempty"`
	Scopes     *OAuth2ScopesMediationPolicyRuleCondition `json:"scopes,omitempty"`
}
