// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// GracePeriodRequiredSoon represents the GracePeriodRequiredSoon schema
// <x-lifecycle-container><x-lifecycle class="ea"></x-lifecycle></x-lifecycle-container>Customizable strings to use with [grace periods](https://developer.okta.com/docs/concepts/policies/#authenticato...
type GracePeriodRequiredSoon struct {
	// The description that's shown on the Sign-In Widget for users who are within an authenticator grace period. This description prompts users to enroll required authenticators before their grace period...
	GracePeriodRequiredSoonDescription string `json:"gracePeriodRequiredSoonDescription,omitempty"`
	// The label of the custom link that's shown on the Sign-In Widget when users are prompted to enroll required authenticators before their grace period ends.
	GracePeriodRequiredSoonCustomLinkLabel string `json:"gracePeriodRequiredSoonCustomLinkLabel,omitempty"`
	// The URL for the custom link that's shown on the Sign-In Widget when users are prompted to enroll required authenticators before their grace period ends.
	GracePeriodRequiredSoonCustomLinkUrl string `json:"gracePeriodRequiredSoonCustomLinkUrl,omitempty"`
}
