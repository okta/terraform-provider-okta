// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CustomTelephonyProviderSettingsTwilioMessagingService represents the CustomTelephonyProviderSettingsTwilioMessagingService schema
type CustomTelephonyProviderSettingsTwilioMessagingService struct {
	// The Twilio Messaging Service SID used for sending SMS messages. You can find this value in your Twilio console.  This method uses Twilio's [Programmable Messaging API](https://www.twilio.com/docs/m...
	TwilioMessageSid string `json:"twilioMessageSid,omitempty"`
}
