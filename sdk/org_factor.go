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
	DuoFactor              = "duo"
	ExternalIdpFactor      = "external_idp"
	FidoU2fFactor          = "fido_u2f"
	FidoWebauthnFactor     = "fido_webauthn"
	GoogleOtpFactor        = "google_otp"
	HotpFactor             = "hotp"
	OktaCallFactor         = "okta_call"
	OktaEmailFactor        = "okta_email"
	OktaOtpFactor          = "okta_otp"
	OktaPasswordFactor     = "okta_password" // Not configurable for OIE
	OktaPushFactor         = "okta_push"
	OktaQuestionFactor     = "okta_question"
	OktaSmsFactor          = "okta_sms"
	OktaVerifyFactor       = "okta_verify"  // OIE only (Combo of OktaOtp, OktaPush)
	OnPremMfaFactor        = "onprem_mfa"   // OIE only
	PhoneNumberFactor      = "phone_number" // OIE only (Combo of OktaSms + OktaCall)
	RsaTokenFactor         = "rsa_token"
	SecurityQuestionFactor = "security_question" // OIE only (Evolution/rename from okta_question)
	SymantecVipFactor      = "symantec_vip"
	WebauthnFactor         = "webauthn" // OIE only (Evolution/rename from fido_webauthn)
	YubikeyTokenFactor     = "yubikey_token"
)

// List of factors that are applicable to Okta Classic Engine
var FactorProviders = []string{
	DuoFactor,
	FidoU2fFactor,
	FidoWebauthnFactor,
	HotpFactor,
	GoogleOtpFactor,
	OktaCallFactor,
	OktaEmailFactor,
	OktaOtpFactor,
	OktaPasswordFactor,
	OktaPushFactor,
	OktaQuestionFactor,
	OktaSmsFactor,
	RsaTokenFactor,
	SymantecVipFactor,
	YubikeyTokenFactor,
}

// List of factors that are applicable to Okta Identity Engine (OIE)
var AuthenticatorProviders = []string{
	// DuoFactor, // not implemented
	ExternalIdpFactor,
	GoogleOtpFactor,
	OktaEmailFactor,
	OktaPasswordFactor, // Note: Not configurable in OIE policies (Handle downstream as necessary)
	OktaVerifyFactor,
	OnPremMfaFactor,
	PhoneNumberFactor,
	RsaTokenFactor,
	SecurityQuestionFactor,
	WebauthnFactor,
	// YubikeyTokenFactor, // not implemented
}

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
