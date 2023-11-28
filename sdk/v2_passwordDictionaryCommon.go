// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PasswordDictionaryCommon struct {
	Exclude *bool `json:"exclude,omitempty"`
}

func NewPasswordDictionaryCommon() *PasswordDictionaryCommon {
	return &PasswordDictionaryCommon{
		Exclude: boolPtr(false),
	}
}

func (a *PasswordDictionaryCommon) IsPolicyInstance() bool {
	return true
}
