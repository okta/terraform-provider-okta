// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// InlineHookRequestObject represents the InlineHookRequestObject schema
// The API request that triggered the inline hook
type InlineHookRequestObject struct {
	// The unique identifier that Okta assigned to the API request
	ID string `json:"id,omitempty"`
	// The IP address of the client that made the API request
	IpAddress string `json:"ipAddress,omitempty"`
	// The HTTP request method of the API request
	Method string `json:"method,omitempty"`
	// The URL of the API endpoint
	Url map[string]interface{} `json:"url,omitempty"`
}
