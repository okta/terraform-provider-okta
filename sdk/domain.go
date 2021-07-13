package sdk

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type Domain struct {
	ID                    string `json:"id,omitempty"`
	Domain                string `json:"domain"`
	CertificateSourceType string `json:"certificateSourceType"`
	ValidationStatus      string `json:"validationStatus,omitempty"`
	DNSRecords            []struct {
		Expiration string   `json:"expiration",omitempty"`
		Fqdn       string   `json:"fqdn,omitempty"`
		Values     []string `json:"values,omitempty"`
		RecordType string   `json:"recordType,omitempty"`
	} `json:"dnsRecords,omitempty"`
	PublicCertificate struct {
		Subject     string    `json:"subject,omitempty"`
		Fingerprint string    `json:"fingerprint,omitempty"`
		Expiration  time.Time `json:"expiration,omitempty"`
	} `json:"publicCertificate,omitempty"`
}

type Certificate struct {
	Type             string `json:"type"`
	PrivateKey       string `json:"privateKey"`
	Certificate      string `json:"certificate"`
	CertificateChain string `json:"certificateChain,omitempty"`
}

func (m *ApiSupplement) CreateDomain(ctx context.Context, body Domain) (*Domain, *okta.Response, error) {
	url := "/api/v1/domains"
	req, err := m.RequestExecutor.
		WithAccept("application/json").
		WithContentType("application/json").
		NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, err
	}
	var domain Domain
	resp, err := m.RequestExecutor.Do(ctx, req, &domain)
	if err != nil {
		return nil, resp, err
	}
	return &domain, resp, nil
}

func (m *ApiSupplement) VerifyDomain(ctx context.Context, id string) (*Domain, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/domains/%s/verify", id)
	req, err := m.RequestExecutor.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var domain Domain
	resp, err := m.RequestExecutor.Do(ctx, req, &domain)
	if err != nil {
		return nil, resp, err
	}
	return &domain, resp, nil
}

func (m *ApiSupplement) GetDomain(ctx context.Context, id string) (*Domain, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/domains/%v", id)
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var domain Domain
	resp, err := m.RequestExecutor.Do(ctx, req, &domain)
	if err != nil {
		return nil, resp, err
	}
	return &domain, resp, nil
}

func (m *ApiSupplement) DeleteDomain(ctx context.Context, id string) (*okta.Response, error) {
	url := fmt.Sprintf("/api/v1/domains/%v", id)
	req, err := m.RequestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := m.RequestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
