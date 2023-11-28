// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type Protocol struct {
	Algorithms  *ProtocolAlgorithms          `json:"algorithms,omitempty"`
	Credentials *IdentityProviderCredentials `json:"credentials,omitempty"`
	Endpoints   *ProtocolEndpoints           `json:"endpoints,omitempty"`
	Issuer      *ProtocolEndpoint            `json:"issuer,omitempty"`
	RelayState  *ProtocolRelayState          `json:"relayState,omitempty"`
	Scopes      []string                     `json:"scopes,omitempty"`
	Settings    *ProtocolSettings            `json:"settings,omitempty"`
	Type        string                       `json:"type,omitempty"`
}
