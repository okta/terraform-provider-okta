// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CustomTelephonyProviderSettingsTwilioCallerId represents the CustomTelephonyProviderSettingsTwilioCallerId schema
type CustomTelephonyProviderSettingsTwilioCallerId struct {
	// The Twilio caller ID that's used for making calls. You can find this value in your Twilio console.  This method uses Twilio's [Programmable Voice API](https://www.twilio.com/docs/voice).
	TwilioCallerId string `json:"twilioCallerId,omitempty"`
}
