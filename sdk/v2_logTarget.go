package sdk

type LogTarget struct {
	AlternateId string      `json:"alternateId,omitempty"`
	DetailEntry interface{} `json:"detailEntry,omitempty"`
	DisplayName string      `json:"displayName,omitempty"`
	Id          string      `json:"id,omitempty"`
	Type        string      `json:"type,omitempty"`
}
