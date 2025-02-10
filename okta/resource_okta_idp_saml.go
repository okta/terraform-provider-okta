package okta

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceIdpSaml() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdpSamlCreate,
		ReadContext:   resourceIdpSamlRead,
		UpdateContext: resourceIdpSamlUpdate,
		DeleteContext: resourceIdpDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Creates a SAML Identity Provider. This resource allows you to create and configure a SAML Identity Provider.",
		Schema: buildIdpSchema(map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"acs_binding": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"acs_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "INSTANCE",
				Description: "The type of ACS. It can be `INSTANCE` or `ORG`. Default: `INSTANCE`",
			},
			"sso_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "URL of binding-specific endpoint to send an AuthnRequest message to IdP.",
			},
			"sso_binding": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     postBindingAlias,
				Description: "The method of making an SSO request. It can be set to `HTTP-POST` or `HTTP-REDIRECT`. Default: `HTTP-POST`",
			},
			"sso_destination": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URI reference indicating the address to which the AuthnRequest message is sent.",
			},
			"name_format": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
				Description: "The name identifier format to use. By default `urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified`.",
			},
			"subject_format": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "The name format.",
			},
			"subject_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional regular expression pattern used to filter untrusted IdP usernames.",
			},
			"issuer": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "URI that identifies the issuer.",
			},
			"issuer_mode": issuerMode,
			"audience": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the signing key.",
			},
			"max_clock_skew": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum allowable clock-skew when processing messages from the IdP.",
			},
			"user_type_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_signature_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The XML digital Signature Algorithm used when signing an `AuthnRequest` message. It can be `SHA-256` or `SHA-1`. Default: `SHA-256`",
				Default:     "SHA-256",
			},
			"request_signature_scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies whether to digitally sign an AuthnRequest messages to the IdP. It can be `REQUEST` or `NONE`. Default: `REQUEST`",
				Default:     "REQUEST",
			},
			"response_signature_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The minimum XML digital signature algorithm allowed when verifying a `SAMLResponse` message or Assertion element. It can be `SHA-256` or `SHA-1`. Default: `SHA-256`",
				Default:     "SHA-256",
			},
			"response_signature_scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies whether to verify a `SAMLResponse` message or Assertion element XML digital signature. It can be `RESPONSE`, `ASSERTION`, or `ANY`. Default: `ANY`",
				Default:     "ANY",
			},
		}),
	}
}

func resourceIdpSamlCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	idp, err := buildIdPSaml(d)
	if err != nil {
		return diag.FromErr(err)
	}
	respIdp, _, err := getOktaClientFromMetadata(meta).IdentityProvider.CreateIdentityProvider(ctx, idp)
	if err != nil {
		return diag.Errorf("failed to create SAML identity provider: %v", err)
	}
	d.SetId(respIdp.Id)
	err = setIdpStatus(ctx, d, getOktaClientFromMetadata(meta), idp.Status)
	if err != nil {
		return diag.Errorf("failed to change SAML identity provider's status: %v", err)
	}
	return resourceIdpSamlRead(ctx, d, meta)
}

func resourceIdpSamlRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	idp, resp, err := getOktaClientFromMetadata(meta).IdentityProvider.GetIdentityProvider(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get SAML identity provider: %v", err)
	}
	if idp == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", idp.Name)
	_ = d.Set("acs_binding", idp.Protocol.Endpoints.Acs.Binding)
	_ = d.Set("acs_type", idp.Protocol.Endpoints.Acs.Type)
	if idp.Protocol.Endpoints.Sso != nil {
		_ = d.Set("sso_binding", idp.Protocol.Endpoints.Sso.Binding)
		_ = d.Set("sso_destination", idp.Protocol.Endpoints.Sso.Destination)
		_ = d.Set("sso_url", idp.Protocol.Endpoints.Sso.Url)
	}
	if idp.Policy.MaxClockSkewPtr != nil {
		_ = d.Set("max_clock_skew", idp.Policy.MaxClockSkewPtr)
	}
	_ = d.Set("provisioning_action", idp.Policy.Provisioning.Action)
	_ = d.Set("deprovisioned_action", idp.Policy.Provisioning.Conditions.Deprovisioned.Action)
	_ = d.Set("profile_master", idp.Policy.Provisioning.ProfileMaster)
	_ = d.Set("suspended_action", idp.Policy.Provisioning.Conditions.Suspended.Action)
	_ = d.Set("subject_match_type", idp.Policy.Subject.MatchType)
	_ = d.Set("subject_match_attribute", idp.Policy.Subject.MatchAttribute)
	_ = d.Set("subject_filter", idp.Policy.Subject.Filter)
	_ = d.Set("username_template", idp.Policy.Subject.UserNameTemplate.Template)
	_ = d.Set("issuer", idp.Protocol.Credentials.Trust.Issuer)
	_ = d.Set("audience", idp.Protocol.Credentials.Trust.Audience)
	_ = d.Set("kid", idp.Protocol.Credentials.Trust.Kid)
	_ = d.Set("name_format", idp.Protocol.Settings.NameFormat)
	syncIdpSamlAlgo(d, idp.Protocol.Algorithms)
	err = syncGroupActions(d, idp.Policy.Provisioning.Groups)
	if err != nil {
		return diag.Errorf("failed to set SAML identity provider properties: %v", err)
	}
	if idp.IssuerMode != "" {
		_ = d.Set("issuer_mode", idp.IssuerMode)
	}
	if idp.Status != "" {
		_ = d.Set("status", idp.Status)
	}
	mapping, resp, err := getProfileMappingBySourceID(ctx, idp.Id, "", meta)
	if err := suppressErrorOn401("resource okta_idp_saml.user_type_id", meta, resp, err); err != nil {
		return diag.Errorf("failed to get SAML identity provider profile mapping: %v", err)
	}
	if mapping != nil {
		_ = d.Set("user_type_id", mapping.Target.Id)
	}
	setMap := map[string]interface{}{
		"subject_format": convertStringSliceToSet(idp.Policy.Subject.Format),
	}
	if idp.Policy.AccountLink != nil {
		_ = d.Set("account_link_action", idp.Policy.AccountLink.Action)
		if idp.Policy.AccountLink.Filter != nil {
			setMap["account_link_group_include"] = convertStringSliceToSet(idp.Policy.AccountLink.Filter.Groups.Include)
		}
	}
	err = setNonPrimitives(d, setMap)
	if err != nil {
		return diag.Errorf("failed to set SAML identity provider properties: %v", err)
	}
	return nil
}

func resourceIdpSamlUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	idp, err := buildIdPSaml(d)
	if err != nil {
		return diag.FromErr(err)
	}
	_, _, err = getOktaClientFromMetadata(meta).IdentityProvider.UpdateIdentityProvider(ctx, d.Id(), idp)
	if err != nil {
		return diag.Errorf("failed to update SAML identity provider: %v", err)
	}
	err = setIdpStatus(ctx, d, getOktaClientFromMetadata(meta), idp.Status)
	if err != nil {
		return diag.Errorf("failed to update SAML identity provider's status: %v", err)
	}
	return resourceIdpSamlRead(ctx, d, meta)
}

func buildIdPSaml(d *schema.ResourceData) (sdk.IdentityProvider, error) {
	if d.Get("subject_match_type").(string) != "CUSTOM_ATTRIBUTE" &&
		len(d.Get("subject_match_attribute").(string)) > 0 {
		return sdk.IdentityProvider{}, errors.New("you can only provide 'subject_match_attribute' with 'subject_match_type' set to 'CUSTOM_ATTRIBUTE'")
	}
	idp := sdk.IdentityProvider{
		Name:       d.Get("name").(string),
		Type:       saml2Idp,
		IssuerMode: d.Get("issuer_mode").(string),
		Policy: &sdk.IdentityProviderPolicy{
			AccountLink:  buildPolicyAccountLink(d),
			Provisioning: buildIdPProvisioning(d),
			Subject: &sdk.PolicySubject{
				Filter:         d.Get("subject_filter").(string),
				Format:         convertInterfaceToStringSet(d.Get("subject_format")),
				MatchType:      d.Get("subject_match_type").(string),
				MatchAttribute: d.Get("subject_match_attribute").(string),
				UserNameTemplate: &sdk.PolicyUserNameTemplate{
					Template: d.Get("username_template").(string),
				},
			},
			MaxClockSkewPtr: int64Ptr(d.Get("max_clock_skew").(int)),
		},
		Protocol: &sdk.Protocol{
			Algorithms: buildAlgorithms(d),
			Endpoints: &sdk.ProtocolEndpoints{
				Acs: &sdk.ProtocolEndpoint{
					// ACS endpoint can only be HTTP-POST
					// https://developer.okta.com/docs/reference/api/idps/#assertion-consumer-service-acs-endpoint-object
					Binding: "HTTP-POST",
					Type:    d.Get("acs_type").(string),
				},
				Sso: &sdk.ProtocolEndpoint{
					Binding:     d.Get("sso_binding").(string),
					Destination: d.Get("sso_destination").(string),
					Url:         d.Get("sso_url").(string),
				},
			},
			Type: saml2Idp,
			Credentials: &sdk.IdentityProviderCredentials{
				Trust: &sdk.IdentityProviderCredentialsTrust{
					Issuer:   d.Get("issuer").(string),
					Kid:      d.Get("kid").(string),
					Audience: d.Get("audience").(string),
				},
			},
			Settings: &sdk.ProtocolSettings{
				NameFormat: d.Get("name_format").(string),
			},
		},
	}
	if d.Get("status") != nil {
		idp.Status = d.Get("status").(string)
	}
	return idp, nil
}

func syncIdpSamlAlgo(d *schema.ResourceData, alg *sdk.ProtocolAlgorithms) {
	if alg != nil {
		if alg.Request != nil && alg.Request.Signature != nil {
			_ = d.Set("request_signature_algorithm", alg.Request.Signature.Algorithm)
			_ = d.Set("request_signature_scope", alg.Request.Signature.Scope)
		}
		if alg.Response != nil && alg.Response.Signature != nil {
			_ = d.Set("response_signature_algorithm", alg.Response.Signature.Algorithm)
			_ = d.Set("response_signature_scope", alg.Response.Signature.Scope)
		}
	}
}
