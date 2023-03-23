package sdk

import "encoding/json"

type PasswordPolicyPasswordSettingsLockout struct {
	AutoUnlockMinutes               int64    `json:"-"`
	AutoUnlockMinutesPtr            *int64   `json:"autoUnlockMinutes,omitempty"`
	MaxAttempts                     int64    `json:"-"`
	MaxAttemptsPtr                  *int64   `json:"maxAttempts,omitempty"`
	ShowLockoutFailures             *bool    `json:"showLockoutFailures,omitempty"`
	UserLockoutNotificationChannels []string `json:"userLockoutNotificationChannels,omitempty"`
}

func NewPasswordPolicyPasswordSettingsLockout() *PasswordPolicyPasswordSettingsLockout {
	return &PasswordPolicyPasswordSettingsLockout{}
}

func (a *PasswordPolicyPasswordSettingsLockout) IsPolicyInstance() bool {
	return true
}

func (a *PasswordPolicyPasswordSettingsLockout) MarshalJSON() ([]byte, error) {
	type Alias PasswordPolicyPasswordSettingsLockout
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.AutoUnlockMinutes != 0 {
		result.AutoUnlockMinutesPtr = Int64Ptr(a.AutoUnlockMinutes)
	}
	if a.MaxAttempts != 0 {
		result.MaxAttemptsPtr = Int64Ptr(a.MaxAttempts)
	}
	return json.Marshal(&result)
}

func (a *PasswordPolicyPasswordSettingsLockout) UnmarshalJSON(data []byte) error {
	type Alias PasswordPolicyPasswordSettingsLockout

	result := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.AutoUnlockMinutesPtr != nil {
		a.AutoUnlockMinutes = *result.AutoUnlockMinutesPtr
		a.AutoUnlockMinutesPtr = result.AutoUnlockMinutesPtr
	}
	if result.MaxAttemptsPtr != nil {
		a.MaxAttempts = *result.MaxAttemptsPtr
		a.MaxAttemptsPtr = result.MaxAttemptsPtr
	}
	return nil
}
