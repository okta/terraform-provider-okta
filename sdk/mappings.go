// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"context"
	"fmt"
	"net/http"
)

// FIXME uses internal api
func (m *APISupplement) ApplyMappings(ctx context.Context, sourceID, targetID string) (*Response, error) {
	url := fmt.Sprintf("/api/internal/v1/mappings/reapply?source=%s&target=%s", sourceID, targetID)
	re := m.cloneRequestExecutor()
	req, err := re.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return nil, err
	}
	return re.Do(ctx, req, nil)
}
