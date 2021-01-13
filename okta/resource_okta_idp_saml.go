package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
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
			"acs_binding": bindingSchema,
			"acs_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "INSTANCE",
				ValidateDiagFunc: stringInSlice([]string{"INSTANCE", "ORG"}),
			},
			"sso_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sso_binding": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringInSlice([]string{postBindingAlias, redirectBindingAlias}),
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
		}),
	}
}

func resourceIdpSamlCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	idp := buildIdPSaml(d)
	_, _, err := getSupplementFromMetadata(m).CreateIdentityProvider(ctx, idp, nil)
	if err != nil {
		return diag.Errorf("failed to create SAML identity provider: %v", err)
	}
	d.SetId(idp.ID)
	err = setIdpStatus(ctx, d, getOktaClientFromMetadata(m), idp.Status)
	if err != nil {
		return diag.Errorf("failed to change SAML identity provider's status: %v", err)
	}
	return resourceIdpSamlRead(ctx, d, m)
}

func resourceIdpSamlRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	idp := &sdk.SAMLIdentityProvider{}
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
	_ = d.Set("profile_master", idp.Policy.Provisioning.ProfileMaster)
	_ = d.Set("suspended_action", idp.Policy.Provisioning.Conditions.Suspended.Action)
	_ = d.Set("subject_match_type", idp.Policy.Subject.MatchType)
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
	setMap := map[string]interface{}{
		"subject_format": convertStringSetToInterface(idp.Policy.Subject.Format),
	}
	if idp.Policy.AccountLink != nil {
		_ = d.Set("account_link_action", idp.Policy.AccountLink.Action)
		if idp.Policy.AccountLink.Filter != nil {
			setMap["account_link_group_include"] = convertStringSetToInterface(idp.Policy.AccountLink.Filter.Groups.Include)
		}
	}
	err = setNonPrimitives(d, setMap)
	if err != nil {
		return diag.Errorf("failed to set SAML identity provider properties: %v", err)
	}
	return nil
}

func resourceIdpSamlUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	idp := buildIdPSaml(d)
	_, _, err := getSupplementFromMetadata(m).UpdateIdentityProvider(ctx, d.Id(), idp, nil)
	if err != nil {
		return diag.Errorf("failed to update SAML identity provider: %v", err)
	}
	err = setIdpStatus(ctx, d, getOktaClientFromMetadata(m), idp.Status)
	if err != nil {
		return diag.Errorf("failed to update SAML identity provider's status: %v", err)
	}
	return resourceIdpSamlRead(ctx, d, m)
}

func buildIdPSaml(d *schema.ResourceData) *sdk.SAMLIdentityProvider {
	return &sdk.SAMLIdentityProvider{
		Name:       d.Get("name").(string),
		Type:       "SAML2",
		IssuerMode: d.Get("issuer_mode").(string),
		Policy: &sdk.SAMLPolicy{
			AccountLink:  NewAccountLink(d),
			Provisioning: NewIdpProvisioning(d),
			Subject: &sdk.SAMLSubject{
				Filter:    d.Get("subject_filter").(string),
				Format:    convertInterfaceToStringSet(d.Get("subject_format")),
				MatchType: d.Get("subject_match_type").(string),
				UserNameTemplate: &okta.ApplicationCredentialsUsernameTemplate{
					Template: d.Get("username_template").(string),
				},
			},
			MaxClockSkew: int64(d.Get("max_clock_skew").(int)),
		},
		Protocol: &sdk.SAMLProtocol{
			Algorithms: NewAlgorithms(d),
			Endpoints: &sdk.SAMLEndpoints{
				Acs: &sdk.ACSSSO{
					Binding: d.Get("acs_binding").(string),
					Type:    d.Get("acs_type").(string),
				},
				Sso: &sdk.IDPSSO{
					Binding:     d.Get("sso_binding").(string),
					Destination: d.Get("sso_destination").(string),
					URL:         d.Get("sso_url").(string),
				},
			},
			Type: "SAML2",
			Credentials: &sdk.SAMLCredentials{
				Trust: &sdk.IDPTrust{
					Issuer:   d.Get("issuer").(string),
					Kid:      d.Get("kid").(string),
					Audience: d.Get("audience").(string),
				},
			},
		},
	}
}
