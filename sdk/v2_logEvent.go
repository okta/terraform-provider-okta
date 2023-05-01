package sdk

import (
	"context"
	"time"

	"github.com/okta/terraform-provider-okta/sdk/query"
)

type LogEventResource resource

type LogEvent struct {
	Actor                 *LogActor                 `json:"actor,omitempty"`
	AuthenticationContext *LogAuthenticationContext `json:"authenticationContext,omitempty"`
	Client                *LogClient                `json:"client,omitempty"`
	DebugContext          *LogDebugContext          `json:"debugContext,omitempty"`
	DisplayMessage        string                    `json:"displayMessage,omitempty"`
	EventType             string                    `json:"eventType,omitempty"`
	LegacyEventType       string                    `json:"legacyEventType,omitempty"`
	Outcome               *LogOutcome               `json:"outcome,omitempty"`
	Published             *time.Time                `json:"published,omitempty"`
	Request               *LogRequest               `json:"request,omitempty"`
	SecurityContext       *LogSecurityContext       `json:"securityContext,omitempty"`
	Severity              string                    `json:"severity,omitempty"`
	Target                []*LogTarget              `json:"target,omitempty"`
	Transaction           *LogTransaction           `json:"transaction,omitempty"`
	Uuid                  string                    `json:"uuid,omitempty"`
	Version               string                    `json:"version,omitempty"`
}

// The Okta System Log API provides read access to your organization’s system log. This API provides more functionality than the Events API
func (m *LogEventResource) GetLogs(ctx context.Context, qp *query.Params) ([]*LogEvent, *Response, error) {
	url := "/api/v1/logs"
	if qp != nil {
		url = url + qp.String()
	}

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var logEvent []*LogEvent

	resp, err := rq.Do(ctx, req, &logEvent)
	if err != nil {
		return nil, resp, err
	}

	return logEvent, resp, nil
}
