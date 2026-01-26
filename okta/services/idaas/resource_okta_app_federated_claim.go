package idaas

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &appFederatedClaim{}
	_ resource.ResourceWithConfigure   = &appFederatedClaim{}
	_ resource.ResourceWithImportState = &appFederatedClaim{}
)

var _ resource.Resource = &appFederatedClaim{}

type appFederatedClaim struct {
	*config.Config
}

func newAppFederatedClaimResource() resource.Resource {
	return &appFederatedClaim{}
}

func (r *appFederatedClaim) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *appFederatedClaim) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Expected format: resource_id/sequence_id",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("app_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

func (r *appFederatedClaim) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_federated_claim"
}

func (r *appFederatedClaim) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "`id` used to specify the app feature ID. Its a combination of `app_id` and `name` separated by a forward slash (/).",
			},
			"app_id": schema.StringAttribute{
				Required:    true,
				Description: "`app_id` used to specify the app ID.",
			},
			"expression": schema.StringAttribute{
				Required:    true,
				Description: "The Okta Expression Language expression to be evaluated at runtime.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the claim to be used in the produced token.",
			},
		},
	}
}

type appFederatedClaimModel struct {
	ID         types.String `tfsdk:"id"`
	AppID      types.String `tfsdk:"app_id"`
	Expression types.String `tfsdk:"expression"`
	Name       types.String `tfsdk:"name"`
}

func (r *appFederatedClaim) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data appFederatedClaimModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	createFederatedClaimResp, _, err := r.OktaIDaaSClient.OktaSDKClientV6().ApplicationSSOFederatedClaimsAPI.CreateFederatedClaim(ctx, data.AppID.ValueString()).FederatedClaimRequestBody(buildFederatedClaimRequestBody(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating app federated claim",
			"Could not create app federated claim, unexpected error: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(updateAppFederatedClaimToState(&data, createFederatedClaimResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func updateAppFederatedClaimToState(a *appFederatedClaimModel, resp *v6okta.FederatedClaim) diag.Diagnostics {
	var diags diag.Diagnostics
	a.ID = types.StringValue(resp.GetId())
	a.Expression = types.StringValue(resp.GetExpression())
	a.Name = types.StringValue(resp.GetName())
	return diags
}

func buildFederatedClaimRequestBody(data appFederatedClaimModel) v6okta.FederatedClaimRequestBody {
	return v6okta.FederatedClaimRequestBody{
		Expression: data.Expression.ValueStringPointer(),
		Name:       data.Name.ValueStringPointer(),
	}
}

func (r *appFederatedClaim) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data appFederatedClaimModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Read API call logic
	getAppFederatedClaimResp, _, err := r.OktaIDaaSClient.OktaSDKClientV6().ApplicationSSOFederatedClaimsAPI.GetFederatedClaim(ctx, data.AppID.ValueString(), data.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading app federated claim",
			"Could not read app federated claim, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(updateReadAppFederatedClaimToState(&data, getAppFederatedClaimResp)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func updateReadAppFederatedClaimToState(a *appFederatedClaimModel, resp *v6okta.FederatedClaimRequestBody) diag.Diagnostics {
	var diags diag.Diagnostics
	a.Expression = types.StringValue(resp.GetExpression())
	a.Name = types.StringValue(resp.GetName())
	return diags
}

func (r *appFederatedClaim) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, data appFederatedClaimModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Read API call logic
	updateAppFederatedClaimResp, _, err := r.OktaIDaaSClient.OktaSDKClientV6().ApplicationSSOFederatedClaimsAPI.ReplaceFederatedClaim(ctx, data.AppID.ValueString(), state.ID.ValueString()).FederatedClaim(updateFederatedClaimBody(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading app federated claim",
			"Could not read app federated claim, unexpected error: "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(updateAppFederatedClaimToState(&data, updateAppFederatedClaimResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func updateFederatedClaimBody(data appFederatedClaimModel) v6okta.FederatedClaim {
	return v6okta.FederatedClaim{
		Expression: data.Expression.ValueStringPointer(),
		Name:       data.Name.ValueStringPointer(),
	}
}

func (r *appFederatedClaim) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data appFederatedClaimModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.OktaIDaaSClient.OktaSDKClientV6().ApplicationSSOFederatedClaimsAPI.DeleteFederatedClaim(ctx, data.AppID.ValueString(), data.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting app federated claim",
			"Could not delete app federated claim, unexpected error: "+err.Error(),
		)
		return
	}
}
