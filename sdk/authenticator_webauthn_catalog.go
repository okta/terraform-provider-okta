package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WebauthnAuthenticator struct {
	AAGUID                       string               `json:"aaguid"`
	ModelName                    string               `json:"modelName"`
	AuthenticatorCharactoristics AuthnCharactoristics `json:"authenticatorCharacteristics"`
}

type AuthnCharactoristics struct {
	PlatformAttached          bool `json:"platformAttached"`
	FipsCompliant             bool `json:"fipsCompliant"`
	HardwareProtected         bool `json:"hardwareProtected"`
	UserVerificationSupported bool `json:"userVerificationSupported"`
}

func (m *APISupplement) ListWebauthnCatalog(ctx context.Context, authn, orgName, domain, token string, client *http.Client) ([]*WebauthnAuthenticator, error) {
	url := fmt.Sprintf("https://%s-admin.%s/api/internal/authenticators/%s/catalog", orgName, domain, authn)
	buff := new(bytes.Buffer)
	encoder := json.NewEncoder(buff)
	encoder.SetEscapeHTML(false)
	req, err := http.NewRequest(http.MethodGet, url, buff)
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
	var waa []*WebauthnAuthenticator
	err = json.Unmarshal(respBody, &waa)
	if err != nil {
		return nil, err
	}
	return waa, nil
}
