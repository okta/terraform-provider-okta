package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &featuresResource{}
	_ resource.ResourceWithConfigure   = &featuresResource{}
	_ resource.ResourceWithImportState = &featuresResource{}
)

func newFeaturesResource() resource.Resource {
	return &featuresResource{}
}

type featuresResource struct {
	*config.Config
}

type featureResourceModel struct {
	ID          types.String `tfsdk:"id"`
	FeatureID   types.String `tfsdk:"feature_id"`
	Mode        types.Bool   `tfsdk:"mode"`
	Name        types.String `tfsdk:"name"`
	Lifecycle   types.String `tfsdk:"life_cycle"`
	Status      types.String `tfsdk:"status"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	Stage       types.Object `tfsdk:"stage"`
}

func (r *featuresResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_feature"
}

func (r *featuresResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages feature",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the resource. This ID is simply the feature.",
				Computed:    true,
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.UseStateForUnknown(),
				// },
			},
			"feature_id": schema.StringAttribute{
				Description: "Okta API for feature only reads and updates therefore the okta_feature resource needs to act as a quasi data source. Do this by setting feature_id",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mode": schema.BoolAttribute{
				Description: "Indicates if you want to force enable or disable a feature. Value is `true` meaning force",
				Optional:    true,
			},
			"life_cycle": schema.StringAttribute{
				Description: "Whether to `ENABLE` or `DISABLE` the feature",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the feature.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Brief description of the feature and what it provides.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "The feature status",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of feature.",
				Computed:    true,
			},
			"stage": schema.ObjectAttribute{
				Description: "Current release cycle stage of a feature.",
				Computed:    true,
				AttributeTypes: map[string]attr.Type{
					"state": types.StringType,
					"value": types.StringType,
				},
			},
		},
	}
}

func (r *featuresResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *featuresResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Warn(ctx, "True create for an okta_features is a no-op as this resource cannot not create, only enabled. Additionally, enabled an okta_features will also enabled all of its dependent")
	var state featureResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	feature, _, err := r.OktaIDaaSClient.OktaSDKClientV5().FeatureAPI.GetFeature(ctx, state.FeatureID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to get feature", err.Error())
		return
	}
	var lifecycle string
	if !state.Lifecycle.IsNull() {
		lifecycle = state.Lifecycle.ValueString()
	} else {
		if feature.GetStatus() == "DISABLED" {
			lifecycle = "DISABLE"
		} else {
			lifecycle = "ENABLE"
		}
	}
	apiReq := r.OktaIDaaSClient.OktaSDKClientV5().FeatureAPI.UpdateFeatureLifecycle(ctx, state.FeatureID.ValueString(), lifecycle)
	if !state.Mode.IsNull() && state.Mode.ValueBool() {
		apiReq = apiReq.Mode("force")
	}

	updatedFeature, _, err := apiReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to update feature", err.Error())
		return
	}

	resp.Diagnostics.Append(mapFeatureToState(updatedFeature, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *featuresResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state featureResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	feature, _, err := r.OktaIDaaSClient.OktaSDKClientV5().FeatureAPI.GetFeature(ctx, state.FeatureID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to get feature", err.Error())
		return
	}
	resp.Diagnostics.Append(mapFeatureToState(feature, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *featuresResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state featureResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	feature, _, err := r.OktaIDaaSClient.OktaSDKClientV5().FeatureAPI.GetFeature(ctx, state.FeatureID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to get feature", err.Error())
		return
	}
	var lifecycle string
	if !state.Lifecycle.IsNull() {
		lifecycle = state.Lifecycle.ValueString()
	} else {
		if feature.GetStatus() == "DISABLED" {
			lifecycle = "DISABLE"
		} else {
			lifecycle = "ENABLE"
		}
	}
	apiReq := r.OktaIDaaSClient.OktaSDKClientV5().FeatureAPI.UpdateFeatureLifecycle(ctx, state.FeatureID.ValueString(), lifecycle)
	if !state.Mode.IsNull() && state.Mode.ValueBool() {
		apiReq = apiReq.Mode("force")
	}

	updatedFeature, _, err := apiReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to update feature", err.Error())
		return
	}

	resp.Diagnostics.Append(mapFeatureToState(updatedFeature, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *featuresResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Warn(ctx, "True delete for an okta_features is a no-op as this resource cannot not delete, only disabled. Additionally, disabled an okta_features will also disabled all of its dependencies")
	var state featureResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, _, err := r.OktaIDaaSClient.OktaSDKClientV5().FeatureAPI.UpdateFeatureLifecycle(ctx, state.FeatureID.ValueString(), "DISABLE").Mode("force").Execute()
	if err != nil {
		resp.Diagnostics.AddError("failed to disabled feature", err.Error())
		return
	}
}

func (r *featuresResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapFeatureToState(data *okta.Feature, state *featureResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	state.ID = types.StringPointerValue(data.Id)
	state.FeatureID = types.StringPointerValue(data.Id)
	state.Name = types.StringPointerValue(data.Name)
	state.Status = types.StringPointerValue(data.Status)
	state.Description = types.StringPointerValue(data.Description)
	state.Type = types.StringPointerValue(data.Type)
	if data.Stage != nil {
		featureStageValue := map[string]attr.Value{
			"state": types.StringPointerValue(data.GetStage().State),
			"value": types.StringPointerValue(data.GetStage().Value),
		}
		featureStageTypes := map[string]attr.Type{
			"state": types.StringType,
			"value": types.StringType,
		}
		featureStage, _ := types.ObjectValue(featureStageTypes, featureStageValue)
		state.Stage = featureStage
	}
	return diags
}
