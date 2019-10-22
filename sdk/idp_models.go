package sdk

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/okta"
)

type (
	IdentityProvider interface {
		IsIDP() bool
	}

	SAMLIdentityProvider struct {
		ID         string        `json:"id,omitempty"`
		IssuerMode string        `json:"issuerMode,omitempty"`
		Name       string        `json:"name,omitempty"`
		Policy     *SAMLPolicy   `json:"policy,omitempty"`
		Protocol   *SAMLProtocol `json:"protocol,omitempty"`
		Type       string        `json:"type,omitempty"`
		Status     string        `json:"status,omitempty"`
	}

	SAMLPolicy struct {
		AccountLink  *AccountLink     `json:"accountLink,omitempty"`
		Provisioning *IDPProvisioning `json:"provisioning,omitempty"`
		Subject      *SAMLSubject     `json:"subject,omitempty"`
		Type         string           `json:"type,omitempty"`
	}

	SAMLEndpoints struct {
		Acs *ACSSSO `json:"acs,omitempty"`
		Sso *IDPSSO `json:"sso,omitempty"`
	}

	IDPProvisioning struct {
		Action        string           `json:"action,omitempty"`
		Conditions    *IDPConditions   `json:"conditions,omitempty"`
		Groups        *IDPGroupsAction `json:"groups,omitempty"`
		ProfileMaster bool             `json:"profileMaster,omitempty"`
	}

	AccountLink struct {
		Action string  `json:"action,omitempty"`
		Filter *Filter `json:"filter,omitempty"`
	}

	Filter struct {
		Groups *Included `json:"groups"`
	}

	Included struct {
		Include []string `json:"include"`
	}

	IDPAction struct {
		Action string `json:"action,omitempty"`
	}

	IDPGroupsAction struct {
		Action              string   `json:"action,omitempty"`
		Assignments         []string `json:"assignments,omitempty"`
		Filter              []string `json:"filter,omitempty"`
		SourceAttributeName string   `json:"sourceAttributeName,omitempty"`
	}

	Signature struct {
		Algorithm string `json:"algorithm,omitempty"`
		Scope     string `json:"scope,omitempty"`
	}

	SAMLProtocol struct {
		Algorithms  *Algorithms      `json:"algorithms,omitempty"`
		Credentials *SAMLCredentials `json:"credentials,omitempty"`
		Endpoints   *SAMLEndpoints   `json:"endpoints,omitempty"`
		Type        string           `json:"type,omitempty"`
	}

	IDPTrust struct {
		Audience string `json:"audience,omitempty"`
		Issuer   string `json:"issuer,omitempty"`
		Kid      string `json:"kid,omitempty"`
	}

	IDPSSO struct {
		Binding     string `json:"binding,omitempty"`
		Destination string `json:"destination,omitempty"`
		URL         string `json:"url,omitempty"`
	}

	ACSSSO struct {
		Binding string `json:"binding,omitempty"`
		Type    string `json:"type,omitempty"`
	}

	Endpoint struct {
		Binding string `json:"binding,omitempty"`
		URL     string `json:"url,omitempty"`
	}

	IDPConditions struct {
		Deprovisioned *IDPAction `json:"deprovisioned,omitempty"`
		Suspended     *IDPAction `json:"suspended,omitempty"`
	}

	SAMLSubject struct {
		Filter           string                                       `json:"filter,omitempty"`
		Format           []string                                     `json:"format,omitempty"`
		MatchType        string                                       `json:"matchType,omitempty"`
		UserNameTemplate *okta.ApplicationCredentialsUsernameTemplate `json:"userNameTemplate,omitempty"`
	}

	Algorithms struct {
		Request  *IDPSignature `json:"request,omitempty"`
		Response *IDPSignature `json:"response,omitempty"`
	}

	IDPSignature struct {
		Signature *Signature `json:"signature,omitempty"`
	}

	SAMLCredentials struct {
		Trust *IDPTrust `json:"trust,omitempty"`
	}

	OIDCIdentityProvider struct {
		ID         string        `json:"id,omitempty"`
		IssuerMode string        `json:"issuerMode,omitempty"`
		Name       string        `json:"name,omitempty"`
		Policy     *OIDCPolicy   `json:"policy,omitempty"`
		Protocol   *OIDCProtocol `json:"protocol,omitempty"`
		Type       string        `json:"type,omitempty"`
		Status     string        `json:"status,omitempty"`
	}

	OIDCPolicy struct {
		AccountLink  *AccountLink     `json:"accountLink,omitempty"`
		MaxClockSkew int64            `json:"maxClockSkew"`
		Provisioning *IDPProvisioning `json:"provisioning,omitempty"`
		Subject      *OIDCSubject     `json:"subject,omitempty"`
	}

	OIDCEndpoints struct {
		Acs           *ACSSSO   `json:"acs,omitempty"`
		Authorization *Endpoint `json:"authorization,omitempty"`
		Jwks          *Endpoint `json:"jwks,omitempty"`
		Token         *Endpoint `json:"token,omitempty"`
		UserInfo      *Endpoint `json:"userInfo,omitempty"`
	}

	OIDCProtocol struct {
		Algorithms  *Algorithms      `json:"algorithms,omitempty"`
		Credentials *OIDCCredentials `json:"credentials,omitempty"`
		Endpoints   *OIDCEndpoints   `json:"endpoints,omitempty"`
		Issuer      *Issuer          `json:"issuer,omitempty"`
		Scopes      []string         `json:"scopes,omitempty"`
		Type        string           `json:"type,omitempty"`
	}

	OIDCCredentials struct {
		Client *OIDCClient `json:"client,omitempty"`
	}

	OIDCClient struct {
		ClientID     string `json:"client_id,omitempty"`
		ClientSecret string `json:"client_secret,omitempty"`
	}

	OIDCSubject struct {
		MatchType        string                                       `json:"matchType,omitempty"`
		MatchAttribute   string                                       `json:"matchAttribute,omitempty"`
		UserNameTemplate *okta.ApplicationCredentialsUsernameTemplate `json:"userNameTemplate,omitempty"`
	}

	Issuer struct {
		URL string `json:"url,omitempty"`
	}

	SigningKey struct {
		Created   string   `json:"created,omitempty"`
		ExpiresAt string   `json:"expiresAt,omitempty"`
		X5C       []string `json:"x5c"`
		Kid       string   `json:"kid"`
		Kty       string   `json:"kty"`
		Use       string   `json:"use"`
		X5T256    string   `json:"x5t#S256"`
		E         string   `json:"e,omitempty"`
		N         string   `json:"n,omitempty"`
	}

	Certificate struct {
		X5C []string `json:"x5c"`
	}

	BasicIdp struct {
		IdentityProvider
		Id       string        `json:"id"`
		Name     string        `json:"name"`
		Type     string        `json:"type,omitempty"`
		Status   string        `json:"status,omitempty"`
		Protocol *SAMLProtocol `json:"protocol,omitempty"`
	}
)

func (i *OIDCIdentityProvider) IsIDP() bool {
	return true
}

func (i *SAMLIdentityProvider) IsIDP() bool {
	return true
}

func (i *BasicIdp) IsIDP() bool {
	return true
}

func GetEndpoint(d *schema.ResourceData, key string) *Endpoint {
	binding := d.Get(fmt.Sprintf("%s_binding", key)).(string)
	url := d.Get(fmt.Sprintf("%s_url", key)).(string)

	if binding != "" && url != "" {
		return &Endpoint{
			Binding: binding,
			URL:     url,
		}
	}
	return nil
}
