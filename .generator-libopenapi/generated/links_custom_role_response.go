// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LinksCustomRoleResponse represents the LinksCustomRoleResponse schema
// Specifies link relations (see [Web Linking](https://www.rfc-editor.org/rfc/rfc8288)) available using the [JSON Hypertext Application Language](https://datatracker.ietf.org/doc/html/draft-kelly-json...
type LinksCustomRoleResponse struct {
	Assignee interface{} `json:"assignee,omitempty"`
	Member interface{} `json:"member,omitempty"`
	Permissions interface{} `json:"permissions,omitempty"`
	Resource-set interface{} `json:"resource-set,omitempty"`
	Role interface{} `json:"role,omitempty"`
}
