package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v4/okta"
)

func resourceEmailTemplateSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEmailTemplateSettingsPut,
		ReadContext:   resourceEmailTemplateSettingsRead,
		UpdateContext: resourceEmailTemplateSettingsPut,
		DeleteContext: resourceFuncNoOp,
		Importer:      createNestedResourceImporter([]string{"brand_id", "template_name"}),
		Description: `Update settings for an email template belonging to a brand in an Okta organization.
		Use this resource to get and set the [settings for an email template](https://developer.okta.com/docs/api/openapi/okta-management/management/tag/CustomTemplates/#tag/CustomTemplates/operation/getEmailSettings) 
        belonging to a brand in an Okta organization.`,
		Schema: map[string]*schema.Schema{
			"brand_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Brand ID",
			},
			"template_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Email template name - Example values: `AccountLockout`,`ADForgotPassword`,`ADForgotPasswordDenied`,`ADSelfServiceUnlock`,`ADUserActivation`,`AuthenticatorEnrolled`,`AuthenticatorReset`,`ChangeEmailConfirmation`,`EmailChallenge`,`EmailChangeConfirmation`,`EmailFactorVerification`,`ForgotPassword`,`ForgotPasswordDenied`,`IGAReviewerEndNotification`,`IGAReviewerNotification`,`IGAReviewerPendingNotification`,`IGAReviewerReassigned`,`LDAPForgotPassword`,`LDAPForgotPasswordDenied`,`LDAPSelfServiceUnlock`,`LDAPUserActivation`,`MyAccountChangeConfirmation`,`NewSignOnNotification`,`OktaVerifyActivation`,`PasswordChanged`,`PasswordResetByAdmin`,`PendingEmailChange`,`RegistrationActivation`,`RegistrationEmailVerification`,`SelfServiceUnlock`,`SelfServiceUnlockOnUnlockedAccount`,`UserActivation`",
			},
			"recipients": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The recipients the emails of this template will be sent to - Valid values: `ALL_USERS`, `ADMINS_ONLY`, `NO_USERS`",
			},
		},
	}
}

func resourceEmailTemplateSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ets, diagErr := etsValues("read", d)
	if diagErr != nil {
		return diagErr
	}

	emailTemplateSettings, resp, err := getOktaV3ClientFromMetadata(m).CustomizationAPI.GetEmailSettings(ctx, ets.brandID, ets.templateName).Execute()
	if err := v3suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get email template settings: %v", err)
	}
	if emailTemplateSettings == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("brand_id", ets.brandID)
	_ = d.Set("template_name", ets.templateName)
	_ = d.Set("recipients", emailTemplateSettings.GetRecipients())

	return nil
}

func resourceEmailTemplateSettingsPut(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ets, diagErr := etsValues("update", d)
	if diagErr != nil {
		return diagErr
	}

	es := okta.EmailSettings{}
	if recipients, ok := d.GetOk("recipients"); ok {
		es.Recipients = recipients.(string)
	}

	_, err := getOktaV3ClientFromMetadata(m).CustomizationAPI.ReplaceEmailSettings(ctx, ets.brandID, ets.templateName).EmailSettings(es).Execute()
	if err != nil {
		return diag.Errorf("failed to update email template settings: %v", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", ets.brandID, ets.templateName))
	return resourceEmailTemplateSettingsRead(ctx, d, m)
}

type etsHelper struct {
	brandID      string
	templateName string
}

func etsValues(action string, d *schema.ResourceData) (*etsHelper, diag.Diagnostics) {
	brandID, ok := d.GetOk("brand_id")
	if !ok {
		return nil, diag.Errorf("brand_id required to %s email template settings", action)
	}

	templateName, ok := d.GetOk("template_name")
	if !ok {
		return nil, diag.Errorf("template name required to %s email template settings", action)
	}

	return &etsHelper{
		brandID:      brandID.(string),
		templateName: templateName.(string),
	}, nil
}
