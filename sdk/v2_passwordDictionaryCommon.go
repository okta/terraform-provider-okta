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
