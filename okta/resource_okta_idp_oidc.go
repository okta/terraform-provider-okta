package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

func resourceIdpOidc() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdpCreate,
		ReadContext:   resourceIdpRead,
		UpdateContext: resourceIdpUpdate,
		DeleteContext: resourceIdpDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		// Note the base schema
		Schema: buildIdpSchema(map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorization_url":     urlSchema,
			"authorization_binding": bindingSchema,
			"token_url":             urlSchema,
			"token_binding":         bindingSchema,
			"user_info_url":         optionalURLSchema,
			"user_info_binding":     optionalBindingSchema,
			"jwks_url":              urlSchema,
			"jwks_binding":          bindingSchema,
			"acs_binding":           bindingSchema,
			"acs_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "INSTANCE",
				ValidateDiagFunc: stringInSlice([]string{"INSTANCE"}),
			},
			"scopes": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"protocol_type": {
				Type:             schema.TypeString,
				Default:          "OIDC",
				Optional:         true,
				ValidateDiagFunc: stringInSlice([]string{"OIDC", "OAUTH2"}),
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_secret": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"issuer_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer_mode": {
				Type:             schema.TypeString,
				Description:      "Indicates whether Okta uses the original Okta org domain URL, or a custom domain URL",
				ValidateDiagFunc: stringInSlice([]string{"ORG_URL", "CUSTOM_URL"}),
				Default:          "ORG_URL",
				Optional:         true,
			},
			"max_clock_skew": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		}),
	}
}

func resourceIdpCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	idp := buildIdPOidc(d)
	_, _, err := getSupplementFromMetadata(m).CreateIdentityProvider(ctx, idp, nil)
	if err != nil {
		return diag.Errorf("failed to create OIDC identity provider: %v", err)
	}
	d.SetId(idp.ID)
	err = setIdpStatus(ctx, d, getOktaClientFromMetadata(m), idp.Status)
	if err != nil {
		return diag.Errorf("failed to change OIDC identity provider's status: %v", err)
	}
	return resourceIdpRead(ctx, d, m)
}

func resourceIdpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	idp := &sdk.OIDCIdentityProvider{}
	_, resp, err := getSupplementFromMetadata(m).GetIdentityProvider(ctx, d.Id(), idp)
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get OIDC identity provider: %v", err)
	}
	if idp.ID == "" {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", idp.Name)
	_ = d.Set("type", idp.Type)
	_ = d.Set("max_clock_skew", idp.Policy.MaxClockSkew)
	_ = d.Set("provisioning_action", idp.Policy.Provisioning.Action)
	_ = d.Set("deprovisioned_action", idp.Policy.Provisioning.Conditions.Deprovisioned.Action)
	_ = d.Set("suspended_action", idp.Policy.Provisioning.Conditions.Suspended.Action)
	_ = d.Set("profile_master", idp.Policy.Provisioning.ProfileMaster)
	_ = d.Set("subject_match_type", idp.Policy.Subject.MatchType)
	_ = d.Set("username_template", idp.Policy.Subject.UserNameTemplate.Template)
	_ = d.Set("issuer_url", idp.Protocol.Issuer.URL)
	_ = d.Set("client_secret", idp.Protocol.Credentials.Client.ClientSecret)
	_ = d.Set("client_id", idp.Protocol.Credentials.Client.ClientID)
	syncEndpoint("authorization", idp.Protocol.Endpoints.Authorization, d)
	syncEndpoint("token", idp.Protocol.Endpoints.Token, d)
	syncEndpoint("user_info", idp.Protocol.Endpoints.UserInfo, d)
	syncEndpoint("jwks", idp.Protocol.Endpoints.Jwks, d)
	syncAlgo(d, idp.Protocol.Algorithms)

	err = syncGroupActions(d, idp.Policy.Provisioning.Groups)
	if err != nil {
		return diag.Errorf("failed to set OIDC identity provider properties: %v", err)
	}
	if idp.Protocol.Endpoints.Acs != nil {
		_ = d.Set("acs_binding", idp.Protocol.Endpoints.Acs.Binding)
		_ = d.Set("acs_type", idp.Protocol.Endpoints.Acs.Type)
	}
	if idp.IssuerMode != "" {
		_ = d.Set("issuer_mode", idp.IssuerMode)
	}
	setMap := map[string]interface{}{
		"scopes": convertStringSetToInterface(idp.Protocol.Scopes),
	}
	if idp.Policy.AccountLink != nil {
		_ = d.Set("account_link_action", idp.Policy.AccountLink.Action)
		if idp.Policy.AccountLink.Filter != nil {
			setMap["account_link_group_include"] = convertStringSetToInterface(idp.Policy.AccountLink.Filter.Groups.Include)
		}
	}
	err = setNonPrimitives(d, setMap)
	if err != nil {
		return diag.Errorf("failed to set OIDC identity provider properties: %v", err)
	}
	return nil
}

func syncEndpoint(key string, e *sdk.Endpoint, d *schema.ResourceData) {
	if e != nil {
		_ = d.Set(key+"_binding", e.Binding)
		_ = d.Set(key+"_url", e.URL)
	}
}

func resourceIdpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	idp := buildIdPOidc(d)
	_, _, err := getSupplementFromMetadata(m).UpdateIdentityProvider(ctx, d.Id(), idp, nil)
	if err != nil {
		return diag.Errorf("failed to update OIDC identity provider: %v", err)
	}
	err = setIdpStatus(ctx, d, getOktaClientFromMetadata(m), idp.Status)
	if err != nil {
		return diag.Errorf("failed to update OIDC identity provider's status: %v", err)
	}
	return resourceIdpRead(ctx, d, m)
}

func buildIdPOidc(d *schema.ResourceData) *sdk.OIDCIdentityProvider {
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
