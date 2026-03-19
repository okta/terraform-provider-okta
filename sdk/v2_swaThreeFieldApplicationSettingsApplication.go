// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type SwaThreeFieldApplicationSettingsApplication struct {
	ButtonSelector     string `json:"buttonSelector,omitempty"`
	ExtraFieldSelector string `json:"extraFieldSelector,omitempty"`
	ExtraFieldValue    string `json:"extraFieldValue,omitempty"`
	LoginUrlRegex      string `json:"loginUrlRegex,omitempty"`
	PasswordSelector   string `json:"passwordSelector,omitempty"`
	TargetURL          string `json:"targetURL,omitempty"`
	UserNameSelector   string `json:"userNameSelector,omitempty"`
}
