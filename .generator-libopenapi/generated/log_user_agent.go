// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// LogUserAgent represents the LogUserAgent schema
// "A user agent is software (a software agent) that is acting on behalf of a user." ([Definition of User Agent](https://developer.mozilla.org/en-US/docs/Glossary/User_agent))  In the Okta event data ...
type LogUserAgent struct {
	// If the client is a web browser, this field identifies the type of web browser (for example, CHROME, FIREFOX)
	Browser string `json:"browser,omitempty"`
	// The operating system that the client runs on (for example, Windows 10)
	Os string `json:"os,omitempty"`
	// A raw string representation of the user agent that is formatted according to [section 5.5.3 of HTTP/1.1 Semantics and Content](https://datatracker.ietf.org/doc/html/rfc7231#section-5.5.3). Both the...
	RawUserAgent string `json:"rawUserAgent,omitempty"`
}
