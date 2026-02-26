package idaas

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
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
			"honor_persistent_name_id": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines if the IdP should persist account linking when the incoming assertion NameID format is urn:oasis:names:tc:SAML:2.0:nameid-format:persistent",
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
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get SAML identity provider: %v", err)
	}
	if idp == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", idp.Name)
	if idp.Protocol != nil {
		if idp.Protocol.Endpoints != nil {
			if idp.Protocol.Endpoints.Acs != nil {
				_ = d.Set("acs_binding", idp.Protocol.Endpoints.Acs.Binding)
				_ = d.Set("acs_type", idp.Protocol.Endpoints.Acs.Type)
			}
			if idp.Protocol.Endpoints.Sso != nil {
				_ = d.Set("sso_binding", idp.Protocol.Endpoints.Sso.Binding)
				_ = d.Set("sso_destination", idp.Protocol.Endpoints.Sso.Destination)
				_ = d.Set("sso_url", idp.Protocol.Endpoints.Sso.Url)
			}
		}
	}

	if idp.Policy.MaxClockSkewPtr != nil {
		_ = d.Set("max_clock_skew", idp.Policy.MaxClockSkewPtr)
	}

	if idp.Policy.Provisioning != nil {
		_ = d.Set("provisioning_action", idp.Policy.Provisioning.Action)
		if idp.Policy.Provisioning.Conditions.Deprovisioned != nil {
			_ = d.Set("deprovisioned_action", idp.Policy.Provisioning.Conditions.Deprovisioned.Action)
		}
		_ = d.Set("profile_master", idp.Policy.Provisioning.ProfileMaster)
		_ = d.Set("suspended_action", idp.Policy.Provisioning.Conditions.Suspended.Action)
	}
	if idp.Policy.Subject != nil {
		_ = d.Set("subject_match_type", idp.Policy.Subject.MatchType)
		_ = d.Set("subject_match_attribute", idp.Policy.Subject.MatchAttribute)
		_ = d.Set("subject_filter", idp.Policy.Subject.Filter)
		if idp.Policy.Subject.UserNameTemplate != nil {
			_ = d.Set("username_template", idp.Policy.Subject.UserNameTemplate.Template)
		}
	}
	if idp.Protocol != nil {
		if idp.Protocol.Credentials != nil {
			if idp.Protocol.Credentials.Trust != nil {
				_ = d.Set("issuer", idp.Protocol.Credentials.Trust.Issuer)
				_ = d.Set("audience", idp.Protocol.Credentials.Trust.Audience)
				_ = d.Set("kid", idp.Protocol.Credentials.Trust.Kid)
			}
		}
		if idp.Protocol.Settings != nil {
			_ = d.Set("name_format", idp.Protocol.Settings.NameFormat)
			_ = d.Set("honor_persistent_name_id", idp.Protocol.Settings.HonorPersistentNameId)
		}
		syncIdpSamlAlgo(d, idp.Protocol.Algorithms)
	}
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
	if err := utils.SuppressErrorOn401("resource okta_idp_saml.user_type_id", meta, resp, err); err != nil {
		return diag.Errorf("failed to get SAML identity provider profile mapping: %v", err)
	}
	if mapping != nil {
		_ = d.Set("user_type_id", mapping.Target.Id)
	}
	setMap := map[string]interface{}{
		"subject_format": utils.ConvertStringSliceToSet(idp.Policy.Subject.Format),
	}
	if idp.Policy.AccountLink != nil {
		_ = d.Set("account_link_action", idp.Policy.AccountLink.Action)
		if idp.Policy.AccountLink.Filter != nil && idp.Policy.AccountLink.Filter.Groups != nil {
			setMap["account_link_group_include"] = utils.ConvertStringSliceToSet(idp.Policy.AccountLink.Filter.Groups.Include)
		}
	}
	err = utils.SetNonPrimitives(d, setMap)
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
		Type:       Saml2Idp,
		IssuerMode: d.Get("issuer_mode").(string),
		Policy: &sdk.IdentityProviderPolicy{
			AccountLink:  buildPolicyAccountLink(d),
			Provisioning: buildIdPProvisioning(d),
			Subject: &sdk.PolicySubject{
				Filter:         d.Get("subject_filter").(string),
				Format:         utils.ConvertInterfaceToStringSet(d.Get("subject_format")),
				MatchType:      d.Get("subject_match_type").(string),
				MatchAttribute: d.Get("subject_match_attribute").(string),
				UserNameTemplate: &sdk.PolicyUserNameTemplate{
					Template: d.Get("username_template").(string),
				},
			},
			MaxClockSkewPtr: utils.Int64Ptr(d.Get("max_clock_skew").(int)),
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
			Type: Saml2Idp,
			Credentials: &sdk.IdentityProviderCredentials{
				Trust: &sdk.IdentityProviderCredentialsTrust{
					Issuer:   d.Get("issuer").(string),
					Kid:      d.Get("kid").(string),
					Audience: d.Get("audience").(string),
				},
			},
			Settings: &sdk.ProtocolSettings{
				NameFormat:            d.Get("name_format").(string),
				HonorPersistentNameId: d.Get("honor_persistent_name_id").(bool),
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
