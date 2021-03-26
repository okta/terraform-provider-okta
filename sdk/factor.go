package sdk

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type Factor struct {
	Id         string `json:"id"`
	Provider   string `json:"provider"`
	FactorType string `json:"factorType"`
	Status     string `json:"status"`
}

// Current available factors for MFA
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

// GetFactor gets a factor by ID.
func (m *ApiSupplement) GetFactor(ctx context.Context, id string) (*Factor, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/org/factors/%s", id)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	var factor Factor
	resp, err := m.RequestExecutor.Do(ctx, req, &factor)
	if err != nil {
		return nil, resp, err
	}
	return &factor, resp, nil
}

// ActivateFactor allows multifactor authentication to use provided factor type
func (m *ApiSupplement) ActivateFactor(ctx context.Context, id string) (*Factor, *okta.Response, error) {
	return m.lifecycleChangeFactor(ctx, id, "activate")
}

// ActivateFactor denies multifactor authentication to use provided factor type
func (m *ApiSupplement) DeactivateFactor(ctx context.Context, id string) (*Factor, *okta.Response, error) {
	return m.lifecycleChangeFactor(ctx, id, "deactivate")
}

func (m *ApiSupplement) lifecycleChangeFactor(ctx context.Context, id, action string) (*Factor, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/org/factors/%s/lifecycle/%s", id, action)
	req, err := m.RequestExecutor.
		WithAccept("application/json").
		WithContentType("application/json").
		NewRequest("POST", url, nil)
	if err != nil {
		return nil, nil, err
	}
	var factor *Factor
	resp, err := m.RequestExecutor.Do(ctx, req, factor)
	if err != nil {
		return nil, resp, err
	}
	return factor, resp, nil
}
