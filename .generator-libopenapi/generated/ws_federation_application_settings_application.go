// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// WsFederationApplicationSettingsApplication represents the WsFederationApplicationSettingsApplication schema
type WsFederationApplicationSettingsApplication struct {
	// The entity ID of the SP. Use the entity ID value exactly as provided by the SP.
	AudienceRestriction string `json:"audienceRestriction"`
	// Identifies the SAML authentication context class for the assertion's authentication statement
	AuthnContextClassRef string `json:"authnContextClassRef"`
	// The group name to include in the WS-Fed response attribute statement. This property is used in conjunction with the `groupFilter` property.  Groups that are filtered through the `groupFilter` expre...
	GroupName string `json:"groupName,omitempty"`
	// Launch URL for the web app
	SiteURL string `json:"siteURL"`
	// You can federate user attributes such as Okta profile fields, LDAP, Active Directory, and Workday values. The SP uses the federated WS-Fed attribute values accordingly.
	AttributeStatements string `json:"attributeStatements,omitempty"`
	// A regular expression that filters for the User Groups you want included with the `groupName` attribute. If the matching User Group has a corresponding AD group, then the attribute statement include...
	GroupFilter string `json:"groupFilter,omitempty"`
	// Specifies the WS-Fed assertion attribute value for filtered groups. This attribute is only applied to Active Directory groups.
	GroupValueFormat string `json:"groupValueFormat"`
	// The username format that you send in the WS-Fed response
	NameIDFormat string `json:"nameIDFormat"`
	// The uniform resource identifier (URI) of the WS-Fed app that's used to share resources securely within a domain. It's the identity that's sent to the Okta IdP when signing in. See [Realm name](http...
	Realm string `json:"realm,omitempty"`
	// Specifies additional username attribute statements to include in the WS-Fed assertion
	UsernameAttribute string `json:"usernameAttribute"`
	// Enables a web app to override the `wReplyURL` URL with a reply parameter.
	WReplyOverride bool `json:"wReplyOverride,omitempty"`
	// The WS-Fed SP endpoint where your users sign in
	WReplyURL string `json:"wReplyURL"`
}
