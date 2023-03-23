package sdk

type LogTransaction struct {
	Detail interface{} `json:"detail,omitempty"`
	Id     string      `json:"id,omitempty"`
	Type   string      `json:"type,omitempty"`
}
