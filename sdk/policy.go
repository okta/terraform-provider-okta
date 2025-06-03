// DO NOT EDIT LOCAL SDK - USE v3 okta-sdk-golang FOR API CALLS THAT DO NOT EXIST IN LOCAL SDK
package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
func PasswordPolicy() SdkPolicy {
	// Initialize a policy with password data
	p := SdkPolicy{}
	p.Type = PasswordPolicyType
	return p
}

// SignOnPolicy returns policy of OKTA_SIGN_ON type
func SignOnPolicy() SdkPolicy {
	p := SdkPolicy{}
	p.Type = SignOnPolicyType
	return p
}

// MfaPolicy returns policy of MFA_ENROLL type
func MfaPolicy() SdkPolicy {
	p := SdkPolicy{}
	p.Type = MfaPolicyType
	return p
}

// ProfileEnrollmentPolicy returns policy of PROFILE_ENROLLMENT type
func ProfileEnrollmentPolicy() SdkPolicy {
	p := SdkPolicy{}
	p.Type = ProfileEnrollmentPolicyType
	return p
}

// Policy wrapper over okta.Policy until all of the public properties are fully supported
type SdkPolicy struct {
	// TODO
	Policy

	Settings *SdkPolicySettings `json:"settings,omitempty"`
}

// MarshalJSON Deal with the embedded struct okta.Policy having its own
// marshaler. okta.Policy doens't support a policy settings fully so we have a
// local implementation of it here.
// https://developer.okta.com/docs/reference/api/policy/#policy-settings-object
func (a *SdkPolicy) MarshalJSON() ([]byte, error) {
	// This technique is derived from
	// https://jhall.io/posts/go-json-tricks-embedded-marshaler/
	policyJSON, err := a.Policy.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var settingsJSON []byte
	if a.Settings != nil {
		type Alias SdkPolicySettings
		type local struct {
			*Alias
		}
		result := local{Alias: (*Alias)(a.Settings)}
		settingsJSON, err = json.Marshal(&result)
		if err != nil {
			return nil, err
		}

		// manipulate a serialized policyJSON with a serialized settingsJSON to have
		// the settings embedded properly.
		separator := ","
		if string(policyJSON) == "{}" {
			separator = ""
		}
		settingsJSON = []byte(fmt.Sprintf("%s\"settings\":%s}", separator, settingsJSON))
	}

	var _json string
	if len(settingsJSON) > 0 {
		_json = fmt.Sprintf("%s%s", policyJSON[:len(policyJSON)-1], settingsJSON)
	} else {
		_json = string(policyJSON)
	}
	return []byte(_json), nil
}

func (a *SdkPolicy) UnmarshalJSON(data []byte) error {
	type Alias SdkPolicy
	result := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	// Need to get around multiple embedded structs issue when unmarshalling so
	// make use of an anonymous struct so only settings are unmarshaled
	settings := struct {
		Settings SdkPolicySettings `json:"settings,omitempty"`
	}{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return err
	}
	a.Settings = &settings.Settings

	return nil
}

// PolicySettings missing from okta-sdk-golang. However, there is a subset okta.PasswordPolicySettings.
type SdkPolicySettings struct {
	// TODO
	Authenticators []*PolicyAuthenticator            `json:"authenticators,omitempty"`
	Delegation     *PasswordPolicyDelegationSettings `json:"delegation,omitempty"`
	Factors        *PolicyFactorsSettings            `json:"factors,omitempty"`
	Password       *PasswordPolicyPasswordSettings   `json:"password,omitempty"`
	Recovery       *PasswordPolicyRecoverySettings   `json:"recovery,omitempty"`
	Type           string                            `json:"type,omitempty"`
}

// PolicyFactorsSettings is not expressed in the okta-sdk-golang yet
type PolicyFactorsSettings struct {
	// TODO
	Duo          *PolicyFactor `json:"duo,omitempty"`
	FidoU2f      *PolicyFactor `json:"fido_u2f,omitempty"`
	FidoWebauthn *PolicyFactor `json:"fido_webauthn,omitempty"`
	Hotp         *PolicyFactor `json:"hotp,omitempty"`
	GoogleOtp    *PolicyFactor `json:"google_otp,omitempty"`
	CustomOtp    *PolicyFactor `json:"custom_otp,omitempty"`
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
	Key         string                          `json:"key,omitempty"`
	ID          string                          `json:"id,omitempty"`
	Enroll      *Enroll                         `json:"enroll,omitempty"`
	Constraints *PolicyAuthenticatorConstraints `json:"constraints,omitempty"`
}

type PolicyAuthenticatorConstraints struct {
	AaguidGroups []string `json:"aaguidGroups,omitempty"`
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
func (m *APISupplement) GetPolicy(ctx context.Context, policyID string) (*SdkPolicy, *Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v", policyID)
	req, err := m.RequestExecutor.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	var policy *SdkPolicy
	resp, err := m.RequestExecutor.Do(ctx, req, &policy)
	if err != nil {
		return nil, resp, err
	}
	return policy, resp, nil
}

// UpdatePolicy updates a policy.
func (m *APISupplement) UpdatePolicy(ctx context.Context, policyID string, body SdkPolicy) (*SdkPolicy, *Response, error) {
	url := fmt.Sprintf("/api/v1/policies/%v", policyID)
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, nil, err
	}
	var policy *SdkPolicy
	resp, err := m.RequestExecutor.Do(ctx, req, &policy)
	if err != nil {
		return nil, resp, err
	}
	return policy, resp, nil
}

// CreatePolicy creates a policy.
func (m *APISupplement) CreatePolicy(ctx context.Context, body SdkPolicy) (*SdkPolicy, *Response, error) {
	url := "/api/v1/policies"
	req, err := m.RequestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, err
	}
	var policy *SdkPolicy
	resp, err := m.RequestExecutor.Do(ctx, req, &policy)
	if err != nil {
		return nil, resp, err
	}
	return policy, resp, nil
}
