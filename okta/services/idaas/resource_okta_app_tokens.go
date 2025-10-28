package idaas

import (
	"context"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &appTokens{}
	_ resource.ResourceWithConfigure   = &appTokens{}
	_ resource.ResourceWithImportState = &appTokens{}
)

var _ resource.Resource = &appTokens{}

type appTokens struct {
	*config.Config
}

func newAppTokensResource() resource.Resource {
	return &appTokens{}
}

func (r *appTokens) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rate_limit_warning_threshold_percentage"
}

func (r *appTokens) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *appTokens) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, tfpath.Root("id"), req, resp)
}

func (r *appTokens) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"token_id": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *appTokens) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data rateLimitWarningThresholdPercentageModel

	// Read Terraform plan Data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	rateLimitWarningThresholdPercentageResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().RateLimitSettingsAPI.ReplaceRateLimitSettingsWarningThreshold(ctx).RateLimitWarningThreshold(buildPerClientRateLimitWarningThresholdPercentage(data)).Execute()
	if err != nil {
		return
	}

	// Example Data value setting
	//data.Id = types.StringValue("example-id")
	data.Id = types.StringValue("rate_limiting_warning_threshold_percentage")
	data.WarningThreshold = types.Int32Value(rateLimitWarningThresholdPercentageResp.GetWarningThreshold())

	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *appTokens) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data rateLimitWarningThresholdPercentageModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getRateLimitWarningThresholdResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().RateLimitSettingsAPI.GetRateLimitSettingsWarningThreshold(ctx).Execute()
	if err != nil {
		return
	}

	data.Id = types.StringValue("rate_limiting_warning_threshold_percentage")
	data.WarningThreshold = types.Int32Value(getRateLimitWarningThresholdResp.GetWarningThreshold())

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *appTokens) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

func (r *appTokens) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Delete Not Supported",
		"This resource cannot be deleted via Terraform.",
	)
}
