// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CustomTelephonyProviderSettingsTelesignServiceCall represents the CustomTelephonyProviderSettingsTelesignServiceCall schema
type CustomTelephonyProviderSettingsTelesignServiceCall struct {
	// The Telesign service identifier used for sending making calls. You can find this value in your Telesign console.  The `telesignVerifyService` method uses Telesign's [Verify API](https://developer.t...
	TelesignService string `json:"telesignService,omitempty"`
}
