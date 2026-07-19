// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CustomTelephonyProviderSettingsTwilioVerify represents the CustomTelephonyProviderSettingsTwilioVerify schema
type CustomTelephonyProviderSettingsTwilioVerify struct {
	// The Twilio Verify Service SID used for sending verification messages or calls. You can find this value in your Twilio console.  This method uses Twilio's [Verify API](https://www.twilio.com/docs/ve...
	TwilioVerifySid string `json:"twilioVerifySid,omitempty"`
}
