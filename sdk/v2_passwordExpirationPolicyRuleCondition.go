// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import "encoding/json"

type PasswordExpirationPolicyRuleCondition struct {
	Number    int64  `json:"-"`
	NumberPtr *int64 `json:"number,omitempty"`
	Unit      string `json:"unit,omitempty"`
}

func NewPasswordExpirationPolicyRuleCondition() *PasswordExpirationPolicyRuleCondition {
	return &PasswordExpirationPolicyRuleCondition{}
}

func (a *PasswordExpirationPolicyRuleCondition) IsPolicyInstance() bool {
	return true
}

func (a *PasswordExpirationPolicyRuleCondition) MarshalJSON() ([]byte, error) {
	type Alias PasswordExpirationPolicyRuleCondition
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.Number != 0 {
		result.NumberPtr = Int64Ptr(a.Number)
	}
	return json.Marshal(&result)
}

func (a *PasswordExpirationPolicyRuleCondition) UnmarshalJSON(data []byte) error {
	type Alias PasswordExpirationPolicyRuleCondition

	result := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.NumberPtr != nil {
		a.Number = *result.NumberPtr
		a.NumberPtr = result.NumberPtr
	}
	return nil
}
