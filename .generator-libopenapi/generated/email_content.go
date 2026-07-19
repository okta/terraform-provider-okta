// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// EmailContent represents the EmailContent schema
type EmailContent struct {
	// The HTML body of the email. May contain [variable references](https://velocity.apache.org/engine/1.7/user-guide.html#references).  <x-lifecycle class="ea"></x-lifecycle> Not required if Custom lang...
	Body string `json:"body"`
	// The email subject. May contain [variable references](https://velocity.apache.org/engine/1.7/user-guide.html#references).  <x-lifecycle class="ea"></x-lifecycle> Not required if Custom languages for...
	Subject string `json:"subject"`
}
