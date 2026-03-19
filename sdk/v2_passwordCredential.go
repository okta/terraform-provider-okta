// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

type PasswordCredential struct {
	Hash  *PasswordCredentialHash `json:"hash,omitempty"`
	Hook  *PasswordCredentialHook `json:"hook,omitempty"`
	Value string                  `json:"value,omitempty"`
}
