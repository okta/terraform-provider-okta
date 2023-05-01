package sdk

import "encoding/json"

type Duration struct {
	Number    int64  `json:"-"`
	NumberPtr *int64 `json:"number,omitempty"`
	Unit      string `json:"unit,omitempty"`
}

func NewDuration() *Duration {
	return &Duration{}
}

func (a *Duration) IsPolicyInstance() bool {
	return true
}

func (a *Duration) MarshalJSON() ([]byte, error) {
	type Alias Duration
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.Number != 0 {
		result.NumberPtr = Int64Ptr(a.Number)
	}
	return json.Marshal(&result)
}

func (a *Duration) UnmarshalJSON(data []byte) error {
	type Alias Duration

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
