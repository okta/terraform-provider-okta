package okta

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
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
		Schema: buildIdpSchema(map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"acs_binding": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: elemInSlice([]string{"HTTP-POST", "HTTP-REDIRECT"}),
				Deprecated:       "This property will be removed in the future, as it can only be set to 'HTTP-POST'",
				DiffSuppressFunc: func(string, string, string, *schema.ResourceData) bool {
					return true
				},
			},
			"acs_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "INSTANCE",
				ValidateDiagFunc: elemInSlice([]string{"INSTANCE", "ORG"}),
			},
			"sso_url": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: stringIsURL(validURLSchemes...),
			},
			"sso_binding": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: elemInSlice([]string{postBindingAlias, redirectBindingAlias}),
				Default:          postBindingAlias,
			},
			"sso_destination": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_format": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
			},
			"subject_format": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"subject_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"issuer": {
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer_mode": issuerMode,
			"audience": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"max_clock_skew": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"user_type_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		}),
	}
}

func resourceIdpSamlCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	idp, err := buildIdPSaml(d)
	if err != nil {
		return diag.FromErr(err)
	}
	respIdp, _, err := getOktaClientFromMetadata(m).IdentityProvider.CreateIdentityProvider(ctx, idp)
	if err != nil {
		return diag.Errorf("failed to create SAML identity provider: %v", err)
	}
	d.SetId(respIdp.Id)
	err = setIdpStatus(ctx, d, getOktaClientFromMetadata(m), idp.Status)
	if err != nil {
		return diag.Errorf("failed to change SAML identity provider's status: %v", err)
	}
	return resourceIdpSamlRead(ctx, d, m)
}

func resourceIdpSamlRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	idp, resp, err := getOktaClientFromMetadata(m).IdentityProvider.GetIdentityProvider(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get SAML identity provider: %v", err)
	}
	if idp == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", idp.Name)
	_ = d.Set("acs_type", idp.Protocol.Endpoints.Acs.Type)
	_ = d.Set("max_clock_skew", idp.Policy.MaxClockSkew)
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
	syncAlgo(d, idp.Protocol.Algorithms)
	err = syncGroupActions(d, idp.Policy.Provisioning.Groups)
	if err != nil {
		return diag.Errorf("failed to set SAML identity provider properties: %v", err)
	}
	if idp.IssuerMode != "" {
		_ = d.Set("issuer_mode", idp.IssuerMode)
	}
	mapping, _, err := getProfileMappingBySourceID(ctx, idp.Id, "", m)
	if err != nil {
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

func resourceIdpSamlUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	idp, err := buildIdPSaml(d)
	if err != nil {
		return diag.FromErr(err)
	}
	_, _, err = getOktaClientFromMetadata(m).IdentityProvider.UpdateIdentityProvider(ctx, d.Id(), idp)
	if err != nil {
		return diag.Errorf("failed to update SAML identity provider: %v", err)
	}
	err = setIdpStatus(ctx, d, getOktaClientFromMetadata(m), idp.Status)
	if err != nil {
		return diag.Errorf("failed to update SAML identity provider's status: %v", err)
	}
	return resourceIdpSamlRead(ctx, d, m)
}

func buildIdPSaml(d *schema.ResourceData) (okta.IdentityProvider, error) {
	if d.Get("subject_match_type").(string) != "CUSTOM_ATTRIBUTE" &&
		len(d.Get("subject_match_attribute").(string)) > 0 {
		return okta.IdentityProvider{}, errors.New("you can only provide 'subject_match_attribute' with 'subject_match_type' set to 'CUSTOM_ATTRIBUTE'")
	}
	return okta.IdentityProvider{
		Name:       d.Get("name").(string),
		Type:       saml2Idp,
		IssuerMode: d.Get("issuer_mode").(string),
		Policy: &okta.IdentityProviderPolicy{
			AccountLink:  buildPolicyAccountLink(d),
			Provisioning: buildIdPProvisioning(d),
			Subject: &okta.PolicySubject{
				Filter:         d.Get("subject_filter").(string),
				Format:         convertInterfaceToStringSet(d.Get("subject_format")),
				MatchType:      d.Get("subject_match_type").(string),
				MatchAttribute: d.Get("subject_match_attribute").(string),
				UserNameTemplate: &okta.PolicyUserNameTemplate{
					Template: d.Get("username_template").(string),
				},
			},
			MaxClockSkew: int64(d.Get("max_clock_skew").(int)),
		},
		Protocol: &okta.Protocol{
			Algorithms: buildAlgorithms(d),
			Endpoints: &okta.ProtocolEndpoints{
				Acs: &okta.ProtocolEndpoint{
					// ACS endpoint can only be HTTP-POST
					// https://developer.okta.com/docs/reference/api/idps/#assertion-consumer-service-acs-endpoint-object
					Binding: "HTTP-POST",
					Type:    d.Get("acs_type").(string),
				},
				Sso: &okta.ProtocolEndpoint{
					Binding:     d.Get("sso_binding").(string),
					Destination: d.Get("sso_destination").(string),
					Url:         d.Get("sso_url").(string),
				},
			},
			Type: saml2Idp,
			Credentials: &okta.IdentityProviderCredentials{
				Trust: &okta.IdentityProviderCredentialsTrust{
					Issuer:   d.Get("issuer").(string),
					Kid:      d.Get("kid").(string),
					Audience: d.Get("audience").(string),
				},
			},
		},
	}, nil
}
