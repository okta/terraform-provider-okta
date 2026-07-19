// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PolicyRuleConditions represents the PolicyRuleConditions schema
type PolicyRuleConditions struct {
	Network interface{} `json:"network,omitempty"`
	RiskScore interface{} `json:"riskScore,omitempty"`
	UserIdentifier interface{} `json:"userIdentifier,omitempty"`
	GrantTypes interface{} `json:"grantTypes,omitempty"`
	IdentityProvider interface{} `json:"identityProvider,omitempty"`
	Platform interface{} `json:"platform,omitempty"`
	Risk interface{} `json:"risk,omitempty"`
	Scopes interface{} `json:"scopes,omitempty"`
	Users interface{} `json:"users,omitempty"`
	UserStatus interface{} `json:"userStatus,omitempty"`
	Apps interface{} `json:"apps,omitempty"`
	Context interface{} `json:"context,omitempty"`
	Device interface{} `json:"device,omitempty"`
	Groups interface{} `json:"groups,omitempty"`
	People interface{} `json:"people,omitempty"`
	AuthContext interface{} `json:"authContext,omitempty"`
	App interface{} `json:"app,omitempty"`
	AuthProvider interface{} `json:"authProvider,omitempty"`
	BeforeScheduledAction interface{} `json:"beforeScheduledAction,omitempty"`
	Clients interface{} `json:"clients,omitempty"`
	MdmEnrollment interface{} `json:"mdmEnrollment,omitempty"`
}
