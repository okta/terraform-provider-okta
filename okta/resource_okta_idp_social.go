package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/terraform-providers/terraform-provider-okta/sdk"
)

func resourceIdpSocial() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdpSocialCreate,
		Read:   resourceIdpSocialRead,
		Update: resourceIdpSocialUpdate,
		Delete: resourceIdpDelete,
		Exists: getIdentityProviderExists(&sdk.SAMLIdentityProvider{}),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		// Note the base schema
		Schema: buildIdpSchema(map[string]*schema.Schema{
			"authorization_url":     optUrlSchema,
			"authorization_binding": optBindingSchema,
			"token_url":             optUrlSchema,
			"token_binding":         optBindingSchema,
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"FACEBOOK", "LINKEDIN", "MICROSOFT", "GOOGLE"},
					false,
				),
			},
			"scopes": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"protocol_type": &schema.Schema{
				Type:         schema.TypeString,
				Default:      "OAUTH2",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"OIDC", "OAUTH2"}, false),
			},
			"client_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"match_type": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "This property was incorrectly added to this resource, you should use \"subject_match_type\"",
			},
			"match_attribute": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "This property was incorrectly added to this resource, you should use \"subject_match_attribute\"",
			},
			"client_secret": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"max_clock_skew": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"issuer_mode": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "Indicates whether Okta uses the original Okta org domain URL, or a custom domain URL",
				ValidateFunc: validation.StringInSlice([]string{"ORG_URL", "CUSTOM_URL"}, false),
				Default:      "ORG_URL",
				Optional:     true,
			},
		}),
	}
}

func resourceIdpSocialCreate(d *schema.ResourceData, m interface{}) error {
	idp := buildidpSocial(d)
	if err := createIdp(m, idp); err != nil {
		return err
	}
	d.SetId(idp.ID)

	if err := setIdpStatus(idp.ID, idp.Status, d.Get("status").(string), m); err != nil {
		return err
	}

	return resourceIdpSocialRead(d, m)
}

func resourceIdpSocialRead(d *schema.ResourceData, m interface{}) error {
	idp := &sdk.OIDCIdentityProvider{}

	if err := fetchIdp(d.Id(), m, idp); err != nil {
		return err
	}

	if idp == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", idp.Name)
	d.Set("max_clock_skew", idp.Policy.MaxClockSkew)
	d.Set("provisioning_action", idp.Policy.Provisioning.Action)
	d.Set("deprovisioned_action", idp.Policy.Provisioning.Conditions.Deprovisioned.Action)
	d.Set("suspended_action", idp.Policy.Provisioning.Conditions.Suspended.Action)
	d.Set("profile_master", idp.Policy.Provisioning.ProfileMaster)
	d.Set("subject_match_type", idp.Policy.Subject.MatchType)
	d.Set("subject_match_attribute", idp.Policy.Subject.MatchAttribute)
	d.Set("username_template", idp.Policy.Subject.UserNameTemplate.Template)
	d.Set("client_id", idp.Protocol.Credentials.Client.ClientID)
	d.Set("client_secret", idp.Protocol.Credentials.Client.ClientSecret)

	if err := syncGroupActions(d, idp.Policy.Provisioning.Groups); err != nil {
		return err
	}

	if idp.IssuerMode != "" {
		d.Set("issuer_mode", idp.IssuerMode)
	}

	setMap := map[string]interface{}{
		"scopes": convertStringSetToInterface(idp.Protocol.Scopes),
	}

	if idp.Policy.AccountLink != nil {
		d.Set("account_link_action", idp.Policy.AccountLink.Action)

		if idp.Policy.AccountLink.Filter != nil {
			setMap["account_link_group_include"] = convertStringSetToInterface(idp.Policy.AccountLink.Filter.Groups.Include)
		}
	}

	return setNonPrimitives(d, setMap)
}

func resourceIdpSocialUpdate(d *schema.ResourceData, m interface{}) error {
	idp := buildidpSocial(d)
	d.Partial(true)

	if err := updateIdp(d.Id(), m, idp); err != nil {
		return err
	}

	d.Partial(false)

	if err := setIdpStatus(idp.ID, idp.Status, d.Get("status").(string), m); err != nil {
		return err
	}

	return resourceIdpSocialRead(d, m)
}

func buildidpSocial(d *schema.ResourceData) *sdk.OIDCIdentityProvider {
	return &sdk.OIDCIdentityProvider{
		Name:       d.Get("name").(string),
		Type:       d.Get("type").(string),
		IssuerMode: d.Get("issuer_mode").(string),
		Policy: &sdk.OIDCPolicy{
			AccountLink:  NewAccountLink(d),
			MaxClockSkew: int64(d.Get("max_clock_skew").(int)),
			Provisioning: NewIdpProvisioning(d),
			Subject: &sdk.OIDCSubject{
				MatchType:      d.Get("subject_match_type").(string),
				MatchAttribute: d.Get("subject_match_attribute").(string),
				UserNameTemplate: &okta.ApplicationCredentialsUsernameTemplate{
					Template: d.Get("username_template").(string),
				},
			},
		},
		Protocol: &sdk.OIDCProtocol{
			Scopes: convertInterfaceToStringSet(d.Get("scopes")),
			Type:   d.Get("protocol_type").(string),
			Credentials: &sdk.OIDCCredentials{
				Client: &sdk.OIDCClient{
					ClientID:     d.Get("client_id").(string),
					ClientSecret: d.Get("client_secret").(string),
				},
			},
		},
	}
}
