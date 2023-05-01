package sdk

type Scope struct {
	AllowedOktaApps []*IframeEmbedScopeAllowedApps `json:"allowedOktaApps,omitempty"`
	StringValue     string                         `json:"stringValue,omitempty"`
	Type            string                         `json:"type,omitempty"`
}
