// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PasswordDictionary struct {
	Common *PasswordDictionaryCommon `json:"common,omitempty"`
}

func NewPasswordDictionary() *PasswordDictionary {
	return &PasswordDictionary{}
}

func (a *PasswordDictionary) IsPolicyInstance() bool {
	return true
}
