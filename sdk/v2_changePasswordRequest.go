package sdk

type ChangePasswordRequest struct {
	NewPassword *PasswordCredential `json:"newPassword,omitempty"`
	OldPassword *PasswordCredential `json:"oldPassword,omitempty"`
}
