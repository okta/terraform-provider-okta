package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
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
		Schema: buildIdpSchema(map[string]*schema.Schema{
			"authorization_url":     optURLSchema,
			"authorization_binding": optBindingSchema,
			"token_url":             optURLSchema,
			"token_binding":         optBindingSchema,
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: elemInSlice([]string{"FACEBOOK", "LINKEDIN", "MICROSOFT", "GOOGLE", "APPLE"}),
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
				ValidateDiagFunc: elemInSlice([]string{"OIDC", "OAUTH2"}),
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
				ValidateDiagFunc: elemInSlice([]string{"ORG_URL", "CUSTOM_URL"}),
				Default:          "ORG_URL",
				Optional:         true,
			},
			"apple_kid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Key ID that you obtained from Apple when you created the private key for the client",
			},
			"apple_private_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The PKCS #8 encoded private key that you created for the client and downloaded from Apple",
			},
			"apple_team_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Team ID associated with your Apple developer account",
			},
		}),
	}
}

func resourceIdpSocialCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	idp := buildIdPSocial(d)
	respIdp, _, err := getOktaClientFromMetadata(m).IdentityProvider.CreateIdentityProvider(ctx, idp)
	if err != nil {
		return diag.Errorf("failed to create social identity provider: %v", err)
	}
	d.SetId(respIdp.Id)
	err = setIdpStatus(ctx, d, getOktaClientFromMetadata(m), idp.Status)
	if err != nil {
		return diag.Errorf("failed to change social identity provider's status: %v", err)
	}
	return resourceIdpSocialRead(ctx, d, m)
}

func resourceIdpSocialRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	idp, resp, err := getOktaClientFromMetadata(m).IdentityProvider.GetIdentityProvider(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get identity provider: %v", err)
	}
	if idp == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("type", idp.Type)
	_ = d.Set("name", idp.Name)
	_ = d.Set("max_clock_skew", idp.Policy.MaxClockSkew)
	_ = d.Set("provisioning_action", idp.Policy.Provisioning.Action)
	_ = d.Set("deprovisioned_action", idp.Policy.Provisioning.Conditions.Deprovisioned.Action)
	_ = d.Set("suspended_action", idp.Policy.Provisioning.Conditions.Suspended.Action)
	_ = d.Set("profile_master", idp.Policy.Provisioning.ProfileMaster)
	_ = d.Set("subject_match_type", idp.Policy.Subject.MatchType)
	_ = d.Set("subject_match_attribute", idp.Policy.Subject.MatchAttribute)
	_ = d.Set("username_template", idp.Policy.Subject.UserNameTemplate.Template)
	_ = d.Set("protocol_type", idp.Protocol.Type)
	_ = d.Set("client_id", idp.Protocol.Credentials.Client.ClientId)
	_ = d.Set("client_secret", idp.Protocol.Credentials.Client.ClientSecret)
	if idp.Type == "APPLE" {
		_ = d.Set("apple_kid", idp.Protocol.Credentials.Signing.Kid)
		_ = d.Set("apple_team_id", idp.Protocol.Credentials.Signing.TeamId)
	}

	err = syncGroupActions(d, idp.Policy.Provisioning.Groups)
	if err != nil {
		return diag.Errorf("failed to set social identity provider properties: %v", err)
	}

	if idp.IssuerMode != "" {
		_ = d.Set("issuer_mode", idp.IssuerMode)
	}

	setMap := map[string]interface{}{
		"scopes": convertStringSliceToSet(idp.Protocol.Scopes),
	}

	if idp.Policy.AccountLink != nil {
		_ = d.Set("account_link_action", idp.Policy.AccountLink.Action)

		if idp.Policy.AccountLink.Filter != nil {
			setMap["account_link_group_include"] = convertStringSliceToSet(idp.Policy.AccountLink.Filter.Groups.Include)
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
	_, _, err := getOktaClientFromMetadata(m).IdentityProvider.UpdateIdentityProvider(ctx, d.Id(), idp)
	if err != nil {
		return diag.Errorf("failed to update social identity provider: %v", err)
	}
	err = setIdpStatus(ctx, d, getOktaClientFromMetadata(m), idp.Status)
	if err != nil {
		return diag.Errorf("failed to update social identity provider's status: %v", err)
	}
	return resourceIdpSocialRead(ctx, d, m)
}

func buildIdPSocial(d *schema.ResourceData) okta.IdentityProvider {
	idp := okta.IdentityProvider{
		Name:       d.Get("name").(string),
		Type:       d.Get("type").(string),
		IssuerMode: d.Get("issuer_mode").(string),
		Policy: &okta.IdentityProviderPolicy{
			AccountLink:  buildPolicyAccountLink(d),
			MaxClockSkew: int64(d.Get("max_clock_skew").(int)),
			Provisioning: buildIdPProvisioning(d),
			Subject: &okta.PolicySubject{
				MatchType:      d.Get("subject_match_type").(string),
				MatchAttribute: d.Get("subject_match_attribute").(string),
				UserNameTemplate: &okta.PolicyUserNameTemplate{
					Template: d.Get("username_template").(string),
				},
			},
		},
		Protocol: &okta.Protocol{
			Scopes: convertInterfaceToStringSet(d.Get("scopes")),
			Type:   d.Get("protocol_type").(string),
			Credentials: &okta.IdentityProviderCredentials{
				Client: &okta.IdentityProviderCredentialsClient{
					ClientId:     d.Get("client_id").(string),
					ClientSecret: d.Get("client_secret").(string),
				},
			},
		},
	}
	if idp.Type == "APPLE" {
		idp.Protocol.Credentials.Signing = &okta.IdentityProviderCredentialsSigning{
			Kid:        "",
			PrivateKey: d.Get("apple_private_key").(string),
			TeamId:     d.Get("apple_team_id").(string),
		}
		if kid, ok := d.GetOk("apple_kid"); ok {
			idp.Protocol.Credentials.Signing.Kid = kid.(string)
		}
	}
	return idp
}
