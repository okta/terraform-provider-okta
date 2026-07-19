// Code generated from OpenAPI spec. DO NOT EDIT.
package models

// CustomTelephonyProviderSettings represents the CustomTelephonyProviderSettings schema
// Settings for custom telephony provider.  These settings vary based on the telephony provider and the type of telephony operation (SMS or Voice). For `sms` and `call`, you can select one method per ...
type CustomTelephonyProviderSettings struct {
	// Method for making voice calls. Choose one method for making voice calls based on your telephony provider setup.
	Call interface{} `json:"call,omitempty"`
	// Method for sending SMS messages. Choose one method for sending SMS messages based on your telephony provider setup.
	Sms interface{} `json:"sms,omitempty"`
}
