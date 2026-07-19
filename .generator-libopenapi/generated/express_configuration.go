// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// ExpressConfiguration represents the ExpressConfiguration schema
// Auth0 tenant details used to enable Express Configuration on an OIN Integration submission. Populates the submission's sign-on, authentication, and app settings so that admins can configure the app...
type ExpressConfiguration struct {
	// The Express Configuration capabilities to enable on the submission. If omitted, all capabilities are automatically configured based on the OIN integration's supported protocols.
	Capabilities []interface{} `json:"capabilities,omitempty"`
	// The URL template that Okta uses to launch the app from the end-user dashboard. Supports template variables `{organization_name}`, `{organization_id}`, and `{connection_name}`.
	InitiateLoginUriTemplate string `json:"initiateLoginUriTemplate,omitempty"`
	// The Auth0 admin login domain that Okta redirects to as part of the consent flow in a web browser. Use the custom domain name if one is configured in Auth0.
	LoginDomain string `json:"loginDomain"`
	// The client ID of the additional client application that Auth0 creates for the OIN administrator consent flow between Okta and Auth0
	OinClientId string `json:"oinClientId"`
	// The Auth0 tenant domain (for example, `example.auth0.com`)
	TenantDomain string `json:"tenantDomain"`
	// The client ID of the SaaS application that end users sign in to
	ApplicationClientId string `json:"applicationClientId"`
}
