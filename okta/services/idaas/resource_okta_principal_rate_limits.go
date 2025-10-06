package idaas

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &principalRateLimits{}
	_ resource.ResourceWithConfigure   = &principalRateLimits{}
	_ resource.ResourceWithImportState = &principalRateLimits{}
)

var _ resource.Resource = &principalRateLimits{}

type principalRateLimits struct {
	*config.Config
}

type principalRateLimitsModel struct {
	Id                           types.String `tfsdk:"id"`
	PrincipalId                  types.String `tfsdk:"principal_id"`
	PrincipalType                types.String `tfsdk:"principal_type"`
	DefaultConcurrencyPercentage types.Int32  `tfsdk:"default_concurrency_percentage"`
	DefaultPercentage            types.Int32  `tfsdk:"default_percentage"`
	CreatedBy                    types.String `tfsdk:"created_by"`
	CreatedDate                  types.String `tfsdk:"created_date"`
	LastUpdate                   types.String `tfsdk:"last_update"`
	LastUpdatedBy                types.String `tfsdk:"last_updated_by"`
	OrgId                        types.String `tfsdk:"org_id"`
}

func newPrincipalRateLimitsResource() resource.Resource {
	return &principalRateLimits{}
}

func (r *principalRateLimits) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_principal_rate_limits"
}

func (r *principalRateLimits) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *principalRateLimits) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, tfpath.Root("id"), req, resp)
}

func (r *principalRateLimits) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The unique identifier of the principle rate limit entity.",
			},
			"principal_id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of the principal. This is the ID of the API token or OAuth 2.0 app.",
			},
			"principal_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"OAUTH_CLIENT",
						"SSWS_TOKEN",
					}...),
				},
				Description: "The type of principal, either an API token or an OAuth 2.0 app.",
			},
			"default_concurrency_percentage": schema.Int32Attribute{
				Computed:    true,
				Optional:    true,
				Description: "The default percentage of a given concurrency limit threshold that the owning principal can consume.",
			},
			"default_percentage": schema.Int32Attribute{
				Computed:    true,
				Optional:    true,
				Description: "The default percentage of a given rate limit threshold that the owning principal can consume.",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "The Okta user ID of the user who created the principle rate limit entity.",
			},
			"created_date": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time the principle rate limit entity was created.",
			},
			"last_update": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time the principle rate limit entity was last updated.",
			},
			"last_updated_by": schema.StringAttribute{
				Computed:    true,
				Description: "The Okta user ID of the user who last updated the principle rate limit entity.",
			},
			"org_id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the Okta org.",
			},
		},
	}
}

func (r *principalRateLimits) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError(
		"Create Not Supported",
		"This resource cannot be created via Terraform.",
	)
}

func (r *principalRateLimits) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data principalRateLimitsModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getPrincipalRateSettingsResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().PrincipalRateLimitAPI.GetPrincipalRateLimitEntity(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to read principal rate limit",
			err.Error(),
		)
		return
	}

	applyPrincipalRateSettingsToState(&data, getPrincipalRateSettingsResp)

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *principalRateLimits) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data principalRateLimitsModel

	// Read Terraform plan Data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	updatePrincipalRateSettingsRespSettingsResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().PrincipalRateLimitAPI.ReplacePrincipalRateLimitEntity(ctx, data.Id.ValueString()).Entity(buildPrincipalRateLimits(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update principal rate limit",
			err.Error(),
		)
		return
	}

	applyPrincipalRateSettingsToState(&data, updatePrincipalRateSettingsRespSettingsResp)

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *principalRateLimits) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Delete Not Supported",
		"This resource cannot be deleted via Terraform.",
	)
}

func buildPrincipalRateLimits(data principalRateLimitsModel) v5okta.PrincipalRateLimitEntity {

	principalRateSettings := v5okta.PrincipalRateLimitEntity{
		PrincipalId:   data.PrincipalId.ValueString(),
		PrincipalType: data.PrincipalType.ValueString(),
	}

	if data.DefaultConcurrencyPercentage.ValueInt32Pointer() != nil {
		principalRateSettings.DefaultConcurrencyPercentage = data.DefaultConcurrencyPercentage.ValueInt32Pointer()
	}

	if data.DefaultPercentage.ValueInt32Pointer() != nil {
		principalRateSettings.DefaultPercentage = data.DefaultPercentage.ValueInt32Pointer()
	}

	return principalRateSettings
}

func applyPrincipalRateSettingsToState(data *principalRateLimitsModel, principalRateLimitSettingsResp *v5okta.PrincipalRateLimitEntity) {
	data.Id = types.StringValue(principalRateLimitSettingsResp.GetId())
	data.PrincipalId = types.StringValue(principalRateLimitSettingsResp.GetPrincipalId())
	data.PrincipalType = types.StringValue(principalRateLimitSettingsResp.GetPrincipalType())
	data.DefaultConcurrencyPercentage = types.Int32Value(principalRateLimitSettingsResp.GetDefaultConcurrencyPercentage())
	data.DefaultPercentage = types.Int32Value(principalRateLimitSettingsResp.GetDefaultPercentage())
	data.CreatedBy = types.StringValue(principalRateLimitSettingsResp.GetCreatedBy())
	data.CreatedDate = types.StringValue(principalRateLimitSettingsResp.GetCreatedDate().Format(time.RFC3339))
	data.LastUpdate = types.StringValue(principalRateLimitSettingsResp.GetLastUpdate().Format(time.RFC3339))
	data.LastUpdatedBy = types.StringValue(principalRateLimitSettingsResp.GetLastUpdatedBy())
	data.OrgId = types.StringValue(principalRateLimitSettingsResp.GetOrgId())
}
