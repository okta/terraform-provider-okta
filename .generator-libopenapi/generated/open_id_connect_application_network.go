// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// OpenIdConnectApplicationNetwork represents the OpenIdConnectApplicationNetwork schema
// The network restrictions of the client
type OpenIdConnectApplicationNetwork struct {
	// If `ZONE` is specified as a connection, then specify the included IP network zones here. Value can be "ALL_IP_ZONES" or an array of zone IDs.
	Include []string `json:"include,omitempty"`
	// The connection type of the network. Can be `ANYWHERE` or `ZONE`.
	Connection string `json:"connection"`
	// If `ZONE` is specified as a connection, then specify the excluded IP network zones here. Value can be "ALL_IP_ZONES" or an array of zone IDs.
	Exclude []string `json:"exclude,omitempty"`
}
