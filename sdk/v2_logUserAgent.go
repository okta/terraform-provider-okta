// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type LogUserAgent struct {
	Browser      string `json:"browser,omitempty"`
	Os           string `json:"os,omitempty"`
	RawUserAgent string `json:"rawUserAgent,omitempty"`
}
