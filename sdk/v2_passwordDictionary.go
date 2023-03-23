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
