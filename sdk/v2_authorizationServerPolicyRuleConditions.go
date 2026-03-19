// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type AuthorizationServerPolicyRuleConditions struct {
	Clients    *ClientPolicyCondition                    `json:"clients,omitempty"`
	GrantTypes *GrantTypePolicyRuleCondition             `json:"grantTypes,omitempty"`
	People     *PolicyPeopleCondition                    `json:"people,omitempty"`
	Scopes     *OAuth2ScopesMediationPolicyRuleCondition `json:"scopes,omitempty"`
}
