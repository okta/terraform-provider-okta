// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AuthenticationMethodObject represents the AuthenticationMethodObject schema
type AuthenticationMethodObject struct {
	// <x-lifecycle-container><x-lifecycle class="oie"></x-lifecycle></x-lifecycle-container>Authenticator ID
	ID string `json:"id,omitempty"`
	// A label that identifies the authenticator
	Key string `json:"key"`
	// Specifies the method used for the authenticator
	Method string `json:"method,omitempty"`
}
