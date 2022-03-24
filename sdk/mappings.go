package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

// FIXME uses internal api
func (m *APISupplement) ApplyMappings(ctx context.Context, sourceID, targetID string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/internal/v1/mappings/reapply?source=%s&target=%s", sourceID, targetID)
	re := m.cloneRequestExecutor()
	req, err := re.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return nil, err
	}
	return re.Do(ctx, req, nil)
}
