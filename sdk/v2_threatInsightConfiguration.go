package sdk

import (
	"context"
	"fmt"
	"time"
)

type ThreatInsightConfigurationResource resource

type ThreatInsightConfiguration struct {
	Links        interface{} `json:"_links,omitempty"`
	Action       string      `json:"action,omitempty"`
	Created      *time.Time  `json:"created,omitempty"`
	ExcludeZones []string    `json:"excludeZones,omitempty"`
	LastUpdated  *time.Time  `json:"lastUpdated,omitempty"`
}

// Gets current ThreatInsight configuration
func (m *ThreatInsightConfigurationResource) GetCurrentConfiguration(ctx context.Context) (*ThreatInsightConfiguration, *Response, error) {
	url := fmt.Sprintf("/api/v1/threats/configuration")

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var threatInsightConfiguration *ThreatInsightConfiguration

	resp, err := rq.Do(ctx, req, &threatInsightConfiguration)
	if err != nil {
		return nil, resp, err
	}

	return threatInsightConfiguration, resp, nil
}

// Updates ThreatInsight configuration
func (m *ThreatInsightConfigurationResource) UpdateConfiguration(ctx context.Context, body ThreatInsightConfiguration) (*ThreatInsightConfiguration, *Response, error) {
	url := fmt.Sprintf("/api/v1/threats/configuration")

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var threatInsightConfiguration *ThreatInsightConfiguration

	resp, err := rq.Do(ctx, req, &threatInsightConfiguration)
	if err != nil {
		return nil, resp, err
	}

	return threatInsightConfiguration, resp, nil
}
