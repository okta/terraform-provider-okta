// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// BaseEmailServer represents the BaseEmailServer schema
type BaseEmailServer struct {
	// Username that's used to access your SMTP server
	Username string `json:"username"`
	// Human-readable name for your SMTP server
	Alias string `json:"alias"`
	AuthType interface{} `json:"authType"`
	// If `true`, all email traffic is routed through your SMTP server
	Enabled bool `json:"enabled"`
	// Hostname or IP address of your SMTP server
	Host string `json:"host"`
	// ID of your SMTP server
	ID string `json:"id,omitempty"`
	// Port number of your SMTP server
	Port int `json:"port"`
}
