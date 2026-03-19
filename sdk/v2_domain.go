// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"context"
	"fmt"
)

type DomainResource resource

type Domain struct {
	CertificateSourceType string                     `json:"certificateSourceType,omitempty"`
	DnsRecords            []*DNSRecord               `json:"dnsRecords,omitempty"`
	Domain                string                     `json:"domain,omitempty"`
	Id                    string                     `json:"id,omitempty"`
	PublicCertificate     *DomainCertificateMetadata `json:"publicCertificate,omitempty"`
	ValidationStatus      string                     `json:"validationStatus,omitempty"`
}

// List all verified custom Domains for the org.
func (m *DomainResource) ListDomains(ctx context.Context) (*DomainListResponse, *Response, error) {
	url := "/api/v1/domains"

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var domainListResponse *DomainListResponse

	resp, err := rq.Do(ctx, req, &domainListResponse)
	if err != nil {
		return nil, resp, err
	}

	return domainListResponse, resp, nil
}

// Creates your domain.
func (m *DomainResource) CreateDomain(ctx context.Context, body Domain) (*Domain, *Response, error) {
	url := "/api/v1/domains"

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var domain *Domain

	resp, err := rq.Do(ctx, req, &domain)
	if err != nil {
		return nil, resp, err
	}

	return domain, resp, nil
}

// Deletes a Domain by &#x60;id&#x60;.
func (m *DomainResource) DeleteDomain(ctx context.Context, domainId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/domains/%v", domainId)

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

// Fetches a Domain by &#x60;id&#x60;.
func (m *DomainResource) GetDomain(ctx context.Context, domainId string) (*Domain, *Response, error) {
	url := fmt.Sprintf("/api/v1/domains/%v", domainId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var domain *Domain

	resp, err := rq.Do(ctx, req, &domain)
	if err != nil {
		return nil, resp, err
	}

	return domain, resp, nil
}

// Creates the Certificate for the Domain.
func (m *DomainResource) CreateCertificate(ctx context.Context, domainId string, body DomainCertificate) (*Response, error) {
	url := fmt.Sprintf("/api/v1/domains/%v/certificate", domainId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Verifies the Domain by &#x60;id&#x60;.
func (m *DomainResource) VerifyDomain(ctx context.Context, domainId string) (*Domain, *Response, error) {
	url := fmt.Sprintf("/api/v1/domains/%v/verify", domainId)

	rq := m.client.CloneRequestExecutor()

	req, err := rq.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var domain *Domain

	resp, err := rq.Do(ctx, req, &domain)
	if err != nil {
		return nil, resp, err
	}

	return domain, resp, nil
}
