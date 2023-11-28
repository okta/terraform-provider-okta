// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
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
