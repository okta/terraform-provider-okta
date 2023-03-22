package sdk

type LogActor struct {
	AlternateId string      `json:"alternateId,omitempty"`
	Detail      interface{} `json:"detail,omitempty"`
	DisplayName string      `json:"displayName,omitempty"`
	Id          string      `json:"id,omitempty"`
	Type        string      `json:"type,omitempty"`
}
