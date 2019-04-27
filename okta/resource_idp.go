package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/okta/okta-sdk-golang/okta"
)

func resourceIdp() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdpCreate,
		Read:   resourceIdpRead,
		Update: resourceIdpUpdate,
		Delete: resourceIdentityProviderDelete,
		Exists: getIdentityProviderExists(&OIDCIdentityProvider{}),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		// Note the base schema
		Schema: buildIdpSchema(map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"OIDC", "FACEBOOK", "LINKEDIN", "MICROSOFT", "GOOGLE"},
					false,
				),
			},
			"authorization_url":     urlSchema,
			"authorization_binding": bindingSchema,
			"token_url":             urlSchema,
			"token_binding":         bindingSchema,
			"user_info_url":         urlSchema,
			"user_info_binding":     bindingSchema,
			"jwks_url":              urlSchema,
			"jwks_binding":          bindingSchema,
			"scopes": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     schema.TypeString,
				Required: true,
			},
			"protocol_type": &schema.Schema{
				Type:         schema.TypeString,
				Default:      "OIDC",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"OIDC", "OAUTH2"}, false),
			},
			"client_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"client_secret": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"max_clock_skew": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		}),
	}
}

func resourceIdpCreate(d *schema.ResourceData, m interface{}) error {
	idp := buildOidcIdp(d)
	if err := createIdp(m, idp); err != nil {
		return err
	}
	d.SetId(idp.ID)

	if err := setIdpStatus(idp.ID, idp.Status, d.Get("status").(string), m); err != nil {
		return err
	}

	return resourceIdpRead(d, m)
}

func resourceIdpRead(d *schema.ResourceData, m interface{}) error {
	var idp *OIDCIdentityProvider

	if err := fetchIdp(d.Id(), m, idp); err != nil {
		return err
	}

	d.Set("name", idp.Name)
	d.Set("account_link_action", idp.Policy.AccountLink.Action)
	d.Set("account_link_filter", idp.Policy.AccountLink.Filter)
	d.Set("max_clock_skew", idp.Policy.MaxClockSkew)
	d.Set("provisioning_action", idp.Policy.Provisioning.Action)
	d.Set("deprovisioned_action", idp.Policy.Provisioning.Conditions.Deprovisioned)
	d.Set("suspended_action", idp.Policy.Provisioning.Conditions.Suspended)
	d.Set("groups_action", idp.Policy.Provisioning.Groups.Action)
	d.Set("subject_match_type", idp.Policy.Subject.MatchType)
	d.Set("username_template", idp.Policy.Subject.UserNameTemplate.Template)
	d.Set("request_algorithm", idp.Protocol.Algorithms.Request.Signature.Algorithm)
	d.Set("request_scope", idp.Protocol.Algorithms.Request.Signature.Scope)
	d.Set("response_algorithm", idp.Protocol.Algorithms.Response.Signature.Algorithm)
	d.Set("response_scope", idp.Protocol.Algorithms.Response.Signature.Scope)

	return nil
}

func resourceIdpUpdate(d *schema.ResourceData, m interface{}) error {
	idp := buildOidcIdp(d)
	d.Partial(true)

	if err := updateIdp(d.Id(), m, idp); err != nil {
		return err
	}

	d.Partial(false)

	if err := setIdpStatus(idp.ID, idp.Status, d.Get("status").(string), m); err != nil {
		return err
	}

	return resourceIdpRead(d, m)
}
func buildOidcIdp(d *schema.ResourceData) *OIDCIdentityProvider {
	return &OIDCIdentityProvider{
		Name: d.Get("name").(string),
		Type: "OIDC",
		Policy: &OIDCPolicy{
			AccountLink: &AccountLink{
				Action: d.Get("account_link_action").(string),
				Filter: d.Get("account_link_filter").(string),
			},
			MaxClockSkew: int64(d.Get("max_clock_skew").(int)),
			Provisioning: &IDPProvisioning{
				Action: d.Get("provisioning_action").(string),
				Conditions: &IDPConditions{
					Deprovisioned: &IDPAction{
						Action: d.Get("deprovisioned_action").(string),
					},
					Suspended: &IDPAction{
						Action: d.Get("suspended_action").(string),
					},
				},
				Groups: &IDPAction{
					Action: d.Get("groups_action").(string),
				},
			},
			Subject: &OIDCSubject{
				MatchType: d.Get("subject_match_type").(string),
				UserNameTemplate: &okta.ApplicationCredentialsUsernameTemplate{
					Template: d.Get("username_template").(string),
				},
			},
		},
		Protocol: &OIDCProtocol{
			Algorithms: &Algorithms{
				Request: &IDPSignature{
					Signature: &Signature{
						Algorithm: d.Get("request_algorithm").(string),
						Scope:     d.Get("request_scope").(string),
					},
				},
				Response: &IDPSignature{
					Signature: &Signature{
						Algorithm: d.Get("response_algorithm").(string),
						Scope:     d.Get("response_scope").(string),
					},
				},
			},
			Endpoints: &OIDCEndpoints{
				Acs: &ACSSSO{
					Binding: d.Get("acs_binding").(string),
					Type:    d.Get("acs_type").(string),
				},
				Authorization: getEndpoint(d, "authorization"),
				Token:         getEndpoint(d, "token"),
				UserInfo:      getEndpoint(d, "user_info"),
				Jwks:          getEndpoint(d, "jwks"),
			},
			Scopes: convertInterfaceToStringSet(d.Get("scopes")),
			Type:   d.Get("protocol_type").(string),
			Credentials: &OIDCCredentials{
				Client: &OIDCClient{
					ClientID:     d.Get("client_id").(string),
					ClientSecret: d.Get("client_secret").(string),
				},
			},
			Issuer: &Issuer{
				URL: d.Get("issuer_url").(string),
			},
		},
	}
}
