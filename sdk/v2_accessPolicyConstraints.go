// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import "encoding/json"

type AccessPolicyConstraints struct {
	Knowledge            *KnowledgeConstraint  `json:"knowledge,omitempty"`
	Possession           *PossessionConstraint `json:"possession,omitempty"`
	AdditionalProperties map[string]any        `json:"-"`
}

// UnmarshalJSON captures all fields, including any not modelled in the struct,
// into AdditionalProperties so they survive a round-trip through the SDK.
func (a *AccessPolicyConstraints) UnmarshalJSON(data []byte) error {
	type plain AccessPolicyConstraints
	var known plain
	if err := json.Unmarshal(data, &known); err != nil {
		return err
	}
	a.Knowledge = known.Knowledge
	a.Possession = known.Possession

	var all map[string]json.RawMessage
	if err := json.Unmarshal(data, &all); err != nil {
		return err
	}
	delete(all, "knowledge")
	delete(all, "possession")
	if len(all) > 0 {
		a.AdditionalProperties = make(map[string]any, len(all))
		for k, v := range all {
			var val any
			if err := json.Unmarshal(v, &val); err != nil {
				return err
			}
			a.AdditionalProperties[k] = val
		}
	}
	return nil
}

// MarshalJSON emits known fields plus any extra fields stored in AdditionalProperties.
func (a AccessPolicyConstraints) MarshalJSON() ([]byte, error) {
	type plain AccessPolicyConstraints
	b, err := json.Marshal(plain(a))
	if err != nil {
		return nil, err
	}
	if len(a.AdditionalProperties) == 0 {
		return b, nil
	}
	var merged map[string]any
	if err := json.Unmarshal(b, &merged); err != nil {
		return nil, err
	}
	for k, v := range a.AdditionalProperties {
		merged[k] = v
	}
	return json.Marshal(merged)
}

func NewAccessPolicyConstraints() *AccessPolicyConstraints {
	return &AccessPolicyConstraints{}
}

func (a *AccessPolicyConstraints) IsPolicyInstance() bool {
	return true
}
