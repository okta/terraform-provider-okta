package sdk

type OAuth2Scope struct {
	Consent         string `json:"consent,omitempty"`
	Default         *bool  `json:"default,omitempty"`
	Description     string `json:"description,omitempty"`
	DisplayName     string `json:"displayName,omitempty"`
	Id              string `json:"id,omitempty"`
	MetadataPublish string `json:"metadataPublish,omitempty"`
	Name            string `json:"name,omitempty"`
	System          *bool  `json:"system,omitempty"`
}
