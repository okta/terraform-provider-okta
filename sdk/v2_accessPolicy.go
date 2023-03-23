package sdk

import (
	"encoding/json"
	"time"
)

type AccessPolicy struct {
	Embedded    interface{}           `json:"_embedded,omitempty"`
	Links       interface{}           `json:"_links,omitempty"`
	Conditions  *PolicyRuleConditions `json:"conditions,omitempty"`
	Created     *time.Time            `json:"created,omitempty"`
	Description string                `json:"description,omitempty"`
	Id          string                `json:"id,omitempty"`
	LastUpdated *time.Time            `json:"lastUpdated,omitempty"`
	Name        string                `json:"name,omitempty"`
	Priority    int64                 `json:"-"`
	PriorityPtr *int64                `json:"priority,omitempty"`
	Status      string                `json:"status,omitempty"`
	System      *bool                 `json:"system,omitempty"`
	Type        string                `json:"type,omitempty"`
}

func NewAccessPolicy() *AccessPolicy {
	return &AccessPolicy{
		Type: "ACCESS_POLICY",
	}
}

func (a *AccessPolicy) IsPolicyInstance() bool {
	return true
}

func (a *AccessPolicy) MarshalJSON() ([]byte, error) {
	type Alias AccessPolicy
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.Priority != 0 {
		result.PriorityPtr = Int64Ptr(a.Priority)
	}
	return json.Marshal(&result)
}

func (a *AccessPolicy) UnmarshalJSON(data []byte) error {
	type Alias AccessPolicy

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
