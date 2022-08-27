package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SecurityNotificationEmails struct {
	SendEmailForNewDeviceEnabled        bool `json:"sendEmailForNewDeviceEnabled"`
	SendEmailForFactorEnrollmentEnabled bool `json:"sendEmailForFactorEnrollmentEnabled"`
	SendEmailForFactorResetEnabled      bool `json:"sendEmailForFactorResetEnabled"`
	SendEmailForPasswordChangedEnabled  bool `json:"sendEmailForPasswordChangedEnabled"`
	ReportSuspiciousActivityEnabled     bool `json:"reportSuspiciousActivityEnabled"`
}

func (m *APISupplement) UpdateSecurityNotificationEmails(ctx context.Context, body SecurityNotificationEmails, orgName, domain, token string, client *http.Client) (*SecurityNotificationEmails, error) {
	url := fmt.Sprintf("https://%s-admin.%s/api/internal/org/settings/security-notification-settings", orgName, domain)
	buff := new(bytes.Buffer)
	encoder := json.NewEncoder(buff)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, url, buff)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "SSWS "+token)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode > http.StatusNoContent {
		return nil, fmt.Errorf("API returned HTTP status %d, err: %s", res.StatusCode, string(respBody))
	}
	var threatInsightConfiguration SecurityNotificationEmails
	err = json.Unmarshal(respBody, &threatInsightConfiguration)
	if err != nil {
		return nil, err
	}
	return &threatInsightConfiguration, nil
}

func (m *APISupplement) GetSecurityNotificationEmails(ctx context.Context, orgName, domain, token string, client *http.Client) (*SecurityNotificationEmails, error) {
	url := fmt.Sprintf("https://%s-admin.%s/api/internal/org/settings/security-notification-settings", orgName, domain)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "SSWS "+token)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode > http.StatusNoContent {
		return nil, fmt.Errorf("API returned HTTP status %d, err: %s", res.StatusCode, string(respBody))
	}
	var threatInsightConfiguration SecurityNotificationEmails
	err = json.Unmarshal(respBody, &threatInsightConfiguration)
	if err != nil {
		return nil, err
	}
	return &threatInsightConfiguration, nil
}
