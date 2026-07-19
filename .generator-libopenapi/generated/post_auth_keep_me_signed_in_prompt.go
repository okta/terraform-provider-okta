// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// PostAuthKeepMeSignedInPrompt represents the PostAuthKeepMeSignedInPrompt schema
type PostAuthKeepMeSignedInPrompt struct {
	// The label on the accept button when prompting for Stay signed in
	AcceptButtonText string `json:"acceptButtonText,omitempty"`
	// The label on the reject button when prompting for Stay signed in
	RejectButtonText string `json:"rejectButtonText,omitempty"`
	// The subtitle on the Sign-In Widget when prompting for Stay signed in
	Subtitle string `json:"subtitle,omitempty"`
	// The title on the Sign-In Widget when prompting for Stay signed in
	Title string `json:"title,omitempty"`
}
