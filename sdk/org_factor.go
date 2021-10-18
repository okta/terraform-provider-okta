package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type OrgFactor struct {
	Id         string `json:"id"`
	Provider   string `json:"provider"`
	FactorType string `json:"factorType"`
	Status     string `json:"status"`
}

// Current available org factors for MFA
const (
	DuoFactor          = "duo"
	FidoU2fFactor      = "fido_u2f"
	FidoWebauthnFactor = "fido_webauthn"
	GoogleOtpFactor    = "google_otp"
	OktaCallFactor     = "okta_call"
	OktaOtpFactor      = "okta_otp"
	OktaPasswordFactor = "okta_password"
	OktaPushFactor     = "okta_push"
	OktaQuestionFactor = "okta_question"
	OktaSmsFactor      = "okta_sms"
	OktaEmailFactor    = "okta_email"
	RsaTokenFactor     = "rsa_token"
	SymantecVipFactor  = "symantec_vip"
	YubikeyTokenFactor = "yubikey_token"
	HotpFactor         = "hotp"
)

// GetOrgFactor gets a factor by ID.
func (m *APISupplement) GetOrgFactor(ctx context.Context, id string) (*OrgFactor, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/org/factors/%s", id)
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var factor *OrgFactor
	resp, err := m.RequestExecutor.Do(ctx, req, &factor)
	if err != nil {
		return nil, resp, err
	}
	return factor, resp, nil
}

// ActivateOrgFactor allows multifactor authentication to use provided factor type
func (m *APISupplement) ActivateOrgFactor(ctx context.Context, id string) (*OrgFactor, *okta.Response, error) {
	return m.lifecycleChangeOrgFactor(ctx, id, "activate")
}

// DeactivateOrgFactor denies multifactor authentication to use provided factor type
func (m *APISupplement) DeactivateOrgFactor(ctx context.Context, id string) (*OrgFactor, *okta.Response, error) {
	return m.lifecycleChangeOrgFactor(ctx, id, "deactivate")
}

func (m *APISupplement) lifecycleChangeOrgFactor(ctx context.Context, id, action string) (*OrgFactor, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/org/factors/%s/lifecycle/%s", id, action)
	req, err := m.RequestExecutor.
		WithAccept("application/json").
		WithContentType("application/json").
		NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var factor *OrgFactor
	resp, err := m.RequestExecutor.Do(ctx, req, factor)
	if err != nil {
		return nil, resp, err
	}
	return factor, resp, nil
}
