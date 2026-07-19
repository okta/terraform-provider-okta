// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CustomTelephonyProviderSettingsTelesignServiceSms represents the CustomTelephonyProviderSettingsTelesignServiceSms schema
type CustomTelephonyProviderSettingsTelesignServiceSms struct {
	// The Telesign service identifier used for sending SMS messages. You can find this value in your Telesign console.  The `telesignVerifyService` method uses Telesign's [Verify API](https://developer.t...
	TelesignService string `json:"telesignService,omitempty"`
}
