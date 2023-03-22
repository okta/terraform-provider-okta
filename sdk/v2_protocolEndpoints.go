package sdk

type ProtocolEndpoints struct {
	Acs           *ProtocolEndpoint `json:"acs,omitempty"`
	Authorization *ProtocolEndpoint `json:"authorization,omitempty"`
	Jwks          *ProtocolEndpoint `json:"jwks,omitempty"`
	Metadata      *ProtocolEndpoint `json:"metadata,omitempty"`
	Slo           *ProtocolEndpoint `json:"slo,omitempty"`
	Sso           *ProtocolEndpoint `json:"sso,omitempty"`
	Token         *ProtocolEndpoint `json:"token,omitempty"`
	UserInfo      *ProtocolEndpoint `json:"userInfo,omitempty"`
}
