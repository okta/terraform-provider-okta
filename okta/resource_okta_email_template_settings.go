package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-sdk-golang/v4/okta"
)

var (
	_ resource.Resource                = &emailTemplateSettingsResource{}
	_ resource.ResourceWithConfigure   = &emailTemplateSettingsResource{}
	_ resource.ResourceWithImportState = &emailTemplateSettingsResource{}
)

func NewEmailTemplateSettingsResource() resource.Resource {
	return &emailTemplateSettingsResource{}
}

type emailTemplateSettingsResource struct {
	*Config
}

type emailTemplateSettingsResourceModel struct {
	ID           types.String `tfsdk:"id"`
	BrandID      types.String `tfsdk:"brand_id"`
	TemplateName types.String `tfsdk:"template_name"`
	Recipients   types.String `tfsdk:"recipients"`
}

func (r *emailTemplateSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_email_template_settings"
}

func (r *emailTemplateSettingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Manages email template settings`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the resource. This is a compound ID of the brand ID and the template name.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"brand_id": schema.StringAttribute{
				Description: "The ID of the brand.",
				Required:    true,
			},
			"template_name": schema.StringAttribute{
				Description: "Email template name - Example values: `AccountLockout`,`ADForgotPassword`,`ADForgotPasswordDenied`,`ADSelfServiceUnlock`,`ADUserActivation`,`AuthenticatorEnrolled`,`AuthenticatorReset`,`ChangeEmailConfirmation`,`EmailChallenge`,`EmailChangeConfirmation`,`EmailFactorVerification`,`ForgotPassword`,`ForgotPasswordDenied`,`IGAReviewerEndNotification`,`IGAReviewerNotification`,`IGAReviewerPendingNotification`,`IGAReviewerReassigned`,`LDAPForgotPassword`,`LDAPForgotPasswordDenied`,`LDAPSelfServiceUnlock`,`LDAPUserActivation`,`MyAccountChangeConfirmation`,`NewSignOnNotification`,`OktaVerifyActivation`,`PasswordChanged`,`PasswordResetByAdmin`,`PendingEmailChange`,`RegistrationActivation`,`RegistrationEmailVerification`,`SelfServiceUnlock`,`SelfServiceUnlockOnUnlockedAccount`,`UserActivation`",
				Required:    true,
			},
			"recipients": schema.StringAttribute{
				Description: "The recipients the emails of this template will be sent to - Valid values: `ALL_USERS`, `ADMINS_ONLY`, `NO_USERS`",
				Required:    true,
			},
		},
	}
}

func (r *emailTemplateSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

// Unable to true create because there must always be a template setting
func (r *emailTemplateSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan emailTemplateSettingsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.put(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update email template settings",
			err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(formatId(plan.BrandID.ValueString(), plan.TemplateName.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *emailTemplateSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state emailTemplateSettingsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	emailSettings, _, err := r.oktaSDKClientV3.CustomizationAPI.GetEmailSettings(ctx, state.BrandID.ValueString(), state.TemplateName.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to read email template settings",
			err.Error(),
		)
		return
	}

	state.Recipients = types.StringValue(emailSettings.Recipients)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Noop for delete because there must always be a template setting
func (r *emailTemplateSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *emailTemplateSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan emailTemplateSettingsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.put(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update email template settings",
			err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(formatId(plan.BrandID.ValueString(), plan.TemplateName.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *emailTemplateSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *emailTemplateSettingsResource) put(ctx context.Context, plan emailTemplateSettingsResourceModel) error {
	emailSettings := buildEmailSettings(plan)
	_, err := r.oktaSDKClientV3.CustomizationAPI.ReplaceEmailSettings(ctx, plan.BrandID.ValueString(), plan.TemplateName.ValueString()).EmailSettings(emailSettings).Execute()
	return err
}

func formatId(brandID string, templateName string) string {
	return fmt.Sprintf("%s/%s", brandID, templateName)
}

func buildEmailSettings(model emailTemplateSettingsResourceModel) okta.EmailSettings {
	emailSettings := okta.EmailSettings{}
	emailSettings.SetRecipients(model.Recipients.ValueString())
	return emailSettings
}
