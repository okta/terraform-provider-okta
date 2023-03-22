package sdk

import "encoding/json"

type InactivityPolicyRuleCondition struct {
	Number    int64  `json:"-"`
	NumberPtr *int64 `json:"number,omitempty"`
	Unit      string `json:"unit,omitempty"`
}

func NewInactivityPolicyRuleCondition() *InactivityPolicyRuleCondition {
	return &InactivityPolicyRuleCondition{}
}

func (a *InactivityPolicyRuleCondition) IsPolicyInstance() bool {
	return true
}

func (a *InactivityPolicyRuleCondition) MarshalJSON() ([]byte, error) {
	type Alias InactivityPolicyRuleCondition
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.Number != 0 {
		result.NumberPtr = Int64Ptr(a.Number)
	}
	return json.Marshal(&result)
}

func (a *InactivityPolicyRuleCondition) UnmarshalJSON(data []byte) error {
	type Alias InactivityPolicyRuleCondition

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
