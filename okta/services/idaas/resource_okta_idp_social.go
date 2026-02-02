package idaas

import (
	"context"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
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
		Description: "Creates a Social Identity Provider. This resource allows you to create and configure a Social Identity Provider.",
		ValidateRawResourceConfigFuncs: []schema.ValidateRawResourceConfigFunc{
			validation.PreferWriteOnlyAttribute(cty.GetAttrPath("client_secret"), cty.GetAttrPath("client_secret_wo")),
		},
		Schema: buildIdpSchema(map[string]*schema.Schema{
			"authorization_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IdP Authorization Server (AS) endpoint to request consent from the user and obtain an authorization code grant.",
			},
			"authorization_binding": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The method of making an authorization request. It can be set to `HTTP-POST` or `HTTP-REDIRECT`.",
			},
			"token_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IdP Authorization Server (AS) endpoint to exchange the authorization code grant for an access token.",
			},
			"token_binding": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The method of making a token request. It can be set to `HTTP-POST` or `HTTP-REDIRECT`.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identity Provider Types: https://developer.okta.com/docs/reference/api/idps/#identity-provider-type",
			},
			"scopes": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "The scopes of the IdP.",
			},
			"protocol_type": {
				Type:        schema.TypeString,
				Default:     "OAUTH2",
				Optional:    true,
				Description: "The type of protocol to use. It can be `OIDC` or `OAUTH2`. Default: `OAUTH2`",
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Unique identifier issued by AS for the Okta IdP instance.",
			},
			"client_secret": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"client_secret_wo"},
				Description:   "Client secret issued by AS for the Okta IdP instance. When set, this secret will be stored in the Terraform state file. For Terraform 1.11+, consider using `client_secret_wo` instead to avoid persisting secrets in state.",
			},
			"client_secret_wo": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				WriteOnly:     true,
				ConflictsWith: []string{"client_secret"},
				Description:   "Write-only client secret issued by AS for the Okta IdP instance for Terraform 1.11+. Unlike `client_secret`, this secret will not be persisted in the Terraform state file, providing improved security. Only use this attribute with Terraform 1.11 or higher.",
			},
			"client_secret_wo_version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Version number for the write-only client secret. Increment this value to trigger an update when changing `client_secret_wo`.",
			},
			"max_clock_skew": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum allowable clock-skew when processing messages from the IdP.",
			},
			"issuer_mode": {
				Type:        schema.TypeString,
				Description: "Indicates whether Okta uses the original Okta org domain URL, or a custom domain URL. It can be `ORG_URL` or `CUSTOM_URL`. Default: `ORG_URL`",
				Default:     "ORG_URL",
				Optional:    true,
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
				Description: "The Key ID that you obtained from Apple when you created the private key for the client. PrivateKey is required when resource is first created. For all consecutive updates, it can be empty/omitted and keeps the existing value if it is empty/omitted. PrivateKey isn't returned when importing this resource.",
			},
			"apple_team_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Team ID associated with your Apple developer account",
			},
			"trust_issuer": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Trust issuer for the Okta IdP instance.",
			},
			"trust_audience": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Trust audience for the Okta IdP instance.",
			},
			"trust_kid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Trust kid for the Okta IdP instance.",
			},
			"trust_revocation": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Trust revocation for the Okta IdP instance.",
			},
			"trust_revocation_cache_lifetime": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Trust revocation cache lifetime for the Okta IdP instance.",
			},
		}),
	}
}

func resourceIdpSocialCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	idp := buildIdPSocial(d)
	respIdp, _, err := getOktaClientFromMetadata(meta).IdentityProvider.CreateIdentityProvider(ctx, idp)
	if err != nil {
		return diag.Errorf("failed to create social identity provider: %v", err)
	}
	d.SetId(respIdp.Id)
	err = setIdpStatus(ctx, d, getOktaClientFromMetadata(meta), idp.Status)
	if err != nil {
		return diag.Errorf("failed to change social identity provider's status: %v", err)
	}
	return resourceIdpSocialRead(ctx, d, meta)
}

func resourceIdpSocialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	idp, resp, err := getOktaClientFromMetadata(meta).IdentityProvider.GetIdentityProvider(ctx, d.Id())
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get identity provider: %v", err)
	}
	if idp == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("type", idp.Type)
	_ = d.Set("name", idp.Name)
	if idp.Policy.MaxClockSkewPtr != nil {
		_ = d.Set("max_clock_skew", idp.Policy.MaxClockSkewPtr)
	}
	_ = d.Set("provisioning_action", idp.Policy.Provisioning.Action)
	if idp.Policy.Provisioning.Conditions != nil {
		_ = d.Set("deprovisioned_action", idp.Policy.Provisioning.Conditions.Deprovisioned.Action)
		_ = d.Set("suspended_action", idp.Policy.Provisioning.Conditions.Suspended.Action)
	}
	_ = d.Set("profile_master", idp.Policy.Provisioning.ProfileMaster)
	_ = d.Set("subject_match_type", idp.Policy.Subject.MatchType)
	_ = d.Set("subject_match_attribute", idp.Policy.Subject.MatchAttribute)
	_ = d.Set("username_template", idp.Policy.Subject.UserNameTemplate.Template)
	_ = d.Set("protocol_type", idp.Protocol.Type)
	if idp.Protocol.Credentials.Client != nil {
		_ = d.Set("client_id", idp.Protocol.Credentials.Client.ClientId)

		// Only set client_secret if client_secret_wo was not used in the config
		// Check if client_secret_wo is configured (not null in raw config)
		woVal, diags := d.GetRawConfigAt(cty.GetAttrPath("client_secret_wo"))
		if len(diags) == 0 && woVal.Type().Equals(cty.String) && !woVal.IsNull() {
			// client_secret_wo is being used, so don't persist client_secret in state
			_ = d.Set("client_secret", "")
		} else {
			// client_secret is being used, persist it in state
			_ = d.Set("client_secret", idp.Protocol.Credentials.Client.ClientSecret)
		}
	}
	if idp.Protocol.Credentials.Trust != nil {
		_ = d.Set("trust_issuer", idp.Protocol.Credentials.Trust.Issuer)
		_ = d.Set("trust_audience", idp.Protocol.Credentials.Trust.Audience)
		_ = d.Set("trust_kid", idp.Protocol.Credentials.Trust.Kid)
		_ = d.Set("trust_revocation", idp.Protocol.Credentials.Trust.Revocation)
		if idp.Protocol.Credentials.Trust.RevocationCacheLifetimePtr != nil {
			_ = d.Set("trust_revocation_cache_lifetime", idp.Protocol.Credentials.Trust.RevocationCacheLifetimePtr)
		}
	}
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
		"scopes": utils.ConvertStringSliceToSet(idp.Protocol.Scopes),
	}

	if idp.Policy.AccountLink != nil {
		_ = d.Set("account_link_action", idp.Policy.AccountLink.Action)

		if idp.Policy.AccountLink.Filter != nil {
			setMap["account_link_group_include"] = utils.ConvertStringSliceToSet(idp.Policy.AccountLink.Filter.Groups.Include)
		}
	}
	err = utils.SetNonPrimitives(d, setMap)
	if err != nil {
		return diag.Errorf("failed to set social identity provider properties: %v", err)
	}
	return nil
}

func resourceIdpSocialUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	idp := buildIdPSocial(d)
	_, _, err := getOktaClientFromMetadata(meta).IdentityProvider.UpdateIdentityProvider(ctx, d.Id(), idp)
	if err != nil {
		return diag.Errorf("failed to update social identity provider: %v", err)
	}
	err = setIdpStatus(ctx, d, getOktaClientFromMetadata(meta), idp.Status)
	if err != nil {
		return diag.Errorf("failed to update social identity provider's status: %v", err)
	}
	return resourceIdpSocialRead(ctx, d, meta)
}

func buildIdPSocial(d *schema.ResourceData) sdk.IdentityProvider {
	idp := sdk.IdentityProvider{
		Name:       d.Get("name").(string),
		Type:       d.Get("type").(string),
		IssuerMode: d.Get("issuer_mode").(string),
		Policy: &sdk.IdentityProviderPolicy{
			AccountLink:     buildPolicyAccountLink(d),
			MaxClockSkewPtr: utils.Int64Ptr(d.Get("max_clock_skew").(int)),
			Provisioning:    buildIdPProvisioning(d),
			Subject: &sdk.PolicySubject{
				MatchType:      d.Get("subject_match_type").(string),
				MatchAttribute: d.Get("subject_match_attribute").(string),
				UserNameTemplate: &sdk.PolicyUserNameTemplate{
					Template: d.Get("username_template").(string),
				},
			},
		},
		Protocol: &sdk.Protocol{
			Scopes: utils.ConvertInterfaceToStringSet(d.Get("scopes")),
			Type:   d.Get("protocol_type").(string),
			Credentials: &sdk.IdentityProviderCredentials{
				Client: buildClientCredentials(d),
			},
		},
	}
	if idp.Type == "APPLE" {
		idp.Protocol.Credentials.Signing = &sdk.IdentityProviderCredentialsSigning{
			Kid:        "",
			PrivateKey: d.Get("apple_private_key").(string),
			TeamId:     d.Get("apple_team_id").(string),
		}
		if kid, ok := d.GetOk("apple_kid"); ok {
			idp.Protocol.Credentials.Signing.Kid = kid.(string)
		}
	}
	if d.Get("status") != nil {
		idp.Status = d.Get("status").(string)
	}
	return idp
}

func buildClientCredentials(d *schema.ResourceData) *sdk.IdentityProviderCredentialsClient {
	// Try to get write-only attribute first, fall back to regular attribute
	var clientSecret string
	woVal, diags := d.GetRawConfigAt(cty.GetAttrPath("client_secret_wo"))
	if len(diags) == 0 && woVal.Type().Equals(cty.String) && !woVal.IsNull() {
		clientSecret = woVal.AsString()
	} else {
		clientSecret = d.Get("client_secret").(string)
	}

	return &sdk.IdentityProviderCredentialsClient{
		ClientId:     d.Get("client_id").(string),
		ClientSecret: clientSecret,
	}
}
