// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// SamlAttributeStatementExpression represents the SamlAttributeStatementExpression schema
// Generic `EXPRESSION` attribute statements
type SamlAttributeStatementExpression struct {
	// The name of the attribute in your app. The attribute name must be unique across all user and group attribute statements.
	Name string `json:"name,omitempty"`
	// The name format of the attribute. Supported values:
	Namespace string `json:"namespace,omitempty"`
	// The type of attribute statements object
	Type string `json:"type,omitempty"`
	// The attribute values (supports [Okta Expression Language](https://developer.okta.com/docs/reference/okta-expression-language/))
	Values []string `json:"values,omitempty"`
}
