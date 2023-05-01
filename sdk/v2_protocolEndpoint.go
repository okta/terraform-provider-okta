package sdk

type ProtocolEndpoint struct {
	Binding     string `json:"binding,omitempty"`
	Destination string `json:"destination,omitempty"`
	Type        string `json:"type,omitempty"`
	Url         string `json:"url,omitempty"`
}
