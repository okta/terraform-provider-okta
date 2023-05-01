package sdk

type ProfileMappingSource struct {
	Links interface{} `json:"_links,omitempty"`
	Id    string      `json:"id,omitempty"`
	Name  string      `json:"name,omitempty"`
	Type  string      `json:"type,omitempty"`
}
