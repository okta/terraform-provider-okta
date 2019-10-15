package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/okta/okta-sdk-golang/okta"
	"github.com/terraform-providers/terraform-provider-okta/sdk"
)

func resourceIdpSaml() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdpSamlCreate,
		Read:   resourceIdpSamlRead,
		Update: resourceIdpSamlUpdate,
		Delete: resourceIdpDelete,
		Exists: getIdentityProviderExists(&sdk.SAMLIdentityProvider{}),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: buildIdpSchema(map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"acs_binding": bindingSchema,
			"acs_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "INSTANCE",
				ValidateFunc: validation.StringInSlice([]string{"INSTANCE", "ORG"}, false),
			},
			"sso_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"sso_binding": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice(
					[]string{postBindingAlias, redirectBindingAlias},
					false,
				),
				Default: postBindingAlias,
			},
			"sso_destination": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_format": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified",
			},
			"subject_format": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"subject_filter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"issuer": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer_mode": issuerMode,
			"audience": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"kid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}

func resourceIdpSamlCreate(d *schema.ResourceData, m interface{}) error {
	idp := buildidpSaml(d)
	if err := createIdp(m, idp); err != nil {
		return err
	}
	d.SetId(idp.ID)

	if err := setIdpStatus(d.Id(), idp.Status, d.Get("status").(string), m); err != nil {
		return err
	}

	return resourceIdpSamlRead(d, m)
}

func resourceIdpSamlRead(d *schema.ResourceData, m interface{}) error {
	idp := &sdk.SAMLIdentityProvider{}
	if err := fetchIdp(d.Id(), m, idp); err != nil {
		return err
	}

	if idp == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", idp.Name)
	d.Set("provisioning_action", idp.Policy.Provisioning.Action)
	d.Set("deprovisioned_action", idp.Policy.Provisioning.Conditions.Deprovisioned.Action)
	d.Set("profile_master", idp.Policy.Provisioning.ProfileMaster)
	d.Set("suspended_action", idp.Policy.Provisioning.Conditions.Suspended.Action)
	d.Set("subject_match_type", idp.Policy.Subject.MatchType)
	d.Set("subject_filter", idp.Policy.Subject.Filter)
	d.Set("username_template", idp.Policy.Subject.UserNameTemplate.Template)
	d.Set("issuer", idp.Protocol.Credentials.Trust.Issuer)
	d.Set("audience", idp.Protocol.Credentials.Trust.Audience)
	d.Set("kid", idp.Protocol.Credentials.Trust.Kid)
	syncAlgo(d, idp.Protocol.Algorithms)

	if err := syncGroupActions(d, idp.Policy.Provisioning.Groups); err != nil {
		return err
	}

	if idp.IssuerMode != "" {
		d.Set("issuer_mode", idp.IssuerMode)
	}

	setMap := map[string]interface{}{
		"subject_format": convertStringSetToInterface(idp.Policy.Subject.Format),
	}

	if idp.Policy.AccountLink != nil {
		d.Set("account_link_action", idp.Policy.AccountLink.Action)

		if idp.Policy.AccountLink.Filter != nil {
			setMap["account_link_group_include"] = convertStringSetToInterface(idp.Policy.AccountLink.Filter.Groups.Include)
		}
	}

	return setNonPrimitives(d, setMap)
}

func resourceIdpSamlUpdate(d *schema.ResourceData, m interface{}) error {
	idp := buildidpSaml(d)
	d.Partial(true)

	if err := updateIdp(d.Id(), m, idp); err != nil {
		return err
	}

	d.Partial(false)

	if err := setIdpStatus(idp.ID, idp.Status, d.Get("status").(string), m); err != nil {
		return err
	}

	return resourceIdpSamlRead(d, m)
}

func buildidpSaml(d *schema.ResourceData) *sdk.SAMLIdentityProvider {
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
