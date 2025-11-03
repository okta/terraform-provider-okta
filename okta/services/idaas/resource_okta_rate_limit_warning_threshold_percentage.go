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
	_ resource.Resource                = &rateLimitWarningThresholdPercentage{}
	_ resource.ResourceWithConfigure   = &rateLimitWarningThresholdPercentage{}
	_ resource.ResourceWithImportState = &rateLimitWarningThresholdPercentage{}
)

var _ resource.Resource = &rateLimitWarningThresholdPercentage{}

type rateLimitWarningThresholdPercentage struct {
	*config.Config
}

type rateLimitWarningThresholdPercentageModel struct {
	Id               types.String `tfsdk:"id"`
	WarningThreshold types.Int32  `tfsdk:"warning_threshold"`
}

func newRateLimitWarningThresholdPercentageResource() resource.Resource {
	return &rateLimitWarningThresholdPercentage{}
}

func (r *rateLimitWarningThresholdPercentage) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rate_limit_warning_threshold_percentage"
}

func (r *rateLimitWarningThresholdPercentage) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *rateLimitWarningThresholdPercentage) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, tfpath.Root("id"), req, resp)
}

func (r *rateLimitWarningThresholdPercentage) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"warning_threshold": schema.Int32Attribute{
				Required:    true,
				Description: "The threshold value (percentage) of a rate limit that, when exceeded, triggers a warning notification. By default, this value is 90 for Workforce orgs and 60 for CIAM orgs.",
			},
		},
	}
}

func (r *rateLimitWarningThresholdPercentage) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data rateLimitWarningThresholdPercentageModel

	// Read Terraform plan Data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	rateLimitWarningThresholdPercentageResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().RateLimitSettingsAPI.ReplaceRateLimitSettingsWarningThreshold(ctx).RateLimitWarningThreshold(buildPerClientRateLimitWarningThresholdPercentage(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to read rate limit warning threshold percentage",
			err.Error(),
		)
		return
	}

	// Example Data value setting
	//data.Id = types.StringValue("example-id")
	data.Id = types.StringValue("rate_limiting_warning_threshold_percentage")
	data.WarningThreshold = types.Int32Value(rateLimitWarningThresholdPercentageResp.GetWarningThreshold())

	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *rateLimitWarningThresholdPercentage) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data rateLimitWarningThresholdPercentageModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getRateLimitWarningThresholdResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().RateLimitSettingsAPI.GetRateLimitSettingsWarningThreshold(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update rate limit warning threshold percentage",
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue("rate_limiting_warning_threshold_percentage")
	data.WarningThreshold = types.Int32Value(getRateLimitWarningThresholdResp.GetWarningThreshold())

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *rateLimitWarningThresholdPercentage) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data rateLimitWarningThresholdPercentageModel

	// Read Terraform plan Data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	rateLimitWarningThresholdPercentageResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().RateLimitSettingsAPI.ReplaceRateLimitSettingsWarningThreshold(ctx).RateLimitWarningThreshold(buildPerClientRateLimitWarningThresholdPercentage(data)).Execute()
	if err != nil {
		return
	}

	// Example Data value setting
	data.Id = types.StringValue("rate_limiting_warning_threshold_percentage")
	data.WarningThreshold = types.Int32Value(rateLimitWarningThresholdPercentageResp.GetWarningThreshold())

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *rateLimitWarningThresholdPercentage) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Delete Not Supported",
		"This resource cannot be deleted via Terraform.",
	)
}

func buildPerClientRateLimitWarningThresholdPercentage(data rateLimitWarningThresholdPercentageModel) v5okta.RateLimitWarningThresholdRequest {

	rateLimitAdminNotificationSettings := v5okta.RateLimitWarningThresholdRequest{
		WarningThreshold: data.WarningThreshold.ValueInt32(),
	}

	return rateLimitAdminNotificationSettings
}
