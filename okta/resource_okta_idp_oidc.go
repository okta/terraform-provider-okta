package okta

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		// Note the base schema
		Schema: buildIdpSchema(map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of OIDC IdP.",
			},
			"authorization_url":     urlSchema,
			"authorization_binding": bindingSchema,
			"token_url":             urlSchema,
			"token_binding":         bindingSchema,
			"user_info_url":         optionalURLSchema,
			"user_info_binding":     optionalBindingSchema,
			"jwks_url":              urlSchema,
			"jwks_binding":          bindingSchema,
			"scopes": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"protocol_type": {
				Type:     schema.TypeString,
				Default:  "OIDC",
				Optional: true,
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier (opens new window)issued by the AS for the Okta IdP instance",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Client secret issued (opens new window)by the AS for the Okta IdP instance",
			},
			"pkce_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Require Proof Key for Code Exchange (PKCE) for additional verification key rotation mode. See: https://developer.okta.com/docs/reference/api/idps/#oauth-2-0-and-openid-connect-client-object",
			},
			"issuer_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer_mode": {
				Type:        schema.TypeString,
				Description: "Indicates whether Okta uses the original Okta org domain URL, custom domain URL, or dynamic. See Identity Provider attributes - issuerMode - https://developer.okta.com/docs/reference/api/idps/#identity-provider-attributes",
				Default:     "ORG_URL",
				Optional:    true,
			},
			"max_clock_skew": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"user_type_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_signature_algorithm": oidcRequestSignatureAlgorithmSchema,
			"request_signature_scope":     oidcRequestSignatureScopeSchema,
		}),
	}
}

func resourceIdpCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	idp, err := buildIdPOidc(d)
	if err != nil {
		return diag.FromErr(err)
	}
	respIdp, _, err := getOktaClientFromMetadata(m).IdentityProvider.CreateIdentityProvider(ctx, idp)
	if err != nil {
		return diag.Errorf("failed to create OIDC identity provider: %v", err)
	}
	d.SetId(respIdp.Id)
	err = setIdpStatus(ctx, d, getOktaClientFromMetadata(m), idp.Status)
	if err != nil {
		return diag.Errorf("failed to change OIDC identity provider's status: %v", err)
	}
	return resourceIdpRead(ctx, d, m)
}

func resourceIdpRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	idp, resp, err := getOktaClientFromMetadata(m).IdentityProvider.GetIdentityProvider(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get OIDC identity provider: %v", err)
	}
	if idp == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", idp.Name)
	_ = d.Set("type", idp.Type)
	if idp.Policy.MaxClockSkewPtr != nil {
		_ = d.Set("max_clock_skew", *idp.Policy.MaxClockSkewPtr)
	}
	_ = d.Set("provisioning_action", idp.Policy.Provisioning.Action)
	_ = d.Set("deprovisioned_action", idp.Policy.Provisioning.Conditions.Deprovisioned.Action)
	_ = d.Set("suspended_action", idp.Policy.Provisioning.Conditions.Suspended.Action)
	_ = d.Set("profile_master", idp.Policy.Provisioning.ProfileMaster)
	_ = d.Set("subject_match_type", idp.Policy.Subject.MatchType)
	_ = d.Set("username_template", idp.Policy.Subject.UserNameTemplate.Template)
	_ = d.Set("issuer_url", idp.Protocol.Issuer.Url)
	_ = d.Set("client_secret", idp.Protocol.Credentials.Client.ClientSecret)
	_ = d.Set("client_id", idp.Protocol.Credentials.Client.ClientId)
	if idp.Protocol.Credentials.Client.PKCERequired != nil {
		_ = d.Set("pkce_required", idp.Protocol.Credentials.Client.PKCERequired)
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
	mapping, _, err := getProfileMappingBySourceID(ctx, idp.Id, "", m)
	if err != nil {
		return diag.Errorf("failed to get identity provider profile mapping: %v", err)
	}
	if mapping != nil {
		_ = d.Set("user_type_id", mapping.Target.Id)
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
		return diag.Errorf("failed to set OIDC identity provider properties: %v", err)
	}
	return nil
}

func resourceIdpUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	idp, err := buildIdPOidc(d)
	if err != nil {
		return diag.FromErr(err)
	}
	_, _, err = getOktaClientFromMetadata(m).IdentityProvider.UpdateIdentityProvider(ctx, d.Id(), idp)
	if err != nil {
		return diag.Errorf("failed to update OIDC identity provider: %v", err)
	}
	err = setIdpStatus(ctx, d, getOktaClientFromMetadata(m), idp.Status)
	if err != nil {
		return diag.Errorf("failed to update OIDC identity provider's status: %v", err)
	}
	return resourceIdpRead(ctx, d, m)
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
		client.PKCERequired = boolPtr(d.Get("pkce_required").(bool))
	}
	idp := sdk.IdentityProvider{
		Name:       d.Get("name").(string),
		Type:       "OIDC",
		IssuerMode: d.Get("issuer_mode").(string),
		Policy: &sdk.IdentityProviderPolicy{
			AccountLink:     buildPolicyAccountLink(d),
			MaxClockSkewPtr: int64Ptr(d.Get("max_clock_skew").(int)),
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
			Algorithms: buildAlgorithms(d),
			Endpoints:  buildProtocolEndpoints(d),
			Scopes:     convertInterfaceToStringSet(d.Get("scopes")),
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
