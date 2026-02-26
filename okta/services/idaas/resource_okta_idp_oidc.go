package idaas

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
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
		Description: "Creates an OIDC Identity Provider. This resource allows you to create and configure an OIDC Identity Provider.",
		// Note the base schema
		Schema: buildIdpSchema(map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of OIDC IdP.",
			},
			"authorization_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IdP Authorization Server (AS) endpoint to request consent from the user and obtain an authorization code grant.",
			},
			"authorization_binding": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The method of making an authorization request. It can be set to `HTTP-POST` or `HTTP-REDIRECT`.",
			},
			"token_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IdP Authorization Server (AS) endpoint to exchange the authorization code grant for an access token.",
			},
			"token_binding": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The method of making a token request. It can be set to `HTTP-POST` or `HTTP-REDIRECT`.",
			},
			"user_info_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Protected resource endpoint that returns claims about the authenticated user.",
			},
			"user_info_binding": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"jwks_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Endpoint where the keys signer publishes its keys in a JWK Set.",
			},
			"jwks_binding": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The method of making a request for the OIDC JWKS. It can be set to `HTTP-POST` or `HTTP-REDIRECT`",
			},
			"scopes": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "The scopes of the IdP.",
			},
			"protocol_type": {
				Type:        schema.TypeString,
				Default:     "OIDC",
				Optional:    true,
				Description: " The type of protocol to use. It can be `OIDC` or `OAUTH2`. Default: `OIDC`",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier issued by AS for the Okta IdP instance.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Client secret issued by AS for the Okta IdP instance.",
			},
			"pkce_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Require Proof Key for Code Exchange (PKCE) for additional verification key rotation mode. See: https://developer.okta.com/docs/reference/api/idps/#oauth-2-0-and-openid-connect-client-object",
			},
			"issuer_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "URI that identifies the issuer.",
			},
			"issuer_mode": {
				Type:        schema.TypeString,
				Description: "Indicates whether Okta uses the original Okta org domain URL, a custom domain URL, or dynamic. It can be `ORG_URL`, `CUSTOM_URL`, or `DYNAMIC`. Default: `ORG_URL`",
				Default:     "ORG_URL",
				Optional:    true,
			},
			"max_clock_skew": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum allowable clock-skew when processing messages from the IdP.",
			},
			"user_type_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User type ID. Can be used as `target_id` in the `okta_profile_mapping` resource.",
			},
			"request_signature_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The HMAC Signature Algorithm used when signing an authorization request. Defaults to `HS256`. It can be `HS256`, `HS384`, `HS512`, `SHA-256`. `RS256`, `RS384`, or `RS512`. NOTE: `SHA-256` an undocumented legacy value and not continue to be valid. See API docs https://developer.okta.com/docs/reference/api/idps/#oidc-request-signature-algorithm-object",
				Default:     "HS256",
			},
			"request_signature_scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies whether to digitally sign an AuthnRequest messages to the IdP. Defaults to `REQUEST`. It can be `REQUEST` or `NONE`.",
				Default:     "REQUEST",
			},
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional regular expression pattern used to filter untrusted IdP usernames.",
			},
		}),
	}
}

func resourceIdpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	idp, err := buildIdPOidc(d)
	if err != nil {
		return diag.FromErr(err)
	}
	respIdp, _, err := getOktaClientFromMetadata(meta).IdentityProvider.CreateIdentityProvider(ctx, idp)
	if err != nil {
		return diag.Errorf("failed to create OIDC identity provider: %v", err)
	}
	d.SetId(respIdp.Id)
	err = setIdpStatus(ctx, d, getOktaClientFromMetadata(meta), idp.Status)
	if err != nil {
		return diag.Errorf("failed to change OIDC identity provider's status: %v", err)
	}
	return resourceIdpRead(ctx, d, meta)
}

func resourceIdpRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	idp, resp, err := getOktaClientFromMetadata(meta).IdentityProvider.GetIdentityProvider(ctx, d.Id())
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get OIDC identity provider: %v", err)
	}
	if idp == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", idp.Name)
	_ = d.Set("type", idp.Type)
	if idp.Policy != nil {
		if idp.Policy.MaxClockSkewPtr != nil {
			_ = d.Set("max_clock_skew", idp.Policy.MaxClockSkewPtr)
		}
		if idp.Policy.Provisioning != nil {
			_ = d.Set("provisioning_action", idp.Policy.Provisioning.Action)
			if idp.Policy.Provisioning.Conditions != nil {
				if idp.Policy.Provisioning.Conditions.Deprovisioned != nil {
					_ = d.Set("deprovisioned_action", idp.Policy.Provisioning.Conditions.Deprovisioned.Action)
				}
				if idp.Policy.Provisioning.Conditions.Suspended != nil {
					_ = d.Set("suspended_action", idp.Policy.Provisioning.Conditions.Suspended.Action)
				}
			}
			if idp.Policy.Provisioning.ProfileMaster != nil {
				_ = d.Set("profile_master", idp.Policy.Provisioning.ProfileMaster)
			}
		}
		if idp.Policy.Subject != nil {
			_ = d.Set("subject_match_type", idp.Policy.Subject.MatchType)
			_ = d.Set("username_template", idp.Policy.Subject.UserNameTemplate.Template)
			_ = d.Set("filter", idp.Policy.Subject.Filter)
		}
	}
	if idp.Protocol != nil {
		if idp.Protocol.Issuer != nil {
			_ = d.Set("issuer_url", idp.Protocol.Issuer.Url)
		}
		if idp.Protocol.Credentials != nil {
			if idp.Protocol.Credentials.Client != nil {
				_ = d.Set("client_id", idp.Protocol.Credentials.Client.ClientId)
				_ = d.Set("client_secret", idp.Protocol.Credentials.Client.ClientSecret)
				if idp.Protocol.Credentials.Client.PKCERequired != nil {
					_ = d.Set("pkce_required", idp.Protocol.Credentials.Client.PKCERequired)
				}
			}
		}
	}
	syncEndpoint("authorization", idp.Protocol.Endpoints.Authorization, d)
	syncEndpoint("token", idp.Protocol.Endpoints.Token, d)
	syncEndpoint("user_info", idp.Protocol.Endpoints.UserInfo, d)
	syncEndpoint("jwks", idp.Protocol.Endpoints.Jwks, d)
	syncIdpOidcAlgo(d, idp.Protocol.Algorithms)
	err = syncGroupActions(d, idp.Policy.Provisioning.Groups)
	if err != nil {
		return diag.Errorf("failed to set OIDC identity provider properties: %v", err)
	}
	if idp.IssuerMode != "" {
		_ = d.Set("issuer_mode", idp.IssuerMode)
	}
	mapping, _, err := getProfileMappingBySourceID(ctx, idp.Id, "", meta)
	if err != nil {
		return diag.Errorf("failed to get identity provider profile mapping: %v", err)
	}
	if mapping != nil {
		_ = d.Set("user_type_id", mapping.Target.Id)
	}
	setMap := map[string]interface{}{
		"scopes": utils.ConvertStringSliceToSet(idp.Protocol.Scopes),
	}
	if idp.Policy.AccountLink != nil {
		_ = d.Set("account_link_action", idp.Policy.AccountLink.Action)
		if idp.Policy.AccountLink.Filter != nil && idp.Policy.AccountLink.Filter.Groups != nil {
			setMap["account_link_group_include"] = utils.ConvertStringSliceToSet(idp.Policy.AccountLink.Filter.Groups.Include)
		}
	}
	err = utils.SetNonPrimitives(d, setMap)
	if err != nil {
		return diag.Errorf("failed to set OIDC identity provider properties: %v", err)
	}
	return nil
}

func resourceIdpUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	idp, err := buildIdPOidc(d)
	if err != nil {
		return diag.FromErr(err)
	}
	_, _, err = getOktaClientFromMetadata(meta).IdentityProvider.UpdateIdentityProvider(ctx, d.Id(), idp)
	if err != nil {
		return diag.Errorf("failed to update OIDC identity provider: %v", err)
	}
	err = setIdpStatus(ctx, d, getOktaClientFromMetadata(meta), idp.Status)
	if err != nil {
		return diag.Errorf("failed to update OIDC identity provider's status: %v", err)
	}
	return resourceIdpRead(ctx, d, meta)
}

func buildIdPOidc(d *schema.ResourceData) (sdk.IdentityProvider, error) {
	if d.Get("subject_match_type").(string) != "CUSTOM_ATTRIBUTE" &&
		len(d.Get("subject_match_attribute").(string)) > 0 {
		return sdk.IdentityProvider{}, errors.New("you can only provide 'subject_match_attribute' with 'subject_match_type' set to 'CUSTOM_ATTRIBUTE'")
	}
	client := &sdk.IdentityProviderCredentialsClient{
		ClientId:     d.Get("client_id").(string),
		ClientSecret: d.Get("client_secret").(string),
	}
	pkceVal := d.GetRawConfig().GetAttr("pkce_required")
	if !pkceVal.IsNull() {
		client.PKCERequired = utils.BoolPtr(d.Get("pkce_required").(bool))
	}
	idp := sdk.IdentityProvider{
		Name:       d.Get("name").(string),
		Type:       "OIDC",
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
				Filter: d.Get("filter").(string),
			},
		},
		Protocol: &sdk.Protocol{
			Algorithms: buildAlgorithms(d),
			Endpoints:  buildProtocolEndpoints(d),
			Scopes:     utils.ConvertInterfaceToStringSet(d.Get("scopes")),
			Type:       d.Get("protocol_type").(string),
			Credentials: &sdk.IdentityProviderCredentials{
				Client: client,
			},
			Issuer: &sdk.ProtocolEndpoint{
				Url: d.Get("issuer_url").(string),
			},
		},
	}
	if d.Get("status") != nil {
		idp.Status = d.Get("status").(string)
	}
	return idp, nil
}

func syncIdpOidcAlgo(d *schema.ResourceData, alg *sdk.ProtocolAlgorithms) {
	if alg != nil {
		if alg.Request != nil && alg.Request.Signature != nil {
			_ = d.Set("request_signature_algorithm", alg.Request.Signature.Algorithm)
			_ = d.Set("request_signature_scope", alg.Request.Signature.Scope)
		}
	}
}
