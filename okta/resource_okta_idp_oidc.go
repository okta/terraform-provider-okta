package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/terraform-providers/terraform-provider-okta/sdk"
)

func resourceIdpOidc() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdpCreate,
		Read:   resourceIdpRead,
		Update: resourceIdpUpdate,
		Delete: resourceIdpDelete,
		Exists: getIdentityProviderExists(&sdk.OIDCIdentityProvider{}),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		// Note the base schema
		Schema: buildIdpSchema(map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorization_url":     urlSchema,
			"authorization_binding": bindingSchema,
			"token_url":             urlSchema,
			"token_binding":         bindingSchema,
			"user_info_url":         optionalUrlSchema,
			"user_info_binding":     optionalBindingSchema,
			"jwks_url":              urlSchema,
			"jwks_binding":          bindingSchema,
			"acs_binding":           bindingSchema,
			"acs_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "INSTANCE",
				ValidateFunc: validation.StringInSlice([]string{"INSTANCE"}, false),
			},
			"scopes": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"issuer_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer_mode": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "Indicates whether Okta uses the original Okta org domain URL, or a custom domain URL",
				ValidateFunc: validation.StringInSlice([]string{"ORG_URL", "CUSTOM_URL"}, false),
				Default:      "ORG_URL",
				Optional:     true,
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
	idp := &sdk.OIDCIdentityProvider{}

	if err := fetchIdp(d.Id(), m, idp); err != nil {
		return err
	}

	d.Set("name", idp.Name)
	d.Set("type", idp.Type)
	d.Set("max_clock_skew", idp.Policy.MaxClockSkew)
	d.Set("provisioning_action", idp.Policy.Provisioning.Action)
	d.Set("deprovisioned_action", idp.Policy.Provisioning.Conditions.Deprovisioned.Action)
	d.Set("suspended_action", idp.Policy.Provisioning.Conditions.Suspended.Action)
	d.Set("profile_master", idp.Policy.Provisioning.ProfileMaster)
	d.Set("subject_match_type", idp.Policy.Subject.MatchType)
	d.Set("username_template", idp.Policy.Subject.UserNameTemplate.Template)
	d.Set("issuer_url", idp.Protocol.Issuer.URL)
	d.Set("client_secret", idp.Protocol.Credentials.Client.ClientSecret)
	d.Set("client_id", idp.Protocol.Credentials.Client.ClientID)
	syncEndpoint("authorization", idp.Protocol.Endpoints.Authorization, d)
	syncEndpoint("token", idp.Protocol.Endpoints.Token, d)
	syncEndpoint("user_info", idp.Protocol.Endpoints.UserInfo, d)
	syncEndpoint("jwks", idp.Protocol.Endpoints.Jwks, d)
	syncAlgo(d, idp.Protocol.Algorithms)

	if err := syncGroupActions(d, idp.Policy.Provisioning.Groups); err != nil {
		return err
	}

	if idp.Protocol.Endpoints.Acs != nil {
		d.Set("acs_binding", idp.Protocol.Endpoints.Acs.Binding)
		d.Set("acs_type", idp.Protocol.Endpoints.Acs.Type)
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

func syncEndpoint(key string, e *sdk.Endpoint, d *schema.ResourceData) {
	if e != nil {
		d.Set(key+"_binding", e.Binding)
		d.Set(key+"_url", e.URL)
	}
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

func buildOidcIdp(d *schema.ResourceData) *sdk.OIDCIdentityProvider {
	return &sdk.OIDCIdentityProvider{
		Name:       d.Get("name").(string),
		Type:       "OIDC",
		IssuerMode: d.Get("issuer_mode").(string),
		Policy: &sdk.OIDCPolicy{
			AccountLink:  NewAccountLink(d),
			MaxClockSkew: int64(d.Get("max_clock_skew").(int)),
			Provisioning: NewIdpProvisioning(d),
			Subject: &sdk.OIDCSubject{
				MatchType: d.Get("subject_match_type").(string),
				UserNameTemplate: &okta.ApplicationCredentialsUsernameTemplate{
					Template: d.Get("username_template").(string),
				},
			},
		},
		Protocol: &sdk.OIDCProtocol{
			Algorithms: NewAlgorithms(d),
			Endpoints:  NewEndpoints(d),
			Scopes:     convertInterfaceToStringSet(d.Get("scopes")),
			Type:       d.Get("protocol_type").(string),
			Credentials: &sdk.OIDCCredentials{
				Client: &sdk.OIDCClient{
					ClientID:     d.Get("client_id").(string),
					ClientSecret: d.Get("client_secret").(string),
				},
			},
			Issuer: &sdk.Issuer{
				URL: d.Get("issuer_url").(string),
			},
		},
	}
}
