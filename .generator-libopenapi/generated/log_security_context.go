// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogSecurityContext represents the LogSecurityContext schema
// The `securityContext` object provides security information that is directly related to the evaluation of the event's IP reputation. IP reputation is a trustworthiness rating that evaluates how like...
type LogSecurityContext struct {
	// The Internet service provider that's used to send the event's request
	Isp string `json:"isp,omitempty"`
	// Specifies whether an event's request is from a known proxy
	IsProxy bool `json:"isProxy,omitempty"`
	BotProtection interface{} `json:"botProtection,omitempty"`
	// The domain name that's associated with the IP address of the inbound event request
	Domain string `json:"domain,omitempty"`
	Risk interface{} `json:"risk,omitempty"`
	// The result of the user behavior detection models associated with the event
	UserBehaviors []interface{} `json:"userBehaviors,omitempty"`
	// The [Autonomous system](https://docs.telemetry.mozilla.org/datasets/other/asn_aggregates/reference) number that's associated with the autonomous system the event request was sourced to
	AsNumber int `json:"asNumber,omitempty"`
	// The organization that is associated with the autonomous system that the event request is sourced to
	AsOrg string `json:"asOrg,omitempty"`
	IpDetails interface{} `json:"ipDetails,omitempty"`
}
