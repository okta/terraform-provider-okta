// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ContentSecurityPolicySetting represents the ContentSecurityPolicySetting schema
type ContentSecurityPolicySetting struct {
	Mode string `json:"mode,omitempty"`
	ReportUri string `json:"reportUri,omitempty"`
	SrcList []string `json:"srcList,omitempty"`
}
