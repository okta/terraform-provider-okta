package okta

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/okta/okta-sdk-golang/okta"
)

type (
	IdentityProvider interface {
		IsIDP() bool
	}

	SAMLIdentityProvider struct {
		ID       string        `json:"id"`
		Name     string        `json:"name"`
		Policy   *SAMLPolicy   `json:"policy"`
		Protocol *SAMLProtocol `json:"protocol"`
		Type     string        `json:"type"`
		Status   string        `json:"status"`
	}

	SAMLPolicy struct {
		AccountLink  *AccountLink     `json:"accountLink"`
		Provisioning *IDPProvisioning `json:"provisioning"`
		Subject      *SAMLSubject     `json:"subject"`
	}

	SAMLEndpoints struct {
		Acs *ACSSSO `json:"acs"`
		Sso *IDPSSO `json:"sso"`
	}

	IDPProvisioning struct {
		Action        string         `json:"action"`
		Conditions    *IDPConditions `json:"conditions"`
		Groups        *IDPAction     `json:"groups"`
		ProfileMaster bool           `json:"profileMaster"`
	}

	AccountLink struct {
		Action string      `json:"action"`
		Filter interface{} `json:"filter"`
	}

	IDPAction struct {
		Action string `json:"action"`
	}

	Signature struct {
		Algorithm string `json:"algorithm"`
		Scope     string `json:"scope"`
	}

	SAMLProtocol struct {
		Algorithms  *Algorithms      `json:"algorithms"`
		Credentials *SAMLCredentials `json:"credentials"`
		Endpoints   *SAMLEndpoints   `json:"endpoints"`
		Type        string           `json:"type"`
	}

	IDPTrust struct {
		Audience string `json:"audience"`
		Issuer   string `json:"issuer"`
		Kid      string `json:"kid"`
	}

	IDPSSO struct {
		Binding     string `json:"binding"`
		Destination string `json:"destination"`
		URL         string `json:"url"`
	}

	ACSSSO struct {
		Binding string `json:"binding"`
		Type    string `json:"type"`
	}

	Endpoint struct {
		Binding string `json:"binding"`
		URL     string `json:"url"`
	}

	IDPConditions struct {
		Deprovisioned *IDPAction `json:"deprovisioned"`
		Suspended     *IDPAction `json:"suspended"`
	}

	SAMLSubject struct {
		Filter           string                                       `json:"filter"`
		Format           []string                                     `json:"format"`
		MatchType        string                                       `json:"matchType"`
		UserNameTemplate *okta.ApplicationCredentialsUsernameTemplate `json:"userNameTemplate"`
	}

	Algorithms struct {
		Request  *IDPSignature `json:"request"`
		Response *IDPSignature `json:"response"`
	}

	IDPSignature struct {
		Signature *Signature `json:"signature"`
	}

	SAMLCredentials struct {
		Trust *IDPTrust `json:"trust"`
	}

	OIDCIdentityProvider struct {
		ID       string        `json:"id"`
		Name     string        `json:"name"`
		Policy   *OIDCPolicy   `json:"policy"`
		Protocol *OIDCProtocol `json:"protocol"`
		Type     string        `json:"type"`
		Status   string        `json:"status"`
	}

	OIDCPolicy struct {
		AccountLink  *AccountLink     `json:"accountLink"`
		MaxClockSkew int64            `json:"maxClockSkew"`
		Provisioning *IDPProvisioning `json:"provisioning"`
		Subject      *OIDCSubject     `json:"subject"`
	}

	OIDCEndpoints struct {
		Acs           *ACSSSO   `json:"acs"`
		Authorization *Endpoint `json:"authorization"`
		Jwks          *Endpoint `json:"jwks"`
		Token         *Endpoint `json:"token"`
		UserInfo      *Endpoint `json:"userInfo"`
	}

	OIDCProtocol struct {
		Algorithms  *Algorithms      `json:"algorithms"`
		Credentials *OIDCCredentials `json:"credentials"`
		Endpoints   *OIDCEndpoints   `json:"endpoints"`
		Issuer      *Issuer          `json:"issuer"`
		Scopes      []string         `json:"scopes"`
		Type        string           `json:"type"`
	}

	OIDCCredentials struct {
		Client *OIDCClient `json:"client"`
	}

	OIDCClient struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}

	OIDCSubject struct {
		MatchType        string                                       `json:"matchType"`
		UserNameTemplate *okta.ApplicationCredentialsUsernameTemplate `json:"userNameTemplate"`
	}

	Issuer struct {
		URL string `json:"url"`
	}
)

func (i *OIDCIdentityProvider) IsIDP() bool {
	return true
}

func (i *SAMLIdentityProvider) IsIDP() bool {
	return true
}

func getEndpoint(d *schema.ResourceData, key string) *Endpoint {
	return &Endpoint{
		Binding: d.Get(fmt.Sprintf("%s_binding", key)).(string),
		URL:     d.Get(fmt.Sprintf("%s_url", key)).(string),
	}
}
