package sdk

import (
	"encoding/json"
	"time"
)

type AccessPolicyRule struct {
	Actions     *AccessPolicyRuleActions    `json:"actions,omitempty"`
	Conditions  *AccessPolicyRuleConditions `json:"conditions,omitempty"`
	Created     *time.Time                  `json:"created,omitempty"`
	Id          string                      `json:"id,omitempty"`
	LastUpdated *time.Time                  `json:"lastUpdated,omitempty"`
	Name        string                      `json:"name,omitempty"`
	Priority    int64                       `json:"-"`
	PriorityPtr *int64                      `json:"priority,omitempty"`
	Status      string                      `json:"status,omitempty"`
	System      *bool                       `json:"system,omitempty"`
	Type        string                      `json:"type,omitempty"`
}

func NewAccessPolicyRule() *AccessPolicyRule {
	return &AccessPolicyRule{
		Status: "ACTIVE",
		System: boolPtr(false),
		Type:   "ACCESS_POLICY",
	}
}

func (a *AccessPolicyRule) IsPolicyInstance() bool {
	return true
}

func (a *AccessPolicyRule) MarshalJSON() ([]byte, error) {
	type Alias AccessPolicyRule
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.Priority != 0 {
		result.PriorityPtr = Int64Ptr(a.Priority)
	}
	return json.Marshal(&result)
}

func (a *AccessPolicyRule) UnmarshalJSON(data []byte) error {
	type Alias AccessPolicyRule

	result := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.PriorityPtr != nil {
		a.Priority = *result.PriorityPtr
		a.PriorityPtr = result.PriorityPtr
	}
	return nil
}
