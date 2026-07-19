// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CustomTelephonyProviderSettingsTwilioPhoneNumber represents the CustomTelephonyProviderSettingsTwilioPhoneNumber schema
type CustomTelephonyProviderSettingsTwilioPhoneNumber struct {
	// The Twilio phone number that's used for sending SMS messages or voice calls. You can find this value in your Twilio console.  This method uses Twilio's [Programmable Messaging API](https://www.twil...
	TwilioPhoneNumber string `json:"twilioPhoneNumber,omitempty"`
}
