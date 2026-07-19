// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// NumberFactorChallengeEmbeddedLinks represents the NumberFactorChallengeEmbeddedLinks schema
// Contains the `challenge` and `correctAnswer` objects for `push` factors that use a number matching challenge
type NumberFactorChallengeEmbeddedLinks struct {
	// Number matching challenge for a `push` factor
	Challenge map[string]interface{} `json:"challenge,omitempty"`
}
