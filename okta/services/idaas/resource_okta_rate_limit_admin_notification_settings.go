package idaas

import (
	"context"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &rateLimitAdminNotificationSettingsResource{}
	_ resource.ResourceWithConfigure   = &rateLimitAdminNotificationSettingsResource{}
	_ resource.ResourceWithImportState = &rateLimitAdminNotificationSettingsResource{}
)

var _ resource.Resource = &rateLimitAdminNotificationSettingsResource{}

type rateLimitAdminNotificationSettingsResource struct {
	*config.Config
}

type rateLimitAdminNotificationSettingsModel struct {
	Id                   types.String `tfsdk:"id"`
	NotificationsEnabled types.Bool   `tfsdk:"notifications_enabled"`
}

func newRateLimitAdminNotificationSettingsResource() resource.Resource {
	return &rateLimitAdminNotificationSettingsResource{}
}

func (r *rateLimitAdminNotificationSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rate_limit_admin_notification_settings"
}

func (r *rateLimitAdminNotificationSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *rateLimitAdminNotificationSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, tfpath.Root("id"), req, resp)
}

func (r *rateLimitAdminNotificationSettingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"notifications_enabled": schema.BoolAttribute{
				Required:    true,
				Description: "Enables or disables admin notifications for rate limiting events.",
			},
		},
	}
}

func (r *rateLimitAdminNotificationSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data rateLimitAdminNotificationSettingsModel

	// Read Terraform plan Data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	rateLimitAdminNotificationSettingsResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().RateLimitSettingsAPI.ReplaceRateLimitSettingsAdminNotifications(ctx).RateLimitAdminNotifications(buildPerClientRateLimitAdminNotifications(data)).Execute()
	if err != nil {
		return
	}

	data.Id = types.StringValue("rate_limiting_admin_notification")
	data.NotificationsEnabled = types.BoolValue(rateLimitAdminNotificationSettingsResp.GetNotificationsEnabled())

	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *rateLimitAdminNotificationSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data rateLimitAdminNotificationSettingsModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getRateLimitAdminNotificationSettingsResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().RateLimitSettingsAPI.GetRateLimitSettingsAdminNotifications(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to read rate limit admin notification settings",
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue("rate_limiting_admin_notification")
	data.NotificationsEnabled = types.BoolValue(getRateLimitAdminNotificationSettingsResp.GetNotificationsEnabled())

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *rateLimitAdminNotificationSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data rateLimitAdminNotificationSettingsModel

	// Read Terraform plan Data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	rateLimitAdminNotificationSettingsResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().RateLimitSettingsAPI.ReplaceRateLimitSettingsAdminNotifications(ctx).RateLimitAdminNotifications(buildPerClientRateLimitAdminNotifications(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update rate limit admin notification settings",
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue("rate_limiting_admin_notification")
	data.NotificationsEnabled = types.BoolValue(rateLimitAdminNotificationSettingsResp.GetNotificationsEnabled())

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *rateLimitAdminNotificationSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Delete Not Supported",
		"This resource cannot be deleted via Terraform.",
	)
}

func buildPerClientRateLimitAdminNotifications(data rateLimitAdminNotificationSettingsModel) v5okta.RateLimitAdminNotifications {

	rateLimitAdminNotificationSettings := v5okta.RateLimitAdminNotifications{
		NotificationsEnabled: data.NotificationsEnabled.ValueBool(),
	}

	return rateLimitAdminNotificationSettings
}
