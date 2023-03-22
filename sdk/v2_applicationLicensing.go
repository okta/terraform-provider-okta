package sdk

import "encoding/json"

type ApplicationLicensing struct {
	SeatCount    int64  `json:"-"`
	SeatCountPtr *int64 `json:"seatCount,omitempty"`
}

func (a *ApplicationLicensing) MarshalJSON() ([]byte, error) {
	type Alias ApplicationLicensing
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.SeatCount != 0 {
		result.SeatCountPtr = Int64Ptr(a.SeatCount)
	}
	return json.Marshal(&result)
}

func (a *ApplicationLicensing) UnmarshalJSON(data []byte) error {
	type Alias ApplicationLicensing

	result := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.SeatCountPtr != nil {
		a.SeatCount = *result.SeatCountPtr
		a.SeatCountPtr = result.SeatCountPtr
	}
	return nil
}
