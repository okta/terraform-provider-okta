package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceOrgConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrgSettingsCreate,
		ReadContext:   resourceOrgSettingsRead,
		UpdateContext: resourceOrgSettingsUpdate,
		DeleteContext: resourceOrgSettingsDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"company_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of org",
			},
			"website": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The org's website",
			},
			"phone_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Support help phone of org",
			},
			"end_user_support_help_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Support link of org",
			},
			"support_phone_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Support help phone of org",
			},
			"address_1": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Primary address of org",
			},
			"address_2": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Secondary address of org",
			},
			"city": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "City of org",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "State of org",
			},
			"country": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Country of org",
			},
			"postal_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Postal code of org",
			},
			"expires_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiration of org",
			},
			"subdomain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subdomain of org",
			},
			"logo": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: logoValid(),
				Description:      "Local path to logo of the org.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == ""
				},
				StateFunc: func(val interface{}) string {
					logoPath := val.(string)
					if logoPath == "" {
						return logoPath
					}
					return fmt.Sprintf("%s (%s)", logoPath, computeFileHash(logoPath))
				},
			},
			"billing_contact_user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User ID representing the billing contact",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == ""
				},
			},
			"technical_contact_user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User ID representing the technical contact",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == ""
				},
			},
			"opt_out_communication_emails": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether the org's users receive Okta Communication emails",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == ""
				},
			},
		},
	}
}

func resourceOrgSettingsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	settings, _, err := getOktaClientFromMetadata(m).OrgSetting.PartialUpdateOrgSetting(ctx, buildOrgSettings(d))
	if err != nil {
		return diag.Errorf("failed to update org settings: %v", err)
	}
	d.SetId(settings.Id)
	logo, ok := d.GetOk("logo")
	if ok {
		_, err := getSupplementFromMetadata(m).UploadOrgLogo(ctx, logo.(string))
		if err != nil {
			return diag.Errorf("failed to upload org logo: %v", err)
		}
	}
	err = updateCommunicationSettings(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	err = updateContactUsers(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceOrgSettingsRead(ctx, d, m)
}

func resourceOrgSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	settings, _, err := getOktaClientFromMetadata(m).OrgSetting.GetOrgSettings(ctx)
	if err != nil {
		return diag.Errorf("failed to get org settings: %v", err)
	}
	setOrgSettings(d, settings)
	comm, _, err := getOktaClientFromMetadata(m).OrgSetting.GetOktaCommunicationSettings(ctx)
	if err != nil {
		return diag.Errorf("failed to get org communication settings: %v", err)
	}
	_ = d.Set("opt_out_communication_emails", comm.OptOutEmailUsers)
	billingContact, _, err := getOktaClientFromMetadata(m).OrgSetting.GetOrgContactUser(ctx, "BILLING")
	if err != nil {
		return diag.Errorf("failed to get billing contact user: %v", err)
	}
	_ = d.Set("billing_contact_user", billingContact.UserId)
	technicalContact, _, err := getOktaClientFromMetadata(m).OrgSetting.GetOrgContactUser(ctx, "TECHNICAL")
	if err != nil {
		return diag.Errorf("failed to get technical contact user: %v", err)
	}
	_ = d.Set("technical_contact_user", technicalContact.UserId)
	return nil
}

func resourceOrgSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, _, err := getOktaClientFromMetadata(m).OrgSetting.UpdateOrgSetting(ctx, buildOrgSettings(d))
	if err != nil {
		return diag.Errorf("failed to update org settings: %v", err)
	}
	logo, ok := d.GetOk("logo")
	if ok {
		_, err := getSupplementFromMetadata(m).UploadOrgLogo(ctx, logo.(string))
		if err != nil {
			return diag.Errorf("failed to upload org logo: %v", err)
		}
	}
	err = updateCommunicationSettings(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	err = updateContactUsers(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceOrgSettingsRead(ctx, d, m)
}

func resourceOrgSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func updateCommunicationSettings(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	comm, _, err := getOktaClientFromMetadata(m).OrgSetting.GetOktaCommunicationSettings(ctx)
	if err != nil {
		return fmt.Errorf("failed to get org communication settings: %v", err)
	}
	o, ok := d.GetOk("opt_out_communication_emails")
	if ok && *comm.OptOutEmailUsers != o.(bool) {
		if o.(bool) {
			_, _, err = getOktaClientFromMetadata(m).OrgSetting.OptOutUsersFromOktaCommunicationEmails(ctx)
		} else {
			_, _, err = getOktaClientFromMetadata(m).OrgSetting.OptInUsersToOktaCommunicationEmails(ctx)
		}
		if err != nil {
			return fmt.Errorf("failed to update org communication settings: %v", err)
		}
	}
	return nil
}

func updateContactUsers(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	billingContact, _, err := getOktaClientFromMetadata(m).OrgSetting.GetOrgContactUser(ctx, "BILLING")
	if err != nil {
		return fmt.Errorf("failed to get billing contact user: %v", err)
	}
	billing, ok := d.GetOk("billing_contact_user")
	if ok && billingContact.UserId != billing.(string) {
		_, _, err := getOktaClientFromMetadata(m).OrgSetting.UpdateOrgContactUser(ctx,
			"BILLING", okta.UserIdString{UserId: billing.(string)})
		if err != nil {
			return fmt.Errorf("failed to update billing contact user: %v", err)
		}
	}
	technicalContact, _, err := getOktaClientFromMetadata(m).OrgSetting.GetOrgContactUser(ctx, "TECHNICAL")
	if err != nil {
		return fmt.Errorf("failed to get technical contact user: %v", err)
	}
	technical, ok := d.GetOk("technical_contact_user")
	if ok && technicalContact.UserId != technical.(string) {
		_, _, err := getOktaClientFromMetadata(m).OrgSetting.UpdateOrgContactUser(ctx,
			"TECHNICAL", okta.UserIdString{UserId: technical.(string)})
		if err != nil {
			return fmt.Errorf("failed to update technical contact user: %v", err)
		}
	}
	return nil
}

func setOrgSettings(d *schema.ResourceData, settings *okta.OrgSetting) {
	_ = d.Set("address_1", settings.Address1)
	_ = d.Set("address_2", settings.Address2)
	_ = d.Set("city", settings.City)
	_ = d.Set("company_name", settings.CompanyName)
	_ = d.Set("country", settings.Country)
	_ = d.Set("end_user_support_help_url", settings.EndUserSupportHelpURL)
	_ = d.Set("phone_number", settings.PhoneNumber)
	_ = d.Set("postal_code", settings.PostalCode)
	_ = d.Set("state", settings.State)
	_ = d.Set("support_phone_number", settings.SupportPhoneNumber)
	_ = d.Set("website", settings.Website)
	_ = d.Set("subdomain", settings.Subdomain)
	if settings.ExpiresAt != nil {
		_ = d.Set("expires_at", settings.ExpiresAt.String())
	}
}

func buildOrgSettings(d *schema.ResourceData) okta.OrgSetting {
	return okta.OrgSetting{
		Address1:              d.Get("address_1").(string),
		Address2:              d.Get("address_2").(string),
		City:                  d.Get("city").(string),
		CompanyName:           d.Get("company_name").(string),
		Country:               d.Get("country").(string),
		EndUserSupportHelpURL: d.Get("end_user_support_help_url").(string),
		PhoneNumber:           d.Get("phone_number").(string),
		PostalCode:            d.Get("postal_code").(string),
		State:                 d.Get("state").(string),
		SupportPhoneNumber:    d.Get("support_phone_number").(string),
		Website:               d.Get("website").(string),
	}
}
