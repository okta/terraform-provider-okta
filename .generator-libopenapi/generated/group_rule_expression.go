// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// GroupRuleExpression represents the GroupRuleExpression schema
// Defines Okta specific [group-rules expression](https://developer.okta.com/docs/reference/okta-expression-language/#expressions-in-group-rules)
type GroupRuleExpression struct {
	// Expression type. Only valid value is '`urn:okta:expression:1.0`'.
	Type string `json:"type,omitempty"`
	// Okta expression that would result in a Boolean value
	Value string `json:"value,omitempty"`
}
