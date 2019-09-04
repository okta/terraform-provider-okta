package okta

import "fmt"

type (
	Factor struct {
		Id         string `json:"id"`
		Provider   string `json:"provider"`
		FactorType string `json:"factorType"`
		Status     string `json:"status"`
		Links      struct {
			Deactivate struct {
				Href  string `json:"href"`
				Hints *Hints `json:"hints"`
			} `json:"deactivate"`
		} `json:"_links"`
	}

	// OrgService allows you to perform actions against resources at the organization level.
	OrgService service
)

// Current available factors for MFA
const (
	DuoFactor          = "duo"
	FidoU2fFactor      = "fido_u2f"
	FidoWebauthnFactor = "fido_webauthn"
	GoogleOtpFactor    = "google_otp"
	OktaCallFactor     = "okta_call"
	OktaOtpFactor      = "okta_otp"
	OktaPushFactor     = "okta_push"
	OktaQuestionFactor = "okta_question"
	OktaSmsFactor      = "okta_sms"
	RsaTokenFactor     = "rsa_token"
	SymantecVipFactor  = "symantec_vip"
	YubikeyTokenFactor = "yubikey_token"
)

// ListFactors lists information around factors for organization.
func (s *OrgService) ListFactors() ([]*Factor, *Response, error) {
	var factorList []*Factor
	req, err := s.client.NewRequest("GET", "org/factors", nil)
	if err != nil {
		return factorList, nil, err
	}
	resp, err := s.client.Do(req, &factorList)

	return factorList, resp, err
}

// ActivateFactor ability to activate factor provider for an organization. For valid providers IDs see API docs
// https://developer.okta.com/docs/api/resources/factor_admin.
func (s *OrgService) ActivateFactor(id string) (*Factor, *Response, error) {
	return s.lifecycleChangeFactor(id, "activate")
}

// DeactivateFactor ability to deactivate factor provider for an organization. For valid provider IDs see API docs
// https://developer.okta.com/docs/api/resources/factor_admin.
func (s *OrgService) DeactivateFactor(id string) (*Factor, *Response, error) {
	return s.lifecycleChangeFactor(id, "deactivate")
}

func (s *OrgService) lifecycleChangeFactor(id, action string) (*Factor, *Response, error) {
	var factor *Factor
	relUrl := fmt.Sprintf("org/factors/%s/lifecycle/%s", id, action)
	req, err := s.client.NewRequest("POST", relUrl, nil)
	if err != nil {
		return factor, nil, err
	}
	resp, err := s.client.Do(req, &factor)

	return factor, resp, err
}
