package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceSecurityNotificationEmails() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityNotificationEmailsCreate,
		ReadContext:   resourceSecurityNotificationEmailsRead,
		UpdateContext: resourceSecurityNotificationEmailsUpdate,
		DeleteContext: resourceSecurityNotificationEmailsDelete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
		Schema: map[string]*schema.Schema{
			"send_email_for_new_device_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Notifies end users about new sign-on activity",
				Default:     true,
			},
			"send_email_for_factor_enrollment_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Notifies end users of any activity on their account related to MFA factor enrollment",
				Default:     true,
			},
			"send_email_for_factor_reset_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Notifies end users that one or more factors have been reset for their account",
				Default:     true,
			},
			"send_email_for_password_changed_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Notifies end users that the password for their account has changed",
				Default:     true,
			},
			"report_suspicious_activity_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Notifies end users about suspicious or unrecognized activity from their account",
				Default:     true,
			},
		},
	}
}

func resourceSecurityNotificationEmailsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	client := c.oktaClient.GetConfig().HttpClient
	emails, err := getAPISupplementFromMetadata(m).UpdateSecurityNotificationEmails(ctx, buildSecurityNotificationEmails(d), c.orgName, c.domain, c.apiToken, client)
	if err != nil {
		return diag.Errorf("failed to update security notification emails: %v", err)
	}
	d.SetId("security_notification_emails")
	_ = d.Set("send_email_for_new_device_enabled", emails.SendEmailForNewDeviceEnabled)
	_ = d.Set("send_email_for_factor_enrollment_enabled", emails.SendEmailForFactorEnrollmentEnabled)
	_ = d.Set("send_email_for_factor_reset_enabled", emails.SendEmailForFactorResetEnabled)
	_ = d.Set("send_email_for_password_changed_enabled", emails.SendEmailForPasswordChangedEnabled)
	_ = d.Set("report_suspicious_activity_enabled", emails.ReportSuspiciousActivityEnabled)
	return nil
}

func resourceSecurityNotificationEmailsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	client := c.oktaClient.GetConfig().HttpClient
	emails, err := getAPISupplementFromMetadata(m).GetSecurityNotificationEmails(ctx, c.orgName, c.domain, c.apiToken, client)
	if err != nil {
		return diag.Errorf("failed to get security notification emails: %v", err)
	}
	d.SetId("security_notification_emails")
	_ = d.Set("send_email_for_new_device_enabled", emails.SendEmailForNewDeviceEnabled)
	_ = d.Set("send_email_for_factor_enrollment_enabled", emails.SendEmailForFactorEnrollmentEnabled)
	_ = d.Set("send_email_for_factor_reset_enabled", emails.SendEmailForFactorResetEnabled)
	_ = d.Set("send_email_for_password_changed_enabled", emails.SendEmailForPasswordChangedEnabled)
	_ = d.Set("report_suspicious_activity_enabled", emails.ReportSuspiciousActivityEnabled)
	return nil
}

func resourceSecurityNotificationEmailsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	client := c.oktaClient.GetConfig().HttpClient
	_, err := getAPISupplementFromMetadata(m).UpdateSecurityNotificationEmails(ctx, buildSecurityNotificationEmails(d), c.orgName, c.domain, c.apiToken, client)
	if err != nil {
		return diag.Errorf("failed to update security notification emails: %v", err)
	}
	return nil
}

func resourceSecurityNotificationEmailsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)
	client := c.oktaClient.GetConfig().HttpClient
	emails := sdk.SecurityNotificationEmails{
		SendEmailForNewDeviceEnabled:        true,
		SendEmailForFactorEnrollmentEnabled: true,
		SendEmailForFactorResetEnabled:      true,
		SendEmailForPasswordChangedEnabled:  true,
		ReportSuspiciousActivityEnabled:     true,
	}
	_, err := getAPISupplementFromMetadata(m).UpdateSecurityNotificationEmails(ctx, emails, c.orgName, c.domain, c.apiToken, client)
	if err != nil {
		return diag.Errorf("failed to set default security notification emails: %v", err)
	}
	return nil
}

func buildSecurityNotificationEmails(d *schema.ResourceData) sdk.SecurityNotificationEmails {
	return sdk.SecurityNotificationEmails{
		SendEmailForNewDeviceEnabled:        d.Get("send_email_for_new_device_enabled").(bool),
		SendEmailForFactorEnrollmentEnabled: d.Get("send_email_for_factor_enrollment_enabled").(bool),
		SendEmailForFactorResetEnabled:      d.Get("send_email_for_factor_reset_enabled").(bool),
		SendEmailForPasswordChangedEnabled:  d.Get("send_email_for_password_changed_enabled").(bool),
		ReportSuspiciousActivityEnabled:     d.Get("report_suspicious_activity_enabled").(bool),
	}
}
