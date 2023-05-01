package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type ApplicationGroupAssignmentResource resource

type ApplicationGroupAssignment struct {
	Embedded    interface{} `json:"_embedded,omitempty"`
	Links       interface{} `json:"_links,omitempty"`
	Id          string      `json:"id,omitempty"`
	LastUpdated *time.Time  `json:"lastUpdated,omitempty"`
	Priority    int64       `json:"-"`
	PriorityPtr *int64      `json:"priority,omitempty"`
	Profile     interface{} `json:"profile,omitempty"`
}

// Removes a group assignment from an application.
func (m *ApplicationGroupAssignmentResource) DeleteApplicationGroupAssignment(ctx context.Context, appId, groupId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/groups/%v", appId, groupId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (a *ApplicationGroupAssignment) MarshalJSON() ([]byte, error) {
	type Alias ApplicationGroupAssignment
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.Priority != 0 {
		result.PriorityPtr = Int64Ptr(a.Priority)
	}
	return json.Marshal(&result)
}

func (a *ApplicationGroupAssignment) UnmarshalJSON(data []byte) error {
	type Alias ApplicationGroupAssignment

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
