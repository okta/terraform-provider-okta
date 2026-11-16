// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"context"
	"fmt"
	"net/http"
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
	SmartCardIdpFactor     = "smart_card_idp"
	CustomAppFactor        = "custom_app" // OIE only
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

// GetOrgFactor gets a factor by ID.
func (m *APISupplement) GetOrgFactor(ctx context.Context, id string) (*OrgFactor, *Response, error) {
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
func (m *APISupplement) ActivateOrgFactor(ctx context.Context, id string) (*OrgFactor, *Response, error) {
	return m.lifecycleChangeOrgFactor(ctx, id, "activate")
}

// DeactivateOrgFactor denies multifactor authentication to use provided factor type
func (m *APISupplement) DeactivateOrgFactor(ctx context.Context, id string) (*OrgFactor, *Response, error) {
	return m.lifecycleChangeOrgFactor(ctx, id, "deactivate")
}

func (m *APISupplement) lifecycleChangeOrgFactor(ctx context.Context, id, action string) (*OrgFactor, *Response, error) {
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
