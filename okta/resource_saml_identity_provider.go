package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/okta/okta-sdk-golang/okta"
)

func resourceSamlIdp() *schema.Resource {
	return &schema.Resource{
		Create: resourceSamlIdpCreate,
		Read:   resourceSamlIdpRead,
		Update: resourceSamlIdpUpdate,
		Delete: resourceIdentityProviderDelete,
		Exists: getIdentityProviderExists(&SAMLIdentityProvider{}),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: buildIdpSchema(map[string]*schema.Schema{
			"sso_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"sso_binding": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice(
					[]string{postBindingAlias, redirectBindingAlias},
					false,
				),
				Default: postBindingAlias,
			},
			"sso_destination": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"subject_format": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     schema.TypeString,
				Optional: true,
			},
			"subject_filter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"provisioning_action": actionSchema,
			"issuer": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"audience": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"kid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		}),
	}
}

func resourceSamlIdpCreate(d *schema.ResourceData, m interface{}) error {
	idp := buildSamlIdp(d)
	if err := createIdp(m, idp); err != nil {
		return err
	}
	d.SetId(idp.ID)

	if err := setIdpStatus(idp.ID, idp.Status, d.Get("status").(string), m); err != nil {
		return err
	}

	return resourceSamlIdpRead(d, m)
}

func resourceSamlIdpRead(d *schema.ResourceData, m interface{}) error {
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

func resourceSamlIdpUpdate(d *schema.ResourceData, m interface{}) error {
	idp := buildSamlIdp(d)
	d.Partial(true)

	if err := updateIdp(d.Id(), m, idp); err != nil {
		return err
	}

	d.Partial(false)

	if err := setIdpStatus(idp.ID, idp.Status, d.Get("status").(string), m); err != nil {
		return err
	}

	return resourceSamlIdpRead(d, m)
}

func buildSamlIdp(d *schema.ResourceData) *SAMLIdentityProvider {
	return &SAMLIdentityProvider{
		Name: d.Get("name").(string),
		Type: "OIDC",
		Policy: &SAMLPolicy{
			AccountLink: &AccountLink{
				Action: d.Get("account_link_action").(string),
				Filter: d.Get("account_link_filter").(string),
			},
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
			Subject: &SAMLSubject{
				Filter:    d.Get("subject_filter").(string),
				Format:    convertInterfaceToStringSet(d.Get("subject_format")),
				MatchType: d.Get("subject_match_type").(string),
				UserNameTemplate: &okta.ApplicationCredentialsUsernameTemplate{
					Template: d.Get("username_template").(string),
				},
			},
		},
		Protocol: &SAMLProtocol{
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
			Endpoints: &SAMLEndpoints{
				Acs: &ACSSSO{
					Binding: d.Get("acs_binding").(string),
					Type:    d.Get("acs_type").(string),
				},
				Sso: &IDPSSO{
					Binding:     d.Get("acs_binding").(string),
					Destination: d.Get("acs_destination").(string),
					URL:         d.Get("acs_url").(string),
				},
			},
			Type: d.Get("protocol_type").(string),
			Credentials: &SAMLCredentials{
				Trust: &IDPTrust{
					Issuer:   d.Get("issuer").(string),
					Audience: d.Get("audience").(string),
					Kid:      d.Get("kid").(string),
				},
			},
		},
	}
}
