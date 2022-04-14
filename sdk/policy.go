package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

const (
	PasswordPolicyType           = "PASSWORD"
	SignOnPolicyType             = "OKTA_SIGN_ON"
	SignOnPolicyRuleType         = "SIGN_ON"
	MfaPolicyType                = "MFA_ENROLL"
	AccessPolicyType             = "ACCESS_POLICY"
	ProfileEnrollmentPolicyType  = "PROFILE_ENROLLMENT"
	IdpDiscoveryType             = "IDP_DISCOVERY"
	OauthAuthorizationPolicyType = "OAUTH_AUTHORIZATION_POLICY"
)

// PasswordPolicy returns policy of PASSWORD type
func PasswordPolicy() Policy {
	// Initialize a policy with password data
	p := Policy{}
	p.Type = PasswordPolicyType
	return p
}

// SignOnPolicy returns policy of OKTA_SIGN_ON type
func SignOnPolicy() Policy {
	p := Policy{}
	p.Type = SignOnPolicyType
	return p
}

// MfaPolicy returns policy of MFA_ENROLL type
func MfaPolicy() Policy {
	p := Policy{}
	p.Type = MfaPolicyType
	return p
}

// ProfileEnrollmentPolicy returns policy of PROFILE_ENROLLMENT type
func ProfileEnrollmentPolicy() Policy {
	p := Policy{}
	p.Type = ProfileEnrollmentPolicyType
	return p
}

// Policy wrapper over okta.Policy until all of the public properties are fully supported
type Policy struct {
	// TODO
	okta.Policy

	Settings *PolicySettings `json:"settings,omitempty"`
}

// PolicySettings missing from okta-sdk-golang. However, there is a subset okta.PasswordPolicySettings.
type PolicySettings struct {
	// TODO
	Authenticators []*PolicyAuthenticator                 `json:"authenticators,omitempty"`
	Delegation     *okta.PasswordPolicyDelegationSettings `json:"delegation,omitempty"`
	Factors        *PolicyFactorsSettings                 `json:"factors,omitempty"`
	Password       *okta.PasswordPolicyPasswordSettings   `json:"password,omitempty"`
	Recovery       *okta.PasswordPolicyRecoverySettings   `json:"recovery,omitempty"`
	Type           string                                 `json:"type,omitempty"`
}

// PolicyFactorsSettings is not expressed in the okta-sdk-golang yet
type PolicyFactorsSettings struct {
	// TODO
	Duo          *PolicyFactor `json:"duo,omitempty"`
	FidoU2f      *PolicyFactor `json:"fido_u2f,omitempty"`
	FidoWebauthn *PolicyFactor `json:"fido_webauthn,omitempty"`
	Hotp         *PolicyFactor `json:"hotp,omitempty"`
	GoogleOtp    *PolicyFactor `json:"google_otp,omitempty"`
	OktaCall     *PolicyFactor `json:"okta_call,omitempty"`
	OktaOtp      *PolicyFactor `json:"okta_otp,omitempty"`
	OktaPassword *PolicyFactor `json:"okta_password,omitempty"`
	OktaPush     *PolicyFactor `json:"okta_push,omitempty"`
	OktaQuestion *PolicyFactor `json:"okta_question,omitempty"`
	OktaSms      *PolicyFactor `json:"okta_sms,omitempty"`
	OktaEmail    *PolicyFactor `json:"okta_email,omitempty"`
	RsaToken     *PolicyFactor `json:"rsa_token,omitempty"`
	SymantecVip  *PolicyFactor `json:"symantec_vip,omitempty"`
	YubikeyToken *PolicyFactor `json:"yubikey_token,omitempty"`
}

type PolicyFactor struct {
	Consent *Consent `json:"consent,omitempty"`
	Enroll  *Enroll  `json:"enroll,omitempty"`
}

type PolicyAuthenticator struct {
	Key    string  `json:"key,omitempty"`
	Enroll *Enroll `json:"enroll,omitempty"`
}

type Consent struct {
	Terms *Terms `json:"terms,omitempty"`
	Type  string `json:"type,omitempty"`
}

type Terms struct {
	Format string `json:"format,omitempty"`
	Value  string `json:"value,omitempty"`
}

type Enroll struct {
	Self string `json:"self,omitempty"`
}

// GetPolicy gets a policy by ID
func (m *APISupplement) GetPolicy(ctx context.Context, policyID string) (*Policy, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v", policyID)
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var policy *Policy
	resp, err := m.RequestExecutor.Do(ctx, req, &policy)
	if err != nil {
		return nil, resp, err
	}
	return policy, resp, nil
}

// UpdatePolicy updates a policy.
func (m *APISupplement) UpdatePolicy(ctx context.Context, policyID string, body Policy) (*Policy, *okta.Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v", policyID)
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, nil, err
	}
	var policy *Policy
	resp, err := m.RequestExecutor.Do(ctx, req, &policy)
	if err != nil {
		return nil, resp, err
	}
	return policy, resp, nil
}

// CreatePolicy creates a policy.
func (m *APISupplement) CreatePolicy(ctx context.Context, body Policy) (*Policy, *okta.Response, error) {
	url := "/api/v1/policies"
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, err
	}
	var policy *Policy
	resp, err := m.RequestExecutor.Do(ctx, req, &policy)
	if err != nil {
		return nil, resp, err
	}
	return policy, resp, nil
}
