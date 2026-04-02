package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &rateLimitAdminNotificationSettingsDataSource{}

func newRateLimitAdminNotificationSettingsDataSource() datasource.DataSource {
	return &rateLimitAdminNotificationSettingsDataSource{}
}

func (d *rateLimitAdminNotificationSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

type rateLimitAdminNotificationSettingsDataSource struct {
	*config.Config
}

func (d *rateLimitAdminNotificationSettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rate_limit_admin_notification_settings"
}

func (d *rateLimitAdminNotificationSettingsDataSource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, tfpath.Root("id"), req, resp)
}

func (d *rateLimitAdminNotificationSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"notifications_enabled": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

type rateLimitAdminNotificationSettingsDataSourceModel struct {
	Id                   types.String `tfsdk:"id"`
	NotificationsEnabled types.Bool   `tfsdk:"notifications_enabled"`
}

func (d *rateLimitAdminNotificationSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data rateLimitAdminNotificationSettingsDataSourceModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getRateLimitAdminNotificationSettingsResp, _, err := d.OktaIDaaSClient.OktaSDKClientV5().RateLimitSettingsAPI.GetRateLimitSettingsAdminNotifications(ctx).Execute()
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
