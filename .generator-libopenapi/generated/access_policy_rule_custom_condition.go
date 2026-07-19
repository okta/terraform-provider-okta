// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AccessPolicyRuleCustomCondition represents the AccessPolicyRuleCustomCondition schema
// Specifies [Okta Expression Language](https://developer.okta.com/docs/reference/okta-expression-language-in-identity-engine/) expressions
type AccessPolicyRuleCustomCondition struct {
	// expression to match
	Condition string `json:"condition"`
}
