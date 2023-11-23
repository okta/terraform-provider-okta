package sdk

import (
	"context"
	"fmt"
	"github.com/okta/terraform-provider-okta/sdk/query"
	"time"
)

type LogStreamResource resource

type LogStream struct {
	Id          string             `json:"id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Type        string             `json:"type,omitempty"`
	Settings    *LogStreamSettings `json:"settings,omitempty"`
	Links       interface{}        `json:"_links,omitempty"`
	Created     *time.Time         `json:"created,omitempty"`
	LastUpdated *time.Time         `json:"lastUpdated,omitempty"`
	Status      string             `json:"status,omitempty"`
}

// Fetches a log stream from your Okta organization by &#x60;id&#x60;.
func (m *LogStreamResource) GetLogStream(ctx context.Context, logStreamId string) (*LogStream, *Response, error) {
	url := fmt.Sprintf("/api/v1/logStreams/%v", logStreamId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var logStream *LogStream

	resp, err := rq.Do(ctx, req, &logStream)
	if err != nil {
		return nil, resp, err
	}

	return logStream, resp, nil
}

// Updates a log stream in your organization.
func (m *LogStreamResource) UpdateLogStream(ctx context.Context, logStreamId string, body LogStream) (*LogStream, *Response, error) {
	url := fmt.Sprintf("/api/v1/logStreams/%v", logStreamId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var LogStream *LogStream

	resp, err := rq.Do(ctx, req, &LogStream)
	if err != nil {
		return nil, resp, err
	}

	return LogStream, resp, nil
}

// Removes log stream.
func (m *LogStreamResource) DeleteLogStream(ctx context.Context, logStreamId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/logStreams/%v", logStreamId)

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

// Enumerates log streams added to your organization with pagination. A subset of log streams can be returned that match a supported filter expression or query.
func (m *LogStreamResource) ListLogStreams(ctx context.Context, qp *query.Params) ([]*LogStream, *Response, error) {
	url := "/api/v1/logStreams"
	if qp != nil {
		url = url + qp.String()
	}

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var LogStream []*LogStream

	resp, err := rq.Do(ctx, req, &LogStream)
	if err != nil {
		return nil, resp, err
	}

	return LogStream, resp, nil
}

// Adds a new log stream to your Okta organization.
func (m *LogStreamResource) CreateLogStream(ctx context.Context, body LogStream) (*LogStream, *Response, error) {
	url := "/api/v1/logStreams"

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var LogStream *LogStream

	resp, err := rq.Do(ctx, req, &LogStream)
	if err != nil {
		return nil, resp, err
	}

	return LogStream, resp, nil
}

// Activate log stream
func (m *LogStreamResource) ActivateLogStream(ctx context.Context, logStreamId string) (*LogStream, *Response, error) {
	url := fmt.Sprintf("/api/v1/logStreams/%v/lifecycle/activate", logStreamId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var LogStream *LogStream

	resp, err := rq.Do(ctx, req, &LogStream)
	if err != nil {
		return nil, resp, err
	}

	return LogStream, resp, nil
}

// Deactivates a log stream.
func (m *LogStreamResource) DeactivateLogStream(ctx context.Context, logStreamId string) (*LogStream, *Response, error) {
	url := fmt.Sprintf("/api/v1/logStreams/%v/lifecycle/deactivate", logStreamId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var LogStream *LogStream

	resp, err := rq.Do(ctx, req, &LogStream)
	if err != nil {
		return nil, resp, err
	}

	return LogStream, resp, nil
}
