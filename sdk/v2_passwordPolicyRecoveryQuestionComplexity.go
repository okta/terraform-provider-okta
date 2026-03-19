// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import "encoding/json"

type PasswordPolicyRecoveryQuestionComplexity struct {
	MinLength    int64  `json:"-"`
	MinLengthPtr *int64 `json:"minLength,omitempty"`
}

func NewPasswordPolicyRecoveryQuestionComplexity() *PasswordPolicyRecoveryQuestionComplexity {
	return &PasswordPolicyRecoveryQuestionComplexity{}
}

func (a *PasswordPolicyRecoveryQuestionComplexity) IsPolicyInstance() bool {
	return true
}

func (a *PasswordPolicyRecoveryQuestionComplexity) MarshalJSON() ([]byte, error) {
	type Alias PasswordPolicyRecoveryQuestionComplexity
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.MinLength != 0 {
		result.MinLengthPtr = Int64Ptr(a.MinLength)
	}
	return json.Marshal(&result)
}

func (a *PasswordPolicyRecoveryQuestionComplexity) UnmarshalJSON(data []byte) error {
	type Alias PasswordPolicyRecoveryQuestionComplexity

	result := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.MinLengthPtr != nil {
		a.MinLength = *result.MinLengthPtr
		a.MinLengthPtr = result.MinLengthPtr
	}
	return nil
}
