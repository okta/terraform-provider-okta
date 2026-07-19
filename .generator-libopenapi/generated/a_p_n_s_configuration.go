// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// APNSConfiguration represents the APNSConfiguration schema
type APNSConfiguration struct {
	// (Optional) File name for Admin Console display
	FileName string `json:"fileName,omitempty"`
	// 10-character Key ID obtained from the Apple developer account
	KeyId string `json:"keyId,omitempty"`
	// 10-character Team ID used to develop the iOS app
	TeamId string `json:"teamId,omitempty"`
	// APNs private authentication token signing key
	TokenSigningKey string `json:"tokenSigningKey,omitempty"`
}
