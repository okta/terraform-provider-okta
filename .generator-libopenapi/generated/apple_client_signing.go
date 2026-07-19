// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// AppleClientSigning represents the AppleClientSigning schema
// Information used to generate the secret JSON Web Token for the token requests to Apple IdP > **Note:** The `privateKey` property is required for a CREATE request. For an UPDATE request, it can be n...
type AppleClientSigning struct {
	// The key ID that you obtained from Apple when you created the private key for the client
	Kid string `json:"kid,omitempty"`
	// The PKCS \#8 encoded private key that you created for the client and downloaded from Apple
	PrivateKey string `json:"privateKey,omitempty"`
	// The Team ID associated with your Apple developer account
	TeamId string `json:"teamId,omitempty"`
}
