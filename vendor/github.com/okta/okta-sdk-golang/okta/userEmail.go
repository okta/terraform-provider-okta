package okta

import ()

type UserEmail struct {
	Value  string `json:"value,omitempty"`
	Status string `json:"status,omitempty"`
	Type   string `json:"type,omitempty"`
}
