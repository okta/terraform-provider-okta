package sdk

type PasswordCredential struct {
	Hash  *PasswordCredentialHash `json:"hash,omitempty"`
	Hook  *PasswordCredentialHook `json:"hook,omitempty"`
	Value string                  `json:"value,omitempty"`
}
