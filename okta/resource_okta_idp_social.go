package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

func resourceIdpSocial() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdpSocialCreate,
		ReadContext:   resourceIdpSocialRead,
		UpdateContext: resourceIdpSocialUpdate,
		DeleteContext: resourceIdpDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		// Note the base schema
		Schema: buildIdpSchema(map[string]*schema.Schema{
			"authorization_url":     optURLSchema,
			"authorization_binding": optBindingSchema,
			"token_url":             optURLSchema,
			"token_binding":         optBindingSchema,
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: stringInSlice([]string{"FACEBOOK", "LINKEDIN", "MICROSOFT", "GOOGLE"}),
			},
			"scopes": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"protocol_type": {
				Type:             schema.TypeString,
				Default:          "OAUTH2",
				Optional:         true,
				ValidateDiagFunc: stringInSlice([]string{"OIDC", "OAUTH2"}),
			},
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"match_type": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "This property was incorrectly added to this resource, you should use \"subject_match_type\"",
			},
			"match_attribute": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "This property was incorrectly added to this resource, you should use \"subject_match_attribute\"",
			},
			"client_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"max_clock_skew": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"issuer_mode": {
				Type:             schema.TypeString,
				Description:      "Indicates whether Okta uses the original Okta org domain URL, or a custom domain URL",
				ValidateDiagFunc: stringInSlice([]string{"ORG_URL", "CUSTOM_URL"}),
				Default:          "ORG_URL",
				Optional:         true,
			},
		}),
	}
}

func resourceIdpSocialCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	idp := buildIdPSocial(d)
	_, _, err := getSupplementFromMetadata(m).CreateIdentityProvider(ctx, idp, nil)
	if err != nil {
		return diag.Errorf("failed to create social identity provider: %v", err)
	}
	d.SetId(idp.ID)
	err = setIdpStatus(ctx, d, getOktaClientFromMetadata(m), idp.Status)
	if err != nil {
		return diag.Errorf("failed to change social identity provider's status: %v", err)
	}
	return resourceIdpSocialRead(ctx, d, m)
}

func resourceIdpSocialRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	idp := &sdk.OIDCIdentityProvider{}
	_, resp, err := getSupplementFromMetadata(m).GetIdentityProvider(ctx, d.Id(), idp)
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get SAML identity provider: %v", err)
	}
	if idp.ID == "" {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", idp.Name)
	_ = d.Set("max_clock_skew", idp.Policy.MaxClockSkew)
	_ = d.Set("provisioning_action", idp.Policy.Provisioning.Action)
	_ = d.Set("deprovisioned_action", idp.Policy.Provisioning.Conditions.Deprovisioned.Action)
	_ = d.Set("suspended_action", idp.Policy.Provisioning.Conditions.Suspended.Action)
	_ = d.Set("profile_master", idp.Policy.Provisioning.ProfileMaster)
	_ = d.Set("subject_match_type", idp.Policy.Subject.MatchType)
	_ = d.Set("subject_match_attribute", idp.Policy.Subject.MatchAttribute)
	_ = d.Set("username_template", idp.Policy.Subject.UserNameTemplate.Template)
	_ = d.Set("client_id", idp.Protocol.Credentials.Client.ClientID)
	_ = d.Set("client_secret", idp.Protocol.Credentials.Client.ClientSecret)

	err = syncGroupActions(d, idp.Policy.Provisioning.Groups)
	if err != nil {
		return diag.Errorf("failed to set social identity provider properties: %v", err)
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
		return diag.Errorf("failed to set social identity provider properties: %v", err)
	}
	return nil
}

func resourceIdpSocialUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	idp := buildIdPSocial(d)
	_, _, err := getSupplementFromMetadata(m).UpdateIdentityProvider(ctx, d.Id(), idp, nil)
	if err != nil {
		return diag.Errorf("failed to update social identity provider: %v", err)
	}
	err = setIdpStatus(ctx, d, getOktaClientFromMetadata(m), idp.Status)
	if err != nil {
		return diag.Errorf("failed to update social identity provider's status: %v", err)
	}
	return resourceIdpSocialRead(ctx, d, m)
}

func buildIdPSocial(d *schema.ResourceData) *sdk.OIDCIdentityProvider {
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
