// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// IdpPolicyRuleActionMatchCriteria represents the IdpPolicyRuleActionMatchCriteria schema
type IdpPolicyRuleActionMatchCriteria struct {
	// You can provide an Okta Expression Language expression with the Login Context that's evaluated with the IdP. For example, the value `login.identifier` refers to the user's username. If the user is ...
	ProviderExpression string `json:"providerExpression,omitempty"`
	// The IdP property that the evaluated string should match to
	PropertyName string `json:"propertyName,omitempty"`
}
